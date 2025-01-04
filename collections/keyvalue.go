package collections

import (
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-kit/generics"
)

// KeyValue is a generic struct that holds a key-value pair.
type KeyValue[K comparable, V any] struct {
	key   K
	value V
}

// NewKeyValue creates a new KeyValue instance with the given key and value.
func NewKeyValue[K comparable, V any](key K, value V) KeyValue[K, V] {
	return KeyValue[K, V]{
		key,
		value,
	}
}

// EmptyKeyValue creates a new KeyValue instance with empty key and value.
func EmptyKeyValue[K comparable, V any]() KeyValue[K, V] {
	return NewKeyValue(generics.Empty[K](), generics.Empty[V]())
}

// Key returns the key of the KeyValue pair.
func (kv KeyValue[K, V]) Key() K {
	return kv.key
}

// Value returns the value of the KeyValue pair.
func (kv KeyValue[K, V]) Value() V {
	return kv.value
}

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
	return k.hasher.Hash(item.key)
}
