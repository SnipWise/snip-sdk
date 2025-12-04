# Conversion Package

The `conversion` package provides utility functions for converting string values to other primitive types in Go.

## Overview

This package offers simple, safe conversion functions that handle errors gracefully. It provides three variants for each conversion type:

1. **Basic functions** (`StringToInt`, `StringToFloat`, `StringToBool`) - Print errors to stdout and return default values
2. **Error-returning functions** (`StringToIntErr`, `StringToFloatErr`, `StringToBoolErr`) - Return errors for explicit error handling
3. **Custom default functions** (`StringToIntOrDefault`, `StringToFloatOrDefault`, `StringToBoolOrDefault`) - Accept custom default values

## Functions

### Basic Conversion Functions

### StringToInt

Converts a string to an integer.

**Signature:**
```go
func StringToInt(str string) int
```

**Parameters:**
- `str` - The string to convert to an integer

**Returns:**
- An integer value, or `0` if conversion fails

**Example:**
```go
package main

import (
    "fmt"
    "github.com/snipwise/snip-sdk/conversion"
)

func main() {
    result := conversion.StringToInt("42")
    fmt.Println(result) // Output: 42

    invalid := conversion.StringToInt("abc")
    fmt.Println(invalid) // Output: 0 (with error message printed)
}
```

### StringToFloat

Converts a string to a 64-bit floating point number.

**Signature:**
```go
func StringToFloat(str string) float64
```

**Parameters:**
- `str` - The string to convert to a float

**Returns:**
- A float64 value, or `0.0` if conversion fails

**Example:**
```go
package main

import (
    "fmt"
    "github.com/snipwise/snip-sdk/conversion"
)

func main() {
    result := conversion.StringToFloat("3.14159")
    fmt.Println(result) // Output: 3.14159

    result2 := conversion.StringToFloat("42")
    fmt.Println(result2) // Output: 42.0

    invalid := conversion.StringToFloat("not-a-number")
    fmt.Println(invalid) // Output: 0.0 (with error message printed)
}
```

### StringToBool

Converts a string to a boolean value.

**Signature:**
```go
func StringToBool(str string) bool
```

**Parameters:**
- `str` - The string to convert to a boolean

**Returns:**
- A boolean value, or `false` if conversion fails

**Accepted Values:**
- `true`: "1", "t", "T", "true", "TRUE", "True"
- `false`: "0", "f", "F", "false", "FALSE", "False"

**Example:**
```go
package main

import (
    "fmt"
    "github.com/snipwise/snip-sdk/conversion"
)

func main() {
    result1 := conversion.StringToBool("true")
    fmt.Println(result1) // Output: true

    result2 := conversion.StringToBool("1")
    fmt.Println(result2) // Output: true

    result3 := conversion.StringToBool("false")
    fmt.Println(result3) // Output: false

    result4 := conversion.StringToBool("0")
    fmt.Println(result4) // Output: false

    invalid := conversion.StringToBool("maybe")
    fmt.Println(invalid) // Output: false (with error message printed)
}
```

---

### Error-Returning Functions

These functions return errors for explicit error handling, following Go's idiomatic error handling pattern.

#### StringToIntErr

Converts a string to an integer and returns an error if conversion fails.

**Signature:**
```go
func StringToIntErr(str string) (int, error)
```

**Example:**
```go
package main

import (
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/conversion"
)

func main() {
    result, err := conversion.StringToIntErr("42")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result) // Output: 42

    _, err = conversion.StringToIntErr("abc")
    if err != nil {
        fmt.Println("Error:", err) // Output: Error: cannot convert to int: ...
    }
}
```

#### StringToFloatErr

Converts a string to a float64 and returns an error if conversion fails.

**Signature:**
```go
func StringToFloatErr(str string) (float64, error)
```

**Example:**
```go
package main

import (
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/conversion"
)

func main() {
    result, err := conversion.StringToFloatErr("3.14159")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result) // Output: 3.14159

    _, err = conversion.StringToFloatErr("not-a-number")
    if err != nil {
        fmt.Println("Error:", err) // Output: Error: cannot convert to float: ...
    }
}
```

#### StringToBoolErr

Converts a string to a boolean and returns an error if conversion fails.

**Signature:**
```go
func StringToBoolErr(str string) (bool, error)
```

**Example:**
```go
package main

import (
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/conversion"
)

func main() {
    result, err := conversion.StringToBoolErr("true")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result) // Output: true

    _, err = conversion.StringToBoolErr("maybe")
    if err != nil {
        fmt.Println("Error:", err) // Output: Error: cannot convert to bool: ...
    }
}
```

---

### Custom Default Functions

These functions accept a custom default value to return when conversion fails. They don't print errors or return error values.

#### StringToIntOrDefault

Converts a string to an integer, returning a custom default value if conversion fails.

**Signature:**
```go
func StringToIntOrDefault(str string, defaultValue int) int
```

**Example:**
```go
package main

import (
    "fmt"
    "os"
    "github.com/snipwise/snip-sdk/conversion"
)

func main() {
    // Use environment variable or default to 8080
    port := conversion.StringToIntOrDefault(os.Getenv("PORT"), 8080)
    fmt.Println(port) // Output: 8080 (if PORT is not set or invalid)

    // Explicit conversion with custom default
    timeout := conversion.StringToIntOrDefault("invalid", 30)
    fmt.Println(timeout) // Output: 30
}
```

#### StringToFloatOrDefault

Converts a string to a float64, returning a custom default value if conversion fails.

**Signature:**
```go
func StringToFloatOrDefault(str string, defaultValue float64) float64
```

**Example:**
```go
package main

import (
    "fmt"
    "github.com/snipwise/snip-sdk/conversion"
)

func main() {
    // Use configuration value or default to 30.5
    timeout := conversion.StringToFloatOrDefault("invalid", 30.5)
    fmt.Println(timeout) // Output: 30.5

    rate := conversion.StringToFloatOrDefault("0.95", 1.0)
    fmt.Println(rate) // Output: 0.95
}
```

#### StringToBoolOrDefault

Converts a string to a boolean, returning a custom default value if conversion fails.

**Signature:**
```go
func StringToBoolOrDefault(str string, defaultValue bool) bool
```

**Example:**
```go
package main

import (
    "fmt"
    "github.com/snipwise/snip-sdk/conversion"
)

func main() {
    // Enable feature by default if config is invalid
    enabled := conversion.StringToBoolOrDefault("invalid", true)
    fmt.Println(enabled) // Output: true

    debug := conversion.StringToBoolOrDefault("false", true)
    fmt.Println(debug) // Output: false
}
```

---

## Error Handling

The package provides three approaches to error handling:

### Basic Functions (with stdout logging)
- `StringToInt`, `StringToFloat`, `StringToBool`
- Print error messages to stdout using `fmt.Println`
- Return default values: `0`, `0.0`, `false`

### Error-Returning Functions (explicit error handling)
- `StringToIntErr`, `StringToFloatErr`, `StringToBoolErr`
- Return `(value, error)` tuple
- Wrap errors using `fmt.Errorf` with `%w` for error unwrapping
- Allow explicit error handling with Go's idiomatic patterns

### Custom Default Functions (silent fallback)
- `StringToIntOrDefault`, `StringToFloatOrDefault`, `StringToBoolOrDefault`
- Don't print errors or return error values
- Accept custom default values
- Ideal for configuration values with sensible defaults

## Usage Notes

- Choose the variant that best fits your use case:
  - Use **basic functions** for simple scripts or when you want visible error logging
  - Use **error-returning functions** for production code requiring explicit error handling
  - Use **custom default functions** for configuration values with sensible fallbacks
- All error-returning functions use `%w` for error wrapping, compatible with `errors.Unwrap`
- Custom default functions are silent and don't produce any output on error

## Installation

```bash
go get github.com/snipwise/snip-sdk/conversion
```

## License

See the main repository for license information.
