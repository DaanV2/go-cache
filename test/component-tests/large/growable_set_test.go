package large_test

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/large"
	test_util "github.com/daanv2/go-cache/test/util"
	"github.com/stretchr/testify/require"
)

func Test_GrowableSet(t *testing.T) {
	sizes := []int{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrency(%v)", size), func(t *testing.T) {
			col, err := large.NewGrowableSet[*test_util.TestItem](test_util.Hasher())
			require.NoError(t, err)

			items := test_util.Generate(size)
			collections.Shuffle(items)

			for _, item := range items {
				v, ok := col.GetOrAdd(item)
				require.True(t, ok)
				require.Equal(t, v, item)
			}
		})
	}
}

func Test_GrowableSet_Concurrency(t *testing.T) {
	sizes := []int{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrency(%v)", size), func(t *testing.T) {
			col, err := large.NewGrowableSet[*test_util.TestItem](test_util.Hasher())
			require.NoError(t, err)

			items := test_util.Generate(size)
			collections.Shuffle(items)

			splitWithOverlap(col, items)
			check := make(map[int]int, size)

			for item := range col.Read() {
				check[item.ID] = check[item.ID] + 1
				require.LessOrEqual(t, check[item.ID], 1)
			}
		})
	}
}

func Benchmark_GrowableSet_Concurrency(t *testing.B) {
	sizes := []int{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrency(%v)", size), func(t *testing.B) {
			for i := 0; i < t.N; i++ {
				col, err := large.NewGrowableSet[*test_util.TestItem](test_util.Hasher())
				require.NoError(t, err)

				items := test_util.Generate(size)
				collections.Shuffle(items)

				splitWithOverlap(col, items)
				check := make(map[int]int, size)

				for item := range col.Read() {
					check[item.ID] = check[item.ID] + 1
					require.LessOrEqual(t, check[item.ID], 1)
				}
			}
		})
	}
}
