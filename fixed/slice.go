package fixed

import (
	"iter"
	"sync"

	"github.com/daanv2/go-kit/generics"
)

type Slice[T any] struct {
	items []T
	lock  sync.RWMutex
}

func NewSlice[T any](amount int) Slice[T] {
	return Slice[T]{
		items: make([]T, 0, amount),
	}
}

func (s *Slice[T]) Cap() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return cap(s.items)
}

func (s *Slice[T]) UnsafeCap() int {
	return cap(s.items)
}

func (s *Slice[T]) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.items)
}

func (s *Slice[T]) UnsafeLen() int {
	return len(s.items)
}

func (s *Slice[T]) SpaceLeft() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return cap(s.items) - len(s.items)
}

func (s *Slice[T]) UnsafeSpaceLeft() int {
	return cap(s.items) - len(s.items)
}

func (s *Slice[T]) IsFull() bool {
	return s.SpaceLeft() == 0
}

func (s *Slice[T]) UnsafeIsFull() bool {
	return s.UnsafeSpaceLeft() == 0
}

func (s *Slice[T]) Get(index int) (T, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.UnsafeGet(index)
}

func (s *Slice[T]) UnsafeGet(index int) (T, bool) {
	if index >= len(s.items) {
		return generics.Empty[T](), false
	}

	return s.items[index], true
}

func (s *Slice[T]) Set(index int, value T) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.UnsafeSet(index, value)
}

func (s *Slice[T]) UnsafeSet(index int, value T) bool {
	if index >= len(s.items) {
		return false
	}

	s.items[index] = value
	return true
}

func (s *Slice[T]) Find(predicate func(v T) bool) (T, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, v := range s.items {
		if predicate(v) {
			return v, true
		}
	}

	return generics.Empty[T](), false
}

// TryAppend will check how much space is left, and attempt to write as much as possible from the given data into its own buffer
//
// If you have 5 items, and there is room for 3, it will return 3, and has added 3 items to its buffer
func (s *Slice[T]) TryAppend(items ...T) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Check if we have space and how much
	space := s.UnsafeSpaceLeft()
	if space <= 0 {
		return 0
	}
	if space < len(items) {
		items = items[:space]
	}

	s.items = append(s.items, items...)
	return len(items)
}

func (s *Slice[T]) Read() iter.Seq[T] {
	return func(yield func(T) bool) {
		s.lock.RLock()
		defer s.lock.RUnlock()

		for _, item := range s.items {
			if !yield(item) {
				return
			}
		}
	}
}
