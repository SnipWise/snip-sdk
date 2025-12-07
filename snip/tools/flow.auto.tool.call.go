package tools

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

func EnableAutoToolCallFlow() ToolsAgentOption {
	return func(toolsAgent *ToolsAgent) {
		initializeAutoToolCallFlow(toolsAgent)
	}
}

func initializeAutoToolCallFlow(toolsAgent *ToolsAgent) {
	// STEP 1: Define tool-calling flow

	autoToolCallingFlow := genkit.DefineFlow(toolsAgent.genKitInstance, toolsAgent.Name+"-tool-calling-flow",
		func(ctx context.Context, req *ToolCallsRequest) (ToolCallsResult, error) {

			// STEP 2: Initialize loop control variables
			stopped := false           // Controls the conversation loop
			lastAssistantMessage := "" // Final AI message

			//totalOfToolsCalls := 0
			toolCallsResults := []map[string]string{}

			history := []*ai.Message{}
			// STEP 3: Start the conversation loop
			// To avoid repeating the first user message in the history
			// we add it here before entering the loop and using prompt
			history = append(history, ai.NewUserTextMessage(req.Prompt))
			// TODO: use agent.Messages as initial history?

			for !stopped { // BEGIN: of loop

				resp, err := genkit.Generate(ctx, toolsAgent.genKitInstance,
					ai.WithModelName("openai/"+toolsAgent.ModelID),
					ai.WithSystem(toolsAgent.SystemInstructions),
					// WithMessages sets the messages. These messages will be sandwiched between the system and user prompts.
					// ai.WithMessages(
					// 	agent.Messages...,
					// ),
					ai.WithMessages(
						history...,
					),
					//ai.WithPrompt(req.Prompt), // NOTE: do not add the prompt again
					ai.WithTools(
						toolsAgent.ToolsIndex...,
					),
					ai.WithToolChoice(ai.ToolChoiceAuto),
					ai.WithReturnToolRequests(true),
				)
				if err != nil {
					return ToolCallsResult{}, err
				}

				// We do not use parallel tool calls
				toolRequests := resp.ToolRequests()
				if len(toolRequests) == 0 {
					// No tool requests, we are done
					stopped = true // Exit the loop
					lastAssistantMessage = resp.Text()
					break // Exit the loop now
				}
				// IMPORTANT: Add the assistant's message with tool requests to history
				// This ensures the model knows it already proposed these tools
				// history = append(history, resp.Message)
				history = append(history, resp.Message)

				for _, req := range toolRequests {

					var tool ai.Tool
					// tool = genkit.LookupTool(agent.genKitInstance, req.Name)

					for _, t := range toolsAgent.ToolsIndex {
						if t.Name() == req.Name {
							// Try to convert ToolRef to Tool
							if toolImpl, ok := t.(ai.Tool); ok {
								tool = toolImpl
								// ✅ Successfully converted to ai.Tool"
								break
							}
							// else: ❌ Failed to convert ToolRef to ai.Tool")
						}
					}

					if toolsAgent.logger.GetLevel() == logger.LevelDebug {
						displayToolRequets(req)
					}

					if tool == nil {
						toolsAgent.logger.Error("❌ Tool %q not found", req.Name)
						continue
					}

					// Execute tool without user confirmation
					runToolExecution := func() {
						output, err := tool.RunRaw(ctx, req.Input)
						if err != nil {
							toolsAgent.logger.Error("❌ Tool %q execution failed: %v", tool.Name(), err)
							// Exit the program on tool execution error
							stopped = true
							return
						}

						// IMPORTANT: or put an argument to execConfirmation() to disable output display
						if toolsAgent.logger.GetLevel() == logger.LevelDebug {
							displayToolCallResult(output)
						}

						part := ai.NewToolResponsePart(&ai.ToolResponse{
							Name:   req.Name,
							Ref:    req.Ref,
							Output: output,
						})

						history = append(history, ai.NewMessage(ai.RoleTool, nil, part))

						toolCallOutput, errOutput := castToToolOutput(output)
						if errOutput != nil {
							toolsAgent.logger.Warn("⚠️ Warning: failed to cast tool output: %v", errOutput)
						}

						if len(toolCallOutput.Content) > 0 {
							toolCallsResults = append(toolCallsResults, map[string]string{
								tool.Name(): toolCallOutput.Content[0].Text,
							})
						} else {
							// Fallback: use the raw output as a string
							outputStr := fmt.Sprintf("%v", output)
							toolCallsResults = append(toolCallsResults, map[string]string{
								tool.Name(): outputStr,
							})
						}

					}
					runToolExecution()

				}

			} // END: of loop
			return ToolCallsResult{
				Text: lastAssistantMessage,
				List: toolCallsResults,
			}, nil
		})

	toolsAgent.toolCallingFlow = autoToolCallingFlow
}
