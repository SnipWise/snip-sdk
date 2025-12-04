package smart

import (
	"context"
	"testing"

	"github.com/firebase/genkit/go/ai"
)

// ============================================================================
// Tests for Agent.GetName
// ============================================================================

func TestAgentGetName(t *testing.T) {
	agent := &Agent{
		Name: "test-agent",
	}

	name := agent.GetName()
	if name != "test-agent" {
		t.Errorf("GetName() = %q, want %q", name, "test-agent")
	}
}

// ============================================================================
// Tests for Agent.Kind
// ============================================================================

func TestAgentKind(t *testing.T) {
	agent := &Agent{}

	kind := agent.Kind()
	if kind != Basic {
		t.Errorf("Kind() = %v, want %v", kind, Basic)
	}
}

// ============================================================================
// Tests for Agent.GetMessages
// ============================================================================

func TestAgentGetMessages(t *testing.T) {
	t.Run("empty messages", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		messages := agent.GetMessages()
		if len(messages) != 0 {
			t.Errorf("GetMessages() length = %d, want 0", len(messages))
		}
	})

	t.Run("with messages", func(t *testing.T) {
		msg1 := ai.NewUserTextMessage("Hello")
		msg2 := ai.NewModelTextMessage("Hi there")

		agent := &Agent{
			Messages: []*ai.Message{msg1, msg2},
		}

		messages := agent.GetMessages()
		if len(messages) != 2 {
			t.Errorf("GetMessages() length = %d, want 2", len(messages))
		}

		if messages[0] != msg1 {
			t.Error("First message doesn't match")
		}
		if messages[1] != msg2 {
			t.Error("Second message doesn't match")
		}
	})
}

// ============================================================================
// Tests for Agent.AddSystemMessage
// ============================================================================

func TestAgentAddSystemMessage(t *testing.T) {
	t.Run("add single message", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		err := agent.AddSystemMessage("You are a helpful assistant")
		if err != nil {
			t.Errorf("AddSystemMessage() unexpected error: %v", err)
		}

		if len(agent.Messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(agent.Messages))
		}

		msg := agent.Messages[0]
		if msg.Role != ai.RoleSystem {
			t.Errorf("Message role = %v, want %v", msg.Role, ai.RoleSystem)
		}

		if len(msg.Content) == 0 {
			t.Fatal("Message content is empty")
		}

		content := msg.Content[0].Text
		if content != "You are a helpful assistant" {
			t.Errorf("Message content = %q, want %q", content, "You are a helpful assistant")
		}
	})

	t.Run("add multiple messages", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		err := agent.AddSystemMessage("First message")
		if err != nil {
			t.Errorf("AddSystemMessage() unexpected error: %v", err)
		}

		err = agent.AddSystemMessage("Second message")
		if err != nil {
			t.Errorf("AddSystemMessage() unexpected error: %v", err)
		}

		if len(agent.Messages) != 2 {
			t.Fatalf("Expected 2 messages, got %d", len(agent.Messages))
		}

		if agent.Messages[0].Content[0].Text != "First message" {
			t.Error("First message content mismatch")
		}
		if agent.Messages[1].Content[0].Text != "Second message" {
			t.Error("Second message content mismatch")
		}
	})

	t.Run("trim whitespace", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		err := agent.AddSystemMessage("  \n  Message with whitespace  \n  ")
		if err != nil {
			t.Errorf("AddSystemMessage() unexpected error: %v", err)
		}

		if len(agent.Messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(agent.Messages))
		}

		content := agent.Messages[0].Content[0].Text
		if content != "Message with whitespace" {
			t.Errorf("Message content = %q, want %q", content, "Message with whitespace")
		}
	})

	t.Run("empty message after trim", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		err := agent.AddSystemMessage("   \n   ")
		if err != nil {
			t.Errorf("AddSystemMessage() unexpected error: %v", err)
		}

		if len(agent.Messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(agent.Messages))
		}

		content := agent.Messages[0].Content[0].Text
		if content != "" {
			t.Errorf("Message content = %q, want empty string", content)
		}
	})

	t.Run("preserve existing messages", func(t *testing.T) {
		existingMsg := ai.NewUserTextMessage("Existing")
		agent := &Agent{
			Messages: []*ai.Message{existingMsg},
		}

		err := agent.AddSystemMessage("New system message")
		if err != nil {
			t.Errorf("AddSystemMessage() unexpected error: %v", err)
		}

		if len(agent.Messages) != 2 {
			t.Fatalf("Expected 2 messages, got %d", len(agent.Messages))
		}

		if agent.Messages[0] != existingMsg {
			t.Error("Existing message was modified or replaced")
		}

		if agent.Messages[1].Content[0].Text != "New system message" {
			t.Error("New message not correctly added")
		}
	})
}

// ============================================================================
// Tests for Agent.GetInfo
// ============================================================================

func TestAgentGetInfo(t *testing.T) {
	t.Run("basic info", func(t *testing.T) {
		agent := &Agent{
			Name:    "test-agent",
			ModelID: "test-model",
			Config: ModelConfig{
				Temperature: 0.7,
				TopP:        0.9,
				MaxTokens:   1000,
			},
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

		if info.Config.TopP != 0.9 {
			t.Errorf("Info.Config.TopP = %f, want 0.9", info.Config.TopP)
		}

		if info.Config.MaxTokens != 1000 {
			t.Errorf("Info.Config.MaxTokens = %d, want 1000", info.Config.MaxTokens)
		}
	})

	t.Run("empty config", func(t *testing.T) {
		agent := &Agent{
			Name:    "minimal-agent",
			ModelID: "minimal-model",
			Config:  ModelConfig{},
		}

		info, err := agent.GetInfo()
		if err != nil {
			t.Errorf("GetInfo() unexpected error: %v", err)
		}

		if info.Name != "minimal-agent" {
			t.Error("Name not correctly returned")
		}

		if info.ModelID != "minimal-model" {
			t.Error("ModelID not correctly returned")
		}

		// Zero values should be present
		if info.Config.Temperature != 0 {
			t.Errorf("Config.Temperature = %f, want 0", info.Config.Temperature)
		}
	})
}

// ============================================================================
// Tests for Agent.Ask (without actual API calls)
// ============================================================================

func TestAgentAsk(t *testing.T) {
	t.Run("chatFlow not initialized", func(t *testing.T) {
		ctx := context.Background()
		agent := &Agent{
			ctx:      ctx,
			chatFlow: nil,
		}

		_, err := agent.Ask("Hello")
		if err == nil {
			t.Error("Ask() expected error when chatFlow is nil, got nil")
		}

		expectedMsg := "chat flow is not initialized"
		if err.Error() != expectedMsg {
			t.Errorf("Ask() error message = %q, want %q", err.Error(), expectedMsg)
		}
	})
}

// ============================================================================
// Tests for Agent.AskStream (without actual API calls)
// ============================================================================

func TestAgentAskStream(t *testing.T) {
	t.Run("chatStreamFlow not initialized", func(t *testing.T) {
		ctx := context.Background()
		agent := &Agent{
			ctx:            ctx,
			chatStreamFlow: nil,
		}

		callback := func(chunk string) error {
			return nil
		}

		_, err := agent.AskStream("Hello", callback)
		if err == nil {
			t.Error("AskStream() expected error when chatStreamFlow is nil, got nil")
		}

		expectedMsg := "chat stream flow is not initialized"
		if err.Error() != expectedMsg {
			t.Errorf("AskStream() error message = %q, want %q", err.Error(), expectedMsg)
		}
	})
}

// ============================================================================
// Tests for AgentKind type
// ============================================================================

func TestAgentKindValues(t *testing.T) {
	tests := []struct {
		name     string
		kind     AgentKind
		expected string
	}{
		{"Basic kind", Basic, "Basic"},
		{"Remote kind", Remote, "Remote"},
		{"Tool kind", Tool, "Tool"},
		{"Intent kind", Intent, "Intent"},
		{"Rag kind", Rag, "Rag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.kind) != tt.expected {
				t.Errorf("AgentKind = %q, want %q", tt.kind, tt.expected)
			}
		})
	}
}
