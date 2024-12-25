package growable_tests

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/daanv2/go-cache/large"
	test_util "github.com/daanv2/go-cache/test/util"
	"github.com/stretchr/testify/require"
)

func Test_GrowableSet(t *testing.T) {
	sizes := []int{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrenty(%v)", size), func(t *testing.T) {
			col, err := large.NewGrowableSet[*test_util.TestItem](test_util.Hasher())
			require.NoError(t, err)

			items := test_util.Generate(size)
			test_util.Shuffle(items)

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
		t.Run(fmt.Sprintf("Concurrenty(%v)", size), func(t *testing.T) {
			col, err := large.NewGrowableSet[*test_util.TestItem](test_util.Hasher())
			require.NoError(t, err)

			items := test_util.Generate(size)
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

func Benchmark_GrowableSet_Concurrency(t *testing.B) {
	sizes := []int{100, 200, 300, 400, 1000, 10000, 20000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Concurrenty(%v)", size), func(t *testing.B) {
			for i := 0; i < t.N; i++ {
				col, err := large.NewGrowableSet[*test_util.TestItem](test_util.Hasher())
				require.NoError(t, err)

				items := test_util.Generate(size)
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

func splitWithOverlap(set *large.GrowableSet[*test_util.TestItem], items []*test_util.TestItem) {
	l := len(items)
	sections := l / max(runtime.GOMAXPROCS(0)*10, 10)
	step := max(sections/2, 1)

	wg := &sync.WaitGroup{}

	for i := 0; i < l; i += step {
		wg.Add(1)
		subitems := items[i:min(l, i+sections)]

		go addToCol(wg, set, subitems)
	}

	wg.Wait()
}

func addToCol(wg *sync.WaitGroup, set *large.GrowableSet[*test_util.TestItem], items []*test_util.TestItem) {
	defer wg.Done()

	for _, item := range items {
		_, _ = set.GetOrAdd(item)
	}
}
