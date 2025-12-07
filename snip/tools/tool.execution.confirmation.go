package tools

type ConfirmationResponse int

const (
	Confirmed ConfirmationResponse = iota
	Denied
	Quit
)
/*
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


*/


type ToolExecutionConfirmation struct {
	Question func(toolName string, toolInput any, toolCallRef string) ConfirmationResponse

	OnConfirmed func(toolName string, toolInput any, toolCallRef string, output any, err error)
	OnDenied    func(toolName string, toolInput any, toolCallRef string)
	OnQuit      func(toolName string, toolInput any, toolCallRef string)
	Default     func()
}

func WithConfirmation(toolExecutionConfirmation ToolExecutionConfirmation) ToolsAgentOption {
	return func(toolsAgent *ToolsAgent) {
		toolsAgent.toolExecutionConfirmation = &toolExecutionConfirmation
	}
}

type ToolExecution struct {
	OnExecuted func(toolName string, toolInput any, toolCallRef string, output any, err error)
}

func WithToolExecution(toolExecution ToolExecution) ToolsAgentOption {
	return func(toolsAgent *ToolsAgent) {
		toolsAgent.toolExecution = &toolExecution
	}
}