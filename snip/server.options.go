package snip

import (
	"github.com/firebase/genkit/go/genkit"
)

// EnableServer configures the agent to expose its flows via HTTP endpoints
func EnableServer(config ConfigHTTP) AgentOption {
	return func(agent *Agent) {
		// Set up HTTP handlers for the flows
		if agent.chatFlowWithMemory != nil && config.ChatFlowHandler == nil {
			config.ChatFlowHandler = genkit.Handler(agent.chatFlowWithMemory)
		}

		if agent.chatStreamFlowWithMemory != nil && config.ChatStreamFlowHandler == nil {
			config.ChatStreamFlowHandler = genkit.Handler(agent.chatStreamFlowWithMemory)
		}

		agent.serverConfig = &config
	}
}
