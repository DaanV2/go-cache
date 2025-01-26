package slices

import (
	"fmt"
	"iter"
	"sync"

	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-kit/generics"
)

// FixedHashed is a fixed size slice, that can be used to store a fixed amount of items
type FixedHashed[T hash.Hashed] struct {
	items     []T          // The items in the slice
	lock      sync.RWMutex // The lock to protect the slic	e
	hashrange hash.Range
}

// Creates a new slice of fixed sized, if its full nothing can be added
func NewFixedHashed[T hash.Hashed](amount int) FixedHashed[T] {
	return FixedHashed[T]{
		items:     make([]T, 0, amount),
		hashrange: hash.NewRange(),
	}
}

// Cap returns the capacity of the slice
func (s *FixedHashed[T]) Cap() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return cap(s.items)
}

// UnsafeCap returns the capacity of the slice without locking
func (s *FixedHashed[T]) UnsafeCap() int {
	return cap(s.items)
}

// Len returns the amount of items in the slice
func (s *FixedHashed[T]) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.items)
}

// UnsafeLen returns the amount of items in the slice without locking
func (s *FixedHashed[T]) UnsafeLen() int {
	return len(s.items)
}

// SpaceLeft returns the amount of space left in the slice
func (s *FixedHashed[T]) SpaceLeft() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return cap(s.items) - len(s.items)
}

// UnsafeSpaceLeft returns the amount of space left in the slice without locking
func (s *FixedHashed[T]) UnsafeSpaceLeft() int {
	return cap(s.items) - len(s.items)
}

// IsFull returns if the slice is full
func (s *FixedHashed[T]) IsFull() bool {
	return s.SpaceLeft() == 0
}

// UnsafeIsFull returns if the slice is full without locking
func (s *FixedHashed[T]) UnsafeIsFull() bool {
	return s.UnsafeSpaceLeft() == 0
}

// Get returns the item at the given index, and if it exists
func (s *FixedHashed[T]) Get(index int) (T, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.UnsafeGet(index)
}

// UnsafeGet returns the item at the given index, and if it exists without locking
func (s *FixedHashed[T]) UnsafeGet(index int) (T, bool) {
	if index >= len(s.items) {
		return generics.Empty[T](), false
	}

	return s.items[index], true
}

// Set will set the item at the given index, and return if it was successful
func (s *FixedHashed[T]) Set(index int, value T) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.UnsafeSet(index, value)
}

// UnsafeSet will set the item at the given index, and return if it was successful without locking
func (s *FixedHashed[T]) UnsafeSet(index int, value T) bool {
	if index >= len(s.items) {
		return false
	}

	old := s.items[index]
	s.items[index] = value

	if old.Hash() != value.Hash() {
		s.Rehash()
	}

	return true
}

// Clear will remove all items from the slice
func (s *FixedHashed[T]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = s.items[:0]
}

// Delete will remove the item at the given index, and return if it was successful
func (s *FixedHashed[T]) Delete(index int) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.UnsafeDelete(index)
}

// UnsafeDelete will remove the item at the given index, and return if it was successful without locking
func (s *FixedHashed[T]) UnsafeDelete(index int) bool {
	if index >= len(s.items) {
		return false
	}

	s.items = append(s.items[:index], s.items[index+1:]...)
	s.Rehash()
	return true
}

// DeleteFunc will remove the all item that matches the predicate, and return the amount of items removed
func (s *FixedHashed[T]) DeleteFunc(predicate func(v T) bool) int {
	s.lock.Lock()
	defer s.lock.Unlock()
	amount := 0

	for i, v := range s.items {
		if predicate(v) {
			s.items = append(s.items[:i], s.items[i+1:]...)
			amount++
		}
	}
	s.Rehash()

	return amount
}

// Find will return the first item that matches the predicate, and if it was found
func (s *FixedHashed[T]) Find(predicate func(v T) bool) (T, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, v := range s.items {
		if predicate(v) {
			return v, true
		}
	}

	return generics.Empty[T](), false
}

// FindIndex will return the index of the first item that matches the predicate, and if it was found
func (s *FixedHashed[T]) FindIndex(predicate func(v T) bool) (int, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i, v := range s.items {
		if predicate(v) {
			return i, true
		}
	}

	return -1, false
}

// TryAppend will check how much space is left, and attempt to write as much as possible from the given data into its own buffer
//
// If you have 5 items, and there is room for 3, it will return 3, and has added 3 items to its buffer
func (s *FixedHashed[T]) TryAppend(items ...T) int {
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

	for _, item := range items {
		s.items = append(s.items, item)
		s.hashrange.Update(item.Hash())
	}

	return len(items)
}

// Read will return a sequence of the items in the slice
func (s *FixedHashed[T]) Read() iter.Seq[T] {
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

func (s *FixedHashed[T]) Rehash() {
	s.lock.Lock()
	defer s.lock.Unlock()

	nr := hash.NewRange()
	for _, item := range s.items {
		nr.Update(item.Hash())
	}
	s.hashrange = nr
}

func (s *FixedHashed[T]) HasHash(v uint64) bool {
	return s.hashrange.Has(v)
}

func (s *FixedHashed[T]) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return fmt.Sprintf("fixed.HashedSlice[%s,%d/%d]", generics.NameOf[T](), len(s.items), cap(s.items))
}

func (s *FixedHashed[T]) GoString() string {
	return s.String()
}
