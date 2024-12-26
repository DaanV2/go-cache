package collections

type Set[T any] interface {
	Readable[T]
	GetOrAdd(item T) (T, bool)
}
