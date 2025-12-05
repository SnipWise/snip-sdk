package smart

type AIRagAgent interface {
	GetName() string
	GetInfo() (RagAgentInfo, error)
	Kind() AgentKind
	AddTextChunksToStore(chunks []TextChunk) (int, error)
	SearchSimilarities(query string) ([]string, error)
}