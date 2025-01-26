package fixed

import hashmark "github.com/daanv2/go-cache/pkg/hash/marked"

// SetItem is a generic item that can be stored in a set.
// It should be seen as immutable.
type SetItem[T comparable] struct {
	Hash  uint64 // The hash of the item marked for empty checks. See [hashmark.MarkedHash]
	Value T
}

func NewSetItem[T comparable](hash uint64, value T) SetItem[T] {
	return SetItem[T]{
		Hash:  hashmark.MarkedHash(hash),
		Value: value,
	}
}

func (s SetItem[T]) GetHash() uint64 {
	return s.Hash
}

func (s SetItem[T]) GetValue() T {
	return s.Value
}

func (s SetItem[T]) IsEmpty() bool {
	return hashmark.IsEmpty(s.Hash)
}
