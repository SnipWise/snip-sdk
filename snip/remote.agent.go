package snip

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/firebase/genkit/go/ai"
)

// Structure for flow input
type RemoteChatRequest struct {
	Data struct {
		Message string `json:"message"`
	} `json:"data"`
}

type RemoteAgent struct {
	ChatStreamEndpoint  string
	ChatEndPoint        string
	InformationEndpoint string
	AddContextEndpoint  string
	GetMessagesEndpoint string
	Name                string
}

func NewRemoteAgent(name string, config ConfigHTTP) *RemoteAgent {
	// Build full URLs from Address and paths
	baseURL := "http://" + config.Address

	// Set default information path if not provided
	informationPath := config.InformationPath
	if informationPath == "" {
		informationPath = DefaultInformationPath
	}

	// Set default add context path if not provided
	addContextPath := config.AddContextPath
	if addContextPath == "" {
		addContextPath = DefaultAddSystemMessagePath
	}

	// Set default get messages path if not provided
	getMessagesPath := config.GetMessagesPath
	if getMessagesPath == "" {
		getMessagesPath = DefaultGetMessagesPath
	}

	return &RemoteAgent{
		ChatStreamEndpoint:  baseURL + config.ChatStreamFlowPath,
		ChatEndPoint:        baseURL + config.ChatFlowPath,
		InformationEndpoint: baseURL + informationPath,
		AddContextEndpoint:  baseURL + addContextPath,
		GetMessagesEndpoint: baseURL + getMessagesPath,
		Name:                name,
	}
}

func (agent *RemoteAgent) GetName() string {
	return agent.Name
}

func (agent *RemoteAgent) Kind() AgentKind {
	return Remote
}

func (agent *RemoteAgent) AddSystemMessage(context string) error {
	// Prepare request
	reqBody := struct {
		Context string `json:"context"`
	}{
		Context: strings.TrimSpace(context),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error creating JSON: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", agent.AddContextEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error during HTTP call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP error: status code %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (agent *RemoteAgent) ReplaceMessagesWith(messages []*ai.Message) error {
	// Remote agents do not support replacing messages directly
	// Message management must be done on the server side
	return fmt.Errorf("ReplaceMessagesWith is not supported for remote agents: message history is managed by the remote server")
}

func (agent *RemoteAgent) ReplaceMessagesWithSystemMessages(systemMessages []string) error {
	// Remote agents do not support replacing messages directly
	// Message management must be done on the server side
	return fmt.Errorf("ReplaceMessagesWithSystemMessages is not supported for remote agents: message history is managed by the remote server")
}

func (agent *RemoteAgent) GetInfo() (AgentInfo, error) {
	// Create HTTP GET request to information endpoint
	req, err := http.NewRequest("GET", agent.InformationEndpoint, nil)
	if err != nil {
		return AgentInfo{}, fmt.Errorf("error creating request: %w", err)
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AgentInfo{}, fmt.Errorf("error during HTTP call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AgentInfo{}, fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AgentInfo{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON response into AgentInfo struct
	var info AgentInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return AgentInfo{}, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return info, nil
}

func (agent *RemoteAgent) GetMessages() []*ai.Message {
	// Create HTTP GET request to get messages endpoint
	req, err := http.NewRequest("GET", agent.GetMessagesEndpoint, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error during HTTP call: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP error: status code %d\n", resp.StatusCode)
		return nil
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil
	}

	// Parse JSON response into []*ai.Message
	var messages []*ai.Message
	if err := json.Unmarshal(body, &messages); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return nil
	}

	return messages
}

func (agent *RemoteAgent) GetCurrentContextSize() int {
	totalContextSize := 0

	// Get messages from remote server
	messages := agent.GetMessages()
	if messages != nil {
		for _, msg := range messages {
			for _, content := range msg.Content {
				totalContextSize += len(content.Text)
			}
		}
	}

	return totalContextSize
}

func (agent *RemoteAgent) AskWithMemory(question string) (ChatResponse, error) {
	// Prepare request
	reqBody := RemoteChatRequest{}
	reqBody.Data.Message = strings.TrimSpace(question)

	// Convert to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("error creating JSON: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", agent.ChatEndPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return ChatResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("error during HTTP call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ChatResponse{}, fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON response
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return ChatResponse{}, fmt.Errorf("error parsing JSON response: %w", err)
	}

	// Extract message from response (try different possible keys)
	// Genkit wraps the response in a "result" object
	if resultObj, ok := result["result"].(map[string]any); ok {
		if response, ok := resultObj["response"].(string); ok {
			return ChatResponse{Text: response}, nil
		}
		if message, ok := resultObj["message"].(string); ok {
			return ChatResponse{Text: message}, nil
		}
		if text, ok := resultObj["text"].(string); ok {
			return ChatResponse{Text: text}, nil
		}
	}

	// Try "response" key (used by ChatResponse struct)
	if response, ok := result["response"].(string); ok {
		return ChatResponse{Text: response}, nil
	}
	// Try "message" key
	if message, ok := result["message"].(string); ok {
		return ChatResponse{Text: message}, nil
	}
	// Try "text" key
	if text, ok := result["text"].(string); ok {
		return ChatResponse{Text: text}, nil
	}
	// Try nested "data.message" or "data.response"
	if data, ok := result["data"].(map[string]any); ok {
		if response, ok := data["response"].(string); ok {
			return ChatResponse{Text: response}, nil
		}
		if message, ok := data["message"].(string); ok {
			return ChatResponse{Text: message}, nil
		}
	}

	return ChatResponse{}, fmt.Errorf("unable to extract message from response")
}

func (agent *RemoteAgent) AskStreamWithMemory(question string, callback func(ChatResponse) error) (ChatResponse, error) {
	// Prepare request
	reqBody := RemoteChatRequest{}
	reqBody.Data.Message = strings.TrimSpace(question)

	// Convert to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error when creating JSON: %v\n", err)
		return ChatResponse{}, err
	}
	// Create HTTP request
	req, err := http.NewRequest("POST", agent.ChatStreamEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error when creating the request: %v\n", err)
		return ChatResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error when HTTP call: %v\n", err)
		return ChatResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP error: status code %d\n", resp.StatusCode)
		resp.Body.Close()
		return ChatResponse{}, fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	// Read the stream
	streamReader := bufio.NewReader(resp.Body)
	fullResponse := ""
	var callbackErr error
	for {
		line, err := streamReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			callbackErr = err
			fmt.Printf("\nError when stream reading: %v\n", err)
			break
		}

		// Read data lines starting with "data: "
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "data: "); ok {
			data := after

			// Check for end of stream
			if data == "[DONE]" {
				break
			}

			// Parse JSON to extract content
			var chunk map[string]any
			if err := json.Unmarshal([]byte(data), &chunk); err == nil {
				var textContent string
				var finishReason string

				// Try to extract from nested message.response structure (Genkit streaming format)
				if messageObj, ok := chunk["message"].(map[string]any); ok {
					if response, ok := messageObj["response"].(string); ok {
						textContent = response
					}
					if fr, ok := messageObj["finish_reason"].(string); ok {
						finishReason = fr
					}
				}

				// Try result.response for final chunk
				if resultObj, ok := chunk["result"].(map[string]any); ok {
					if response, ok := resultObj["response"].(string); ok {
						textContent = response
					}
					if fr, ok := resultObj["finish_reason"].(string); ok {
						finishReason = fr
					}
				}

				// Fallback: try direct "message" or "text" fields (legacy support)
				if textContent == "" {
					if message, ok := chunk["message"].(string); ok {
						textContent = message
					} else if text, ok := chunk["text"].(string); ok {
						textContent = text
					}
				}

				// Send callback if we have content or finish reason
				if textContent != "" || finishReason != "" {
					fullResponse += textContent
					if cbErr := callback(ChatResponse{
						Text:         textContent,
						FinishReason: finishReason,
					}); cbErr != nil {
						callbackErr = cbErr
					}
				}
			} else {
				// If not JSON, print as is
				fullResponse += data
				if cbErr := callback(ChatResponse{Text: data}); cbErr != nil {
					callbackErr = cbErr
				}
			}
		}
	}
	resp.Body.Close()

	// Note: Genkit already sends a final chunk with finish_reason in the stream,
	// so we don't need to send an additional one here (unlike local agents)

	return ChatResponse{Text: fullResponse}, callbackErr
}

// Ask is an alias for AskWithMemory for RemoteAgent
// Remote agents delegate memory management to the server
func (agent *RemoteAgent) Ask(question string) (ChatResponse, error) {
	return agent.AskWithMemory(question)
}

// AskStream is an alias for AskStreamWithMemory for RemoteAgent
// Remote agents delegate memory management to the server
func (agent *RemoteAgent) AskStream(question string, callback func(ChatResponse) error) (ChatResponse, error) {
	return agent.AskStreamWithMemory(question, callback)
}

// CompressContext is not supported for remote agents
// Context compression must be performed on the server side
func (agent *RemoteAgent) CompressContext() (ChatResponse, error) {
	return ChatResponse{}, fmt.Errorf("CompressContext is not supported for remote agents: compression must be performed on the server side")
}

// CompressContextStream is not supported for remote agents
// Context compression must be performed on the server side
func (agent *RemoteAgent) CompressContextStream(callback func(ChatResponse) error) (ChatResponse, error) {
	return ChatResponse{}, fmt.Errorf("CompressContextStream is not supported for remote agents: compression must be performed on the server side")
}
