package large

import (
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-locks"
	optimal "github.com/daanv2/go-optimal"
	"github.com/daanv2/go-optimal/pkg/cpu"
)

// WithBucketSize sets the amount of buckets that the set will use
func WithBucketSize(amount int) options.Option[SetBase] {
	return options.NewFunction[SetBase](func(option *SetBase) {
		option.bucket_size = amount
	})
}

// WithItemLocks sets the locks that the set will use
func WithItemLocks(pool *locks.Pool) options.Option[SetBase] {
	return options.NewFunction(func(option *SetBase) {
		option.items_lock = pool
	})
}


// WithCacheTarget sets the cache target for the set
func WithCacheTarget[T any](target cpu.CacheKind) options.Option[SetBase] {
	return options.NewFunction(func(option *SetBase) {
		option.bucket_size = optimal.SliceSizeFor[T](target)
	})
}
