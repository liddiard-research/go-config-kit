package config_test

import (
	"os"
	"testing"

	"github.com/liddiard-research/go-config-kit"
	"github.com/stretchr/testify/require"
)

func TestEnv(t *testing.T) {
	t.Run("Resolve environment variable", func(t *testing.T) {
		t.Setenv("NAME", "Foo")

		var output string

		require.NoError(t, config.Env("NAME").Into(&output)())
		require.Equal(t, "Foo", output)
	})

	t.Run("Resolve empty environment variable", func(t *testing.T) {
		t.Setenv("NAME", "")

		output := "existing"

		require.NoError(t, config.Env("NAME").Into(&output)())
		require.Equal(t, "", output)
	})

	t.Run("Return not found when environment variable is missing", func(t *testing.T) {
		unsetenv(t, "NAME")

		var output string

		err := config.Env("NAME").Into(&output)()

		require.Error(t, err)
		require.ErrorIs(t, err, config.ErrNotFound)
		require.ErrorIs(t, err, config.ErrMandatoryEnvironmentVariable)
	})
}

func TestEnvMandatory(t *testing.T) {
	t.Run("Resolve mandatory environment variable", func(t *testing.T) {
		t.Setenv("NAME", "Foo")

		var output string

		require.NoError(t,
			config.Env("NAME").
				Mandatory().
				Into(&output)(),
		)

		require.Equal(t, "Foo", output)
	})

	t.Run("Return error when mandatory environment variable is missing", func(t *testing.T) {
		unsetenv(t, "NAME")

		var output string

		err := config.Env("NAME").
			Mandatory().
			Into(&output)()

		require.Error(t, err)
		require.ErrorIs(t, err, config.ErrMandatoryEnvironmentVariable)
		require.ErrorIs(t, err, config.ErrNotFound)
	})

	t.Run("Resolve mandatory environment variable using function", func(t *testing.T) {
		t.Setenv("NAME", "Foo")

		var output string

		require.NoError(t,
			config.Mandatory(
				config.Env("NAME"),
			).
				Into(&output)(),
		)

		require.Equal(t, "Foo", output)
	})

	t.Run("Return error using mandatory function when environment variable is missing", func(t *testing.T) {
		unsetenv(t, "NAME")

		var output string

		err := config.Mandatory(
			config.Env("NAME"),
		).
			Into(&output)()

		require.Error(t, err)
		require.ErrorIs(t, err, config.ErrMandatoryEnvironmentVariable)
		require.ErrorIs(t, err, config.ErrNotFound)
	})
}

func TestEnvOptional(t *testing.T) {
	t.Run("Resolve optional environment variable", func(t *testing.T) {
		t.Setenv("NAME", "Foo")

		var output string

		require.NoError(t,
			config.Env("NAME").
				Optional().
				Into(&output)(),
		)

		require.Equal(t, "Foo", output)
	})

	t.Run("Ignore missing optional environment variable", func(t *testing.T) {
		unsetenv(t, "NAME")

		output := "existing"

		require.NoError(t,
			config.Env("NAME").
				Optional().
				Into(&output)(),
		)

		require.Equal(t, "existing", output)
	})

	t.Run("Resolve empty optional environment variable", func(t *testing.T) {
		t.Setenv("NAME", "")

		output := "existing"

		require.NoError(t,
			config.Env("NAME").
				Optional().
				Into(&output)(),
		)

		require.Equal(t, "", output)
	})
}

func unsetenv(t *testing.T, key string) {
	t.Helper()

	value, exists := os.LookupEnv(key)

	t.Cleanup(func() {
		if exists {
			require.NoError(t, os.Setenv(key, value))
			return
		}

		require.NoError(t, os.Unsetenv(key))
	})

	require.NoError(t, os.Unsetenv(key))
}
