package hashmark_test

import (
	"testing"

	hashmark "github.com/daanv2/go-cache/pkg/hash/marked"
	"github.com/stretchr/testify/require"
)

func Test_HashMarked(t *testing.T) {
	v := []uint64{
		1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024,
	}

	c := hashmark.Empty()
	require.True(t, hashmark.IsEmpty(c))

	for _, i := range v {
		for _, j := range v {
			a := hashmark.MarkedHash(i)
			b := hashmark.MarkedHash(j)

			require.True(t, hashmark.IsFilled(a))
			require.True(t, hashmark.IsFilled(b))

			if i == j {
				require.Equal(t, a, b)
				require.True(t, hashmark.Equal(a, b))
			} else {
				require.NotEqual(t, a, b)
				require.False(t, hashmark.Equal(a, b))
			}

			require.NotEqual(t, a, c)
			require.False(t, hashmark.Equal(a, c))
			require.NotEqual(t, b, c)
			require.False(t, hashmark.Equal(b, c))
		}
	}
}
