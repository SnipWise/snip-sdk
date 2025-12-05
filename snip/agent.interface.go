package snip

import (
	"github.com/firebase/genkit/go/ai"
)

type AIAgent interface {
	// Methods with memory management (conversation history is maintained)
	AskWithMemory(question string) (ChatResponse, error)
	AskStreamWithMemory(question string, callback func(ChatResponse) error) (ChatResponse, error)

	// Methods without memory management (stateless, each request is independent)
	Ask(question string) (ChatResponse, error)
	AskStream(question string, callback func(ChatResponse) error) (ChatResponse, error)

	GetName() string
	GetMessages() []*ai.Message

	ReplaceMessagesWith(messages []*ai.Message) error

	ReplaceMessagesWithSystemMessages(systemMessages []string) error

	GetCurrentContextSize() int

	GetInfo() (AgentInfo, error)
	Kind() AgentKind
	AddSystemMessage(context string) error

	// Context compression methods (require EnableContextCompression option)
	CompressContext() (ChatResponse, error)
	CompressContextStream(callback func(ChatResponse) error) (ChatResponse, error)
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
