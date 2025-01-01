package large

import (
	"fmt"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-kit/generics"
)

// BuckettedSet is a set of items, that uses a pre-defined amount of buckets, each item generates an hash, from which a bucket can be specified
type BuckettedMap[K comparable, V any] struct {
	BuckettedSet[collections.KeyValue[K, V]]
}

// NewBuckettedMap creates a new BuckettedMap with the specified capacity, hasher, and options.
// The BuckettedMap is a concurrent map that uses a bucketing strategy to reduce contention.
func NewBuckettedMap[K comparable, V any](capacity uint64, keyhasher hash.Hasher[K], opts ...options.Option[Options]) (*BuckettedMap[K, V], error) {
	set, err := NewBuckettedSet[collections.KeyValue[K, V]](
		capacity,
		collections.KeyValueHasher[K, V](keyhasher),
		opts...,
	)

	if err != nil {
		return nil, err
	}

	return &BuckettedMap[K, V]{
		BuckettedSet: *set,
	}, nil
}

// Get retrieves the value for the specified key from the BuckettedMap.
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

// Set will add or update the value for the specified key in the BuckettedMap. It returns true if the value was added, false if it was updated.
func (m *BuckettedMap[K, V]) Set(key K, item V) bool {
	kv := collections.NewKeyValue(key, item)
	setitem := NewSetItem(kv, m.hasher.Hash(kv))
	bucket := m.bucketIndex(setitem)
	return m.sets[bucket].updateOrAdd(setitem)
}

// Append adds all items from the specified Rangeable to the BuckettedMap.
func (m *BuckettedMap[K, V]) Append(other collections.Rangeable[collections.KeyValue[K, V]]) {
	other.Range(func(item collections.KeyValue[K, V]) bool {
		m.Set(item.Key(), item.Value())
		return true
	})
}

// AppendParralel adds all items from the specified ParralelRangeable to the BuckettedMap.
func (m *BuckettedMap[K, V]) AppendParralel(other collections.ParralelRangeable[collections.KeyValue[K, V]]) {
	other.RangeParralel(func(item collections.KeyValue[K, V]) bool {
		m.Set(item.Key(), item.Value())
		return true
	})
}

func (m *BuckettedMap[K, V]) Grow(new_capacity uint64) {
	m.BuckettedSet.Grow(new_capacity)
}

func (s *BuckettedMap[K, V]) String() string {
	return fmt.Sprintf("large.BuckettedMap[%s,%s]", generics.NameOf[K](), generics.NameOf[V]())
}

func (s *BuckettedMap[K, V]) GoString() string {
	return s.String()
}
