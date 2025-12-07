package snip

import (
	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/text"
)

type AIRagAgent interface {
	GetName() string
	GetInfo() (agents.RagAgentInfo, error)
	Kind() agents.AgentKind
	AddTextChunksToStore(chunks []text.TextChunk) (int, error)
	SearchSimilarities(query string) ([]string, error)
}
