package constraints

type Comparable[T any] interface {
	Compare(other T) int
}

type Equivalent[T any] interface {
	Equal(other T) bool
}
