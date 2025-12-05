package snip


type ChatRequest struct {
	UserMessage string `json:"message"`
}

// Structure for final flow output
type ChatResponse struct {
	Text string `json:"response"`
	FinishReason string `json:"finish_reason,omitempty"`
	FinishMessage string `json:"finish_message,omitempty"`
}
