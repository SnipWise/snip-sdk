package snip

import (
	"github.com/firebase/genkit/go/ai"
	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/tools"
)

type AIToolsAgent interface {
	GetName() string
	GetInfo() (agents.ToolsAgentInfo, error)
	Kind() agents.AgentKind
	// QUESTION: should this return ai.ToolRef or tools.ToolDefinition?
	SetTools(tools []ai.ToolRef)
	GetTools() []ai.ToolRef
	RunToolCalls(prompt string) (tools.ToolCallsResult, error)
}
