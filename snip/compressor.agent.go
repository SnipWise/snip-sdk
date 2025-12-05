package snip

import (
	"context"
	"fmt"
	"strings"

	"github.com/firebase/genkit/go/ai"
)

//

type CompressorAgent struct {
	agent             *Agent
	compressionPrompt string
}

func NewCompressorAgent(ctx context.Context, agentConfig AgentConfig, modelConfig ModelConfig) (*CompressorAgent, error) {
	// TODO: create an index of compression prompt
	compressionPrompt := `You are a context compression specialist. Your task is to analyze the conversation history and compress it while preserving all essential information.

	## Instructions:
	1. **Preserve Critical Information**: Keep all important facts, decisions, code snippets, file paths, function names, and technical details
	2. **Remove Redundancy**: Eliminate repetitive discussions, failed attempts, and conversational fluff
	3. **Maintain Chronology**: Keep the logical flow and order of important events
	4. **Summarize Discussions**: Convert long discussions into concise summaries with key takeaways
	5. **Keep Context**: Ensure the compressed version provides enough context for continuing the conversation

	## Output Format:
	Return a compressed version of the conversation that:
	- Uses clear, concise language
	- Groups related topics together
	- Highlights key decisions and outcomes
	- Preserves technical accuracy
	- Maintains references to files, functions, and code

	## Compression Guidelines:
	- Remove: Greetings, acknowledgments, verbose explanations, failed attempts
	- Keep: Facts, code, decisions, file paths, function signatures, error messages, requirements
	- Summarize: Long discussions into bullet points with essential information

	Please compress the following conversation history:
	`

	agent, err := NewAgent(
		ctx,
		agentConfig,
		modelConfig,
		EnableChatFlow(),
		EnableChatStreamFlow(),
	)
	if err != nil {
		return nil, err
	}

	return &CompressorAgent{
		agent:             agent,
		compressionPrompt: compressionPrompt,
	}, nil
}

func (c *CompressorAgent) GetName() string {
	return c.agent.Name
}

func (c *CompressorAgent) GetKind() AgentKind {
	return Compressor
}

func (c *CompressorAgent) GetInfo() (AgentInfo, error) {
	return c.agent.GetInfo()
}

// GetCompressionPrompt returns the current compression prompt
func (c *CompressorAgent) GetCompressionPrompt() string {
	return c.compressionPrompt
}

// SetCompressionPrompt sets a new compression prompt
func (c *CompressorAgent) SetCompressionPrompt(prompt string) {
	c.compressionPrompt = prompt
}

func (c *CompressorAgent) CompressText(text string) (ChatResponse, error) {

	prompt := c.compressionPrompt + "\n\n" + text

	response, err := c.agent.Ask(prompt)
	if err != nil {
		return ChatResponse{}, err
	}

	return response, nil

}

// CompressTextStream compresses the given text using streaming
func (c *CompressorAgent) CompressTextStream(text string, callback func(ChatResponse) error) (ChatResponse, error) {

	prompt := c.compressionPrompt + "\n\n" + text

	response, err := c.agent.AskStream(prompt, callback)
	if err != nil {
		return ChatResponse{}, err
	}

	return response, nil

}

// CompressMessages compresses a list of messages into a summary
func (c *CompressorAgent) CompressMessages(messages []*ai.Message) (ChatResponse, error) {
	// Convert messages to text format
	var textBuilder strings.Builder
	for _, msg := range messages {
		textBuilder.WriteString(fmt.Sprintf("%s: ", msg.Role))
		for _, content := range msg.Content {
			textBuilder.WriteString(content.Text)
			textBuilder.WriteString("\n")
		}
	}

	text := textBuilder.String()
	return c.CompressText(text)
}

// CompressMessagesStream compresses a list of messages into a summary using streaming
func (c *CompressorAgent) CompressMessagesStream(messages []*ai.Message, callback func(ChatResponse) error) (ChatResponse, error) {
	// Convert messages to text format
	var textBuilder strings.Builder
	for _, msg := range messages {
		textBuilder.WriteString(fmt.Sprintf("%s: ", msg.Role))
		for _, content := range msg.Content {
			textBuilder.WriteString(content.Text)
			textBuilder.WriteString("\n")
		}
	}

	text := textBuilder.String()
	return c.CompressTextStream(text, callback)
}


type CompressionPrompts struct {
	Minimalist string
	Structured string
	UltraShort string
	ContinuityFocus string
}

var DefaultCompressionPrompts = CompressionPrompts{
	//recommended
	Minimalist: `Summarize the conversation history concisely, preserving key facts, decisions, and context needed for continuation.`,
	Structured: `Compress this conversation into a brief summary including:
		- Main topics discussed
		- Key decisions/conclusions
		- Important context for next exchanges
		Keep it under 200 words.
	`,
	UltraShort: `Summarize this conversation: extract key facts, decisions, and essential context only.`,
	ContinuityFocus: `Create a compact summary of this conversation that preserves all information needed to continue the discussion naturally.`,
}
