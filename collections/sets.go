package collections

// Set is a generic interface that represents a collection of unique items of any type T.
// It embeds the Readable interface, which provides read-only access to the items in the set.
// The GetOrAdd method attempts to retrieve the specified item from the set, and if it does not exist,
// it adds the item to the set. The method returns the item and a boolean indicating whether the item
// was already present in the set (true) or was newly added (false).
type Set[T any] interface {
	Readable[T]
	// GetOrAdd attempts to retrieve the specified item from the set, and if it does not exist,
	GetOrAdd(item T) (T, bool)
}

type Map[K comparable, V any] interface {
	Get(key K) (KeyValue[K, V], bool)
}
