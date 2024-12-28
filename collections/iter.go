package collections

import "iter"

// Readable is a generic interface that represents a readable collection of items of any type T.
// It defines a single method, Read, which returns an iterator sequence of type T.
type Readable[T any] interface {
	// Read returns an iterator sequence of type T.
	Read() iter.Seq[T]
}
