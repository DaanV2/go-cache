package large

import (
	"fmt"
	"iter"
	"runtime"
	"sync"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/pkg/constraints"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/pkg/iterators"
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-kit/generics"
)

// BuckettedSet is a set of items, that uses a pre-defined amount of buckets, each item generates an hash, from which a bucket can be specified
type BuckettedSet[T constraints.Equivalent[T]] struct {
	hasher hash.Hasher[T]
	sets   []*GrowableSet[T]
	base   Options
}

// NewBuckettedSet creates a new BuckettedSet with the specified capacity, hasher, and options.
func NewBuckettedSet[T constraints.Equivalent[T]](capacity uint64, hasher hash.Hasher[T], opts ...options.Option[Options]) (*BuckettedSet[T], error) {
	base, err := CreateOptions[T](opts...)
	if err != nil {
		return nil, err
	}

	buckets := base.bucket_amount
	if buckets == 0 {
		buckets = base.BucketAmount(capacity)
	}

	set := &BuckettedSet[T]{
		hasher: hasher,
		sets:   make([]*GrowableSet[T], 0, buckets),
		base:   base,
	}

	for range buckets {
		s, err := NewGrowableSetFrom(hasher, base)
		if err != nil {
			return nil, err
		}

		set.sets = append(set.sets, s)
	}

	return set, nil
}

// GetOrAdd will return the item if it exists, otherwise it will add the item to the set
func (s *BuckettedSet[T]) GetOrAdd(item T) (T, bool) {
	setitem := collections.NewHashItem[T](s.hasher.Hash(item), item)
	bucket := s.bucketIndex(setitem)
	return s.sets[bucket].getOrAdd(setitem)
}

// UpdateOrAdd will update the item if it exists, otherwise it will add the item to the set, and return true if it had to add it
func (s *BuckettedSet[T]) UpdateOrAdd(item T) bool {
	setitem := collections.NewHashItem[T](s.hasher.Hash(item), item)
	return s.updateOrAdd(setitem)
}

func (s *BuckettedSet[T]) updateOrAdd(item collections.HashItem[T]) bool {
	bucket := s.bucketIndex(item)
	return s.sets[bucket].updateOrAdd(item)
}

// bucketIndex returns the index of the bucket that the item should be placed in
func (s *BuckettedSet[T]) bucketIndex(item collections.HashItem[T]) uint64 {
	return item.Hash() % uint64(len(s.sets))
}

// Read will return a sequence of all items in the set
func (s *BuckettedSet[T]) Read() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, b := range s.sets {
			for item := range b.Read() {
				if !yield(item) {
					return
				}
			}
		}
	}
}

// Range will iterate over all items in the set
func (s *BuckettedSet[T]) Range(yield func(item T) bool) {
	iterators.RangeCol(s, yield)
}

// RangeParralel will iterate over all items in the set in parallel
func (s *BuckettedSet[T]) RangeParralel(yield func(item T) bool) {
	iterators.RangeColParralel(s.sets, yield)
}

func (s *BuckettedSet[T]) String() string {
	return fmt.Sprintf("large.BuckettedSet[%s]", generics.NameOf[T]())
}

func (s *BuckettedSet[T]) GoString() string {
	return s.String()
}

// Grow will increase the capacity of the set
func (m *BuckettedSet[T]) Grow(new_capacity uint64) {
	buckets := m.base.BucketAmount(new_capacity)
	current := uint64(len(m.sets))
	if current >= buckets {
		return
	}

	diff := buckets - current
	// Add the new buckets
	for range diff {
		s, err := NewGrowableSet(m.hasher)
		if err != nil {
			return
		}

		m.sets = append(m.sets, s)
	}

	loops := runtime.GOMAXPROCS(0)
	wg := &sync.WaitGroup{}
	setsCh := make(chan *GrowableSet[T], loops)

	// Start the workers
	for i := 0; i < loops; i++ {
		wg.Add(1)
		go workerGrow(wg, setsCh, m)
	}

	// Remove the old buckets and rehash the items
	for i := range current {
		s := m.sets[i]
		news, err := NewGrowableSet(m.hasher)
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

func workerGrow[T constraints.Equivalent[T]](wg *sync.WaitGroup, process <-chan *GrowableSet[T], receiver *BuckettedSet[T]) {
	defer wg.Done()
	for s := range process {
		s.Range(func(item T) bool {
			_ = receiver.UpdateOrAdd(item)
			return true
		})
	}
}
