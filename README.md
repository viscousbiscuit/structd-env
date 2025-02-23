# structd-env

A **simple** and **blazingly fast** 🔥 package for managing environment variables in a structured way.

This package:

- **Has zero external dependencies**
- **Leverages generics and minimal reflection** to map environment variables to a struct
- **Caches results** to eliminate reflection overhead on subsequent calls
- **Supports a `.env.json` file** for local development
- **Handles multiple casing styles** automatically (snake_case, PascalCase, dash-case, camelCase)
- **Allows custom struct tags** (`env` tag) for fine-tuned variable mapping

### Supported Types

The package supports the following types for environment variable mapping:

- `bool`
- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`

## Installation

```sh
go get github.com/viscousbiscuit/structd-env

```

```go
package main

import (
	"fmt"
	structdEnv "github.com/viscousbiscuit/structd-env"
)

type MyType struct {
	FirstName    string
	LastName     string `env:"LAST_NAME"`
	BusinessName string
	Age          uint    `env:"AGE"`
	NetWorth     float64
	Active       bool    `env:"IS_ACTIVE"`
}

func main() {
 	se, err := structdEnv.GetInstance[MyTestType]()
    if err != nil {
        fmt.Println(err)
    }
    env := se.Get()

	fmt.Println("First Name:", env.FirstName)
	fmt.Println("Last Name:", env.LastName)
	fmt.Println("Age:", env.Age)
	fmt.Println("Net Worth:", env.NetWorth)
	fmt.Println("Active:", env.Active)
}

```

By default, the package attempts to match environment variables to struct fields using a best-effort approach. If multiple environment variables exist with the same name but different casing, use the `env` tag to specify which one to use.

The `env` tag allows you to explicitly define the environment variable name for a struct field. If no `env` tag is provided, the field name will be used as the default.

## `.env.json` Support

If an `.env.json` file is present, its values will override environment variables. Only top-level key-value pairs are supported—nested objects are not.

### Example `.env.json` file:

```json
{
  "FIRST_NAME": "Tom",
  "LAST_NAME": "Smith",
  "AGE": 30
}
```
