package iterators

import (
	"iter"
	"sync"
)

type Collection[T any] interface {
	// Read is an iterator over sequences of individual values.
	// When called as seq(yield), seq calls yield(v) for each value v in the sequence, stopping early if yield returns false. See the iter package documentation for more details.
	Read() iter.Seq[T]
}

// RangeCol goes over the collection and calls the function, stopping early if yield returns false
func RangeCol[C Collection[T], T any](col C, yield func(item T) bool) {
	for item := range col.Read() {
		if !yield(item) {
			return
		}
	}
}

func RangeColParralel[C Collection[T], T any](colls []C, yield func(item T) bool) {
	wg := &sync.WaitGroup{}
	// This can be done unsafely, its faster. but the only state being written to it, is false
	// This means go rountines might go an extra couple of loops before exiting, but it doesn't block as much per item
	continu := new(bool)
	*continu = false

	for _, col := range colls {
		wg.Add(1)
		go rangeColP[C, T](wg, continu, col, yield)
	}

	wg.Wait()
}

func rangeColP[C Collection[T], T any](wg *sync.WaitGroup, continu *bool, col C, yield func(item T) bool) {
	defer wg.Done()
	i := 0

	for item := range col.Read() {
		// Every 16 iterations, check the state has been marked stopped
		i++
		if i%16 == 0 {
			if !(*continu) {
				break
			}
		}

		if !yield(item) {
			*continu = false
			return
		}
	}
}
