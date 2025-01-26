package test_util

import (
	"crypto/sha1"
	"encoding/binary"

	"github.com/daanv2/go-cache/pkg/hash"
	"golang.org/x/exp/constraints"
)

type intHasher[T constraints.Integer] struct{}

// Hash implements hash.Hasher.
func (i intHasher[T]) Hash(item T) uint64 {
	b := binary.LittleEndian.AppendUint64(nil, uint64(item))
	v := sha1.Sum(b)

	return binary.LittleEndian.Uint64(v[:])
}

func CheapIntHasher[T constraints.Integer]() hash.Hasher[T] {
	return intHasher[T]{}
}
