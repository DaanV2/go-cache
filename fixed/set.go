package fixed

import (
	"iter"
	"sync"

	"github.com/daanv2/go-cache/pkg/bloomfilters"
)

// Set is a fixed size slice, that can be used to store a fixed amount of items
type Set[T comparable] struct {
	amount    uint64
	hashrange *bloomfilters.Cheap
	items     []SetItem[T] // The items in the slice
	lock      sync.RWMutex // The lock to protect the slice
}

func NewSet[T comparable](amount uint64) Set[T] {
	return Set[T]{
		amount:    amount,
		items:     make([]SetItem[T], amount),
		hashrange: bloomfilters.NewCheap(amount),
		lock:      sync.RWMutex{},
	}
}

func (s *Set[T]) Cap() int {
	return cap(s.items)
}

func (s *Set[T]) Len() int {
	return len(s.items)
}

func (s *Set[T]) HasHash(hash uint64) bool {
	return s.hashrange.Has(hash)
}

func (s *Set[T]) index(item SetItem[T]) uint64 {
	return item.Hash % s.amount
}

func (s *Set[T]) Get(item SetItem[T]) (SetItem[T], bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.hashrange.Has(item.Hash) {
		return s.get(item)
	}

	return item, false
}

func (s *Set[T]) get(item SetItem[T]) (SetItem[T], bool) {
	sindex := s.index(item)

	sub := s.items[sindex:]
	for _, v := range sub {
		if item == v {
			return v, true
		}
	}

	sub = s.items[:sindex]
	for _, v := range sub {
		if item == v {
			return v, true
		}
	}

	return item, false
}

// Set Add the given item to the set, if equivalant item was overriden, or empty space filled, true is returned
func (s *Set[T]) Set(item SetItem[T]) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.set(item)
}

func (s *Set[T]) set(item SetItem[T]) bool {
	sindex := s.index(item)

	sub := s.items[sindex:]
	for i, v := range sub {
		if (v.Hash == item.Hash && v.Value == item.Value) || v.IsEmpty() {
			sub[i] = item
			s.hashrange.Set(item.Hash)
			return true
		}
	}

	sub = s.items[:sindex]
	for i, v := range sub {
		if (v.Hash == item.Hash && v.Value == item.Value) || v.IsEmpty() {
			sub[i] = item
			s.hashrange.Set(item.Hash)
			return true
		}
	}

	return false
}

func (s *Set[T]) Update(item SetItem[T]) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.update(item)
}

func (s *Set[T]) update(item SetItem[T]) bool {
	sindex := s.index(item)
	sub := s.items[sindex:]
	for i, v := range sub {
		if v.Hash == item.Hash && v.Value == item.Value {
			sub[i] = item
			return true
		}
	}

	sub = s.items[:sindex]
	for i, v := range sub {
		if v.Hash == item.Hash && v.Value == item.Value {
			sub[i] = item
			return true
		}
	}

	return false
}

func (s *Set[T]) Read() iter.Seq[SetItem[T]] {
	return func(yield func(SetItem[T]) bool) {
		s.lock.RLock()
		defer s.lock.RUnlock()

		for _, v := range s.items {
			if v.IsEmpty() {
				continue
			}

			if !yield(v) {
				return
			}
		}
	}
}
