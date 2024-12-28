package test_util

import (
	"github.com/daanv2/go-cache/pkg/hash"
	"golang.org/x/exp/constraints"
)

type intHasher[T constraints.Integer] struct{}

// Hash implements hash.Hasher.
func (i intHasher[T]) Hash(item T) uint64 {
	return uint64(item)
}

func CheapIntHasher[T constraints.Integer]() hash.Hasher[T] {
	return intHasher[T]{}
}
