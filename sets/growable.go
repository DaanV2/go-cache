package sets

import (
	"errors"
	"fmt"
	"iter"
	"sync"

	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/pkg/iterators"
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-kit/generics"
	optimal "github.com/daanv2/go-optimal"
)

// GrowableSet is a set that grows as needed.
type GrowableSet[T comparable] struct {
	Options
	hasher      hash.Hasher[T]
	buckets     []*Fixed[T]
	bucket_lock sync.RWMutex
}

// NewGrowableSet creates a new instance of GrowableSet with the provided hasher and options.
// The hasher is used to hash the elements in the set, and options can be used to configure the set.
func NewGrowableSet[T comparable](hasher hash.Hasher[T], opts ...options.Option[Options]) (*GrowableSet[T], error) {
	o := []options.Option[Options]{
		WithBucketSize(uint64(optimal.SliceSize[T]())),
	}
	o = append(o, opts...)

	base, err := CreateOptions[T](o...)
	if err != nil {
		return nil, err
	}

	return NewGrowableSetFrom(hasher, base)
}

func NewGrowableSetFrom[T comparable](hasher hash.Hasher[T], base Options) (*GrowableSet[T], error) {
	// Validate
	if hasher == nil {
		return nil, errors.New("hasher is nil")
	}
	if base.bucket_size <= 1 {
		return nil, errors.New("bucket size is too small <= 1")
	}

	return &GrowableSet[T]{
		Options:     base,
		hasher:      hasher,
		buckets:     make([]*Fixed[T], 0),
		bucket_lock: sync.RWMutex{},
	}, nil
}

// GetOrAdd returns the item if it exists in the set, otherwise it adds it and returns it.
func (s *GrowableSet[T]) GetOrAdd(item T) (T, bool) {
	setitem := NewSetItem[T](s.hasher.Hash(item), item)

	return s.getOrAdd(setitem)
}

// UpdateOrAdd updates the item if it exists in the set, otherwise it adds it. Returns true if it had to add it instead of update.
func (s *GrowableSet[T]) UpdateOrAdd(item T) bool {
	setitem := NewSetItem[T](s.hasher.Hash(item), item)

	return s.updateOrAdd(setitem)
}

func (s *GrowableSet[T]) getOrAdd(item SetItem[T]) (T, bool) {
	item_lock := s.items_lock.GetLock(item.Hash)

	item_lock.Lock()
	defer item_lock.Unlock()

	l := len(s.buckets)
	switch l {
	case 0:
		s.set(item)
		return item.Value, true
	case 1:
		if s.buckets[0].Set(item) {
			return item.Value, true
		}

	default:
		// Find it
		v, ok := s.Find(item)
		if ok {
			return v.Value, false
		}
	}

	s.set(item)
	return item.Value, true
}

// updateOrAdd TODO. return true if it had to add it instead of update
func (s *GrowableSet[T]) updateOrAdd(item SetItem[T]) bool {
	item_lock := s.items_lock.GetLock(item.Hash)

	item_lock.Lock()
	defer item_lock.Unlock()

	// Find it
	ok := s.updateIf(item)
	if ok {
		return false
	}

	s.set(item)
	return true
}

func (s *GrowableSet[T]) updateIf(item SetItem[T]) bool {
	s.bucket_lock.RLock()
	defer s.bucket_lock.RUnlock()

	l := len(s.buckets)
	switch l {
	case 0:
		break
	case 1:
		return s.buckets[0].Set(item)
	default:
		// Try to find it
		for _, bucket := range s.buckets {
			if !bucket.HasHash(item.Hash) {
				continue
			}

			ok := bucket.Update(item)
			if ok {
				return true
			}
		}
	}

	return false
}

func (s *GrowableSet[T]) set(item SetItem[T]) {
	s.bucket_lock.Lock()
	defer s.bucket_lock.Unlock()

	l := len(s.buckets) - 1
	if l >= 0 {
		if s.buckets[l].Set(item) {
			return
		}
	}

	for {
		b := NewFixed[T](s.Options.bucket_size)
		s.buckets = append(s.buckets, &b)
		if s.buckets[len(s.buckets)-1].Set(item) {
			return
		}
	}
}

func (s *GrowableSet[T]) Find(item SetItem[T]) (SetItem[T], bool) {
	s.bucket_lock.RLock()
	defer s.bucket_lock.RUnlock()

	// Try to find it
	for _, bucket := range s.buckets {
		if !bucket.HasHash(item.Hash) {
			continue
		}

		v, ok := bucket.Get(item)
		if ok {
			return v, true
		}
	}

	return item, false
}

// Read returns an iterator that reads the items in the set.
func (s *GrowableSet[T]) Read() iter.Seq[T] {
	return func(yield func(T) bool) {
		s.bucket_lock.RLock()
		defer s.bucket_lock.RUnlock()

		for _, bucket := range s.buckets {
			for v := range bucket.Read() {
				if !yield(v.Value) {
					return
				}
			}
		}
	}
}

// Range calls the yield function for each item in the set.
func (s *GrowableSet[T]) Range(yield func(item T) bool) {
	iterators.RangeCol(s, yield)
}

func (s *GrowableSet[T]) String() string {
	return fmt.Sprintf("large.GrowableSet[%s,%v]", generics.NameOf[T](), len(s.buckets))
}

func (s *GrowableSet[T]) GoString() string {
	return s.String()
}
