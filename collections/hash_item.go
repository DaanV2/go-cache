package collections

import (
	"github.com/daanv2/go-kit/generics"
)

type HashItem[T comparable] struct {
	hash uint64
	item T
}

func NewHashItem[T comparable](hash uint64, item T) HashItem[T] {
	return HashItem[T]{hash, item}
}

func (s HashItem[T]) Hash() uint64 {
	return s.hash
}

func (s HashItem[T]) Value() T {
	return s.item
}

func (s HashItem[T]) IsEmpty() bool {
	return s == generics.Empty[HashItem[T]]()
}
