package large

import (
	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-kit/generics"
)

// BuckettedSet is a set of items, that uses a pre-defined amount of buckets, each item generates an hash, from which a bucket can be specified
type BuckettedMap[K comparable, V any] struct {
	BuckettedSet[collections.KeyValue[K, V]]
}

func NewBuckettedMap[K comparable, V any](cap uint64, hasher hash.Hasher[K], opts ...options.Option[SetBase]) (*BuckettedMap[K, V], error) {
	set, err := NewBuckettedSet[collections.KeyValue[K, V]](
		cap,
		collections.KeyValueHasher[K, V](hasher),
		opts...,
	)

	if err != nil {
		return nil, err
	}

	return &BuckettedMap[K, V]{
		BuckettedSet: *set,
	}, nil
}

func (m *BuckettedMap[K, V]) Get(key K) (collections.KeyValue[K, V], bool) {
	kv := collections.NewKeyValue(key, generics.Empty[V]())
	setitem := NewSetItem(kv, m.hasher.Hash(kv))
	bucket := m.bucketIndex(setitem)
	v, ok := m.sets[bucket].find(setitem)
	if ok {
		return v.Value(), true
	}

	return collections.EmptyKeyValue[K, V](), false
}

// Set, TODO, returns true if added, false is something got updated
func (m *BuckettedMap[K, V]) Set(key K, item V) bool {
	kv := collections.NewKeyValue(key, item)
	setitem := NewSetItem(kv, m.hasher.Hash(kv))
	bucket := m.bucketIndex(setitem)
	return m.sets[bucket].updateOrAdd(setitem)
}
