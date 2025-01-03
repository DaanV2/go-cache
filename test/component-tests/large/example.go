package large

import (
	"github.com/daanv2/go-cache/large"
	"github.com/daanv2/go-cache/pkg/hash"
)

func Example(size uint64) (*large.BuckettedMap[int, string], error) {
	col, err := large.NewBuckettedMap[int, string](size*10, hash.IntegerHasher[int](hash.MD5))

	return col, err
}
