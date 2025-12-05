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

func (agent *RemoteAgent) Ask(question string) (string, error) {
	// Prepare request
	reqBody := RemoteChatRequest{}
	reqBody.Data.Message = strings.TrimSpace(question)

	// Convert to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error creating JSON: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", agent.ChatEndPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error during HTTP call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON response
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	// Extract message from response (try different possible keys)
	// Genkit wraps the response in a "result" object
	if resultObj, ok := result["result"].(map[string]any); ok {
		if response, ok := resultObj["response"].(string); ok {
			return response, nil
		}
		if message, ok := resultObj["message"].(string); ok {
			return message, nil
		}
		if text, ok := resultObj["text"].(string); ok {
			return text, nil
		}
	}

	// Try "response" key (used by ChatResponse struct)
	if response, ok := result["response"].(string); ok {
		return response, nil
	}
	// Try "message" key
	if message, ok := result["message"].(string); ok {
		return message, nil
	}
	// Try "text" key
	if text, ok := result["text"].(string); ok {
		return text, nil
	}
	// Try nested "data.message" or "data.response"
	if data, ok := result["data"].(map[string]any); ok {
		if response, ok := data["response"].(string); ok {
			return response, nil
		}
		if message, ok := data["message"].(string); ok {
			return message, nil
		}
	}

	return "", fmt.Errorf("unable to extract message from response")
}

func (agent *RemoteAgent) AskStream(question string, callback func(string) error) (string, error) {
	// Prepare request
	reqBody := RemoteChatRequest{}
	reqBody.Data.Message = strings.TrimSpace(question)

	// Convert to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error when creating JSON: %v\n", err)
		//continue
		return "", err
	}
	// Create HTTP request
	req, err := http.NewRequest("POST", agent.ChatStreamEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error when creating the request: %v\n", err)
		//continue
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error when HTTP call: %v\n", err)
		//continue
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP error: status code %d\n", resp.StatusCode)
		resp.Body.Close()
		//continue
		return "", fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
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
				// Display content if available (try "message" then "text")
				if message, ok := chunk["message"].(string); ok {
					//fmt.Print(message)
					fullResponse += message
					callback(message)
				} else if text, ok := chunk["text"].(string); ok {
					//fmt.Print(text)
					fullResponse += text
					callback(text)
				}
			} else {
				// If not JSON, print as is
				//fmt.Print(data)
				fullResponse += data
				callback(data)
			}
			callbackErr = err
		}
	}
	resp.Body.Close()
	return fullResponse, callbackErr
}
