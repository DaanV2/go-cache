package collections

import (
	"github.com/daanv2/go-cache/pkg/constraints"
	"github.com/daanv2/go-kit/generics"
)

type HashItem[T constraints.Equivalent[T]] struct {
	hash uint64
	item T
}

func NewHashItem[T constraints.Equivalent[T]](hash uint64, item T) HashItem[T] {
	return HashItem[T]{ hash, item }
}

func (s HashItem[T]) Equal(other HashItem[T]) bool {
	return s.hash == other.hash && s.item.Equal(other.item)
}

func (s HashItem[T]) Hash() uint64 {
	return s.hash
}

func (s HashItem[T]) Value() T {
	return s.item
}

func (s HashItem[T]) IsEmpty() bool {
	return s.Equal(generics.Empty[HashItem[T]]())
}
