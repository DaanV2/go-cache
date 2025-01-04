package large_test

import (
	"runtime"
	"sync"

	"github.com/daanv2/go-cache/collections"
)

func pumpConcurrent[T any](set collections.Set[T], items []T) {
	procs := runtime.GOMAXPROCS(0)*10
	buffer := max(procs, 1000)

	wg := &sync.WaitGroup{}
	pump := make(chan T, buffer)

	go func (items []T, pump chan <- T)  {
		for _, item := range items {
			pump <- item
		}
		close(pump)
	}(items, pump)

	for range procs {
		wg.Add(1)

		go pumpIntoSynced(wg, set, pump)
	}

	wg.Wait()
}

func pumpInto[T any](set collections.Set[T], pump <- chan T) {
	for item := range pump {
		_, _ = set.GetOrAdd(item)
	}
}

func pumpIntoSynced[T any](wg *sync.WaitGroup, set collections.Set[T], pump <- chan T) {
	defer wg.Done()

	pumpInto(set, pump)
}
