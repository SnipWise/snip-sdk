package smart

import "fmt"

// AgentConfig represents the core configuration parameters for creating an agent
type RagAgentConfig struct {
	// Name is the identifier for the agent
	Name string

	// ModelID specifies which language model to use
	ModelID string

	// EngineURL is the base URL for the model inference engine
	EngineURL string
}

// Validate checks if the AgentConfig has all required fields
func (ac *RagAgentConfig) Validate() error {
	if ac.Name == "" {
		return fmt.Errorf("agent name is required")
	}
	if ac.ModelID == "" {
		return fmt.Errorf("model ID is required")
	}
	if ac.EngineURL == "" {
		return fmt.Errorf("engine URL is required")
	}
	return nil
}
