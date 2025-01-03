package large_test

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/large"
	test_util "github.com/daanv2/go-cache/test/util"
	"github.com/stretchr/testify/require"
)

func Benchmark_GrowableSet_Concurrency(t *testing.B) {
	sizes := []int{100, 200, 300, 400, 1000, 10000, 20000}

	test_util.Case1(sizes, func(size int) {
		items := test_util.Generate(size)
		collections.Shuffle(items)

		t.Run(fmt.Sprintf("Size(%v)", size), func(t *testing.B) {
			for i := 0; i < t.N; i++ {
				col, err := large.NewGrowableSet[*test_util.TestItem](test_util.Hasher())
				require.NoError(t, err)

				splitWithOverlap(col, items)
				check := make(map[int]int, size)

				for item := range col.Read() {
					check[item.ID] = check[item.ID] + 1
					require.LessOrEqual(t, check[item.ID], 1)
				}
			}
		})
	})
}
