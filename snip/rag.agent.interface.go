package snip

import "github.com/snipwise/snip-sdk/snip/agents"

// Structure for agent information endpoint
type RagAgentInfo struct {
	Name               string `json:"name"`
	ModelID            string `json:"model_id"`
	EmbeddingDimension int    `json:"embedding_dimension"`
	StoreName          string `json:"store_name"`
	StorePath          string `json:"store_path"`
	NumberOfDocuments  int    `json:"number_of_documents"`
}

type TextChunk struct {
	Content  string
	Metadata map[string]any
}

type AIRagAgent interface {
	GetName() string
	GetInfo() (RagAgentInfo, error)
	Kind() agents.AgentKind
	AddTextChunksToStore(chunks []TextChunk) (int, error)
	SearchSimilarities(query string) ([]string, error)
}