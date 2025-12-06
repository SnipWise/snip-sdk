package agents

import "github.com/snipwise/snip-sdk/snip/models"

// Structure for agent information endpoint
type AgentInfo struct {
	Name    string             `json:"name"`
	ModelID string             `json:"model_id"`
	Config  models.ModelConfig `json:"config"`
}
