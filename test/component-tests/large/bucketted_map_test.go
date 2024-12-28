package large_test

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/large"
	"github.com/daanv2/go-cache/pkg/hash"
	test_util "github.com/daanv2/go-cache/test/util"
	"github.com/daanv2/go-optimal/pkg/cpu"
	"github.com/stretchr/testify/require"
)

func Test_BuckettedMap(t *testing.T) {
	sizes := []uint64{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrency(%v)", size), func(t *testing.T) {
			col, err := large.NewBuckettedMap[int, string](size*10, hash.IntegerHasher[int](hash.MD5))
			require.NoError(t, err)

			items := test_util.Generate(int(size))
			collections.Shuffle(items)

			for _, item := range items {
				ok := col.Set(item.ID, item.Data)
				require.True(t, ok)

				v, ok := col.Get(item.ID)
				require.True(t, ok)
				require.Equal(t, v.Value(), item.Data)
			}

			for _, item := range items {
				data := item.Data + "updated"
				ok := col.Set(item.ID, data)
				require.False(t, ok)

				v, ok := col.Get(item.ID)
				require.True(t, ok)
				require.Equal(t, v.Value(), data)
			}
		})
	}
}

func Test_BuckettedMap_Grow(t *testing.T) {
	sizes := []uint64{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrency(%v)", size), func(t *testing.T) {
			col, err := large.NewBuckettedMap[int, string](size*10, hash.IntegerHasher[int](hash.MD5))
			require.NoError(t, err)

			items := test_util.Generate(int(size))
			collections.Shuffle(items)

			for _, item := range items {
				ok := col.Set(item.ID, item.Data)
				require.True(t, ok)

				v, ok := col.Get(item.ID)
				require.True(t, ok)
				require.Equal(t, v.Value(), item.Data)
			}

			col.Grow(size * 2)

			for _, item := range items {
				v, ok := col.Get(item.ID)
				require.True(t, ok)
				require.Equal(t, v.Value(), item.Data)
			}
		})
	}
}

func Test_BuckettedMap_Concurrency(t *testing.T) {
	sizes := []uint64{100, 200, 300, 400, 1000, 10000, 20000}
	target := []cpu.CacheKind{cpu.CacheL1, cpu.CacheL2, cpu.CacheL3}

	test_util.Case2(sizes, target, func(size uint64, cache cpu.CacheKind) {
		t.Run(fmt.Sprintf("Concurrency(%v)->%s", size, cache), func(t *testing.T) {
			col, err := large.NewBuckettedMap[int, string](
				size*10,
				 hash.IntegerHasher[int](hash.MD5),
				  large.WithCacheTarget[collections.KeyValue[int, string]](cache),
				)
			require.NoError(t, err)

			items := make([]collections.KeyValue[int, string], 0, int(size))
			for _, item := range test_util.Generate(int(size)) {
				items = append(items, collections.NewKeyValue[int, string](item.ID, item.Data))
			}
			collections.Shuffle(items)

			splitWithOverlap(col, items)
			check := make(map[int]int, size)

			for item := range col.Read() {
				check[item.Key()] = check[item.Key()] + 1
				require.LessOrEqual(t, check[item.Key()], 1)
			}
		})
	})
}

func Benchmark_BuckettedMap_Concurrency(t *testing.B) {
	sizes := []uint64{100, 200, 300, 400, 1000, 10000, 20000}
	target := []cpu.CacheKind{cpu.CacheL1, cpu.CacheL2, cpu.CacheL3}

	test_util.Case2(sizes, target, func(size uint64, cache cpu.CacheKind) {
		t.Run(fmt.Sprintf("Concurrency(%v)->%s", size, cache), func(t *testing.B) {
			t.ReportMetric(float64(size), "size")
			t.ReportMetric(float64(cache)+1, "target")

			for i := 0; i < t.N; i++ {
				col, err := large.NewBuckettedMap[int, string](
					size*10,
					 hash.IntegerHasher[int](hash.MD5),
					  large.WithCacheTarget[collections.KeyValue[int, string]](cache),
					)
				require.NoError(t, err)

				items := make([]collections.KeyValue[int, string], 0, int(size))
				for _, item := range test_util.Generate(int(size)) {
					items = append(items, collections.NewKeyValue[int, string](item.ID, item.Data))
				}
				collections.Shuffle(items)

				splitWithOverlap(col, items)
			}
		})
	})
}
