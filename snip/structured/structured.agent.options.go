package structured

import (
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

// StructuredAgentOption defines a functional option for configuring StructuredAgent
type StructuredAgentOption[O any] func(*StructuredAgent[O])

// WithLogger sets a custom logger for the agent
func WithLogger[O any](log logger.Logger) StructuredAgentOption[O] {
	return func(structuredAgent *StructuredAgent[O]) {
		structuredAgent.logger = log
	}
}

// WithVerbose enables verbose logging (INFO level) with agent name prefix
func WithVerbose[O any](verbose bool) StructuredAgentOption[O] {
	return func(structuredAgent *StructuredAgent[O]) {
		if verbose {
			structuredAgent.logger = logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, structuredAgent.Name)
		} else {
			structuredAgent.logger = &logger.NoOpLogger{}
		}
	}
}

// WithLogLevel sets the log level for the agent
func WithLogLevel[O any](level logger.LogLevel) StructuredAgentOption[O] {
	return func(structuredAgent *StructuredAgent[O]) {
		structuredAgent.logger = logger.NewConsoleLoggerWithPrefix(level, structuredAgent.Name)
	}
}
