package hash_test

import (
	"math"
	"testing"

	"github.com/daanv2/go-cache/pkg/hash"
	"github.com/stretchr/testify/require"
)

func Test_NewRange(t *testing.T) {
	r := hash.NewRange()

	// New ranges should not work for any value.
	require.False(t, r.Has(0))
	require.False(t, r.Has(math.MaxUint64))
	require.False(t, r.Has(1337))
}

func Test_Range_Update(t *testing.T) {
	r := hash.NewRange()

	r.Update(1337)
	require.True(t, r.Has(1337))
	require.False(t, r.Has(0))
	require.False(t, r.Has(math.MaxUint64))
}

func Test_Range_Update_Multiple(t *testing.T) {
	r := hash.NewRange()

	r.Update(1337)
	r.Update(42)
	r.Update(9001)
	require.True(t, r.Has(1337))
	require.True(t, r.Has(42))
	require.True(t, r.Has(9001))
	require.False(t, r.Has(0))
	require.False(t, r.Has(math.MaxUint64))
}
