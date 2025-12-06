package env

import (
	"os"
	"testing"
)

// ============================================================================
// Tests for GetEnvOrDefault
// ============================================================================

func TestGetEnvOrDefault(t *testing.T) {
	t.Run("environment variable exists", func(t *testing.T) {
		// Set a test environment variable
		key := "TEST_ENV_VAR_EXISTS"
		expectedValue := "test_value_123"
		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key) // Clean up after test

		result := GetEnvOrDefault(key, "default_value")
		if result != expectedValue {
			t.Errorf("GetEnvOrDefault(%q, \"default_value\") = %q, want %q", key, result, expectedValue)
		}
	})

	t.Run("environment variable does not exist", func(t *testing.T) {
		key := "TEST_ENV_VAR_DOES_NOT_EXIST"
		defaultValue := "my_default_value"

		// Make sure the env var doesn't exist
		os.Unsetenv(key)

		result := GetEnvOrDefault(key, defaultValue)
		if result != defaultValue {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q", key, defaultValue, result, defaultValue)
		}
	})

	t.Run("environment variable is empty string", func(t *testing.T) {
		key := "TEST_ENV_VAR_EMPTY"
		defaultValue := "default_when_empty"

		// Set env var to empty string
		os.Setenv(key, "")
		defer os.Unsetenv(key)

		result := GetEnvOrDefault(key, defaultValue)
		// Empty string should return default value
		if result != defaultValue {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q (empty env var should use default)", key, defaultValue, result, defaultValue)
		}
	})

	t.Run("environment variable with whitespace", func(t *testing.T) {
		key := "TEST_ENV_VAR_WHITESPACE"
		envValue := "  value_with_spaces  "
		defaultValue := "default_value"

		os.Setenv(key, envValue)
		defer os.Unsetenv(key)

		result := GetEnvOrDefault(key, defaultValue)
		// Should return the env value as-is (with whitespace)
		if result != envValue {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q", key, defaultValue, result, envValue)
		}
	})

	t.Run("default value is empty string", func(t *testing.T) {
		key := "TEST_ENV_VAR_DEFAULT_EMPTY"
		defaultValue := ""

		os.Unsetenv(key)

		result := GetEnvOrDefault(key, defaultValue)
		if result != defaultValue {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q", key, defaultValue, result, defaultValue)
		}
	})

	t.Run("environment variable with special characters", func(t *testing.T) {
		key := "TEST_ENV_VAR_SPECIAL_CHARS"
		envValue := "value@#$%^&*()!~`"
		defaultValue := "default"

		os.Setenv(key, envValue)
		defer os.Unsetenv(key)

		result := GetEnvOrDefault(key, defaultValue)
		if result != envValue {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q", key, defaultValue, result, envValue)
		}
	})

	t.Run("environment variable with newlines", func(t *testing.T) {
		key := "TEST_ENV_VAR_NEWLINES"
		envValue := "line1\nline2\nline3"
		defaultValue := "default"

		os.Setenv(key, envValue)
		defer os.Unsetenv(key)

		result := GetEnvOrDefault(key, defaultValue)
		if result != envValue {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q", key, defaultValue, result, envValue)
		}
	})

	t.Run("numeric values as strings", func(t *testing.T) {
		key := "TEST_ENV_VAR_NUMERIC"
		envValue := "12345"
		defaultValue := "999"

		os.Setenv(key, envValue)
		defer os.Unsetenv(key)

		result := GetEnvOrDefault(key, defaultValue)
		if result != envValue {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q", key, defaultValue, result, envValue)
		}
	})

	t.Run("boolean-like values as strings", func(t *testing.T) {
		key := "TEST_ENV_VAR_BOOLEAN"
		envValue := "true"
		defaultValue := "false"

		os.Setenv(key, envValue)
		defer os.Unsetenv(key)

		result := GetEnvOrDefault(key, defaultValue)
		if result != envValue {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q", key, defaultValue, result, envValue)
		}
	})

	t.Run("unicode characters", func(t *testing.T) {
		key := "TEST_ENV_VAR_UNICODE"
		envValue := "Hello 世界 🌍"
		defaultValue := "default"

		os.Setenv(key, envValue)
		defer os.Unsetenv(key)

		result := GetEnvOrDefault(key, defaultValue)
		if result != envValue {
			t.Errorf("GetEnvOrDefault(%q, %q) = %q, want %q", key, defaultValue, result, envValue)
		}
	})
}

// ============================================================================
// Benchmark for GetEnvOrDefault
// ============================================================================

func BenchmarkGetEnvOrDefault_Exists(b *testing.B) {
	key := "BENCHMARK_ENV_EXISTS"
	os.Setenv(key, "benchmark_value")
	defer os.Unsetenv(key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEnvOrDefault(key, "default")
	}
}

func BenchmarkGetEnvOrDefault_NotExists(b *testing.B) {
	key := "BENCHMARK_ENV_NOT_EXISTS"
	os.Unsetenv(key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEnvOrDefault(key, "default")
	}
}
