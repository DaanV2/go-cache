package benchmarks

import (
	"runtime"
	"sync"
	"testing"

	"github.com/daanv2/go-cache/pkg/collections"
	"github.com/daanv2/go-optimal/pkg/cpu"
)

type GetOrAdd[T any] interface {
	GetOrAdd(item T) (T, bool)
}

type Map[K, V any] interface {
	Set(key K, value V) bool
}

type KeyValue[K, V any] interface {
	GetKey() K
	GetValue() V
}

func CacheTargets() []cpu.CacheKind {
	return []cpu.CacheKind{cpu.CacheL1}
}

var (
	procs  = runtime.GOMAXPROCS(0) * 10
	buffer = max(procs, 1000)
)

func PumpConcurrentSet[T any](set GetOrAdd[T], items []T) {
	wg := &sync.WaitGroup{}
	pump := make(chan T, buffer)

	go transferThenClose(pump, items)

	for range procs {
		wg.Add(1)

		go PumpIntoSyncedSet(wg, set, pump)
	}

	wg.Wait()
}

func PumpConcurrentMap[K, V comparable](set Map[K, V], items []KeyValue[K, V]) {
	wg := &sync.WaitGroup{}
	pump := make(chan KeyValue[K, V], buffer)

	go transferThenClose(pump, items)

	for range procs {
		wg.Add(1)

		go PumpIntoSyncedMap[K, V](wg, set, pump)
	}

	wg.Wait()
}

func transferThenClose[T any](pump chan<- T, items []T) {
	l := len(items)
	step := max(l*procs, 10)
	wg := &sync.WaitGroup{}

	for i := 0; i < l; i += step {
		wg.Add(1)

		go transfer(wg, pump, items[i:min(i+step, l)])
	}

	wg.Wait()
	close(pump)
}

func transfer[T any](wg *sync.WaitGroup, pump chan<- T, items []T) {
	defer wg.Done()

	for _, item := range items {
		pump <- item
	}
}

func PumpSyncSet[T any](set GetOrAdd[T], items []T) {
	for _, item := range items {
		_, _ = set.GetOrAdd(item)
	}
}

func PumpSyncMap[K, V comparable](set Map[K, V], items []KeyValue[K, V]) {
	for _, item := range items {
		_ = set.Set(item.GetKey(), item.GetValue())
	}
}

func PumpIntoSet[T any](set GetOrAdd[T], pump <-chan T) {
	for item := range pump {
		_, _ = set.GetOrAdd(item)
	}
}

func PumpIntoMap[K, V comparable](set Map[K, V], pump <-chan KeyValue[K, V]) {
	for item := range pump {
		_ = set.Set(item.GetKey(), item.GetValue())
	}
}

func Validate[T comparable](t *testing.B, set collections.Readable[T], items []T) {
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
	perS := (inserts / ns) * 1e9

	t.ReportMetric(perS, "inserts/s")
	t.ReportMetric(inserts, "inserts")
}

func PumpIntoSyncedSet[T any](wg *sync.WaitGroup, set GetOrAdd[T], pump <-chan T) {
	defer wg.Done()

	PumpIntoSet(set, pump)
}

func PumpIntoSyncedMap[K, V comparable](wg *sync.WaitGroup, set Map[K, V], pump <-chan KeyValue[K, V]) {
	defer wg.Done()

	PumpIntoMap(set, pump)
}
