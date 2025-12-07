package tools

import (
	"encoding/json"
	"fmt"

	"github.com/firebase/genkit/go/ai"
)

func castToToolOutput(output any) (ToolOutput, error) {
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		//log.Printf("Failed to marshal tool output: %v\n", err)
		return ToolOutput{
			Content: []ContentItem{{
				Text: err.Error(),
				Type: "text",
			}},
		}, err
	}

	var result ToolOutput
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		//log.Printf("Failed to unmarshal tool output: %v\n", err)
		return ToolOutput{
			Content: []ContentItem{{
				Text: err.Error(),
				Type: "text",
			}},
		}, err
	}
	return result, nil
}

func displayToolCallResult(output any) {
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		fmt.Println("🤖 Tool output:", output)
		return
	}
	var result ToolOutput
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		fmt.Println("🤖 Tool output:", output)
		return
	}
	if len(result.Content) > 0 {
		fmt.Println("🤖 Tool output:", result.Content[0].Text)
	} else {
		fmt.Println("🤖 Tool output:", output)
	}
}

func displayToolRequets(toolRequest *ai.ToolRequest) {
	jsonInput, err := json.Marshal(toolRequest.Input)
	if err != nil {
		fmt.Println("🛠️ Tool request:", toolRequest.Name, toolRequest.Ref, toolRequest.Input)
	}
	fmt.Println("🛠️ Tool request:", toolRequest.Name, "args:", string(jsonInput))
}
