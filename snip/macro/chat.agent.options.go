package macro

import (
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

type MacroAgentOption func(*MacroAgent)

// WithLogger sets a custom logger for the agent
func WithLogger(log logger.Logger) MacroAgentOption {
	return func(macroAgent *MacroAgent) {
		macroAgent.logger = log
	}
}

// WithVerbose enables verbose logging (INFO level) with agent name prefix
func WithVerbose(verbose bool) MacroAgentOption {
	return func(macroAgent *MacroAgent) {
		if verbose {
			macroAgent.logger = logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, macroAgent.Name)
		} else {
			macroAgent.logger = &logger.NoOpLogger{}
		}
	}
}

// WithLogLevel sets the log level for the agent
func WithLogLevel(level logger.LogLevel) MacroAgentOption {
	return func(macroAgent *MacroAgent) {
		macroAgent.logger = logger.NewConsoleLoggerWithPrefix(level, macroAgent.Name)
	}
}
