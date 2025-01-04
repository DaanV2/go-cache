package benchmarks

import (
	"runtime"
	"sync"
	"testing"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-optimal/pkg/cpu"
)

func CacheTargets() []cpu.CacheKind {
	return []cpu.CacheKind{cpu.CacheL1}
}

var (
	procs  = runtime.GOMAXPROCS(0) * 10
	buffer = max(procs, 1000)
)

func PumpConcurrent[T any](set collections.Set[T], items []T) {
	wg := &sync.WaitGroup{}
	pump := make(chan T, buffer)

	go transferThenClose(pump, items)

	for range procs {
		wg.Add(1)

		go PumpIntoSynced(wg, set, pump)
	}

	wg.Wait()
}

func transferThenClose[T any](pump chan <- T, items []T) {
	l := len(items)
	step := max(l * procs, 10)
	wg := &sync.WaitGroup{}

	for i := 0; i < l; i += step {
		wg.Add(1)

		go transfer(wg, pump, items[i:min(i + step, l)])
	}

	wg.Wait()
	close(pump)
}

func transfer[T any](wg *sync.WaitGroup, pump chan <- T, items []T) {
	defer wg.Done()

	for _, item := range items {
		pump <- item
	}
}

func PumpSync[T any](set collections.Set[T], items []T) {
	for _, item := range items {
		_, _ = set.GetOrAdd(item)
	}
}

func PumpInto[T any](set collections.Set[T], pump <-chan T) {
	for item := range pump {
		_, _ = set.GetOrAdd(item)
	}
}

func Validate[T comparable](t *testing.B, set collections.Set[T], items []T) {
	check := make(map[T]bool, len(items))

	for item := range set.Read() {
		_, ok := check[item]
		if ok {
			t.Logf("Duplicate of %v", item)
			t.Fail()
		}

		check[item] = false
	}

	for _, item := range items {
		v, ok := check[item]
		if !ok {
			t.Logf("Item was not defined %v", item)
			t.Fail()
		}
		if v {
			t.Log("Item was already inserted", item)
			t.Fail()
		}

		check[item] = true
	}
}

func ReportAdd(t *testing.B, size uint64) {
	inserts := float64(size) * float64(t.N)
	ns := float64(t.Elapsed().Nanoseconds())

	t.ReportMetric(inserts / ns, "inserts/ns")
	t.ReportMetric(inserts, "inserts")
	t.ReportMetric(ns, "ns")
}

func PumpIntoSynced[T any](wg *sync.WaitGroup, set collections.Set[T], pump <-chan T) {
	defer wg.Done()

	PumpInto(set, pump)
}
