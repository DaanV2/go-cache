package collections_test

import (
	"fmt"
	"testing"

	"github.com/daanv2/go-cache/collections"
	"github.com/daanv2/go-kit/generics"
	"github.com/stretchr/testify/require"
)

func Test_KeyValue_Value(t *testing.T) {
	test_keyvalue_item(t, 1, 100)
	test_keyvalue_item(t, "key", "value")
	test_keyvalue_item(t, struct{ A int }{A: 1}, struct{ B string }{B: "test"})
}

func test_keyvalue_item[K comparable, V comparable](t *testing.T, key K, value V) {
	title := fmt.Sprintf("key(%s), value(%s)", generics.NameOf[K](), generics.NameOf[V]())
	t.Run(title, func(t *testing.T) {
		item := collections.NewKeyValue[K, V](key, value)

		require.Equal(t, key, item.Key())
		require.Equal(t, value, item.Value())
	})

	t.Run(title+"-empty", func(t *testing.T) {
		item := collections.EmptyKeyValue[K, V]()

		require.Equal(t, generics.Empty[K](), item.Key())
		require.Equal(t, generics.Empty[V](), item.Value())
	})
}
