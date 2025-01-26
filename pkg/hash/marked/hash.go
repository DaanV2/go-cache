package hashmark

const (
	empty_hash_mark  uint64 = 0
	filled_hash_mark uint64 = 0b11 << 62
	mask             uint64 = 0b11 << 62
)

// MarkedHash marks a hash value as filled using the upper 2 bits.
// This means a hash value that is empty will not equal a hash value that is filled.
func MarkedHash(v uint64) uint64 {
	return v | filled_hash_mark
}

// IsEmpty checks if a hash value is empty.
func IsEmpty(v uint64) bool {
	return !IsFilled(v)
}

// IsFilled checks if a hash value is filled.
func IsFilled(v uint64) bool {
	return v >= filled_hash_mark
}

// Equal checks if two hash values are equal.
func Equal(a, b uint64) bool {
	return a == b
}

// Empty returns an empty hash value.
func Empty() uint64 {
	return empty_hash_mark
}
