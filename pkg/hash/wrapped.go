package hash

import (
	go_hash "hash"

	"github.com/daanv2/go-cache/pkg/binary"
)

var _ HashBuilder = &WrappedHasher[go_hash.Hash]{}

type WrappedHasher[T go_hash.Hash] struct {
	base T
}

func NewWrappedHasher[T go_hash.Hash](base T) *WrappedHasher[T] {
	return &WrappedHasher[T]{base}
}

func (w *WrappedHasher[T]) Sum() uint64 {
	return binary.Uint64(w.base.Sum(nil))
}

// Write implements HashBuilder.
func (w *WrappedHasher[T]) Write(data []byte) error {
	_, err := w.base.Write(data)
	return err
}
