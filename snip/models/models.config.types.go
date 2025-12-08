package models

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

// ModelConfig represents the configuration for chat completion
type ModelConfig struct {
	// Temperature controls randomness in responses (0.0-2.0).
	// Higher values (e.g., 1.0) make output more random, lower values (e.g., 0.2) make it more focused and deterministic.
	Temperature float64

	// TopP controls nucleus sampling (0.0-1.0).
	// An alternative to temperature, it considers tokens with top_p probability mass.
	// For example, 0.1 means only tokens comprising the top 10% probability mass are considered.
	TopP float64

	// MaxTokens sets the maximum number of tokens to generate in the completion.
	// The total length of input tokens and generated tokens is limited by the model's context length.
	MaxTokens int64

	// FrequencyPenalty reduces repetition of token sequences (-2.0 to 2.0).
	// Positive values penalize tokens that have already appeared, decreasing likelihood of verbatim repetition.
	FrequencyPenalty float64

	// PresencePenalty encourages talking about new topics (-2.0 to 2.0).
	// Positive values penalize tokens that have appeared at all, increasing likelihood of new topics.
	PresencePenalty float64

	// Stop defines sequences where the API will stop generating further tokens.
	// The returned text will not contain the stop sequence.
	Stop []string

	// Seed enables deterministic sampling when set.
	// If specified, the system will make a best effort to sample deterministically for repeated requests with the same seed.
	Seed *int64

	// ReasoningEffort controls the reasoning effort for reasoning models.
	// Valid values: "low", "medium", "high". Only applicable to reasoning models like o1.
	ReasoningEffort string

	// ParallelToolCalls enables parallel function calling during tool use.
	// When true, the model can call multiple tools simultaneously.
	// When false, tools are called sequentially.
	ParallelToolCalls *bool

	// LogitBias modifies the likelihood of specified tokens appearing in the completion.
	// Maps token IDs to bias values from -100 to 100. Values of -100 ban the token, while 100 strongly increases likelihood.
	//LogitBias map[string]int64

	// N specifies how many chat completion choices to generate for each input message.
	// Note: Because this parameter generates many completions, it can quickly consume your token quota.
	//N int64

	// User provides a unique identifier representing your end-user.
	// This helps OpenAI monitor and detect abuse.
	//User string
}

// ToOpenAIParams converts ModelConfig to OpenAI ChatCompletionNewParams
func (c ModelConfig) ToOpenAIParams() *openai.ChatCompletionNewParams {
	params := &openai.ChatCompletionNewParams{}
	if c.Temperature != 0 {
		params.Temperature = openai.Float(c.Temperature)
	}
	if c.TopP != 0 {
		params.TopP = openai.Float(c.TopP)
	}
	if c.MaxTokens != 0 {
		params.MaxTokens = openai.Int(c.MaxTokens)
	}
	if c.FrequencyPenalty != 0 {
		params.FrequencyPenalty = openai.Float(c.FrequencyPenalty)
	}
	if c.PresencePenalty != 0 {
		params.PresencePenalty = openai.Float(c.PresencePenalty)
	}
	if len(c.Stop) > 0 {
		// Note: Stop parameter handling depends on OpenAI SDK version
		// For now, we'll skip this as it requires union type handling
		// Users can set Stop directly on the params if needed
	}
	if c.Seed != nil {
		params.Seed = openai.Int(*c.Seed)
	}
	if c.ReasoningEffort != "" {
		params.ReasoningEffort = shared.ReasoningEffort(c.ReasoningEffort)
	}
	if c.ParallelToolCalls != nil {
		params.ParallelToolCalls = openai.Bool(*c.ParallelToolCalls)
	}
	// if len(c.LogitBias) > 0 {
	// 	params.LogitBias = c.LogitBias
	// }
	// if c.N != 0 {
	// 	params.N = openai.Int(c.N)
	// }
	// if c.User != "" {
	// 	params.User = openai.String(c.User)
	// }
	return params
}