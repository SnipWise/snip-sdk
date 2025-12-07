package tools

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func EnableToolCallFlow() ToolsAgentOption {
	return func(toolsAgent *ToolsAgent) {
		initializeToolCallFlow(toolsAgent)
	}
}

// TODO: EnableToolCallFlow(Confirmation) -> QUESTION: How to pass the confirmation function?

func initializeToolCallFlow(toolsAgent *ToolsAgent) {

	// STEP 1: Define tool-calling flow

	toolCallingFlow := genkit.DefineFlow(toolsAgent.genKitInstance, toolsAgent.Name+"-tool-calling-flow",
		func(ctx context.Context, req *ToolCallsRequest) (ToolCallsResult, error) {

			// STEP 2: Initialize loop control variables
			stopped := false           // Controls the conversation loop
			lastAssistantMessage := "" // Final AI message

			//totalOfToolsCalls := 0
			toolCallsResults := []map[string]any{}

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

					if tool == nil {
						toolsAgent.logger.Error("❌ Tool %q not found", req.Name)
						continue
					}

					// Execute tool with user confirmation
					runToolExecutionWithConfirmation := func() {
						for {

							response := toolsAgent.toolExecutionConfirmation.Question(req.Name, req.Input, req.Ref)

							switch response {
							case Quit:
								if toolsAgent.toolExecutionConfirmation != nil && toolsAgent.toolExecutionConfirmation.OnQuit != nil {
									toolsAgent.toolExecutionConfirmation.OnQuit(req.Name, req.Input, req.Ref)
								} else {
									// No OnQuit handler defined
									toolsAgent.logger.Warn("✋ No OnQuit handler defined")
									toolsAgent.logger.Debug("🛑 Tool execution aborted by user.")
								}
								stopped = true
								return
							case Confirmed:
								output, err := tool.RunRaw(ctx, req.Input)

								if err != nil {
									toolsAgent.logger.Error("❌ Tool %q execution failed: %v", tool.Name(), err)
									// Exit the program on tool execution error
									stopped = true
									return
								}

								part := ai.NewToolResponsePart(&ai.ToolResponse{
									Name:   req.Name,
									Ref:    req.Ref,
									Output: output,
								})

								history = append(history, ai.NewMessage(ai.RoleTool, nil, part))

								// Store the raw output (not converted to string) so it can be transformed later
								toolCallsResults = append(toolCallsResults, map[string]any{
									tool.Name(): output,
								})

								if toolsAgent.toolExecutionConfirmation != nil && toolsAgent.toolExecutionConfirmation.OnConfirmed != nil {
									toolsAgent.toolExecutionConfirmation.OnConfirmed(req.Name, req.Input, req.Ref, output, err)
								} else {
									// No OnConfirmed handler defined
									toolsAgent.logger.Warn("✋ No OnConfirmed handler defined")
									toolsAgent.logger.Debug("✅ Tool executed successfully: %s", tool.Name())
								}
								
								return
							case Denied:
								// Skip tool execution
								// Add tool response indicating the tool was not executed
								part := ai.NewToolResponsePart(&ai.ToolResponse{
									Name:   req.Name,
									Ref:    req.Ref,
									Output: map[string]any{"error": "Tool execution cancelled by user"},
								})
								history = append(history, ai.NewMessage(ai.RoleTool, nil, part))

								if toolsAgent.toolExecutionConfirmation != nil && toolsAgent.toolExecutionConfirmation.OnDenied != nil {
									toolsAgent.toolExecutionConfirmation.OnDenied(req.Name, req.Input, req.Ref)
								} else {
									// No OnDenied handler defined
									toolsAgent.logger.Warn("✋ No OnDenied handler defined")
									toolsAgent.logger.Debug("⏩ Tool execution denied by user.")
								}
								return
							default:
								// Handle default case
								if toolsAgent.toolExecutionConfirmation != nil && toolsAgent.toolExecutionConfirmation.Default != nil {
									toolsAgent.toolExecutionConfirmation.Default()
								} else {
									// No Default handler defined
									toolsAgent.logger.Warn("✋ No Default handler defined")
									toolsAgent.logger.Debug("No valid response received. Skipping tool execution.")
								}

								//continue
								return
							}

						}

					}
					runToolExecutionWithConfirmation()

				}

			} // END: of loop
			return ToolCallsResult{
				Text: lastAssistantMessage,
				List: toolCallsResults,
			}, nil
		})
	toolsAgent.toolCallingFlow = toolCallingFlow
}
