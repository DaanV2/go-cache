package fixed_test

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/fixed"
	test_util "github.com/daanv2/go-cache/test/util"
	"github.com/daanv2/go-optimal"
	"github.com/stretchr/testify/require"
)

func Benchmark_Set_Get(b *testing.B) {
	sizes := []uint64{100, 200, 500, 1000, uint64(optimal.SliceSize[collections.HashItem[*test_util.TestItem]]())}
	hasher := test_util.CheapIntHasher[int]()

	test_util.Case1(sizes, func(size uint64) {
		items := make([]collections.HashItem[*test_util.TestItem], 0, size)
		for _, item := range test_util.Generate(int(size)) {
			items = append(items, collections.NewHashItem(hasher.Hash(item.ID), item))
		}

		col := fixed.NewSet[*test_util.TestItem](size)

		for _, item := range items {
			ok := col.Set(item)
			require.True(b, ok)
		}

		for _, item := range items {
			v, ok := col.Get(item)
			require.True(b, ok, "item not found: %s", item.Value)
			require.Equal(b, v, item, "item not equal: %s != %s", v.Value, item.Value)
		}

		b.Run(fmt.Sprintf("Get(%v)", size), func(t *testing.B) {
			t.ReportMetric(float64(size), "size")

			for i := 0; i < t.N; i++ {
				for _, item := range items {
					v, ok := col.Get(item)
					if !ok || v.Value == nil {
						t.Fail()
					}
				}
			}
		})
	})
}
