package fixed_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/daanv2/go-cache/fixed"
	"github.com/stretchr/testify/require"
)

func Test_Slice(t *testing.T) {
	col := fixed.NewSlice[int32](7)

	require.Equal(t, col.Cap(), 7)
	require.Equal(t, col.UnsafeCap(), 7)
	require.Equal(t, col.Len(), 0)
	require.Equal(t, col.UnsafeLen(), 0)

	require.Equal(t, col.TryAppend(1), 1)
	require.Equal(t, col.TryAppend(2), 1)

	require.Equal(t, col.Cap(), 7)
	require.Equal(t, col.UnsafeCap(), 7)
	require.Equal(t, col.Len(), 2)
	require.Equal(t, col.UnsafeLen(), 2)

	require.Equal(t, col.TryAppend(3, 4, 5, 6), 4)

	require.Equal(t, col.Cap(), 7)
	require.Equal(t, col.UnsafeCap(), 7)
	require.Equal(t, col.Len(), 6)
	require.Equal(t, col.UnsafeLen(), 6)

	require.Equal(t, col.TryAppend(3, 4, 5, 6), 1)

	require.Equal(t, col.Cap(), 7)
	require.Equal(t, col.UnsafeCap(), 7)
	require.Equal(t, col.Len(), 7)
	require.Equal(t, col.UnsafeLen(), 7)
}

func Test_Slice_Parralel(t *testing.T) {
	col := fixed.NewSlice[int32](90)
	wg := sync.WaitGroup{}

	adds := make(chan int, 200)
	done := make(chan struct{})
	total := &atomic.Uint64{}

	go func() {
		for s := range adds {
			total.Add(uint64(s))
		}

		done <- struct{}{}
	}()

	for range 10 {
		wg.Add(1)

		// Add a whole bunch that should overload
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			for range 10 {
				adds <- col.TryAppend(1, 2, 3, 4)
			}
		}(&wg)
	}

	wg.Wait()
	close(adds)
	<- done

	require.EqualValues(t, col.Cap(), 90)
	require.EqualValues(t, col.Len(), 90)
	require.EqualValues(t, total.Load(), 90)
}
