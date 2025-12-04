package smart

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/firebase/genkit/go/ai"
)

// ============================================================================
// Tests for NewRemoteAgent
// ============================================================================

func TestNewRemoteAgent(t *testing.T) {
	t.Run("with all custom paths", func(t *testing.T) {
		config := ConfigHTTP{
			Address:            "localhost:8080",
			ChatFlowPath:       "/custom/chat",
			ChatStreamFlowPath: "/custom/stream",
			InformationPath:    "/custom/info",
			AddContextPath:     "/custom/context",
			GetMessagesPath:    "/custom/messages",
		}

		agent := NewRemoteAgent("test-agent", config)

		if agent.Name != "test-agent" {
			t.Errorf("Name = %q, want %q", agent.Name, "test-agent")
		}

		expectedChatEndpoint := "http://localhost:8080/custom/chat"
		if agent.ChatEndPoint != expectedChatEndpoint {
			t.Errorf("ChatEndPoint = %q, want %q", agent.ChatEndPoint, expectedChatEndpoint)
		}

		expectedStreamEndpoint := "http://localhost:8080/custom/stream"
		if agent.ChatStreamEndpoint != expectedStreamEndpoint {
			t.Errorf("ChatStreamEndpoint = %q, want %q", agent.ChatStreamEndpoint, expectedStreamEndpoint)
		}

		expectedInfoEndpoint := "http://localhost:8080/custom/info"
		if agent.InformationEndpoint != expectedInfoEndpoint {
			t.Errorf("InformationEndpoint = %q, want %q", agent.InformationEndpoint, expectedInfoEndpoint)
		}

		expectedContextEndpoint := "http://localhost:8080/custom/context"
		if agent.AddContextEndpoint != expectedContextEndpoint {
			t.Errorf("AddContextEndpoint = %q, want %q", agent.AddContextEndpoint, expectedContextEndpoint)
		}

		expectedMessagesEndpoint := "http://localhost:8080/custom/messages"
		if agent.GetMessagesEndpoint != expectedMessagesEndpoint {
			t.Errorf("GetMessagesEndpoint = %q, want %q", agent.GetMessagesEndpoint, expectedMessagesEndpoint)
		}
	})

	t.Run("with default paths", func(t *testing.T) {
		config := ConfigHTTP{
			Address:            "localhost:9000",
			ChatFlowPath:       "/api/chat",
			ChatStreamFlowPath: "/api/chat-stream",
			// InformationPath not set, should use default
			// AddContextPath not set, should use default
			// GetMessagesPath not set, should use default
		}

		agent := NewRemoteAgent("default-agent", config)

		expectedInfoEndpoint := "http://localhost:9000" + DefaultInformationPath
		if agent.InformationEndpoint != expectedInfoEndpoint {
			t.Errorf("InformationEndpoint = %q, want %q", agent.InformationEndpoint, expectedInfoEndpoint)
		}

		expectedContextEndpoint := "http://localhost:9000" + DefaultAddSystemMessagePath
		if agent.AddContextEndpoint != expectedContextEndpoint {
			t.Errorf("AddContextEndpoint = %q, want %q", agent.AddContextEndpoint, expectedContextEndpoint)
		}

		expectedMessagesEndpoint := "http://localhost:9000" + DefaultGetMessagesPath
		if agent.GetMessagesEndpoint != expectedMessagesEndpoint {
			t.Errorf("GetMessagesEndpoint = %q, want %q", agent.GetMessagesEndpoint, expectedMessagesEndpoint)
		}
	})
}

// ============================================================================
// Tests for RemoteAgent.GetName
// ============================================================================

func TestRemoteAgentGetName(t *testing.T) {
	agent := &RemoteAgent{
		Name: "remote-test-agent",
	}

	name := agent.GetName()
	if name != "remote-test-agent" {
		t.Errorf("GetName() = %q, want %q", name, "remote-test-agent")
	}
}

// ============================================================================
// Tests for RemoteAgent.Kind
// ============================================================================

func TestRemoteAgentKind(t *testing.T) {
	agent := &RemoteAgent{}

	kind := agent.Kind()
	if kind != Remote {
		t.Errorf("Kind() = %v, want %v", kind, Remote)
	}
}

// ============================================================================
// Tests for RemoteAgent.AddSystemMessage
// ============================================================================

func TestRemoteAgentAddSystemMessage(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify method
			if r.Method != "POST" {
				t.Errorf("Request method = %q, want POST", r.Method)
			}

			// Verify content type
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", contentType)
			}

			// Parse request body
			var reqBody struct {
				Context string `json:"context"`
			}
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}

			// Verify context
			expectedContext := "Test context"
			if reqBody.Context != expectedContext {
				t.Errorf("Context = %q, want %q", reqBody.Context, expectedContext)
			}

			// Send success response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		}))
		defer server.Close()

		agent := &RemoteAgent{
			AddContextEndpoint: server.URL,
		}

		err := agent.AddSystemMessage("Test context")
		if err != nil {
			t.Errorf("AddSystemMessage() unexpected error: %v", err)
		}
	})

	t.Run("trim whitespace", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var reqBody struct {
				Context string `json:"context"`
			}
			json.NewDecoder(r.Body).Decode(&reqBody)

			// Context should be trimmed
			if reqBody.Context != "Trimmed context" {
				t.Errorf("Context = %q, want %q", reqBody.Context, "Trimmed context")
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		}))
		defer server.Close()

		agent := &RemoteAgent{
			AddContextEndpoint: server.URL,
		}

		err := agent.AddSystemMessage("  \n  Trimmed context  \n  ")
		if err != nil {
			t.Errorf("AddSystemMessage() unexpected error: %v", err)
		}
	})

	t.Run("http error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
		}))
		defer server.Close()

		agent := &RemoteAgent{
			AddContextEndpoint: server.URL,
		}

		err := agent.AddSystemMessage("Test context")
		if err == nil {
			t.Error("AddSystemMessage() expected error for HTTP 500, got nil")
		}
	})

	t.Run("invalid endpoint", func(t *testing.T) {
		agent := &RemoteAgent{
			AddContextEndpoint: "http://invalid-host-that-does-not-exist:99999",
		}

		err := agent.AddSystemMessage("Test context")
		if err == nil {
			t.Error("AddSystemMessage() expected error for invalid endpoint, got nil")
		}
	})
}

// ============================================================================
// Tests for RemoteAgent.GetInfo
// ============================================================================

func TestRemoteAgentGetInfo(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify method
			if r.Method != "GET" {
				t.Errorf("Request method = %q, want GET", r.Method)
			}

			// Send response
			response := AgentInfo{
				Name:    "test-agent",
				ModelID: "test-model",
				Config: ModelConfig{
					Temperature: 0.7,
					TopP:        0.9,
					MaxTokens:   1000,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		agent := &RemoteAgent{
			InformationEndpoint: server.URL,
		}

		info, err := agent.GetInfo()
		if err != nil {
			t.Errorf("GetInfo() unexpected error: %v", err)
		}

		if info.Name != "test-agent" {
			t.Errorf("Info.Name = %q, want %q", info.Name, "test-agent")
		}

		if info.ModelID != "test-model" {
			t.Errorf("Info.ModelID = %q, want %q", info.ModelID, "test-model")
		}

		if info.Config.Temperature != 0.7 {
			t.Errorf("Info.Config.Temperature = %f, want 0.7", info.Config.Temperature)
		}
	})

	t.Run("http error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		agent := &RemoteAgent{
			InformationEndpoint: server.URL,
		}

		_, err := agent.GetInfo()
		if err == nil {
			t.Error("GetInfo() expected error for HTTP 404, got nil")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid json`))
		}))
		defer server.Close()

		agent := &RemoteAgent{
			InformationEndpoint: server.URL,
		}

		_, err := agent.GetInfo()
		if err == nil {
			t.Error("GetInfo() expected error for invalid JSON, got nil")
		}
	})
}

// ============================================================================
// Tests for RemoteAgent.Ask
// ============================================================================

func TestRemoteAgentAsk(t *testing.T) {
	t.Run("successful request with result wrapper", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify method
			if r.Method != "POST" {
				t.Errorf("Request method = %q, want POST", r.Method)
			}

			// Parse request
			var reqBody RemoteChatRequest
			json.NewDecoder(r.Body).Decode(&reqBody)

			if reqBody.Data.Message != "Test question" {
				t.Errorf("Message = %q, want %q", reqBody.Data.Message, "Test question")
			}

			// Send response (Genkit format)
			response := map[string]any{
				"result": map[string]any{
					"response": "Test answer",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		agent := &RemoteAgent{
			ChatEndPoint: server.URL,
		}

		answer, err := agent.Ask("Test question")
		if err != nil {
			t.Errorf("Ask() unexpected error: %v", err)
		}

		if answer != "Test answer" {
			t.Errorf("Ask() = %q, want %q", answer, "Test answer")
		}
	})

	t.Run("successful request with direct response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]any{
				"response": "Direct answer",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		agent := &RemoteAgent{
			ChatEndPoint: server.URL,
		}

		answer, err := agent.Ask("Test question")
		if err != nil {
			t.Errorf("Ask() unexpected error: %v", err)
		}

		if answer != "Direct answer" {
			t.Errorf("Ask() = %q, want %q", answer, "Direct answer")
		}
	})

	t.Run("trim question", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var reqBody RemoteChatRequest
			json.NewDecoder(r.Body).Decode(&reqBody)

			if reqBody.Data.Message != "Trimmed question" {
				t.Errorf("Message = %q, want %q", reqBody.Data.Message, "Trimmed question")
			}

			response := map[string]any{"response": "Answer"}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		agent := &RemoteAgent{
			ChatEndPoint: server.URL,
		}

		_, err := agent.Ask("  \n  Trimmed question  \n  ")
		if err != nil {
			t.Errorf("Ask() unexpected error: %v", err)
		}
	})

	t.Run("http error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		agent := &RemoteAgent{
			ChatEndPoint: server.URL,
		}

		_, err := agent.Ask("Test question")
		if err == nil {
			t.Error("Ask() expected error for HTTP 500, got nil")
		}
	})

	t.Run("no extractable message", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]any{
				"unknown": "field",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		agent := &RemoteAgent{
			ChatEndPoint: server.URL,
		}

		_, err := agent.Ask("Test question")
		if err == nil {
			t.Error("Ask() expected error for no extractable message, got nil")
		}

		if !strings.Contains(err.Error(), "unable to extract message") {
			t.Errorf("Ask() error = %q, want error containing 'unable to extract message'", err.Error())
		}
	})
}

// ============================================================================
// Tests for RemoteAgent.GetMessages
// ============================================================================

func TestRemoteAgentGetMessages(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify method
			if r.Method != "GET" {
				t.Errorf("Request method = %q, want GET", r.Method)
			}

			// Send response with messages
			messages := []*ai.Message{
				ai.NewUserTextMessage("Hello"),
				ai.NewModelTextMessage("Hi there"),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(messages)
		}))
		defer server.Close()

		agent := &RemoteAgent{
			GetMessagesEndpoint: server.URL,
		}

		messages := agent.GetMessages()
		if messages == nil {
			t.Fatal("GetMessages() returned nil")
		}

		if len(messages) != 2 {
			t.Errorf("GetMessages() length = %d, want 2", len(messages))
		}
	})

	t.Run("empty messages", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}))
		defer server.Close()

		agent := &RemoteAgent{
			GetMessagesEndpoint: server.URL,
		}

		messages := agent.GetMessages()
		if messages == nil {
			t.Fatal("GetMessages() returned nil")
		}

		if len(messages) != 0 {
			t.Errorf("GetMessages() length = %d, want 0", len(messages))
		}
	})

	t.Run("http error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		agent := &RemoteAgent{
			GetMessagesEndpoint: server.URL,
		}

		messages := agent.GetMessages()
		if messages != nil {
			t.Errorf("GetMessages() expected nil for HTTP error, got %v", messages)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid json`))
		}))
		defer server.Close()

		agent := &RemoteAgent{
			GetMessagesEndpoint: server.URL,
		}

		messages := agent.GetMessages()
		if messages != nil {
			t.Errorf("GetMessages() expected nil for invalid JSON, got %v", messages)
		}
	})
}
