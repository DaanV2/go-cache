package collections

import "iter"


type Readable[T any] interface {
	Read() iter.Seq[T]
}
