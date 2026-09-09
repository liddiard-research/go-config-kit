package config

import (
	"errors"
	"fmt"
)

var (
	ErrAmbiguousGroup = errors.New("one-of group is ambiguous")
)

type Mapper[T any] func(T) error

type Resolution func() error

func Pack(resolutions ...Resolution) error {
	var errs []error

	for _, resolution := range resolutions {
		if err := resolution(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return &MultiError{
		Header: "Error in pack:",
		Errs:   errs,
	}
}

func OneOfGroup(name string, resolutions ...Resolution) Resolution {
	return func() error {
		var errs []error

		matched := 0
		resolved := 0

		for _, resolution := range resolutions {
			err := resolution()

			switch {
			case err == nil:
				matched++
				resolved++

			case errors.Is(err, ErrNotFound):
				errs = append(errs, err)

			default:
				matched++
				errs = append(errs, err)
			}
		}

		switch {
		case matched == 0:
			return &GroupError{
				Name: name,
				Errs: errs,
			}

		case matched > 1:
			return &GroupError{
				Name: name,
				Errs: append(
					[]error{
						fmt.Errorf(
							"%w: %d configuration sources matched",
							ErrAmbiguousGroup,
							matched,
						),
					},
					errs...,
				),
			}

		case resolved == 1:
			return nil

		default:
			return &GroupError{
				Name: name,
				Errs: errs,
			}
		}
	}
}
