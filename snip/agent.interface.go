package snip

import "github.com/snipwise/snip-sdk/snip/agents"

import (
	"github.com/firebase/genkit/go/ai"
)

type AIAgent interface {
	// Methods with memory management (conversation history is maintained)
	AskWithMemory(question string) (agents.ChatResponse, error)
	AskStreamWithMemory(question string, callback func(agents.ChatResponse) error) (agents.ChatResponse, error)

	// Methods without memory management (stateless, each request is independent)
	Ask(question string) (agents.ChatResponse, error)
	AskStream(question string, callback func(agents.ChatResponse) error) (agents.ChatResponse, error)

	GetName() string
	GetMessages() []*ai.Message

	ReplaceMessagesWith(messages []*ai.Message) error

	ReplaceMessagesWithSystemMessages(systemMessages []string) error

	GetCurrentContextSize() int

	GetInfo() (agents.AgentInfo, error)
	Kind() agents.AgentKind
	AddSystemMessage(context string) error

	// Context compression methods (require EnableContextCompression option)
	CompressContext() (agents.ChatResponse, error)
	CompressContextStream(callback func(agents.ChatResponse) error) (agents.ChatResponse, error)
}

// TODO: add helpers to handle the messages

// TODO:
// type IntentTraits interface {
// 	// AIAgent
// 	// DetermineIntent(question string) (string, error)
// }

// TODO:
// type RagTraits interface {
// 	// AIAgent
// 	// RetrieveRelevantDocuments(question string) ([]Document, error)
// }
