package logger

import (
	"os"
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"WARN", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"ERROR", LevelError},
		{"none", LevelNone},
		{"NONE", LevelNone},
		{"invalid", LevelInfo}, // Default to INFO
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLogLevel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseLogLevel(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{LevelNone, "NONE"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("LogLevel.String() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestNoOpLogger(t *testing.T) {
	logger := &NoOpLogger{}

	// These should not panic
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	logger.SetLevel(LevelDebug)

	if logger.GetLevel() != LevelNone {
		t.Errorf("NoOpLogger.GetLevel() = %v, want %v", logger.GetLevel(), LevelNone)
	}
}

func TestConsoleLogger(t *testing.T) {
	logger := NewConsoleLogger(LevelInfo)

	if logger.GetLevel() != LevelInfo {
		t.Errorf("ConsoleLogger.GetLevel() = %v, want %v", logger.GetLevel(), LevelInfo)
	}

	logger.SetLevel(LevelWarn)
	if logger.GetLevel() != LevelWarn {
		t.Errorf("After SetLevel, ConsoleLogger.GetLevel() = %v, want %v", logger.GetLevel(), LevelWarn)
	}

	// These should not panic
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestConsoleLoggerWithPrefix(t *testing.T) {
	logger := NewConsoleLoggerWithPrefix(LevelDebug, "TestAgent")

	// These should not panic and should include the prefix
	logger.Debug("debug with prefix")
	logger.Info("info with prefix")
	logger.Warn("warn with prefix")
	logger.Error("error with prefix")
}

func TestConsoleLoggerFiltering(t *testing.T) {
	// Set level to WARN - should only log WARN and ERROR
	logger := NewConsoleLogger(LevelWarn)

	// These should not be logged (we can't easily test output without capturing it)
	logger.Debug("should not appear")
	logger.Info("should not appear")

	// These should be logged
	logger.Warn("should appear")
	logger.Error("should appear")
}

func TestGetLoggerFromEnv(t *testing.T) {
	// Test with no env var set
	os.Unsetenv("SNIP_LOG_LEVEL")
	logger := GetLoggerFromEnv()
	if _, ok := logger.(*NoOpLogger); !ok {
		t.Errorf("GetLoggerFromEnv() with no env var should return NoOpLogger")
	}

	// Test with "none"
	os.Setenv("SNIP_LOG_LEVEL", "none")
	logger = GetLoggerFromEnv()
	if _, ok := logger.(*NoOpLogger); !ok {
		t.Errorf("GetLoggerFromEnv() with SNIP_LOG_LEVEL=none should return NoOpLogger")
	}

	// Test with "info"
	os.Setenv("SNIP_LOG_LEVEL", "info")
	logger = GetLoggerFromEnv()
	if consoleLogger, ok := logger.(*ConsoleLogger); ok {
		if consoleLogger.GetLevel() != LevelInfo {
			t.Errorf("GetLoggerFromEnv() with SNIP_LOG_LEVEL=info should have level INFO")
		}
	} else {
		t.Errorf("GetLoggerFromEnv() with SNIP_LOG_LEVEL=info should return ConsoleLogger")
	}

	// Clean up
	os.Unsetenv("SNIP_LOG_LEVEL")
}

func TestGetLoggerFromEnvWithPrefix(t *testing.T) {
	os.Setenv("SNIP_LOG_LEVEL", "debug")
	logger := GetLoggerFromEnvWithPrefix("MyAgent")

	if consoleLogger, ok := logger.(*ConsoleLogger); ok {
		if consoleLogger.prefix != "MyAgent" {
			t.Errorf("GetLoggerFromEnvWithPrefix() prefix = %s, want MyAgent", consoleLogger.prefix)
		}
		if consoleLogger.GetLevel() != LevelDebug {
			t.Errorf("GetLoggerFromEnvWithPrefix() level = %v, want DEBUG", consoleLogger.GetLevel())
		}
	} else {
		t.Errorf("GetLoggerFromEnvWithPrefix() should return ConsoleLogger")
	}

	// Clean up
	os.Unsetenv("SNIP_LOG_LEVEL")
}

func TestConsoleLoggerWithArgs(t *testing.T) {
	logger := NewConsoleLogger(LevelInfo)

	// Test with formatted arguments
	logger.Info("Hello %s, number: %d", "world", 42)
	logger.Error("Error code: %d, message: %s", 500, "Internal Server Error")
}
