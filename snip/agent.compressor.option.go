package snip

// EnableContextCompression configures an agent to use a compressor agent for context compression
// The compressor agent is used to compress conversation history when needed
// Example:
//   agent, err := snip.NewAgent(ctx, config, modelConfig,
//     snip.EnableContextCompression(myCompressorAgent),
//   )
func EnableContextCompression(compressorAgent AICompressorAgent) AgentOption {
	return func(agent *Agent) {
		agent.compressorAgent = compressorAgent
	}
}
