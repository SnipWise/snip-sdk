package agents

import "github.com/snipwise/snip-sdk/snip/models"

type ToolsAgentInfo struct {
	Name    string             `json:"name"`
	ModelID string             `json:"model_id"`
	Config  models.ModelConfig `json:"config"`

	// TODO: to be completed
}
