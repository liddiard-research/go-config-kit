package config_test

import (
	"testing"

	. "github.com/liddiard-research/go-config-kit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOneOfGroup(t *testing.T) {
	value := func(v string) *string {
		return &v
	}

	run := func(value1 *string, value2 *string) (uint, error) {
		var output uint

		resolver := func(value *string) StrResolver {
			return NewStrResolver(func(m Mapper[string]) error {
				if value == nil {
					return ErrNotFound
				}

				return m(*value)
			})
		}

		err := OneOfGroup(
			"Age",
			resolver(value1).As[uint]().Into(&output),
			resolver(value2).As[uint]().Into(&output),
		)()

		return output, err
	}

	t.Run("error when no resolver matches", func(t *testing.T) {
		_, err := run(nil, nil)

		require.Error(t, err)
	})

	t.Run("first resolver succeeds", func(t *testing.T) {
		output, err := run(value("20"), nil)

		require.NoError(t, err)
		assert.Equal(t, uint(20), output)
	})

	t.Run("second resolver succeeds", func(t *testing.T) {
		output, err := run(nil, value("30"))

		require.NoError(t, err)
		assert.Equal(t, uint(30), output)
	})

	t.Run("error when multiple resolvers succeed", func(t *testing.T) {
		_, err := run(value("20"), value("30"))

		require.ErrorIs(t, err, ErrAmbiguousGroup)
	})

	t.Run("error when resolver matches but fails", func(t *testing.T) {
		_, err := run(value("invalid"), nil)

		var groupErr *GroupError
		require.ErrorAs(t, err, &groupErr)
		require.NotErrorIs(t, err, ErrAmbiguousGroup)
	})

	t.Run("error when one resolver succeeds and another fails", func(t *testing.T) {
		_, err := run(value("invalid"), value("30"))

		require.ErrorIs(t, err, ErrAmbiguousGroup)
	})
}
