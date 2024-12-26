package growable_tests

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/large"
	test_util "github.com/daanv2/go-cache/test/util"
	"github.com/stretchr/testify/require"
)

func Test_BuckettedSet(t *testing.T) {
	sizes := []uint64{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrenty(%v)", size), func(t *testing.T) {
			col, err := large.NewBuckettedSet[*test_util.TestItem](size*10, test_util.Hasher())
			require.NoError(t, err)

			items := test_util.Generate(int(size))
			test_util.Shuffle(items)

			for _, item := range items {
				v, ok := col.GetOrAdd(item)
				require.True(t, ok)
				require.Equal(t, v, item)
			}
		})
	}
}

func Test_BuckettedSet_Concurrency(t *testing.T) {
	sizes := []uint64{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrenty(%v)", size), func(t *testing.T) {
			col, err := large.NewBuckettedSet[*test_util.TestItem](size*10, test_util.Hasher())
			require.NoError(t, err)

			items := test_util.Generate(int(size))
			test_util.Shuffle(items)

			splitWithOverlap(col, items)
			check := make(map[int]int, size)

			for item := range col.Read() {
				check[item.ID] = check[item.ID] + 1
				require.LessOrEqual(t, check[item.ID], 1)
			}
		})
	}
}

func Benchmark_BuckettedSet_Concurrency(t *testing.B) {
	sizes := []uint64{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrenty(%v)", size), func(t *testing.B) {
			for i := 0; i < t.N; i++ {
				col, err := large.NewBuckettedSet[*test_util.TestItem](size*10, test_util.Hasher())
				require.NoError(t, err)

				items := test_util.Generate(int(size))
				test_util.Shuffle(items)

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
