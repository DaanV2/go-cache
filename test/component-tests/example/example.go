package example_tests

import (
	"github.com/daanv2/go-cache/maps"
	"github.com/daanv2/go-cache/pkg/hash"
)

func Example(size uint64) (*maps.Bucketted[int, string], error) {
	col, err := maps.NewBuckettedMap[int, string](size*10, hash.IntegerHasher[int](hash.MD5))

	return col, err
}
