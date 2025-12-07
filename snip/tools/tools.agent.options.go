package tools

import (
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

// AgentOption defines a functional option for configuring LocalAIAgent
type ToolsAgentOption func(*ToolsAgent)

// WithLogger sets a custom logger for the agent
func WithLogger(log logger.Logger) ToolsAgentOption {
	return func(a *ToolsAgent) {
		a.logger = log
	}
}

// WithVerbose enables verbose logging (INFO level) with agent name prefix
func WithVerbose(verbose bool) ToolsAgentOption {
	return func(a *ToolsAgent) {
		if verbose {
			a.logger = logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, a.Name)
		} else {
			a.logger = &logger.NoOpLogger{}
		}
	}
}

// WithLogLevel sets the log level for the agent
func WithLogLevel(level logger.LogLevel) ToolsAgentOption {
	return func(a *ToolsAgent) {
		a.logger = logger.NewConsoleLoggerWithPrefix(level, a.Name)
	}
}
