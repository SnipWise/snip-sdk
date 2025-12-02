package smart

import (
	"github.com/firebase/genkit/go/ai"
)

type AIAgent interface {
	Ask(question string) (string, error)
	AskStream(question string, callback func(string) error) (string, error)
	GetName() string
	GetMessages() []*ai.Message
	GetInfo() (AgentInfo, error)
	Kind() AgentKind
	AddContextToMessages(context string) error
}
// TODO: add helpers to handle the messages

// TODO:
type IntentTraits interface {
	// AIAgent
	// DetermineIntent(question string) (string, error)
}

// TODO:
type RagTraits interface {
	// AIAgent
	// RetrieveRelevantDocuments(question string) ([]Document, error)
}