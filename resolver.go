package config

import "errors"

type Resolver[T any, Self any] struct {
	resolve func(Mapper[T]) error
	wrap    func(Resolver[T, Self]) Self
}

func (r Resolver[T, Self]) MapString(
	mapper func(T) (string, error),
) StrResolver {
	return NewStrResolver(func(next Mapper[string]) error {
		return r.resolve(func(value T) error {
			mapped, err := mapper(value)
			if err != nil {
				return err
			}

			return next(mapped)
		})
	})
}

func (r Resolver[T, Self]) Map[U any](
	mapper func(T) (U, error),
) TypedResolver[U] {
	return NewTypedResolver(func(next Mapper[U]) error {
		return r.resolve(func(value T) error {
			mapped, err := mapper(value)
			if err != nil {
				return err
			}

			return next(mapped)
		})
	})
}

func (r Resolver[T, Self]) Into(dst *T) Resolution {
	return func() error {
		return r.resolve(func(value T) error {
			*dst = value
			return nil
		})
	}
}

func (r Resolver[T, Self]) Mandatory() Self {
	next := r

	next.resolve = func(mapper Mapper[T]) error {
		err := r.resolve(mapper)

		if errors.Is(err, ErrNotFound) {
			return err
		}

		return err
	}

	return r.wrap(next)
}

func (r Resolver[T, Self]) Optional() Self {
	next := r

	next.resolve = func(mapper Mapper[T]) error {
		err := r.resolve(mapper)

		if errors.Is(err, ErrNotFound) {
			return nil
		}

		return err
	}

	return r.wrap(next)
}

func Mandatory[R interface {
	Mandatory() R
}](resolver R) R {
	return resolver.Mandatory()
}
