package chat

import "github.com/snipwise/snip-sdk/snip"

// EnableContextCompression configures an agent to use a compressor agent for context compression
// The compressor agent is used to compress conversation history when needed
// Example:
//   agent, err := chatagent.NewChatAgent(ctx, config, modelConfig,
//     chatagent.EnableContextCompression(myCompressorAgent),
//   )
func EnableContextCompression(compressorAgent snip.AICompressorAgent) AgentOption {
	return func(agent *ChatAgent) {
		agent.compressorAgent = compressorAgent
	}
}
