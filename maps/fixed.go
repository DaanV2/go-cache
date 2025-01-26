package maps

import (
	"iter"
	"sync"

	"github.com/daanv2/go-cache/pkg/hash"
)

// Fixed is a fixed size slice, that can be used to store a fixed amount of items
type Fixed[K, V comparable] struct {
	amount    uint64
	hashrange hash.Range
	items     []KeyValue[K, V] // The items in the slice
	lock      sync.RWMutex     // The lock to protect the slice
}

func NewFixed[K, V comparable](amount uint64) Fixed[K, V] {
	return Fixed[K, V]{
		amount:    amount,
		items:     make([]KeyValue[K, V], amount),
		hashrange: hash.NewRange(),
		lock:      sync.RWMutex{},
	}
}

func (s *Fixed[K, V]) Cap() int {
	return cap(s.items)
}

func (s *Fixed[K, V]) Len() int {
	return len(s.items)
}

func (s *Fixed[K, V]) HasHash(hash uint64) bool {
	return s.hashrange.Has(hash)
}

func (s *Fixed[K, V]) index(item KeyValue[K, V]) uint64 {
	return item.Hash % s.amount
}

func (s *Fixed[K, V]) Get(item KeyValue[K, V]) (KeyValue[K, V], bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.hashrange.Has(item.Hash) {
		return s.get(item)
	}

	return item, false
}

func (s *Fixed[K, V]) get(item KeyValue[K, V]) (KeyValue[K, V], bool) {
	sindex := s.index(item)

	sub := s.items[sindex:]
	for _, v := range sub {
		if sameKey(item, v) {
			return v, true
		}
	}

	sub = s.items[:sindex]
	for _, v := range sub {
		if sameKey(item, v) {
			return v, true
		}
	}

	return item, false
}

// Fixed Add the given item to the set, if equivalant item was overriden, or empty space filled, true is returned
func (s *Fixed[K, V]) Set(item KeyValue[K, V]) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.set(item)
}

func (s *Fixed[K, V]) set(item KeyValue[K, V]) bool {
	sindex := s.index(item)

	sub := s.items[sindex:]
	for i, spot := range sub {
		if spot.IsEmpty() || sameKey(item, spot) {
			sub[i] = item
			s.hashrange.Update(item.Hash)
			return true
		}
	}

	sub = s.items[:sindex]
	for i, spot := range sub {
		if spot.IsEmpty() || sameKey(item, spot) {
			sub[i] = item
			s.hashrange.Update(item.Hash)
			return true
		}
	}

	return false
}

func (s *Fixed[K, V]) Update(item KeyValue[K, V]) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.update(item)
}

func (s *Fixed[K, V]) update(item KeyValue[K, V]) bool {
	if !s.hashrange.Has(item.Hash) {
		return false
	}

	sindex := s.index(item)
	sub := s.items[sindex:]
	for i, v := range sub {
		if sameKey(item, v) {
			sub[i] = item
			return true
		}
	}

	sub = s.items[:sindex]
	for i, v := range sub {
		if sameKey(item, v) {
			sub[i] = item
			return true
		}
	}

	return false
}

func (s *Fixed[K, V]) Read() iter.Seq[KeyValue[K, V]] {
	return func(yield func(KeyValue[K, V]) bool) {
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

func sameKey[K, V comparable](a, b KeyValue[K, V]) bool {
	return a.Hash == b.Hash &&
		a.Key == b.Key
}
