package test_util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Int_Hasher(t *testing.T) {
	sizes := []int{10, 50, 100, 250, 500, 1000, 2000}
	hasher := CheapIntHasher[int]()

	Case1(sizes, func(size int) {
		t.Run(fmt.Sprintf("Size(%d)", size), func(t *testing.T) {
			check := make(map[uint64]struct{}, size)

			for i, item := range Generate(size) {
				h := hasher.Hash(item.ID)
				_, ok := check[h]
				require.False(t, ok, "hash collision", i, h)
				check[h] = struct{}{}
			}
		})
	})
}
