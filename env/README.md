# Env Package

The `env` package provides utility functions for working with environment variables in Go.

## Overview

This package offers simple, convenient functions for retrieving environment variables with fallback default values when the variable is not set or is empty.

## Functions

### GetEnvOrDefault

Retrieves an environment variable value, returning a custom default value if the variable is not set or is empty.

**Signature:**
```go
func GetEnvOrDefault(key, defaultValue string) string
```

**Parameters:**
- `key` - The name of the environment variable to retrieve
- `defaultValue` - The value to return if the environment variable is not set or is empty

**Returns:**
- The environment variable value if set and non-empty, otherwise the default value

**Example:**
```go
package main

import (
    "fmt"
    "github.com/snipwise/snip-sdk/env"
)

func main() {
    // Get PORT environment variable or default to "8080"
    port := env.GetEnvOrDefault("PORT", "8080")
    fmt.Println("Server will run on port:", port)
    // Output: Server will run on port: 8080 (if PORT is not set)

    // Get database host or use localhost
    dbHost := env.GetEnvOrDefault("DB_HOST", "localhost")
    fmt.Println("Database host:", dbHost)
    // Output: Database host: localhost (if DB_HOST is not set)

    // Get application environment
    appEnv := env.GetEnvOrDefault("APP_ENV", "development")
    fmt.Println("Running in:", appEnv)
    // Output: Running in: development (if APP_ENV is not set)
}
```

## Common Use Cases

### Server Configuration

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/snipwise/snip-sdk/env"
)

func main() {
    host := env.GetEnvOrDefault("SERVER_HOST", "0.0.0.0")
    port := env.GetEnvOrDefault("SERVER_PORT", "8080")
    addr := fmt.Sprintf("%s:%s", host, port)

    fmt.Printf("Starting server on %s\n", addr)
    http.ListenAndServe(addr, nil)
}
```

### Database Configuration

```go
package main

import (
    "fmt"
    "github.com/snipwise/snip-sdk/env"
)

type DBConfig struct {
    Host     string
    Port     string
    Database string
    User     string
    SSLMode  string
}

func main() {
    config := DBConfig{
        Host:     env.GetEnvOrDefault("DB_HOST", "localhost"),
        Port:     env.GetEnvOrDefault("DB_PORT", "5432"),
        Database: env.GetEnvOrDefault("DB_NAME", "myapp"),
        User:     env.GetEnvOrDefault("DB_USER", "postgres"),
        SSLMode:  env.GetEnvOrDefault("DB_SSLMODE", "disable"),
    }

    connStr := fmt.Sprintf(
        "host=%s port=%s dbname=%s user=%s sslmode=%s",
        config.Host, config.Port, config.Database, config.User, config.SSLMode,
    )

    fmt.Println("Connection string:", connStr)
}
```

### Application Settings

```go
package main

import (
    "fmt"
    "github.com/snipwise/snip-sdk/env"
)

func main() {
    // Application settings with sensible defaults
    logLevel := env.GetEnvOrDefault("LOG_LEVEL", "info")
    apiKey := env.GetEnvOrDefault("API_KEY", "")
    timeout := env.GetEnvOrDefault("REQUEST_TIMEOUT", "30s")

    fmt.Printf("Log Level: %s\n", logLevel)
    fmt.Printf("API Key set: %v\n", apiKey != "")
    fmt.Printf("Timeout: %s\n", timeout)
}
```

### Feature Flags

```go
package main

import (
    "fmt"
    "github.com/snipwise/snip-sdk/env"
)

func main() {
    // Use environment variables for feature flags
    enableCache := env.GetEnvOrDefault("FEATURE_CACHE", "true")
    enableMetrics := env.GetEnvOrDefault("FEATURE_METRICS", "false")
    enableDebug := env.GetEnvOrDefault("FEATURE_DEBUG", "false")

    fmt.Printf("Cache enabled: %s\n", enableCache)
    fmt.Printf("Metrics enabled: %s\n", enableMetrics)
    fmt.Printf("Debug enabled: %s\n", enableDebug)
}
```

## Behavior

- **Empty string handling**: If an environment variable is set but empty (`""`), the function returns the default value
- **Non-existent variable**: If an environment variable doesn't exist, the function returns the default value
- **Set and non-empty**: If an environment variable is set and non-empty, its value is returned

## Usage Notes

- This function only works with string values. For typed values (int, bool, float), combine with the `conversion` package:

```go
import (
    "github.com/snipwise/snip-sdk/env"
    "github.com/snipwise/snip-sdk/conversion"
)

// Get integer from environment
port := conversion.StringToIntOrDefault(
    env.GetEnvOrDefault("PORT", "8080"),
    8080,
)

// Get boolean from environment
debug := conversion.StringToBoolOrDefault(
    env.GetEnvOrDefault("DEBUG", "false"),
    false,
)

// Get float from environment
timeout := conversion.StringToFloatOrDefault(
    env.GetEnvOrDefault("TIMEOUT", "30.5"),
    30.5,
)
```

- For sensitive values like API keys or passwords, consider using empty string as default and checking if the value was actually set:

```go
apiKey := env.GetEnvOrDefault("API_KEY", "")
if apiKey == "" {
    log.Fatal("API_KEY environment variable must be set")
}
```

## Installation

```bash
go get github.com/snipwise/snip-sdk/env
```

## Related Packages

- [`conversion`](../conversion/README.md) - String conversion utilities that work well with environment variables

## License

See the main repository for license information.
