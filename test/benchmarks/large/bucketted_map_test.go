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

func Benchmark_BuckettedMap(t *testing.B) {
	sizes := []uint64{100, 200, 300, 400}
	target := []cpu.CacheKind{cpu.CacheL1, cpu.CacheL2, cpu.CacheL3}

	test_util.Case2(sizes, target, func(size uint64, cache cpu.CacheKind) {
		items := make([]collections.KeyValue[int, string], 0, int(size))
		for _, item := range test_util.Generate(int(size)) {
			items = append(items, collections.NewKeyValue[int, string](item.ID, item.Data))
		}
		collections.Shuffle(items)

		t.Run(fmt.Sprintf("Single/Size(%v)/Cache(%s)", size, cache), func(t *testing.B) {
			t.ReportMetric(float64(size), "size")
			t.ReportMetric(float64(cache)+1, "target")

			for i := 0; i < t.N; i++ {
				col, err := large.NewBuckettedMap[int, string](
					size,
					hash.IntegerHasher[int](hash.MD5),
					large.WithCacheTarget[collections.KeyValue[int, string]](cache),
					large.WithMaxBucketSize(100),
				)
				require.NoError(t, err)

				for _, item := range items {
					col.Set(item.Key(), item.Value())
				}

				for _, item := range items {
					v, ok := col.Get(item.Key())
					if !ok {
						t.Log("item not found", item)
						t.Fail()
					}
					if v.Value() != item.Value() {
						t.Log("item not equal", v, item)
						t.Fail()
					}
					if v.Key() != item.Key() {
						t.Log("key not equal", v, item)
						t.Fail()
					}
				}
			}
		})

		t.Run(fmt.Sprintf("Concurrency/Size(%v)/Cache(%s)", size, cache), func(t *testing.B) {
			t.ReportMetric(float64(size), "size")
			t.ReportMetric(float64(cache)+1, "target")

			for i := 0; i < t.N; i++ {
				col, err := large.NewBuckettedMap[int, string](
					size,
					hash.IntegerHasher[int](hash.MD5),
					large.WithCacheTarget[collections.KeyValue[int, string]](cache),
					large.WithMaxBucketSize(100),
				)
				require.NoError(t, err)

				splitWithOverlap(col, items)

				for _, item := range items {
					v, ok := col.Get(item.Key())
					if !ok {
						t.Fail()
					}
					if v.Value() != item.Value() {
						t.Fail()
					}
					if v.Key() != item.Key() {
						t.Fail()
					}
				}
			}
		})
	})
}

func Benchmark_BuckettedMap_Set(t *testing.B) {
	sizes := []uint64{100, 200, 300, 400, 1000}
	bucketSizes := []uint64{1, 2, 3, 4}

	test_util.Case2(sizes, bucketSizes, func(size uint64, bucketSize uint64) {
		items := make([]collections.KeyValue[int, string], 0, int(size))
		for _, item := range test_util.Generate(int(size)) {
			items = append(items, collections.NewKeyValue[int, string](item.ID, item.Data))
		}
		collections.Shuffle(items)

		t.Run(fmt.Sprintf("Size(%v)/BucketSize(x%v)", size, bucketSize), func(t *testing.B) {
			t.ReportMetric(float64(size), "size")
			t.ReportMetric(float64(bucketSize), "bucket")

			for i := 0; i < t.N; i++ {
				col, err := large.NewBuckettedMap[int, string](
					size,
					hash.IntegerHasher[int](hash.MD5),
					large.WithBucketSize(size * bucketSize),
					large.WithMaxBucketSize(size * bucketSize),
				)
				require.NoError(t, err)

				splitWithOverlap(col, items)

				for _, item := range items {
					v, ok := col.Get(item.Key())
					if !ok {
						t.Fail()
					}
					if v.Value() != item.Value() {
						t.Fail()
					}
					if v.Key() != item.Key() {
						t.Fail()
					}
				}
			}
		})
	})
}
