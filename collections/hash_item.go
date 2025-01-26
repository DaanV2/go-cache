package collections

import hashmark "github.com/daanv2/go-cache/pkg/hash/marked"

type HashItem[T comparable] struct {
	Hash  uint64
	Value T
}

func NewHashItem[T comparable](hash uint64, value T) HashItem[T] {
	return HashItem[T]{
		Hash:  hashmark.MarkedHash(hash),
		Value: value,
	}
}

func (s HashItem[T]) GetHash() uint64 {
	return s.Hash
}

func (s HashItem[T]) GetValue() T {
	return s.Value
}

func (s HashItem[T]) IsEmpty() bool {
	return hashmark.IsEmpty(s.Hash)
}
