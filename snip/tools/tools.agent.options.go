package tools

import (
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

// AgentOption defines a functional option for configuring LocalAIAgent
type ToolsAgentOption func(*ToolsAgent)

// WithLogger sets a custom logger for the agent
func WithLogger(log logger.Logger) ToolsAgentOption {
	return func(toolsAgent *ToolsAgent) {
		toolsAgent.logger = log
	}
}

// WithVerbose enables verbose logging (INFO level) with agent name prefix
func WithVerbose(verbose bool) ToolsAgentOption {
	return func(toolsAgent *ToolsAgent) {
		if verbose {
			toolsAgent.logger = logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, toolsAgent.Name)
		} else {
			toolsAgent.logger = &logger.NoOpLogger{}
		}
	}
}

// WithLogLevel sets the log level for the agent
func WithLogLevel(level logger.LogLevel) ToolsAgentOption {
	return func(toolsAgent *ToolsAgent) {
		toolsAgent.logger = logger.NewConsoleLoggerWithPrefix(level, toolsAgent.Name)
	}
}
