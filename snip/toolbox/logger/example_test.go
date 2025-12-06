package logger_test

import (
	"fmt"
	"os"

	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

// Example of using the NoOp logger (default, no output)
func ExampleNoOpLogger() {
	log := &logger.NoOpLogger{}

	log.Debug("This will not appear")
	log.Info("This will not appear")
	log.Warn("This will not appear")
	log.Error("This will not appear")

	fmt.Println("NoOp logger produces no output")
	// Output: NoOp logger produces no output
}

// Example of using console logger with different levels
func ExampleConsoleLogger() {
	// Create a console logger at INFO level
	log := logger.NewConsoleLogger(logger.LevelInfo)

	// This won't appear (below INFO level)
	log.Debug("Debug message - won't appear")

	// These will appear
	log.Info("Application started")
	log.Warn("Using default configuration")

	// Change level at runtime
	log.SetLevel(logger.LevelError)

	// Now only errors will appear
	log.Info("This won't appear anymore")
	log.Error("Critical error occurred")
}

// Example of using console logger with prefix
func ExampleNewConsoleLoggerWithPrefix() {
	log := logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, "MyApp")

	log.Info("Server started on port 8080")
	log.Warn("Cache miss for key: user-123")
	log.Error("Database connection failed")
}

// Example of using environment variable
func ExampleGetLoggerFromEnv() {
	// Set environment variable
	os.Setenv("SNIP_LOG_LEVEL", "info")
	defer os.Unsetenv("SNIP_LOG_LEVEL")

	log := logger.GetLoggerFromEnv()

	log.Info("Logger configured from environment")
	log.Debug("This won't appear (level is INFO)")
}

// Example of using environment variable with prefix
func ExampleGetLoggerFromEnvWithPrefix() {
	// Set environment variable
	os.Setenv("SNIP_LOG_LEVEL", "debug")
	defer os.Unsetenv("SNIP_LOG_LEVEL")

	log := logger.GetLoggerFromEnvWithPrefix("ChatAgent")

	log.Debug("Detailed debug information")
	log.Info("Processing request")
	log.Warn("Slow response detected")
	log.Error("Request failed")
}

// Example of formatted logging
func ExampleConsoleLogger_formatting() {
	log := logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, "API")

	userID := "user-123"
	duration := 250

	log.Info("Request received from user: %s", userID)
	log.Warn("Request took %d ms (threshold: 200ms)", duration)
	log.Error("Failed to authenticate user %s", userID)
}

// Example showing different log levels
func ExampleLogLevel() {
	levels := []logger.LogLevel{
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelNone,
	}

	for _, level := range levels {
		fmt.Printf("Level: %s\n", level.String())
	}

	// Output:
	// Level: DEBUG
	// Level: INFO
	// Level: WARN
	// Level: ERROR
	// Level: NONE
}

// Example of parsing log level from string
func ExampleParseLogLevel() {
	levels := []string{"debug", "INFO", "warn", "ERROR", "invalid"}

	for _, levelStr := range levels {
		level := logger.ParseLogLevel(levelStr)
		fmt.Printf("%s -> %s\n", levelStr, level.String())
	}

	// Output:
	// debug -> DEBUG
	// INFO -> INFO
	// warn -> WARN
	// ERROR -> ERROR
	// invalid -> INFO
}
