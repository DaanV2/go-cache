package fixed_test

import (
	"testing"

	"github.com/daanv2/go-cache/fixed"
	"github.com/stretchr/testify/require"
)

func Test_Set(t *testing.T) {
	amount := uint64(41)
	col := fixed.NewSet[uint64](amount + 10)

	// Can set
	for i := range amount {
		ok := col.Set(fixed.NewSetItem[uint64](i, i))
		require.True(t, ok, i)
	}

	// Can get
	for i := range amount {
		item, ok := col.Get(fixed.NewSetItem[uint64](i, i))
		require.True(t, ok, i)
		require.EqualValues(t, item.Value, i)
	}

	// Can set again
	for i := range amount {
		ok := col.Set(fixed.NewSetItem[uint64](i, i))
		require.True(t, ok, i)
	}

	// Check for duplicates
	check := make(map[uint64]bool, amount)
	for item := range col.Read() {
		v, ok := check[item.Value]
		require.False(t, ok, "item was duplicated %v", item)
		require.False(t, v, "item was duplicated %v", item)

		check[item.Value] = true
	}
}
