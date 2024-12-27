package hash

import (
	"github.com/daanv2/go-kit/generics"
	"golang.org/x/exp/constraints"
)

// StringHasher, See [MD5], [Sha1] or [Sha256] for bashHash
func StringHasher(basehash func() HashBuilder) Hasher[string] {
	return NewFunctionHasher(basehash, func(data string) []byte {
		return []byte(data)
	})
}

// IntegerHasher, See [MD5], [Sha1] or [Sha256] for bashHash
func IntegerHasher[T constraints.Integer](basehash func() HashBuilder) Hasher[T] {
	return NewFunctionHasher(basehash, func(item T) []byte {
		l := generics.SizeOf[T]()
		d := make([]byte, 0, l)

		for l > 0 {
			d = append(d, byte(l))
			l = l >> 8
		}

		return d
	})
}
