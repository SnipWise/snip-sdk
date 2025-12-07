package tools

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	oai "github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"

	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/models"
	openaihelpers "github.com/snipwise/snip-sdk/snip/openai-helpers"
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

// ToolsAgent is an implementation of an AI agent that can utilize various tools to perform tasks.
// IMPORTANT: ToolsAgent is not a chat agent.
type ToolsAgent struct {
	ctx                context.Context
	Name               string
	SystemInstructions string
	ModelID            string
	Messages           []*ai.Message

	Config models.ModelConfig

	ToolsIndex      []ai.ToolRef
	toolCallingFlow *core.Flow[*ToolCallsRequest, ToolCallsResult, struct{}]

	genKitInstance *genkit.Genkit

	// flow(s) for the agent

	logger logger.Logger
}

func NewToolsAgent(
	ctx context.Context,
	toolsAgentConfig agents.AgentConfig,
	modelConfig models.ModelConfig,
	//toolsIndex []ai.ToolRef,
	opts ...ToolsAgentOption,
) (*ToolsAgent, error) {

	oaiPlugin := &oai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(toolsAgentConfig.EngineURL),
		},
	}
	genKitInstance := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	// Check if model is available
	if !openaihelpers.IsModelAvailable(ctx, toolsAgentConfig.EngineURL, toolsAgentConfig.ModelID) {
		return nil, fmt.Errorf("model %s is not available at %s", toolsAgentConfig.ModelID, toolsAgentConfig.EngineURL)
	}

	toolsAgent := &ToolsAgent{
		Name:               toolsAgentConfig.Name,
		SystemInstructions: toolsAgentConfig.SystemInstructions,
		ModelID:            toolsAgentConfig.ModelID,

		ToolsIndex: []ai.ToolRef{},

		Messages: []*ai.Message{},

		Config: modelConfig,

		ctx:            ctx,
		genKitInstance: genKitInstance,

		logger: logger.GetLoggerFromEnvWithPrefix(toolsAgentConfig.Name), // Default logger from env
	}

	// Apply options
	for _, opt := range opts {
		opt(toolsAgent)
	}

	// Log model availability
	toolsAgent.logger.Info("✅ Model %s is available at %s", toolsAgentConfig.ModelID, toolsAgentConfig.EngineURL)

	return toolsAgent, nil
}


func (toolsAgent *ToolsAgent) SetTools(tools []ai.ToolRef) { 
	toolsAgent.ToolsIndex = tools
}

func (toolsAgent *ToolsAgent) GetTools() []ai.ToolRef {
	return toolsAgent.ToolsIndex
}


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

// GetName returns the name of the ToolsAgent.
func (toolsAgent *ToolsAgent) GetName() string {
	return toolsAgent.Name
}

// GetMessages returns the message history of the ToolsAgent.
func (toolsAgent *ToolsAgent) GetMessages() []*ai.Message {
	return toolsAgent.Messages
}

// GetCurrentContextSize returns the current context size in characters.
func (toolsAgent *ToolsAgent) GetCurrentContextSize() int {
	totalContextSize := len(toolsAgent.SystemInstructions)
	for _, msg := range toolsAgent.Messages {
		for _, content := range msg.Content {
			totalContextSize += len(content.Text)
		}
	}
	return totalContextSize
}

// Kind returns the kind of the agent.
func (toolsAgent *ToolsAgent) Kind() agents.AgentKind {
	return agents.Tool
}

// GetInfo returns the ToolsAgent information.
func (toolsAgent *ToolsAgent) GetInfo() (agents.ToolsAgentInfo, error) {
	return agents.ToolsAgentInfo{
		Name:    toolsAgent.Name,
		Config:  toolsAgent.Config,
		ModelID: toolsAgent.ModelID,
	}, nil
}
