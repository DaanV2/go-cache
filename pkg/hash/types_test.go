package hash_test

import (
	"testing"

	"github.com/daanv2/go-cache/pkg/hash"
)



func Benchmark_Strings(b *testing.B) {
	strs := []string {
		"57dd5c23-677f-47d5-9277-298bef455ae1",
		"6565dcd5-96b9-4666-812f-a51201d69e84",
		"9b019c3b-07db-4f30-806a-5d2a25ade872",
		"f212118f-af4c-4cfa-b87e-7afaa380a9a2",
	}

	string_benchmark(b, "Sha1", hash.Sha1, strs)
	string_benchmark(b, "MD5", hash.MD5, strs)
	string_benchmark(b, "Sha256", hash.Sha256, strs)
}

func string_benchmark(b *testing.B, name string, basehash func() hash.HashBuilder, strs []string) {
	b.Run("strings->"+name, func(b *testing.B) {
		hasher := hash.StringHasher(basehash)

		for i := 0; i < b.N; i++ {
			for _, s := range strs {
				hash := hasher.Hash(s)
				if hash == 0 {
					b.Fail()
				}
			}
		}
	})
}
