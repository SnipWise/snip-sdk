package chatserver

import (
	"github.com/firebase/genkit/go/genkit"
)

// EnableServer configures the agent to expose its flows via HTTP endpoints
func EnableServer(config ConfigHTTP) ChatAgentServerOption {
	return func(cas *ChatAgentServer) {
		// Set up HTTP handlers for the flows
		if cas.agent.GetChatFlowWithMemory() != nil && config.ChatFlowHandler == nil {
			config.ChatFlowHandler = genkit.Handler(cas.agent.GetChatFlowWithMemory())
		}

		if cas.agent.GetChatStreamFlowWithMemory() != nil && config.ChatStreamFlowHandler == nil {
			config.ChatStreamFlowHandler = genkit.Handler(cas.agent.GetChatStreamFlowWithMemory())
		}

		cas.serverConfig = &config
	}
}
