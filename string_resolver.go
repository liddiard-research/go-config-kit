package config

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type StrResolver struct {
	Resolver[string, StrResolver]
}

func NewStrResolver(
	resolve func(Mapper[string]) error,
) StrResolver {
	var makeResolver func(func(Mapper[string]) error) StrResolver

	makeResolver = func(
		resolve func(Mapper[string]) error,
	) StrResolver {
		base := Resolver[string, StrResolver]{
			resolve: resolve,

			wrap: func(
				r Resolver[string, StrResolver],
			) StrResolver {
				return StrResolver{
					Resolver: r,
				}
			},
		}

		return StrResolver{
			Resolver: base,
		}
	}

	return makeResolver(resolve)
}

func (r StrResolver) As[U any]() TypedResolver[U] {
	return NewTypedResolver(func(mapper Mapper[U]) error {
		return r.resolve(func(value string) error {
			converted, err := parse[U](value)
			if err != nil {
				return err
			}

			return mapper(converted)
		})
	})
}

func (r StrResolver) IntoAs[U any](dst *U) Resolution {
	parsed := r.As[U]()
	return parsed.Into(dst)
}

func (r StrResolver) Into[U ~string](dst *U) Resolution {
	return func() error {
		return r.resolve(func(value string) error {
			*dst = U(value)
			return nil
		})
	}
}

func (r StrResolver) Json[U any]() TypedResolver[U] {
	return NewTypedResolver(func(mapper Mapper[U]) error {
		return r.resolve(func(value string) error {
			var dst U
			if err := json.Unmarshal([]byte(value), &dst); err != nil {
				return err
			}

			return mapper(dst)
		})
	})
}

var durationType = reflect.TypeFor[time.Duration]()

func parse[T any](input string) (T, error) {
	var output T

	if unmarshaler, ok := any(&output).(encoding.TextUnmarshaler); ok {
		if err := unmarshaler.UnmarshalText([]byte(input)); err != nil {
			return output, err
		}

		return output, nil
	}

	value := reflect.ValueOf(&output).Elem()

	if value.Type() == durationType {
		parsed, err := time.ParseDuration(input)
		if err != nil {
			return output, err
		}

		value.SetInt(int64(parsed))
		return output, nil
	}

	switch value.Kind() {
	case reflect.String:
		value.SetString(input)

	case reflect.Bool:
		parsed, err := strconv.ParseBool(input)
		if err != nil {
			return output, err
		}

		value.SetBool(parsed)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(
			input,
			10,
			value.Type().Bits(),
		)
		if err != nil {
			return output, err
		}

		value.SetInt(parsed)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		parsed, err := strconv.ParseUint(
			input,
			10,
			value.Type().Bits(),
		)
		if err != nil {
			return output, err
		}

		value.SetUint(parsed)

	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(
			input,
			value.Type().Bits(),
		)
		if err != nil {
			return output, err
		}

		value.SetFloat(parsed)

	case reflect.Complex64, reflect.Complex128:
		parsed, err := strconv.ParseComplex(
			input,
			value.Type().Bits(),
		)
		if err != nil {
			return output, err
		}

		value.SetComplex(parsed)

	default:
		return output, fmt.Errorf(
			"cannot parse string as %s",
			value.Type(),
		)
	}

	return output, nil
}
