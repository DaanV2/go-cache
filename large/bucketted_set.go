package large

import (
	"iter"

	"github.com/daanv2/go-cache/pkg/constraints"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/pkg/options"
)

type BuckettedSet[T constraints.Equivalent[T]] struct {
	hasher hash.Hasher[T]
	sets   []*GrowableSet[T]
}

func NewBuckettedSet[T constraints.Equivalent[T]](cap uint64, hasher hash.Hasher[T], opts ...options.Option[SetBase]) (*BuckettedSet[T], error) {
	base := NewSetBase[T]()
	err := options.Apply(&base, opts...)
	if err != nil {
		return nil, err
	}

	buckets := max(cap/(uint64(base.bucket_size)*4), 10)
	set := &BuckettedSet[T]{
		hasher: hasher,
		sets:   make([]*GrowableSet[T], 0, buckets),
	}

	for range buckets {
		s, err := NewGrowableSet(hasher, opts...)
		if err != nil {
			return nil, err
		}

		set.sets = append(set.sets, s)
	}

	return set, nil
}

func (s *BuckettedSet[T]) GetOrAdd(item T) (T, bool) {
	setitem := NewSetItem[T](item, s.hasher.Hash(item))

	bucket := setitem.hash % uint64(len(s.sets))
	return s.sets[bucket].getOrAdd(setitem)
}

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
