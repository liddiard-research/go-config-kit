package config_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/liddiard-research/go-config-kit"
	"github.com/stretchr/testify/require"
)

func strResolverFixture(value string) config.StrResolver {
	return config.NewStrResolver(func(mapper config.Mapper[string]) error {
		return mapper(value)
	})
}

func TestStrResolver(t *testing.T) {
	t.Run("IntoAs string", func(t *testing.T) {
		var output string

		require.NoError(t, strResolverFixture("hello world").IntoAs(&output)())
		require.Equal(t, "hello world", output)
	})

	t.Run("IntoAs defined string type", func(t *testing.T) {
		type Table string

		var output Table

		require.NoError(t, strResolverFixture("hello world").IntoAs(&output)())
		require.Equal(t, Table("hello world"), output)
	})

	t.Run("As string", func(t *testing.T) {
		var output string

		require.NoError(t,
			strResolverFixture("hello world").
				As[string]().
				Into(&output)(),
		)

		require.Equal(t, "hello world", output)
	})

	t.Run("As defined string type", func(t *testing.T) {
		type Table string

		var output Table

		require.NoError(t,
			strResolverFixture("users").
				As[Table]().
				Into(&output)(),
		)

		require.Equal(t, Table("users"), output)
	})

	t.Run("As bool", func(t *testing.T) {
		var output bool

		require.NoError(t, strResolverFixture("true").IntoAs(&output)())
		require.True(t, output)
	})

	t.Run("As int", func(t *testing.T) {
		var output int

		require.NoError(t, strResolverFixture("-42").IntoAs(&output)())
		require.Equal(t, -42, output)
	})

	t.Run("As int8", func(t *testing.T) {
		var output int8

		require.NoError(t, strResolverFixture("-42").IntoAs(&output)())
		require.Equal(t, int8(-42), output)
	})

	t.Run("As int16", func(t *testing.T) {
		var output int16

		require.NoError(t, strResolverFixture("-420").IntoAs(&output)())
		require.Equal(t, int16(-420), output)
	})

	t.Run("As int32", func(t *testing.T) {
		var output int32

		require.NoError(t, strResolverFixture("-42000").IntoAs(&output)())
		require.Equal(t, int32(-42000), output)
	})

	t.Run("As int64", func(t *testing.T) {
		var output int64

		require.NoError(t, strResolverFixture("-42000").IntoAs(&output)())
		require.Equal(t, int64(-42000), output)
	})

	t.Run("As defined int type", func(t *testing.T) {
		type Count int

		var output Count

		require.NoError(t, strResolverFixture("42").IntoAs(&output)())
		require.Equal(t, Count(42), output)
	})

	t.Run("As uint", func(t *testing.T) {
		var output uint

		require.NoError(t, strResolverFixture("42").IntoAs(&output)())
		require.Equal(t, uint(42), output)
	})

	t.Run("As uint8", func(t *testing.T) {
		var output uint8

		require.NoError(t, strResolverFixture("255").IntoAs(&output)())
		require.Equal(t, uint8(255), output)
	})

	t.Run("As uint16", func(t *testing.T) {
		var output uint16

		require.NoError(t, strResolverFixture("65535").IntoAs(&output)())
		require.Equal(t, uint16(65535), output)
	})

	t.Run("As uint32", func(t *testing.T) {
		var output uint32

		require.NoError(t, strResolverFixture("42000").IntoAs(&output)())
		require.Equal(t, uint32(42000), output)
	})

	t.Run("As uint64", func(t *testing.T) {
		var output uint64

		require.NoError(t, strResolverFixture("42000").IntoAs(&output)())
		require.Equal(t, uint64(42000), output)
	})

	t.Run("As defined uint type", func(t *testing.T) {
		type Age uint

		var output Age

		require.NoError(t, strResolverFixture("20").IntoAs(&output)())
		require.Equal(t, Age(20), output)
	})

	t.Run("As float32", func(t *testing.T) {
		var output float32

		require.NoError(t, strResolverFixture("12.5").IntoAs(&output)())
		require.Equal(t, float32(12.5), output)
	})

	t.Run("As float64", func(t *testing.T) {
		var output float64

		require.NoError(t, strResolverFixture("12.5").IntoAs(&output)())
		require.Equal(t, 12.5, output)
	})

	t.Run("As complex64", func(t *testing.T) {
		var output complex64

		require.NoError(t, strResolverFixture("1+2i").IntoAs(&output)())
		require.Equal(t, complex64(1+2i), output)
	})

	t.Run("As complex128", func(t *testing.T) {
		var output complex128

		require.NoError(t, strResolverFixture("1+2i").IntoAs(&output)())
		require.Equal(t, complex128(1+2i), output)
	})

	t.Run("As duration", func(t *testing.T) {
		var output time.Duration

		require.NoError(t, strResolverFixture("1h30m").IntoAs(&output)())
		require.Equal(t, 90*time.Minute, output)
	})

	t.Run("As text unmarshaler", func(t *testing.T) {
		var output testEnvironment

		require.NoError(t, strResolverFixture("production").IntoAs(&output)())
		require.Equal(t, testEnvironmentProduction, output)
	})
}

type testEnvironment string

const (
	testEnvironmentDevelopment testEnvironment = "development"
	testEnvironmentProduction  testEnvironment = "production"
)

func (e *testEnvironment) UnmarshalText(value []byte) error {
	switch string(value) {
	case string(testEnvironmentDevelopment):
		*e = testEnvironmentDevelopment
	case string(testEnvironmentProduction):
		*e = testEnvironmentProduction
	default:
		return fmt.Errorf("invalid environment %q", value)
	}

	return nil
}

func TestStrResolverErrors(t *testing.T) {
	t.Run("Invalid bool", func(t *testing.T) {
		var output bool

		require.Error(t, strResolverFixture("invalid").IntoAs(&output)())
	})

	t.Run("Invalid int", func(t *testing.T) {
		var output int

		require.Error(t, strResolverFixture("invalid").IntoAs(&output)())
	})

	t.Run("Invalid uint", func(t *testing.T) {
		var output uint

		require.Error(t, strResolverFixture("-1").IntoAs(&output)())
	})

	t.Run("Integer overflow", func(t *testing.T) {
		var output int8

		require.Error(t, strResolverFixture("128").IntoAs(&output)())
	})

	t.Run("Unsigned integer overflow", func(t *testing.T) {
		var output uint8

		require.Error(t, strResolverFixture("256").IntoAs(&output)())
	})

	t.Run("Invalid float", func(t *testing.T) {
		var output float64

		require.Error(t, strResolverFixture("invalid").IntoAs(&output)())
	})

	t.Run("Invalid complex", func(t *testing.T) {
		var output complex128

		require.Error(t, strResolverFixture("invalid").IntoAs(&output)())
	})

	t.Run("Invalid duration", func(t *testing.T) {
		var output time.Duration

		require.Error(t, strResolverFixture("invalid").IntoAs(&output)())
	})

	t.Run("Text unmarshaler error", func(t *testing.T) {
		var output testEnvironment

		require.Error(t, strResolverFixture("invalid").IntoAs(&output)())
	})

	t.Run("Unsupported type", func(t *testing.T) {
		var output struct {
			Name string
		}

		require.Error(t, strResolverFixture("hello").IntoAs(&output)())
	})
}
