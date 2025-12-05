# Snip SDK (for Docker Model Runner)

> S.N.I.P. Smart Neural Intelligence Partner

## Testing

The project includes comprehensive unit tests for all packages.

### Run All Tests

```bash
# Run all tests
go test ./...

# Run all tests with verbose output
go test ./... -v

# Run tests without cache
go test ./... -count=1
```

### Run Tests for Specific Package

```bash
# Test the smart package
go test ./smart -v

# Test the conversion package
go test ./conversion -v

# Test the files package
go test ./files -v

# Test the env package
go test ./env -v
```

### Run Tests with Coverage

```bash
# Generate coverage report for all packages
go test ./... -cover

# Generate detailed coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Benchmarks

```bash
# Run all benchmarks
go test ./... -bench=.

# Run benchmarks for specific package
go test ./env -bench=.
go test ./files -bench=.
go test ./conversion -bench=.
```

### Test Statistics

The project currently includes:
- **200+ unit tests** covering all packages
- **~3000 lines of test code**
- Tests for edge cases, error handling, and integration scenarios
- Benchmarks for performance-critical functions

## Examples

See the [samples](samples/) directory for example usage of the SDK.

