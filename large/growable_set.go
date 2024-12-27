package large

import (
	"errors"
	"iter"
	"sync"

	"github.com/daanv2/go-cache/fixed"
	"github.com/daanv2/go-cache/pkg/constraints"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/pkg/iterators"
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-locks"
	optimal "github.com/daanv2/go-optimal"
)

type SetBase struct {
	bucket_size int
	bucket_lock *sync.RWMutex
	items_lock  *locks.Pool
}

func NewSetBase[T any]() SetBase {
	return SetBase{
		bucket_size: optimal.SliceSize[T](),
		bucket_lock: &sync.RWMutex{},
		items_lock:  locks.NewPool(),
	}
}

type GrowableSet[T constraints.Equivalent[T]] struct {
	SetBase
	hasher  hash.Hasher[T]
	buckets []*fixed.Slice[SetItem[T]]
}

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

func (s *GrowableSet[T]) GetOrAdd(item T) (T, bool) {
	setitem := NewSetItem[T](item, s.hasher.Hash(item))

	return s.getOrAdd(setitem)
}

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
		v, ok := s.buckets[i].Find(func(v SetItem[T]) bool {
			return v.Equals(item)
		})
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
		index, ok := s.buckets[i].FindIndex(func(v SetItem[T]) bool {
			return v.Equals(item)
		})
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

func (s *GrowableSet[T]) Range(yield func(item T) bool) {
	iterators.RangeCol(s, yield)
}
