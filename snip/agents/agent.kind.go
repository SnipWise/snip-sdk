package agents

// AgentKind represents the type of agent
type AgentKind string

const (
	Basic  AgentKind = "Basic"
	Chat  AgentKind = "Chat"
	Remote AgentKind = "Remote"
	Tool  AgentKind = "Tools"
	Intent AgentKind = "Intent"
	Rag    AgentKind = "Rag"
	Compressor AgentKind = "Compressor"
)

