package tools

func EnableYoloToolCallFlow() ToolsAgentOption {
	return func(toolsAgent *ToolsAgent) {
		initializeYoloToolCallFlow(toolsAgent)
	}
}

// TODO: to be implemented
func initializeYoloToolCallFlow(toolsAgent *ToolsAgent) {}
