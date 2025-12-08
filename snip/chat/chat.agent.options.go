package chat

import (
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

// AgentServerOption defines a functional option for configuring the agent
type ChatAgentOption func(*ChatAgent)

// WithLogger sets a custom logger for the agent
func WithLogger(log logger.Logger) ChatAgentOption {
	return func(a *ChatAgent) {
		a.logger = log
	}
}

// WithVerbose enables verbose logging (INFO level) with agent name prefix
func WithVerbose(verbose bool) ChatAgentOption {
	return func(a *ChatAgent) {
		if verbose {
			a.logger = logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, a.Name)
		} else {
			a.logger = &logger.NoOpLogger{}
		}
	}
}

// WithLogLevel sets the log level for the agent
func WithLogLevel(level logger.LogLevel) ChatAgentOption {
	return func(a *ChatAgent) {
		a.logger = logger.NewConsoleLoggerWithPrefix(level, a.Name)
	}
}
