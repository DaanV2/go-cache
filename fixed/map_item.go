package fixed

import (
	"github.com/daanv2/go-cache/pkg/hash"
	hashmark "github.com/daanv2/go-cache/pkg/hash/marked"
	"github.com/daanv2/go-kit/generics"
)

type KeyValue[K comparable, V any] struct {
	Hash  uint64 // The hash of the key marked for empty checks. See [hashmark.MarkedHash]
	Key   K
	Value V
}

// NewKeyValue creates a new KeyValue instance with the given key and value.
func NewKeyValue[K comparable, V any](hash uint64, key K, value V) KeyValue[K, V] {
	return KeyValue[K, V]{
		hashmark.MarkedHash(hash),
		key,
		value,
	}
}

// NewKey creates a new KeyValue instance with the given key.
func NewKey[K comparable, V any](hash uint64, key K) KeyValue[K, V] {
	return KeyValue[K, V]{
		Hash: hashmark.MarkedHash(hash),
		Key:  key,
	}
}

// EmptyKeyValue creates a new KeyValue instance with empty key and value.
func EmptyKeyValue[K comparable, V any]() KeyValue[K, V] {
	return NewKeyValue(0, generics.Empty[K](), generics.Empty[V]())
}

// Key returns the key of the KeyValue pair.
func (kv KeyValue[K, V]) GetKey() K {
	return kv.Key
}

// Value returns the value of the KeyValue pair.
func (kv KeyValue[K, V]) GetValue() V {
	return kv.Value
}

func (kv KeyValue[K, V]) GetHash() uint64 {
	return kv.Hash
}

func (kv KeyValue[K, V]) IsEmpty() bool {
	return hashmark.IsEmpty(kv.Hash)
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
	return k.hasher.Hash(item.Key)
}
