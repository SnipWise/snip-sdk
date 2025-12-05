package snip

import "github.com/firebase/genkit/go/ai"

// CompressorAgentInterface defines the contract for compression agents
type AICompressorAgent interface {
	// GetName returns the name of the agent
	GetName() string

	// GetKind returns the kind/type of the agent
	GetKind() AgentKind

	// GetInfo returns information about the agent
	GetInfo() (AgentInfo, error)

	// GetCompressionPrompt returns the current compression prompt
	GetCompressionPrompt() string

	// SetCompressionPrompt sets a new compression prompt
	SetCompressionPrompt(prompt string)

	// CompressText compresses the given text using the compression prompt
	CompressText(text string) (ChatResponse, error)

	// CompressTextStream compresses the given text using streaming
	CompressTextStream(text string, callback func(ChatResponse) error) (ChatResponse, error)

	// CompressMessages compresses a list of messages into a summary
	CompressMessages(messages []*ai.Message) (ChatResponse, error)

	// CompressMessagesStream compresses a list of messages into a summary using streaming
	CompressMessagesStream(messages []*ai.Message, callback func(ChatResponse) error) (ChatResponse, error)
}
