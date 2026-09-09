package config

type TypedResolver[T any] struct {
	Resolver[T, TypedResolver[T]]
}

func NewTypedResolver[T any](
	resolve func(Mapper[T]) error,
) TypedResolver[T] {
	var wrap func(Resolver[T, TypedResolver[T]]) TypedResolver[T]

	wrap = func(
		resolver Resolver[T, TypedResolver[T]],
	) TypedResolver[T] {
		resolver.wrap = wrap

		return TypedResolver[T]{
			Resolver: resolver,
		}
	}

	return wrap(Resolver[T, TypedResolver[T]]{
		resolve: resolve,
	})
}
