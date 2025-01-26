package maps

import "github.com/daanv2/go-cache/pkg/hash"

// kvHasher is a struct that implements the hash.Hasher interface for KeyValue pairs.
type kvHasher[K comparable, V any] struct {
	hasher hash.Hasher[K]
}

// KeyValueHasher creates a new hasher for KeyValue pairs using the provided hasher for the key type.
func KeyValueHasher[K comparable, V any](hasher hash.Hasher[K]) hash.Hasher[KeyValue[K, V]] {
	return &kvHasher[K, V]{hasher}
}

// Hash computes the hash value of the given KeyValue pair using the hasher for the key type.
func (k *kvHasher[K, V]) Hash(item KeyValue[K, V]) uint64 {
	return k.hasher.Hash(item.Key)
}
