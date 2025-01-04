package large_test

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/large"
	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/daanv2/go-cache/test/benchmarks"
	test_util "github.com/daanv2/go-cache/test/util"
	"github.com/daanv2/go-optimal"
	"github.com/stretchr/testify/require"
)

func Benchmark_BuckettedMap_BucketSize(t *testing.B) {
	sizes := []uint64{50_000}
	baseBucket := uint64(optimal.SliceSize[collections.KeyValue[int, string]]())

	bucketSizes := []uint64{
		baseBucket - 25, baseBucket - 10,
		baseBucket,
		baseBucket + 50, baseBucket + 100,
	}

	test_util.Case2(sizes, bucketSizes, func(size uint64, bucketSize uint64) {
		// Setup
		items := make([]collections.KeyValue[int, string], 0, int(size))
		for _, item := range test_util.Generate(int(size)) {
			items = append(items, collections.NewKeyValue[int, string](item.ID, item.Data))
		}
		collections.Shuffle(items)

		t.Run(fmt.Sprintf("(%d)/Size(%d)/Single", bucketSize, size), func(t *testing.B) {
			// Report
			t.ReportMetric(float64(size), "size")
			t.ReportMetric(float64(bucketSize), "items/bucket")

			for i := 0; i < t.N; i++ {
				col, err := large.NewBuckettedMap[int, string](
					size,
					hash.IntegerHasher[int](hash.MD5),
					large.WithBucketSize(bucketSize),
				)
				require.NoError(t, err)

				benchmarks.PumpSync(col, items)
			}

			benchmarks.ReportAdd(t, size)
		})

		t.Run(fmt.Sprintf("(%d)/Size(%d)/Single/Reuse", bucketSize, size), func(t *testing.B) {
			// Report
			t.ReportMetric(float64(size), "size")
			t.ReportMetric(float64(bucketSize), "items/bucket")

			col, err := large.NewBuckettedMap[int, string](
				size,
				hash.IntegerHasher[int](hash.MD5),
				large.WithBucketSize(bucketSize),
			)
			require.NoError(t, err)

			for i := 0; i < t.N; i++ {
				benchmarks.PumpSync(col, items)
			}

			benchmarks.ReportAdd(t, size)
		})

		t.Run(fmt.Sprintf("(%d)/Size(%d)/Concur", bucketSize, size), func(t *testing.B) {
			// Report
			t.ReportMetric(float64(size), "size")
			t.ReportMetric(float64(bucketSize), "items/bucket")

			for i := 0; i < t.N; i++ {
				col, err := large.NewBuckettedMap[int, string](
					size,
					hash.IntegerHasher[int](hash.MD5),
					large.WithBucketSize(bucketSize),
				)
				require.NoError(t, err)

				benchmarks.PumpConcurrent(col, items)
			}

			benchmarks.ReportAdd(t, size)
		})

		t.Run(fmt.Sprintf("(%d)/Size(%d)/Concur/Reuse", bucketSize, size), func(t *testing.B) {
			// Report
			t.ReportMetric(float64(size), "size")
			t.ReportMetric(float64(bucketSize), "items/bucket")

			col, err := large.NewBuckettedMap[int, string](
				size,
				hash.IntegerHasher[int](hash.MD5),
				large.WithBucketSize(bucketSize),
			)
			require.NoError(t, err)

			for i := 0; i < t.N; i++ {
				benchmarks.PumpConcurrent(col, items)
			}

			benchmarks.ReportAdd(t, size)
		})
	})
}
