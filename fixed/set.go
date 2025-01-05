package fixed

import (
	"iter"
	"sync"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/pkg/hash"
)

// Set is a fixed size slice, that can be used to store a fixed amount of items
type Set[T comparable] struct {
	amount    uint64
	hashrange hash.Range
	items     []collections.HashItem[T] // The items in the slice
	lock      sync.RWMutex              // The lock to protect the slice
}

func NewSet[T comparable](amount uint64) Set[T] {
	return Set[T]{
		amount:    amount,
		items:     make([]collections.HashItem[T], amount),
		hashrange: hash.NewRange(),
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

func (s *Set[T]) index(item collections.HashItem[T]) uint64 {
	return item.Hash() % s.amount
}

func (s *Set[T]) indexH(hash uint64) uint64 {
	return hash % s.amount
}

func (s *Set[T]) Get(item collections.HashItem[T]) (collections.HashItem[T], bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.hashrange.Has(item.Hash()) {
		return s.get(item)
	}

	return item, false
}

func (s *Set[T]) get(item collections.HashItem[T]) (collections.HashItem[T], bool) {
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
func (s *Set[T]) Set(item collections.HashItem[T]) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.set(item)
}

func (s *Set[T]) set(item collections.HashItem[T]) bool {
	sindex := s.index(item)

	sub := s.items[sindex:]
	for i, v := range sub {
		if item == v || v.IsEmpty() {
			sub[i] = item
			s.hashrange.Update(item.Hash())
			return true
		}
	}

	sub = s.items[:sindex]
	for i, v := range sub {
		if item == v || v.IsEmpty() {
			sub[i] = item
			s.hashrange.Update(item.Hash())
			return true
		}
	}

	return false
}

func (s *Set[T]) Update(item collections.HashItem[T]) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.update(item)
}

func (s *Set[T]) update(item collections.HashItem[T]) bool {
	if !s.hashrange.Has(item.Hash()) {
		return false
	}

	sindex := s.index(item)
	sub := s.items[sindex:]
	for i, v := range sub {
		if v == item {
			sub[i] = item
			return true
		}
	}

	sub = s.items[:sindex]
	for i, v := range sub {
		if v == item {
			sub[i] = item
			return true
		}
	}

	return false
}

func (s *Set[T]) Read() iter.Seq[collections.HashItem[T]] {
	return func(yield func(collections.HashItem[T]) bool) {
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
