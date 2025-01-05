package large

import (
	"errors"
	"fmt"
	"iter"
	"sync"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/fixed"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/pkg/iterators"
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-kit/generics"
	optimal "github.com/daanv2/go-optimal"
)

// GrowableMap is a set that grows as needed.
type GrowableMap[K, V comparable] struct {
	Options
	hasher      hash.Hasher[K]
	buckets     []*fixed.Map[K, V]
	bucket_lock sync.RWMutex
}

// NewGrowableMap creates a new instance of GrowableMap with the provided hasher and options.
// The hasher is used to hash the elements in the set, and options can be used to configure the set.
func NewGrowableMap[K, V comparable](hasher hash.Hasher[K], opts ...options.Option[Options]) (*GrowableMap[K, V], error) {
	o := []options.Option[Options]{
		WithBucketSize(uint64(optimal.SliceSize[V]())),
	}
	o = append(o, opts...)

	base, err := CreateOptions[V](o...)
	if err != nil {
		return nil, err
	}

	return NewGrowableMapFrom[K, V](hasher, base)
}

func NewGrowableMapFrom[K, V comparable](hasher hash.Hasher[K], base Options) (*GrowableMap[K, V], error) {
	// Validate
	if hasher == nil {
		return nil, errors.New("hasher is nil")
	}
	if base.bucket_size <= 1 {
		return nil, errors.New("bucket size is too small <= 1")
	}

	return &GrowableMap[K, V]{
		Options:     base,
		hasher:      hasher,
		buckets:     make([]*fixed.Map[K, V], 0),
		bucket_lock: sync.RWMutex{},
	}, nil
}

// GetOrAdd returns the item if it exists in the set, otherwise it adds it and returns it.
func (s *GrowableMap[K, V]) GetOrAdd(item collections.KeyValue[K, V]) (collections.HashItem[collections.KeyValue[K, V]], bool) {
	setitem := collections.NewHashItem[collections.KeyValue[K, V]](s.hasher.Hash(item.Key()), item)

	return s.getOrAdd(setitem)
}

// UpdateOrAdd updates the item if it exists in the set, otherwise it adds it. Returns true if it had to add it instead of update.
func (s *GrowableMap[K, V]) UpdateOrAdd(item collections.KeyValue[K, V]) bool {
	setitem := collections.NewHashItem[collections.KeyValue[K, V]](s.hasher.Hash(item.Key()), item)

	return s.updateOrAdd(setitem)
}

func (s *GrowableMap[K, V]) getOrAdd(item collections.HashItem[collections.KeyValue[K, V]]) (collections.HashItem[collections.KeyValue[K, V]], bool) {
	item_lock := s.items_lock.GetLock(item.Hash())

	item_lock.Lock()
	defer item_lock.Unlock()

	// Find it
	v, ok := s.Find(item)
	if ok {
		return v, false
	}

	s.set(item)
	return item, true
}

// updateOrAdd TODO. return true if it had to add it instead of update
func (s *GrowableMap[K, V]) updateOrAdd(item collections.HashItem[collections.KeyValue[K, V]]) bool {
	item_lock := s.items_lock.GetLock(item.Hash())

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

func (s *GrowableMap[K, V]) updateIf(item collections.HashItem[collections.KeyValue[K, V]]) bool {
	s.bucket_lock.RLock()
	defer s.bucket_lock.RUnlock()

	// Try to find it
	for i := range s.buckets {
		if !s.buckets[i].HasHash(item.Hash()) {
			continue
		}

		ok := s.buckets[i].Update(item)
		if ok {
			return true
		}
	}

	return false
}

func (s *GrowableMap[K, V]) set(item collections.HashItem[collections.KeyValue[K, V]]) {
	s.bucket_lock.Lock()
	defer s.bucket_lock.Unlock()

	l := len(s.buckets) - 1
	if l >= 0 {
		if s.buckets[l].Set(item) {
			return
		}
	}

	for {
		b := fixed.NewMap[K, V](s.Options.bucket_size)
		s.buckets = append(s.buckets, &b)
		if s.buckets[len(s.buckets)-1].Set(item) {
			return
		}
	}
}

func (s *GrowableMap[K, V]) Find(item collections.HashItem[collections.KeyValue[K, V]]) (collections.HashItem[collections.KeyValue[K, V]], bool) {
	s.bucket_lock.RLock()
	defer s.bucket_lock.RUnlock()

	// Try to find it
	for i := range s.buckets {
		v, ok := s.buckets[i].Get(item)
		if ok {
			return v, true
		}
	}

	return item, false
}

// Read returns an iterator that reads the items in the set.
func (s *GrowableMap[K, V]) Read() iter.Seq[collections.HashItem[collections.KeyValue[K, V]]] {
	return func(yield func(collections.HashItem[collections.KeyValue[K, V]]) bool) {
		s.bucket_lock.RLock()
		defer s.bucket_lock.RUnlock()

		for _, bucket := range s.buckets {
			for v := range bucket.Read() {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Range calls the yield function for each item in the set.
func (s *GrowableMap[K, V]) Range(yield func(item collections.HashItem[collections.KeyValue[K, V]]) bool) {
	iterators.RangeCol(s, yield)
}

func (s *GrowableMap[K, V]) String() string {
	return fmt.Sprintf("large.GrowableMap[%s,%s,%v]", generics.NameOf[K](), generics.NameOf[V](), len(s.buckets))
}

func (s *GrowableMap[K, V]) GoString() string {
	return s.String()
}
