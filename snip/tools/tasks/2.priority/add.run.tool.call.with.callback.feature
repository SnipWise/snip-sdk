Add a callback trigerred after each tool call in the tool-calling flow.

// RunToolCalls runs the tool-calling flow with the given prompt.
func (toolsAgent *ToolsAgent) RunToolCalls(prompt string) (ToolCallsResult, error) {
	resp, err := toolsAgent.toolCallingFlow.Run(toolsAgent.ctx, &ToolCallsRequest{
		Prompt: prompt,
	})
	if err != nil {
		return ToolCallsResult{}, err
	}
	return resp, nil
}