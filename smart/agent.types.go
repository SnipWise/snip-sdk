package smart

// AgentKind represents the type of agent
type AgentKind string

const (
	Basic  AgentKind = "Basic"
	Remote AgentKind = "Remote"
	Tool  AgentKind = "Tool"
	Intent AgentKind = "Intent"
	Rag    AgentKind = "Rag"
)

// Structure for agent information endpoint
type AgentInfo struct {
	Name    string      `json:"name"`
	ModelID string      `json:"model_id"`
	Config  ModelConfig `json:"config"`
}
