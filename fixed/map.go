package fixed

import (
	"iter"
	"sync"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/pkg/hash"
)

// Map is a fixed size slice, that can be used to store a fixed amount of items
type Map[K, V comparable] struct {
	amount    uint64
	hashrange hash.Range
	items     []collections.HashItem[collections.KeyValue[K, V]] // The items in the slice
	lock      sync.RWMutex              // The lock to protect the slice
}

func NewMap[K, V comparable](amount uint64) Map[K, V] {
	return Map[K, V]{
		amount:    amount,
		items:     make([]collections.HashItem[collections.KeyValue[K, V]], amount),
		hashrange: hash.NewRange(),
		lock:      sync.RWMutex{},
	}
}

func (s *Map[K, V]) Cap() int {
	return cap(s.items)
}

func (s *Map[K, V]) Len() int {
	return len(s.items)
}

func (s *Map[K, V]) HasHash(hash uint64) bool {
	return s.hashrange.Has(hash)
}

func (s *Map[K, V]) index(item collections.HashItem[collections.KeyValue[K, V]]) uint64 {
	return item.Hash() % s.amount
}

func (s *Map[K, V]) indexH(hash uint64) uint64 {
	return hash % s.amount
}

func (s *Map[K, V]) Get(item collections.HashItem[collections.KeyValue[K, V]]) (collections.HashItem[collections.KeyValue[K, V]], bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.hashrange.Has(item.Hash()) {
		return s.get(item)
	}

	return item, false
}

func (s *Map[K, V]) get(item collections.HashItem[collections.KeyValue[K, V]]) (collections.HashItem[collections.KeyValue[K, V]], bool) {
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

// Map Add the given item to the set, if equivalant item was overriden, or empty space filled, true is returned
func (s *Map[K, V]) Set(item collections.HashItem[collections.KeyValue[K, V]]) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.set(item)
}

func (s *Map[K, V]) set(item collections.HashItem[collections.KeyValue[K, V]]) bool {
	sindex := s.index(item)

	sub := s.items[sindex:]
	for i, v := range sub {
		if v.IsEmpty() || sameKey(item, v) {
			sub[i] = item
			s.hashrange.Update(item.Hash())
			return true
		}
	}

	sub = s.items[:sindex]
	for i, v := range sub {
		if item == v || sameKey(item, v) {
			sub[i] = item
			s.hashrange.Update(item.Hash())
			return true
		}
	}

	return false
}

func (s *Map[K, V]) Update(item collections.HashItem[collections.KeyValue[K, V]]) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.update(item)
}

func (s *Map[K, V]) update(item collections.HashItem[collections.KeyValue[K, V]]) bool {
	if !s.hashrange.Has(item.Hash()) {
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

func (s *Map[K, V]) Read() iter.Seq[collections.HashItem[collections.KeyValue[K, V]]] {
	return func(yield func(collections.HashItem[collections.KeyValue[K, V]]) bool) {
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

func sameKey[K, V comparable](a, b collections.HashItem[collections.KeyValue[K, V]]) bool {
	return a.Hash() == b.Hash() &&
		a.Value().Key() == b.Value().Key()
}
