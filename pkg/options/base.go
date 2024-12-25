package options

type Option[T any] interface {
	apply(option *T) error
}

var _ Option[struct{}] = NewFunction[struct{}](nil)

type OptionFN[T any] func(option *T) error

func (o OptionFN[T]) apply(option *T) error {
	return o(option)
}

func NewFunction[T any](modify func(option *T)) OptionFN[T] {
	return func(option *T) error {
		modify(option)
		return nil
	}
}

func NewFunctionE[T any](modify func(option *T) error) OptionFN[T] {
	return OptionFN[T](modify)
}

func Apply[T any](options *T, modifies ...Option[T]) error {
	for _, modify := range modifies {
		if err := modify.apply(options); err != nil {
			return err
		}
	}

	return nil
}
