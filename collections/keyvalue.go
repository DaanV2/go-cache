package collections

import (
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-kit/generics"
)

type KeyValue[K comparable, V any] struct {
	key   K
	value V
}

func NewKeyValue[K comparable, V any](key K, value V) KeyValue[K, V] {
	return KeyValue[K, V]{
		key,
		value,
	}
}

func EmptyKeyValue[K comparable, V any]() KeyValue[K, V] {
	return NewKeyValue(generics.Empty[K](), generics.Empty[V]())
}

func (kv KeyValue[K, V]) Key() K {
	return kv.key
}

func (kv KeyValue[K, V]) Value() V {
	return kv.value
}

func (kv KeyValue[K, V]) Equals(other KeyValue[K, V]) bool {
	return kv.key == other.key
}

type kvHasher[K comparable, V any] struct {
	hasher hash.Hasher[K]
}

func KeyValueHasher[K comparable, V any](hasher hash.Hasher[K]) hash.Hasher[KeyValue[K, V]] {
	return &kvHasher[K, V]{hasher}
}

func (k *kvHasher[K, V]) Hash(item KeyValue[K, V]) uint64 {
	return k.hasher.Hash(item.key)
}
