package collections

import (
	"github.com/daanv2/go-kit/generics"
)

type HashItem[T comparable] struct {
	Hash uint64
	Value T
}

func NewHashItem[T comparable](hash uint64, value T) HashItem[T] {
	return HashItem[T]{hash, value}
}

func (s HashItem[T]) GetHash() uint64 {
	return s.Hash
}

func (s HashItem[T]) GetValue() T {
	return s.Value
}

func (s HashItem[T]) IsEmpty() bool {
	return s == generics.Empty[HashItem[T]]()
}
