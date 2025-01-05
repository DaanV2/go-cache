package large

import (
	"fmt"
	"iter"
	"runtime"
	"sync"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/pkg/iterators"
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-kit/generics"
)

// BuckettedSet is a set of items, that uses a pre-defined amount of buckets, each item generates an hash, from which a bucket can be specified
type BuckettedMap[K, V comparable] struct {
	hasher hash.Hasher[K]
	sets   []*GrowableMap[K, V]
	base   Options
}

// NewBuckettedMap creates a new BuckettedMap with the specified capacity, hasher, and options.
// The BuckettedMap is a concurrent map that uses a bucketing strategy to reduce contention.
func NewBuckettedMap[K, V comparable](capacity uint64, keyhasher hash.Hasher[K], opts ...options.Option[Options]) (*BuckettedMap[K, V], error) {
	base, err := CreateOptions[collections.HashItem[collections.KeyValue[K, V]]](opts...)
	if err != nil {
		return nil, err
	}

	buckets := base.bucket_amount
	if buckets == 0 {
		buckets = base.BucketAmount(capacity)
	}

	set := &BuckettedMap[K, V]{
		hasher: keyhasher,
		sets:   make([]*GrowableMap[K, V], 0, buckets),
		base:   base,
	}

	for range buckets {
		s, err := NewGrowableMapFrom[K, V](keyhasher, base)
		if err != nil {
			return nil, err
		}

		set.sets = append(set.sets, s)
	}

	return set, nil
}

// Get retrieves the value for the specified key from the BuckettedMap.
func (m *BuckettedMap[K, V]) Get(key K) (collections.KeyValue[K, V], bool) {
	kv := collections.NewKeyValue(key, generics.Empty[V]())
	setitem := collections.NewHashItem(m.hasher.Hash(key), kv)
	bucket := m.bucketIndex(setitem)
	v, ok := m.sets[bucket].Find(setitem)
	if ok {
		return v.Value(), true
	}

	return collections.EmptyKeyValue[K, V](), false
}

// Set will add or update the value for the specified key in the BuckettedMap. It returns true if the value was added, false if it was updated.
func (m *BuckettedMap[K, V]) Set(key K, item V) bool {
	kv := collections.NewKeyValue(key, item)
	setitem := collections.NewHashItem(m.hasher.Hash(key), kv)
	bucket := m.bucketIndex(setitem)
	return m.sets[bucket].updateOrAdd(setitem)
}

// Append adds all items from the specified Rangeable to the BuckettedMap.
func (m *BuckettedMap[K, V]) Append(other collections.Rangeable[collections.KeyValue[K, V]]) {
	other.Range(func(item collections.KeyValue[K, V]) bool {
		m.Set(item.Key(), item.Value())
		return true
	})
}

// AppendParralel adds all items from the specified ParralelRangeable to the BuckettedMap.
func (m *BuckettedMap[K, V]) AppendParralel(other collections.ParralelRangeable[collections.KeyValue[K, V]]) {
	other.RangeParralel(func(item collections.KeyValue[K, V]) bool {
		m.Set(item.Key(), item.Value())
		return true
	})
}

// bucketIndex returns the index of the bucket that the item should be placed in
func (s *BuckettedMap[K, V]) bucketIndex(item collections.HashItem[collections.KeyValue[K, V]]) uint64 {
	return item.Hash() % uint64(len(s.sets))
}

// Read will return a sequence of all items in the set
func (s *BuckettedMap[K, V]) Read() iter.Seq[collections.HashItem[collections.KeyValue[K, V]]] {
	return func(yield func(collections.HashItem[collections.KeyValue[K, V]]) bool) {
		for _, b := range s.sets {
			for item := range b.Read() {
				if !yield(item) {
					return
				}
			}
		}
	}
}

// Keys will return a sequence of all items in the set
func (s *BuckettedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, b := range s.sets {
			for item := range b.Read() {
				if !yield(item.Value().Key()) {
					return
				}
			}
		}
	}
}

// Values will return a sequence of all items in the set
func (s *BuckettedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, b := range s.sets {
			for item := range b.Read() {
				if !yield(item.Value().Value()) {
					return
				}
			}
		}
	}
}

// KeyValues will return a sequence of all items in the set
func (s *BuckettedMap[K, V]) KeyValues() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, b := range s.sets {
			for item := range b.Read() {
				v := item.Value()
				if !yield(v.Key(), v.Value()) {
					return
				}
			}
		}
	}
}

// Range will iterate over all items in the set
func (s *BuckettedMap[K, V]) Range(yield func(item collections.HashItem[collections.KeyValue[K, V]]) bool) {
	iterators.RangeCol(s, yield)
}

// RangeParralel will iterate over all items in the set in parallel
func (s *BuckettedMap[K, V]) RangeParralel(yield func(item collections.HashItem[collections.KeyValue[K, V]]) bool) {
	iterators.RangeColParralel(s.sets, yield)
}

func (s *BuckettedMap[K, V]) String() string {
	return fmt.Sprintf("large.BuckettedSet[%s]", generics.NameOf[collections.HashItem[collections.KeyValue[K, V]]]())
}

func (s *BuckettedMap[K, V]) GoString() string {
	return s.String()
}

// Grow will increase the capacity of the set
func (m *BuckettedMap[K, V]) Grow(new_capacity uint64) {
	buckets := m.base.BucketAmount(new_capacity)
	current := uint64(len(m.sets))
	if current >= buckets {
		return
	}

	diff := buckets - current
	// Add the new buckets
	for range diff {
		s, err := NewGrowableMap[K, V](m.hasher)
		if err != nil {
			return
		}

		m.sets = append(m.sets, s)
	}

	loops := runtime.GOMAXPROCS(0)
	wg := &sync.WaitGroup{}
	setsCh := make(chan *GrowableMap[K, V], loops)

	// Start the workers
	for i := 0; i < loops; i++ {
		wg.Add(1)
		go workerMapGrow(wg, setsCh, m)
	}

	// Remove the old buckets and rehash the items
	for i := range current {
		s := m.sets[i]
		news, err := NewGrowableMap[K, V](m.hasher)
		if err != nil {
			return
		}
		m.sets[i] = news

		// Add the items to the new bucket
		setsCh <- s
	}

	// Close the channel and wait for the workers to finish
	close(setsCh)
	wg.Wait()
}

func workerMapGrow[K, V comparable](wg *sync.WaitGroup, process <-chan *GrowableMap[K, V], receiver *BuckettedMap[K, V]) {
	defer wg.Done()
	for s := range process {
		s.Range(func(item collections.HashItem[collections.KeyValue[K, V]]) bool {
			_ = receiver.Set(item.Value().Key(), item.Value().Value())
			return true
		})
	}
}
