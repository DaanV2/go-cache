package maps

import (
	"github.com/daanv2/go-cache/pkg/options"
	"github.com/daanv2/go-locks"
	optimal "github.com/daanv2/go-optimal"
	"github.com/daanv2/go-optimal/pkg/cpu"
)

// Options is the base struct for all sets.
type Options struct {
	bucket_size      uint64
	items_lock       *locks.Pool
	bucket_amount    uint64
	bucket_amount_fn func(uint64) uint64
}

// CreateOptions creates a new instance of SetBase with the default bucket size.
func CreateOptions[T any](opts ...options.Option[Options]) (Options, error) {
	op := Options{
		bucket_size: uint64(optimal.SliceSize[T]()),
		items_lock:  locks.NewPool(),
		bucket_amount: 0,
		bucket_amount_fn: nil,
	}

	err := options.Apply(&op, opts...)

	return op, err
}

func (o Options) BucketAmount(capacity uint64) uint64 {
	if o.bucket_amount_fn != nil {
		return o.bucket_amount_fn(capacity)
	}

	amount := capacity / uint64(max(o.bucket_size, 1))

	return max(amount, 10)
}

// WithBucketSize sets the size of buckets that the set will use
func WithBucketSize(size uint64) options.Option[Options] {
	return options.NewFunction[Options](func(option *Options) {
		option.bucket_size = size
	})
}

// WithMaxBucketSize sets the size of buckets if it is larger than the a certain size
func WithMaxBucketSize(size uint64) options.Option[Options] {
	return options.NewFunction[Options](func(option *Options) {
		option.bucket_size = min(option.bucket_size, size)
	})
}

// WithItemLocks sets the locks that the set will use
func WithItemLocks(pool *locks.Pool) options.Option[Options] {
	return options.NewFunction(func(option *Options) {
		option.items_lock = pool
	})
}

// WithCacheTarget sets the cache target for the set
func WithCacheTarget[T any](target cpu.CacheKind) options.Option[Options] {
	return options.NewFunction(func(option *Options) {
		option.bucket_size = uint64(optimal.SliceSizeFor[T](target))
	})
}

// WithBucketAmount sets the amount of buckets that the set will use
func WithBucketAmount(amount uint64) options.Option[Options] {
	return options.NewFunction[Options](func(option *Options) {
		option.bucket_amount = amount
	})
}

// WithBucketAmount sets the amount of buckets that the set will use
func WithBucketFunction(calc func(capacity uint64) uint64) options.Option[Options] {
	return options.NewFunction[Options](func(option *Options) {
		option.bucket_amount_fn = calc
	})
}
