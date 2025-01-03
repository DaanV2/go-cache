package large_test

import (
	"runtime"
	"sync"

	"github.com/daanv2/go-cache/collections"
)

func splitWithOverlap[T any](set collections.Set[T], items []T) {
	l := len(items)
	sections := l / max(runtime.GOMAXPROCS(0)*10, 10)
	step := max(sections/2, 1)

	wg := &sync.WaitGroup{}

	for i := 0; i < l; i += step {
		wg.Add(1)
		subitems := items[i:min(l, i+sections)]

		go addToColsynced(wg, set, subitems)
	}

	wg.Wait()
}

func addToCol[T any](set collections.Set[T], items []T) {
	for _, item := range items {
		_, _ = set.GetOrAdd(item)
	}
}

func addToColsynced[T any](wg *sync.WaitGroup, set collections.Set[T], items []T) {
	defer wg.Done()

	addToCol(set, items)
}
