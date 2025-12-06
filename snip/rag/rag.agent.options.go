package rag

import (
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

// AgentOption defines a functional option for configuring LocalAIAgent
type RagAgentOption func(*RagAgent)

// WithLogger sets a custom logger for the agent
func WithLogger(log logger.Logger) RagAgentOption {
	return func(a *RagAgent) {
		a.logger = log
	}
}

// WithVerbose enables verbose logging (INFO level) with agent name prefix
func WithVerbose(verbose bool) RagAgentOption {
	return func(a *RagAgent) {
		if verbose {
			a.logger = logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, a.Name)
		} else {
			a.logger = &logger.NoOpLogger{}
		}
	}
}

// WithLogLevel sets the log level for the agent
func WithLogLevel(level logger.LogLevel) RagAgentOption {
	return func(a *RagAgent) {
		a.logger = logger.NewConsoleLoggerWithPrefix(level, a.Name)
	}
}
