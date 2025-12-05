package snip

// Structure for agent information endpoint
type AgentInfo struct {
	Name    string      `json:"name"`
	ModelID string      `json:"model_id"`
	Config  ModelConfig `json:"config"`
}
