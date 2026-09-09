# go-config-kit

Composable, typed configuration resolution for Go.

The current implementation resolves values from environment variables, supports typed conversion, and aggregates configuration errors.

## Installation

```bash
go get github.com/liddiard-research/go-config-kit@latest
```

Requires Go 1.27 or later.

## Usage

```go
package main

import (
	"log"
	"time"

	"github.com/liddiard-research/go-config-kit"
)

type Config struct {
	DatabaseURL string
	Port        int
	Timeout     time.Duration
	Debug       bool
}

func main() {
	cfg, err := makeConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func makeConfig() (Config, error) {
	var cfg Config

	err := config.Pack(
		config.Env("DATABASE_URL").Into(&cfg.DatabaseURL),
		config.Env("PORT").As[int]().Into(&cfg.Port),
		config.Env("TIMEOUT").As[time.Duration]().Into(&cfg.Timeout),
		config.Env("DEBUG").Optional().As[bool]().Into(&cfg.Debug),
	)
	
	return cfg, err
}
```

Environment variables are required by default. Use `Optional()` to ignore a missing value and preserve the destination's existing value. An environment variable set to an empty string is considered present.

## Typed values

Use `As[T]()` or `IntoAs()` to parse an environment variable into a supported type:

```go
var port uint16
var timeout time.Duration

err := config.Pack(
	config.Env("PORT").As[uint16]().Into(&port),
	config.Env("TIMEOUT").IntoAs(&timeout),
)
```

Supported values include strings, booleans, integers, unsigned integers, floating-point numbers, complex numbers, `time.Duration`, and types implementing `encoding.TextUnmarshaler`.

## JSON values

Use `Json[T]()` to decode a JSON environment variable:

```go
type Database struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var database Database

err := config.Pack(
	config.Env("DATABASE").Json[Database]().Into(&database),
)
```

```bash
export DATABASE='{"host":"localhost","port":5432}'
```

## Alternative variables

`OneOfGroup` requires exactly one resolver in the group to match:

```go
var databaseURL string

err := config.Pack(
	config.OneOfGroup(
		"database URL",
		config.Env("DATABASE_URL").Into(&databaseURL),
		config.Env("LEGACY_DATABASE_URL").Into(&databaseURL),
	),
)
```

The group fails when neither variable is present, when a matched value is invalid, or when multiple variables are present.

## Errors

`Pack` executes every resolution and returns all failures as a single error. The returned errors support standard `errors.Is` and `errors.As` inspection.

```go
if errors.Is(err, config.ErrNotFound) {
	// At least one required value was not found.
}

if errors.Is(err, config.ErrAmbiguousGroup) {
	// Multiple values matched a one-of group.
}
```

## Licence

Released under the [MIT Licence](LICENSE).

Copyright © 2026 Liddiard Research Limited.
