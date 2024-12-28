package large

import (
	"errors"
	"iter"

	"github.com/daanv2/go-cache/fixed"
	"github.com/daanv2/go-cache/pkg/constraints"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/pkg/iterators"
	"github.com/daanv2/go-cache/pkg/options"
)

// GrowableSet is a set that grows as needed.
type GrowableSet[T constraints.Equivalent[T]] struct {
	SetBase
	hasher  hash.Hasher[T]
	buckets []*fixed.Slice[SetItem[T]]
}

// NewGrowableSet creates a new instance of GrowableSet with the provided hasher and options.
// The hasher is used to hash the elements in the set, and options can be used to configure the set.
func NewGrowableSet[T constraints.Equivalent[T]](hasher hash.Hasher[T], opts ...options.Option[SetBase]) (*GrowableSet[T], error) {
	base := NewSetBase[T]()
	err := options.Apply(&base, opts...)
	if err != nil {
		return nil, err
	}

	// Validate
	if hasher == nil {
		return nil, errors.New("hasher is nil")
	}
	if base.bucket_size <= 1 {
		return nil, errors.New("bucket size is too small <= 1")
	}

	set := &GrowableSet[T]{
		SetBase: base,
		hasher:  hasher,
		buckets: make([]*fixed.Slice[SetItem[T]], 0),
	}

	return set, err
}

// GetOrAdd returns the item if it exists in the set, otherwise it adds it and returns it.
func (s *GrowableSet[T]) GetOrAdd(item T) (T, bool) {
	setitem := NewSetItem[T](item, s.hasher.Hash(item))

	return s.getOrAdd(setitem)
}

// UpdateOrAdd updates the item if it exists in the set, otherwise it adds it. Returns true if it had to add it instead of update.
func (s *GrowableSet[T]) UpdateOrAdd(item T) bool {
	setitem := NewSetItem[T](item, s.hasher.Hash(item))

	return s.updateOrAdd(setitem)
}

func (s *GrowableSet[T]) getOrAdd(item SetItem[T]) (T, bool) {
	item_lock := s.items_lock.GetLock(item.hash)

	item_lock.Lock()
	defer item_lock.Unlock()

	// Find it
	v, ok := s.find(item)
	if ok {
		return v.Value(), false
	}

	s.set(item)
	return item.Value(), true
}

// updateOrAdd TODO. return true if it had to add it instead of update
func (s *GrowableSet[T]) updateOrAdd(item SetItem[T]) bool {
	item_lock := s.items_lock.GetLock(item.hash)

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

func (s *GrowableSet[T]) find(item SetItem[T]) (SetItem[T], bool) {
	s.bucket_lock.RLock()
	defer s.bucket_lock.RUnlock()

	// Try to find it
	for i := range s.buckets {
		v, ok := s.buckets[i].Find(item.Equals)
		if ok {
			return v, true
		}
	}

	return item, false
}

func (s *GrowableSet[T]) updateIf(item SetItem[T]) bool {
	s.bucket_lock.RLock()
	defer s.bucket_lock.RUnlock()

	// Try to find it
	for i := range s.buckets {
		index, ok := s.buckets[i].FindIndex(item.Equals)
		if ok {
			_ = s.buckets[i].Set(index, item)
			return true
		}
	}

	return false
}

func (s *GrowableSet[T]) set(item SetItem[T]) {
	s.bucket_lock.Lock()
	defer s.bucket_lock.Unlock()

	l := len(s.buckets)
	for i := l - 1; i > 0; i-- {
		if !s.buckets[i].IsFull() {
			if s.buckets[i].TryAppend(item) > 0 {
				return
			}
		}
	}

	for {
		b := fixed.NewSlice[SetItem[T]](s.SetBase.bucket_size)
		s.buckets = append(s.buckets, &b)
		if s.buckets[len(s.buckets)-1].TryAppend(item) > 0 {
			return
		}
	}
}

// Read returns an iterator that reads the items in the set.
func (s *GrowableSet[T]) Read() iter.Seq[T] {
	return func(yield func(T) bool) {
		s.bucket_lock.RLock()
		defer s.bucket_lock.RUnlock()

		for _, bucket := range s.buckets {
			for v := range bucket.Read() {
				if !yield(v.Value()) {
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
