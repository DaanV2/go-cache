package large

import "github.com/daanv2/go-cache/pkg/constraints"

type SetItem[T constraints.Equivalent[T]] struct {
	hash uint64
	item T
}

func NewSetItem[T constraints.Equivalent[T]](item T, hash uint64) SetItem[T] {
	return SetItem[T]{
		hash: hash,
		item: item,
	}
}

func (s *SetItem[T]) Value() T {
	return s.item
}

func (s *SetItem[T]) Hash() uint64 {
	return s.hash
}

func (s *SetItem[T]) Equals(other SetItem[T]) bool {
	return s.hash == other.hash &&
		s.item.Equals(other.item)
}
