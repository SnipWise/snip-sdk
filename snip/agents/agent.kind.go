package agents

// AgentKind represents the type of agent
type AgentKind string

const (
	Basic  AgentKind = "Basic"
	Chat  AgentKind = "Chat"
	ChatServer AgentKind = "ChatServer"
	Remote AgentKind = "Remote"
	Tool  AgentKind = "Tools"
	Intent AgentKind = "Intent"
	Rag    AgentKind = "Rag"
	Compressor AgentKind = "Compressor"
	Structured AgentKind = "Structured"
	Macro AgentKind = "Macro"
)

