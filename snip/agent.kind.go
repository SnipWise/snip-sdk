package snip

// AgentKind represents the type of agent
type AgentKind string

const (
	Basic  AgentKind = "Basic"
	Remote AgentKind = "Remote"
	Tool  AgentKind = "Tool"
	Intent AgentKind = "Intent"
	Rag    AgentKind = "Rag"
	Compressor AgentKind = "Compressor"
)

