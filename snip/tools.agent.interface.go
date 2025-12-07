package snip

import (
	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/tools"
)

type AIToolsAgent interface {
	GetName() string
	GetInfo() (agents.ToolsAgentInfo, error)
	Kind() agents.AgentKind
	RunToolCalls(prompt string) (tools.ToolCallsResult, error)
}
