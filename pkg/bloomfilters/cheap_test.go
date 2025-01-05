package bloomfilters_test

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/pkg/bloomfilters"
	"github.com/stretchr/testify/require"
)

func Test_Cheap(t *testing.T) {
	amounts := []uint64{
		64,
		128,
		127,
		129,
	}

	for _, amount := range amounts {
		t.Run(fmt.Sprintf("Set->Has(%v)", amount), func(t *testing.T) {
			filter := bloomfilters.NewCheap(uint64(amount))

			for i := range amount * 2 {
				filter.Set(i)

				ok := filter.Has(i)
				require.True(t, ok, i)
			}
		})
	}
}

func Test_Cheap_Duplicates(t *testing.T) {
	amount := uint64(128)
	filter := bloomfilters.NewCheap(uint64(amount))

	half := amount / 2

	for i := range half {
		filter.Set(i)

		ok := filter.Has(i)
		require.True(t, ok, i)
	}

	correct := 0
	for j := range half {
		i := j + 64
		if filter.Has(i) {
			correct++
		}
	}

	require.Less(t, correct, 64)
}
