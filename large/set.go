package large

import (
	"errors"

	"github.com/daanv2/go-cache/pkg/options"
	optimal "github.com/daanv2/go-optimal"
)

type SetItem[T any] struct {
	hash uint64
}

type Set[T any] struct {
	options SetOptions
}

func NewSet[T any](opts ...options.Option[SetOptions]) (*Set[T], error) {
	opt := SetOptions{
		Subitems: optimal.SliceSize[T](),
	}
	err := options.Apply(&opt, opts...)
	if err != nil {
		return nil, err
	}

	// Validate
	if opt.Subitems <= 1 {
		return nil, errors.New("sub items size is too small < 1")
	}

	return &Set[T]{
		options: opt,
	}, err
}
