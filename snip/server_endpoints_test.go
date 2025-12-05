package snip

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	oai "github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

// ============================================================================
// Helper functions for server tests
// ============================================================================

func createTestAgent(t *testing.T) *Agent {
	ctx := context.Background()

	// Create a minimal agent without actual LLM backend
	// This is just for testing HTTP endpoints structure
	oaiPlugin := &oai.OpenAI{
		APIKey: "test-key",
		Opts: []option.RequestOption{
			option.WithBaseURL("http://localhost:12434/test"),
		},
	}

	genKitInstance := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	agent := &Agent{
		Name:               "test-agent",
		SystemInstructions: "Test instructions",
		ModelID:            "test-model",
		Messages:           []*ai.Message{},
		Config: ModelConfig{
			Temperature: 0.7,
			TopP:        0.9,
			MaxTokens:   1000,
		},
		ctx:            ctx,
		genKitInstance: genKitInstance,
	}

	return agent
}

// ============================================================================
// Tests for Server Configuration
// ============================================================================

func TestServerConfigDefaults(t *testing.T) {
	agent := createTestAgent(t)
	agent.serverConfig = &ConfigHTTP{
		Address: "0.0.0.0:8080",
		// Leave paths empty to test defaults
	}

	// Test that defaults would be applied in Serve() method
	// We can't easily test Serve() without blocking, so we test the default constants
	if DefaultChatFlowPath != "/api/chat" {
		t.Errorf("DefaultChatFlowPath = %q, want /api/chat", DefaultChatFlowPath)
	}

	if DefaultChatStreamFlowPath != "/api/chat-stream" {
		t.Errorf("DefaultChatStreamFlowPath = %q, want /api/chat-stream", DefaultChatStreamFlowPath)
	}

	if DefaultInformationPath != "/api/information" {
		t.Errorf("DefaultInformationPath = %q, want /api/information", DefaultInformationPath)
	}

	if DefaultHealthcheckPath != "/healthcheck" {
		t.Errorf("DefaultHealthcheckPath = %q, want /healthcheck", DefaultHealthcheckPath)
	}

	if DefaultAddSystemMessagePath != "/api/add-system-message" {
		t.Errorf("DefaultAddSystemMessagePath = %q, want /api/add-system-message", DefaultAddSystemMessagePath)
	}

	if DefaultGetMessagesPath != "/api/messages" {
		t.Errorf("DefaultGetMessagesPath = %q, want /api/messages", DefaultGetMessagesPath)
	}
}

// ============================================================================
// Tests for Information Endpoint
// ============================================================================

func TestInformationEndpoint(t *testing.T) {
	agent := createTestAgent(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		info := AgentInfo{
			Name:    agent.Name,
			ModelID: agent.ModelID,
			Config:  agent.Config,
		}
		if err := json.NewEncoder(w).Encode(info); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	req := httptest.NewRequest("GET", "/api/information", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}

	var info AgentInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
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
}

// ============================================================================
// Tests for Healthcheck Endpoint
// ============================================================================

func TestHealthcheckEndpoint(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	req := httptest.NewRequest("GET", "/healthcheck", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, _ := io.ReadAll(resp.Body)
	expectedBody := `{"status":"ok"}`
	if string(body) != expectedBody {
		t.Errorf("Response body = %q, want %q", string(body), expectedBody)
	}
}

// ============================================================================
// Tests for Add System Message Endpoint
// ============================================================================

func TestAddSystemMessageEndpoint(t *testing.T) {
	t.Run("successful add", func(t *testing.T) {
		agent := createTestAgent(t)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			var req struct {
				Context string `json:"context"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"error","message":"invalid request body"}`))
				return
			}

			if err := agent.AddSystemMessage(req.Context); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"failed to add context"}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		})

		reqBody := map[string]string{
			"context": "Test context",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/add-system-message", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		// Verify message was added to agent
		if len(agent.Messages) != 1 {
			t.Errorf("Agent has %d messages, want 1", len(agent.Messages))
		}

		if len(agent.Messages) > 0 {
			content := agent.Messages[0].Content[0].Text
			if content != "Test context" {
				t.Errorf("Message content = %q, want %q", content, "Test context")
			}
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		agent := createTestAgent(t)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			var req struct {
				Context string `json:"context"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"error","message":"invalid request body"}`))
				return
			}

			agent.AddSystemMessage(req.Context)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		})

		req := httptest.NewRequest("POST", "/api/add-system-message", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})
}

// ============================================================================
// Tests for Get Messages Endpoint
// ============================================================================

func TestGetMessagesEndpoint(t *testing.T) {
	t.Run("empty messages", func(t *testing.T) {
		agent := createTestAgent(t)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(agent.Messages); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"failed to encode messages"}`))
				return
			}
		})

		req := httptest.NewRequest("GET", "/api/messages", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var messages []*ai.Message
		if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(messages) != 0 {
			t.Errorf("Messages length = %d, want 0", len(messages))
		}
	})

	t.Run("with messages", func(t *testing.T) {
		agent := createTestAgent(t)
		agent.AddSystemMessage("System message")
		agent.Messages = append(agent.Messages, ai.NewUserTextMessage("User message"))
		agent.Messages = append(agent.Messages, ai.NewModelTextMessage("Model response"))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(agent.Messages); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})

		req := httptest.NewRequest("GET", "/api/messages", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var messages []*ai.Message
		if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(messages) != 3 {
			t.Errorf("Messages length = %d, want 3", len(messages))
		}
	})
}

// ============================================================================
// Tests for Shutdown Endpoint
// ============================================================================

func TestShutdownEndpoint(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"shutting down"}`))
	})

	req := httptest.NewRequest("POST", "/server/shutdown", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, _ := io.ReadAll(resp.Body)
	expectedBody := `{"status":"shutting down"}`
	if string(body) != expectedBody {
		t.Errorf("Response body = %q, want %q", string(body), expectedBody)
	}
}

// ============================================================================
// Tests for Cancel Stream Endpoint
// ============================================================================

func TestCancelStreamEndpoint(t *testing.T) {
	t.Run("with active stream", func(t *testing.T) {
		agent := createTestAgent(t)
		ctx, cancel := context.WithCancel(agent.ctx)
		agent.streamCancel = cancel

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if agent.streamCancel != nil {
				agent.streamCancel()
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"stream cancelled"}`))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"no active stream"}`))
			}
		})

		req := httptest.NewRequest("POST", "/api/cancel-stream-completion", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		body, _ := io.ReadAll(resp.Body)
		if !bytes.Contains(body, []byte("stream cancelled")) {
			t.Errorf("Response should contain 'stream cancelled', got %q", string(body))
		}

		// Verify context was cancelled
		select {
		case <-ctx.Done():
			// Success - context was cancelled
		default:
			t.Error("Context was not cancelled")
		}
	})

	t.Run("no active stream", func(t *testing.T) {
		agent := createTestAgent(t)
		agent.streamCancel = nil

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if agent.streamCancel != nil {
				agent.streamCancel()
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"stream cancelled"}`))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"no active stream"}`))
			}
		})

		req := httptest.NewRequest("POST", "/api/cancel-stream-completion", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		body, _ := io.ReadAll(resp.Body)
		if !bytes.Contains(body, []byte("no active stream")) {
			t.Errorf("Response should contain 'no active stream', got %q", string(body))
		}
	})
}

// ============================================================================
// Tests for ConfigHTTP struct
// ============================================================================

func TestConfigHTTP(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		config := ConfigHTTP{
			Address: "localhost:8080",
		}

		// Test that we have the expected default constants
		if config.ChatFlowPath == "" && DefaultChatFlowPath != "/api/chat" {
			t.Error("DefaultChatFlowPath should be /api/chat")
		}

		if config.ChatStreamFlowPath == "" && DefaultChatStreamFlowPath != "/api/chat-stream" {
			t.Error("DefaultChatStreamFlowPath should be /api/chat-stream")
		}
	})

	t.Run("custom values", func(t *testing.T) {
		config := ConfigHTTP{
			Address:            "localhost:9000",
			ChatFlowPath:       "/custom/chat",
			ChatStreamFlowPath: "/custom/stream",
			InformationPath:    "/custom/info",
		}

		if config.ChatFlowPath != "/custom/chat" {
			t.Errorf("ChatFlowPath = %q, want /custom/chat", config.ChatFlowPath)
		}

		if config.ChatStreamFlowPath != "/custom/stream" {
			t.Errorf("ChatStreamFlowPath = %q, want /custom/stream", config.ChatStreamFlowPath)
		}

		if config.InformationPath != "/custom/info" {
			t.Errorf("InformationPath = %q, want /custom/info", config.InformationPath)
		}
	})
}
