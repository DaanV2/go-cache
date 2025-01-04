package large_test

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-cache/large"
	test_util "github.com/daanv2/go-cache/test/util"
	"github.com/stretchr/testify/require"
)

func Test_BuckettedSet(t *testing.T) {
	sizes := []uint64{100, 200, 300, 400, 1000, 10000, 20000}

	test_util.Case1(sizes, func(size uint64) {
		col, err := large.NewBuckettedSet[*test_util.TestItem](size*10, test_util.Hasher())
		require.NoError(t, err)

		items := test_util.Generate(int(size))
		collections.Shuffle(items)

		t.Run(fmt.Sprintf("Size(%d)", size), func(t *testing.T) {
			for _, item := range items {
				v, ok := col.GetOrAdd(item)
				require.True(t, ok)
				require.Equal(t, v, item)
			}
		})
	})
}

func Test_BuckettedSet_Concurrency(t *testing.T) {
	sizes := []uint64{100, 200, 300, 400, 1000, 10000, 20000}

	test_util.Case1(sizes, func(size uint64) {
		col, err := large.NewBuckettedSet[*test_util.TestItem](size*10, test_util.Hasher())
		require.NoError(t, err)

		items := test_util.Generate(int(size))
		collections.Shuffle(items)

		t.Run(fmt.Sprintf("Concurrency(%v)", size), func(t *testing.T) {
			pumpConcurrent(col, items)
			check := make(map[int]int, size)

			for item := range col.Read() {
				check[item.ID] = check[item.ID] + 1
				require.LessOrEqual(t, check[item.ID], 1)
			}
		})
	})
}
