package large

import (
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-locks"
)

func WithBucketSize(amount int) options.Option[SetBase] {
	return options.NewFunction[SetBase](func(option *SetBase) {
		option.bucket_size = amount
	})
}

func WithItemLocks[T any](pool *locks.Pool) options.Option[SetBase] {
	return options.NewFunction(func(option *SetBase) {
		option.items_lock = pool
	})
}
