package hash

type Hasher[T any] interface {
	Hash(item T) uint64
}
