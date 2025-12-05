package smart

// Structure for agent information endpoint
type RagAgentInfo struct {
	Name               string `json:"name"`
	ModelID            string `json:"model_id"`
	EmbeddingDimension int    `json:"embedding_dimension"`
	StoreName		 string `json:"store_name"`
	StorePath		 string `json:"store_path"`
	NumberOfDocuments  int    `json:"number_of_documents"`

	//Config  ModelConfig `json:"config"`
}
