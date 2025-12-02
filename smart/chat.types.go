package smart

type ChatRequest struct {
	UserMessage string `json:"message"`
}

// Structure for final flow output
type ChatResponse struct {
	Response string `json:"response"`
}

