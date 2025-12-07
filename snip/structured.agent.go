package snip

import (
	"github.com/snipwise/snip-sdk/snip/agents"
)

// StructuredAgentInterface defines the interface for agents that generate structured data
type StructuredAgentInterface[O any] interface {
	GenerateStructuredData(text string) (*O, error)
	Kind() agents.AgentKind
}
