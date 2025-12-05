package snip

import (
	"github.com/firebase/genkit/go/genkit"
)

// EnableServer configures the agent to expose its flows via HTTP endpoints
func EnableServer(config ConfigHTTP) AgentOption {
	return func(agent *Agent) {
		// Set up HTTP handlers for the flows
		if agent.chatFlow != nil && config.ChatFlowHandler == nil {
			config.ChatFlowHandler = genkit.Handler(agent.chatFlow)
		}

		if agent.chatStreamFlow != nil && config.ChatStreamFlowHandler == nil {
			config.ChatStreamFlowHandler = genkit.Handler(agent.chatStreamFlow)
		}

		agent.serverConfig = &config
	}
}
