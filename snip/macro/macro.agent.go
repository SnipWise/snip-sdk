package macro

import (
	"fmt"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/snipwise/snip-sdk/snip"
	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

type MacroAgent struct {
	Name string

	chatAgent       snip.AIChatAgent
	compressorAgent snip.AICompressorAgent

	logger logger.Logger
}

func NewMacroAgent(
	name string,
	chatAgent snip.AIChatAgent,
	compressorAgent snip.AICompressorAgent,
	opts ...MacroAgentOption,

) (*MacroAgent, error) {

	macroAgent := &MacroAgent{
		Name:            name,
		chatAgent:       chatAgent,
		compressorAgent: compressorAgent,
		logger:          &logger.NoOpLogger{},
	}

	for _, opt := range opts {
		opt(macroAgent)
	}

	return macroAgent, nil
}

func (macroAgent *MacroAgent) GetName() string {
	return macroAgent.Name
}

func (macroAgent *MacroAgent) Kind() agents.AgentKind {
	return agents.Macro
}

// +++ Chat Agent methods delegated +++

func (macroAgent *MacroAgent) GetMessages() []*ai.Message {
	return macroAgent.chatAgent.GetMessages()
}

func (macroAgent *MacroAgent) GetCurrentContextSize() int {
	return macroAgent.chatAgent.GetCurrentContextSize()
}

func (macroAgent *MacroAgent) AddSystemMessage(context string) error {
	return macroAgent.chatAgent.AddSystemMessage(context)
}

func (macroAgent *MacroAgent) ReplaceMessagesWith(messages []*ai.Message) error {
	return macroAgent.chatAgent.ReplaceMessagesWith(messages)
}

func (macroAgent *MacroAgent) ReplaceMessagesWithSystemMessages(systemMessages []string) error {
	return macroAgent.chatAgent.ReplaceMessagesWithSystemMessages(systemMessages)
}

func (macroAgent *MacroAgent) AskWithMemory(question string) (agents.ChatResponse, error) {
	return macroAgent.chatAgent.AskWithMemory(question)
}

func (macroAgent *MacroAgent) AskStreamWithMemory(question string, callback func(agents.ChatResponse) error) (agents.ChatResponse, error) {
	return macroAgent.chatAgent.AskStreamWithMemory(question, callback)
}

func (macroAgent *MacroAgent) Ask(question string) (agents.ChatResponse, error) {
	return macroAgent.chatAgent.Ask(question)
}

func (macroAgent *MacroAgent) AskStream(question string, callback func(agents.ChatResponse) error) (agents.ChatResponse, error) {
	return macroAgent.chatAgent.AskStream(question, callback)
}

// func (macroAgent *MacroAgent) GetChatFlowWithMemory() *core.Flow[*agents.ChatRequest, *agents.ChatResponse, struct{}] {
// 	return macroAgent.chatAgent.GetChatFlowWithMemory()
// }

// func (macroAgent *MacroAgent) GetChatStreamFlowWithMemory() *core.Flow[*agents.ChatRequest, *agents.ChatResponse, agents.ChatResponse] {
// 	return macroAgent.chatAgent.GetChatStreamFlowWithMemory()
// }

// func (macroAgent *MacroAgent) GetStreamCancel() func() {
// 	return macroAgent.chatAgent.GetStreamCancel()
// }

// +++ Compressor Agent methods delegated +++

// CompressContext compresses the conversation history using the configured compressor agent
// Returns an error if no compressor agent is configured
// The compression result is returned as a ChatResponse
// After compression, the agent's messages are replaced with a single system message containing the compressed context
func (macroAgent *MacroAgent) CompressContext() (agents.ChatResponse, error) {
	if macroAgent.compressorAgent == nil {
		return agents.ChatResponse{}, fmt.Errorf("no compressor agent configured, use EnableContextCompression option")
	}

	response, err := macroAgent.compressorAgent.CompressMessages(macroAgent.chatAgent.GetMessages())
	if err != nil {
		return agents.ChatResponse{}, err
	}

	// Replace the agent's messages with the compressed context
	compressedMessages := []*ai.Message{
		ai.NewSystemTextMessage(strings.TrimSpace(response.Text)),
	}
	if err := macroAgent.chatAgent.ReplaceMessagesWith(compressedMessages); err != nil {
		return agents.ChatResponse{}, err
	}

	return response, nil
}

// CompressContextStream compresses the conversation history using streaming with the configured compressor agent
// Returns an error if no compressor agent is configured
// The callback function is called for each streamed chunk
// The final compression result is returned as a ChatResponse
// After compression, the agent's messages are replaced with a single system message containing the compressed context
func (macroAgent *MacroAgent) CompressContextStream(callback func(agents.ChatResponse) error) (agents.ChatResponse, error) {
	if macroAgent.compressorAgent == nil {
		return agents.ChatResponse{}, fmt.Errorf("no compressor agent configured, use EnableContextCompression option")
	}

	response, err := macroAgent.compressorAgent.CompressMessagesStream(macroAgent.chatAgent.GetMessages(), callback)
	if err != nil {
		return agents.ChatResponse{}, err
	}

	// Replace the agent's messages with the compressed context
	compressedMessages := []*ai.Message{
		ai.NewSystemTextMessage(strings.TrimSpace(response.Text)),
	}
	if err := macroAgent.chatAgent.ReplaceMessagesWith(compressedMessages); err != nil {
		return agents.ChatResponse{}, err
	}

	return response, nil
}
