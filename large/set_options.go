package large

import (
	"github.com/daanv2/go-cache/pkg/options"
)

type SetOptions struct {
	Subitems int
}

func WithSubBufferSize(amount int) options.Option[SetOptions] {
	return options.NewFunction[SetOptions](func(option *SetOptions) {
		option.Subitems = amount
	})
}
