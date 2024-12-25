package constraints

type Comparable[T any] interface {
	Compare(other T) int
}

type Equivalent[T any] interface {
	Equals(other T) bool
}
