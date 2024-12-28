package hash

var _ Hasher[struct{}] = &FunctionHasher[struct{}]{}

type FunctionHasher[T any] struct {
	builder func() HashBuilder
	toBytes func(item T) []byte
}

// Hash implements Hasher.
func (f *FunctionHasher[T]) Hash(item T) uint64 {
	builder := f.builder()
	_ = builder.Write(f.toBytes(item))
	return builder.Sum()
}

func NewFunctionHasher[T any](builder func() HashBuilder, toBytes func(item T) []byte) *FunctionHasher[T] {
	return &FunctionHasher[T]{builder, toBytes}
}
