package fixed_test

import (
	"testing"

	"github.com/daanv2/go-cache/fixed"
	"github.com/stretchr/testify/require"
)

func Test_Map(t *testing.T) {
	amount := uint64(41)
	col := fixed.NewMap[uint64, uint64](amount + 10)

	newItem := func(id, v uint64) fixed.KeyValue[uint64, uint64] {
		return fixed.NewKeyValue(
			id,
			id,
			v,
		)
	}

	// Can set
	for i := range amount {
		ok := col.Set(newItem(i, i))
		require.True(t, ok, i)
	}

	// Can get
	for i := range amount {
		item, ok := col.Get(newItem(i, i))
		require.True(t, ok, i)
		require.EqualValues(t, item.Value, i)
	}

	// Can set again
	for i := range amount {
		ok := col.Set(newItem(i, i))
		require.True(t, ok, i)
	}

	// Check for duplicates
	check := make(map[uint64]bool, amount)
	for item := range col.Read() {
		v, ok := check[item.Key]
		require.False(t, ok, "item was duplicated %v", item)
		require.False(t, v, "item was duplicated %v", item)

		check[item.Key] = true
	}
}
