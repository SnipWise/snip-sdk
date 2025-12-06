package snip

import "github.com/snipwise/snip-sdk/snip/agents"

import "github.com/firebase/genkit/go/ai"

// AICompressorAgent defines the contract for compression agents
type AICompressorAgent interface {
	// GetName returns the name of the agent
	GetName() string

	// GetKind returns the kind/type of the agent
	GetKind() agents.AgentKind

	// GetInfo returns information about the agent
	GetInfo() (agents.AgentInfo, error)

	// GetCompressionPrompt returns the current compression prompt
	GetCompressionPrompt() string

	// SetCompressionPrompt sets a new compression prompt
	SetCompressionPrompt(prompt string)

	// CompressText compresses the given text using the compression prompt
	CompressText(text string) (agents.ChatResponse, error)

	// CompressTextStream compresses the given text using streaming
	CompressTextStream(text string, callback func(agents.ChatResponse) error) (agents.ChatResponse, error)

	// CompressMessages compresses a list of messages into a summary
	CompressMessages(messages []*ai.Message) (agents.ChatResponse, error)

	// CompressMessagesStream compresses a list of messages into a summary using streaming
	CompressMessagesStream(messages []*ai.Message, callback func(agents.ChatResponse) error) (agents.ChatResponse, error)
}
