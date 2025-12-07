package tools

type ToolCallsRequest struct {
	Prompt string `json:"prompt"`
}
type ToolCallsResult struct {
	Text string              `json:"text"`
	List []map[string]string `json:"list"`
}

type ContentItem struct {
	Text string `json:"text"`
	Type string `json:"type"`
}
type ToolOutput struct {
	Content []ContentItem `json:"content"`
}

// ToolDefinition defines a tool that can be used by the ToolsAgent.
// It includes the tool's name, description, and the function that implements the tool's behavior.
// The `name` is the identifier the model uses to request the tool.
// The `description` helps the model understand when to use the tool.
// The function `fn` implements the tool's logic,
// taking  an input of type `In`, and returning an output of type `Out`.
// The input and output types determine the `inputSchema` and `outputSchema`
// in the tool's definition,
// which guide the model on how to provide input and interpret output.
type ToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Function    any    `json:"-"`
}
