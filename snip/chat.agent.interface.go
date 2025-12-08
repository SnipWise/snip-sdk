package snip

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/snipwise/snip-sdk/snip/agents"
)

type AIChatAgent interface {
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
	//CompressContext() (agents.ChatResponse, error)
	//CompressContextStream(callback func(agents.ChatResponse) error) (agents.ChatResponse, error)

	GetChatFlowWithMemory() *core.Flow[*agents.ChatRequest, *agents.ChatResponse, struct{}]
	GetChatStreamFlowWithMemory() *core.Flow[*agents.ChatRequest, *agents.ChatResponse, agents.ChatResponse]

	GetStreamCancel() context.CancelFunc
}

