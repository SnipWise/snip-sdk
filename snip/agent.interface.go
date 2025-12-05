package snip

import (
	"github.com/firebase/genkit/go/ai"
)

type AIAgent interface {
	Ask(question string) (ChatResponse, error)
	AskStream(question string, callback func(ChatResponse) error) (ChatResponse, error)
	GetName() string
	GetMessages() []*ai.Message
	ReplaceMessagesWith(messages []*ai.Message) error

	GetCurrentContextSize() int

	GetInfo() (AgentInfo, error)
	Kind() AgentKind
	AddSystemMessage(context string) error
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