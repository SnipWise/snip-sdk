package chat

import (
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

// AgentOption defines a functional option for configuring LocalAIAgent
type AgentOption func(*ChatAgent)

// WithLogger sets a custom logger for the agent
func WithLogger(log logger.Logger) AgentOption {
	return func(a *ChatAgent) {
		a.logger = log
	}
}

// WithVerbose enables verbose logging (INFO level) with agent name prefix
func WithVerbose(verbose bool) AgentOption {
	return func(a *ChatAgent) {
		if verbose {
			a.logger = logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, a.Name)
		} else {
			a.logger = &logger.NoOpLogger{}
		}
	}
}

// WithLogLevel sets the log level for the agent
func WithLogLevel(level logger.LogLevel) AgentOption {
	return func(a *ChatAgent) {
		a.logger = logger.NewConsoleLoggerWithPrefix(level, a.Name)
	}
}
