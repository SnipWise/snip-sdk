package agents

import "github.com/firebase/genkit/go/ai"

type ChatRequest struct {
	UserMessage string `json:"message"`
}

// Structure for final flow output
type ChatResponse struct {
	Text          string     `json:"response"`
	Content       []*ai.Part `json:"content,omitempty"`
	Role          ai.Role   `json:"role,omitempty"`
	FinishReason  string `json:"finish_reason,omitempty"`
	FinishMessage string `json:"finish_message,omitempty"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

func (chatResponse *ChatResponse) IsEmpty() bool {
	return chatResponse.Text == ""
}

func (chatResponse *ChatResponse) IsFinishReasonStop() bool {
	return chatResponse.FinishReason == "stop"
}

func (chatResponse *ChatResponse) IsFinishReasonLength() bool {
	return chatResponse.FinishReason == "length"
}

func (chatResponse *ChatResponse) IsFinishReasonContentFilter() bool {
	return chatResponse.FinishReason == "content_filter"
}

func (chatResponse *ChatResponse) IsFinishReasonUnknown() bool {
	return chatResponse.FinishReason == "unknown"
}
