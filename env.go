package config

import (
	"fmt"
	"os"
)

var (
	ErrMandatoryEnvironmentVariable = NewNotFoundError("missing required environment variable")
)

func Env(key string) StrResolver {
	return NewStrResolver(func(mapper Mapper[string]) error {
		value, ok := os.LookupEnv(key)
		if !ok {
			return fmt.Errorf("%w: %s", ErrMandatoryEnvironmentVariable, key)
		}

		return mapper(value)
	})
}
