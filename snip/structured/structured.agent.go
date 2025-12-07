package structured

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/openai/openai-go/option"
	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/models"
	openaihelpers "github.com/snipwise/snip-sdk/snip/openai-helpers"
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"

	oai "github.com/firebase/genkit/go/plugins/compat_oai/openai"
)

type StructuredAgent[O any] struct {
	ctx                context.Context
	Name               string
	SystemInstructions string
	ModelID            string

	Messages []*ai.Message

	Config models.ModelConfig

	genKitInstance *genkit.Genkit

	logger logger.Logger

	structuredFlow *core.Flow[*agents.ChatRequest, *O, struct{}]
}


// NewStructuredAgent creates a new StructuredAgent that generates structured data of type O.
func NewStructuredAgent[O any](
	ctx context.Context,
	structuredAgentConfig agents.AgentConfig,
	modelConfig models.ModelConfig,
	opts ...StructuredAgentOption[O],
) (*StructuredAgent[O], error) {

	oaiPlugin := &oai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(structuredAgentConfig.EngineURL),
		},
	}

	genKitInstance := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	// Check if model is available
	if !openaihelpers.IsModelAvailable(ctx, structuredAgentConfig.EngineURL, structuredAgentConfig.ModelID) {
		return nil, fmt.Errorf("model %s is not available at %s", structuredAgentConfig.ModelID, structuredAgentConfig.EngineURL)
	}

	structuredAgent := &StructuredAgent[O]{
		Name:               structuredAgentConfig.Name,
		SystemInstructions: structuredAgentConfig.SystemInstructions,
		ModelID:            structuredAgentConfig.ModelID,
		Messages:           []*ai.Message{},
		Config:             modelConfig,

		ctx:            ctx,
		genKitInstance: genKitInstance,

		logger: logger.GetLoggerFromEnvWithPrefix(structuredAgentConfig.Name), // Default logger from env

	}
	// Apply all options (can override logger)
	for _, opt := range opts {
		opt(structuredAgent)
	}

	// Log model availability
	structuredAgent.logger.Info("✅ Model %s is available at %s", structuredAgentConfig.ModelID, structuredAgentConfig.EngineURL)

	structuredFlow := genkit.DefineFlow(genKitInstance, structuredAgent.Name+"-structured-flow",
		func(ctx context.Context, input *agents.ChatRequest) (*O, error) {

			structuredOutput, modelResponse, err := genkit.GenerateData[O](ctx, genKitInstance,
				ai.WithModelName("openai/"+structuredAgent.ModelID),
				ai.WithSystem(structuredAgent.SystemInstructions),
				ai.WithPrompt(input.UserMessage),
				ai.WithConfig(structuredAgent.Config.ToOpenAIParams()),
			)
			if err != nil {
				return nil, err
			}
			// export SNIP_LOG_LEVEL=debug to see model response
			structuredAgent.logger.Debug("📝 model response")
			structuredAgent.logger.Debug(modelResponse.Text())
			return structuredOutput, nil

		})
	structuredAgent.structuredFlow = structuredFlow

	return structuredAgent, nil
}


// GenerateStructuredData generates structured data of type O based on the input text.
func (structuredAgent *StructuredAgent[O]) GenerateStructuredData(text string) (*O, error) {
	result, err := structuredAgent.structuredFlow.Run(structuredAgent.ctx, &agents.ChatRequest{
		UserMessage: text,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Kind returns the kind of the agent.
func (agent *StructuredAgent[O]) Kind() agents.AgentKind {
	return agents.Structured
}