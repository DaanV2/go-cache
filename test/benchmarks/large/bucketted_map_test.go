package large_test

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/large"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/test/benchmarks"
	test_util "github.com/daanv2/go-cache/test/util"
	"github.com/daanv2/go-optimal/pkg/cpu"
	"github.com/stretchr/testify/require"
)

func Benchmark_BuckettedMap_Add(t *testing.B) {
	t.SkipNow()
	sizes := []uint64{50_000, 100_000}
	target := benchmarks.CacheTargets()

	test_util.Case2(sizes, target, func(size uint64, cache cpu.CacheKind) {
		// Setup
		items := make([]collections.KeyValue[int, string], 0, int(size))
		for _, item := range test_util.Generate(int(size)) {
			items = append(items, collections.NewKeyValue[int, string](item.ID, item.Data))
		}
		collections.Shuffle(items)

		t.Run(fmt.Sprintf("Size(%d)/Cache(%s)/Single", size, cache), func(t *testing.B) {
			// Report
			t.ReportMetric(float64(size), "size")

			for i := 0; i < t.N; i++ {
				col, err := large.NewBuckettedMap[int, string](
					size,
					hash.IntegerHasher[int](hash.MD5),
					large.WithCacheTarget[collections.KeyValue[int, string]](cache),
				)
				require.NoError(t, err)

				benchmarks.PumpSync(col, items)
			}

			benchmarks.ReportAdd(t, size)
		})

		t.Run(fmt.Sprintf("Size(%d)/Cache(%s)/Single/Reuse", size, cache), func(t *testing.B) {
			// Report
			t.ReportMetric(float64(size), "size")

			col, err := large.NewBuckettedMap[int, string](
				size,
				hash.IntegerHasher[int](hash.MD5),
				large.WithCacheTarget[collections.KeyValue[int, string]](cache),
			)
			require.NoError(t, err)

			for i := 0; i < t.N; i++ {
				benchmarks.PumpSync(col, items)
			}

			benchmarks.ReportAdd(t, size)
		})

		t.Run(fmt.Sprintf("Size(%d)/Cache(%s)/Concur", size, cache), func(t *testing.B) {
			// Report
			t.ReportMetric(float64(size), "size")

			for i := 0; i < t.N; i++ {
				col, err := large.NewBuckettedMap[int, string](
					size,
					hash.IntegerHasher[int](hash.MD5),
					large.WithCacheTarget[collections.KeyValue[int, string]](cache),
				)
				require.NoError(t, err)

				benchmarks.PumpConcurrent(col, items)
			}

			benchmarks.ReportAdd(t, size)
		})

		t.Run(fmt.Sprintf("Size(%d)/Cache(%s)/Concur/Reuse", size, cache), func(t *testing.B) {
			// Report
			t.ReportMetric(float64(size), "size")

			col, err := large.NewBuckettedMap[int, string](
				size,
				hash.IntegerHasher[int](hash.MD5),
				large.WithCacheTarget[collections.KeyValue[int, string]](cache),
			)
			require.NoError(t, err)

			for i := 0; i < t.N; i++ {
				benchmarks.PumpConcurrent(col, items)
			}

			benchmarks.ReportAdd(t, size)
		})
	})
}

func Benchmark_BuckettedMap_Get(t *testing.B) {
	t.SkipNow()
	sizes := []uint64{100, 200, 300, 400, 1000, 5000, 10_000, 50_000, 100_000}
	target := benchmarks.CacheTargets()

	test_util.Case2(sizes, target, func(size uint64, cache cpu.CacheKind) {
		t.Run(fmt.Sprintf("Size(%d)/Cache(%s)", size, cache), func(b *testing.B) {
			// Report
			t.ReportMetric(float64(size), "size")

			// Setup
			items := make([]collections.KeyValue[int, string], 0, int(size))
			for _, item := range test_util.Generate(int(size)) {
				items = append(items, collections.NewKeyValue[int, string](item.ID, item.Data))
			}
			collections.Shuffle(items)
			col, err := large.NewBuckettedMap[int, string](
				size,
				hash.IntegerHasher[int](hash.MD5),
				large.WithCacheTarget[collections.KeyValue[int, string]](cache),
			)
			require.NoError(t, err)

			benchmarks.PumpSync(col, items)

			t.Run("Single", func(t *testing.B) {
				for i := 0; i < t.N; i++ {

				}
			})

			t.Run("Concur", func(t *testing.B) {
				for i := 0; i < t.N; i++ {

				}
			})
		})
	})
}
