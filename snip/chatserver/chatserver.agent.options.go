package chatserver

import "github.com/snipwise/snip-sdk/snip/toolbox/logger"

// ChatAgentServerOption defines a functional option for configuring the agent
type ChatAgentServerOption func(*ChatAgentServer)

// WithLogger sets a custom logger for the agent
func WithLogger(log logger.Logger) ChatAgentServerOption {
	return func(chatAgentServer *ChatAgentServer) {
		chatAgentServer.logger = log
	}
}

// WithVerbose enables verbose logging (INFO level) with agent name prefix
func WithVerbose(verbose bool) ChatAgentServerOption {
	return func(chatAgentServer *ChatAgentServer) {
		if verbose {
			chatAgentServer.logger = logger.NewConsoleLoggerWithPrefix(logger.LevelInfo, chatAgentServer.agent.Name)
		} else {
			chatAgentServer.logger = &logger.NoOpLogger{}
		}
	}
}

// WithLogLevel sets the log level for the agent
func WithLogLevel(level logger.LogLevel) ChatAgentServerOption {
	return func(chatAgentServer *ChatAgentServer) {
		chatAgentServer.logger = logger.NewConsoleLoggerWithPrefix(level, chatAgentServer.agent.Name)
	}
}
