package snip

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
// Tests for Agent.GetCurrentContextSize
// ============================================================================

func TestAgentGetCurrentContextSize(t *testing.T) {
	t.Run("empty messages", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		size := agent.GetCurrentContextSize()
		if size != 0 {
			t.Errorf("GetCurrentContextSize() = %d, want 0", size)
		}
	})

	t.Run("with messages", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{
				ai.NewUserTextMessage("Hello"),
				ai.NewModelTextMessage("Hi there"),
				ai.NewUserTextMessage("How are you?"),
			},
		}

		// Expected size = len("Hello") + len("Hi there") + len("How are you?")
		// = 5 + 8 + 12 = 25
		size := agent.GetCurrentContextSize()
		expectedSize := 25
		if size != expectedSize {
			t.Errorf("GetCurrentContextSize() = %d, want %d", size, expectedSize)
		}
	})

	t.Run("nil messages slice", func(t *testing.T) {
		agent := &Agent{
			Messages: nil,
		}

		size := agent.GetCurrentContextSize()
		if size != 0 {
			t.Errorf("GetCurrentContextSize() = %d, want 0", size)
		}
	})

	t.Run("after adding messages", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		// Initial size should be 0
		if agent.GetCurrentContextSize() != 0 {
			t.Errorf("Initial GetCurrentContextSize() = %d, want 0", agent.GetCurrentContextSize())
		}

		// Add a message
		agent.AddSystemMessage("System message")
		// Expected size = len("System message") = 14
		expectedSize := 14
		if agent.GetCurrentContextSize() != expectedSize {
			t.Errorf("After AddSystemMessage, GetCurrentContextSize() = %d, want %d", agent.GetCurrentContextSize(), expectedSize)
		}

		// Add more messages
		agent.Messages = append(agent.Messages, ai.NewUserTextMessage("User message"))
		agent.Messages = append(agent.Messages, ai.NewModelTextMessage("Model response"))

		// Expected size = len("System message") + len("User message") + len("Model response")
		// = 14 + 12 + 14 = 40
		expectedSize = 40
		if agent.GetCurrentContextSize() != expectedSize {
			t.Errorf("After adding 3 messages, GetCurrentContextSize() = %d, want %d", agent.GetCurrentContextSize(), expectedSize)
		}
	})
}

// ============================================================================
// Tests for Agent.ReplaceMessagesWith
// ============================================================================

func TestAgentReplaceMessages(t *testing.T) {
	t.Run("replace empty messages with new messages", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		newMessages := []*ai.Message{
			ai.NewUserTextMessage("Hello"),
			ai.NewModelTextMessage("Hi there"),
		}

		err := agent.ReplaceMessagesWith(newMessages)
		if err != nil {
			t.Errorf("ReplaceMessagesWith() unexpected error: %v", err)
		}

		if len(agent.Messages) != 2 {
			t.Fatalf("Expected 2 messages, got %d", len(agent.Messages))
		}

		if agent.Messages[0] != newMessages[0] {
			t.Error("First message doesn't match")
		}
		if agent.Messages[1] != newMessages[1] {
			t.Error("Second message doesn't match")
		}
	})

	t.Run("replace existing messages", func(t *testing.T) {
		oldMessages := []*ai.Message{
			ai.NewUserTextMessage("Old message 1"),
			ai.NewUserTextMessage("Old message 2"),
			ai.NewUserTextMessage("Old message 3"),
		}

		agent := &Agent{
			Messages: oldMessages,
		}

		newMessages := []*ai.Message{
			ai.NewSystemTextMessage("System instruction"),
			ai.NewUserTextMessage("New user message"),
		}

		err := agent.ReplaceMessagesWith(newMessages)
		if err != nil {
			t.Errorf("ReplaceMessagesWith() unexpected error: %v", err)
		}

		if len(agent.Messages) != 2 {
			t.Fatalf("Expected 2 messages, got %d", len(agent.Messages))
		}

		if agent.Messages[0].Role != ai.RoleSystem {
			t.Errorf("First message role = %v, want %v", agent.Messages[0].Role, ai.RoleSystem)
		}

		if agent.Messages[1].Role != ai.RoleUser {
			t.Errorf("Second message role = %v, want %v", agent.Messages[1].Role, ai.RoleUser)
		}

		// Verify old messages are completely replaced
		for _, oldMsg := range oldMessages {
			found := false
			for _, currentMsg := range agent.Messages {
				if currentMsg == oldMsg {
					found = true
					break
				}
			}
			if found {
				t.Error("Old message still present after replacement")
			}
		}
	})

	t.Run("replace with empty slice", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{
				ai.NewUserTextMessage("Message 1"),
				ai.NewUserTextMessage("Message 2"),
			},
		}

		newMessages := []*ai.Message{}

		err := agent.ReplaceMessagesWith(newMessages)
		if err != nil {
			t.Errorf("ReplaceMessagesWith() unexpected error: %v", err)
		}

		if len(agent.Messages) != 0 {
			t.Fatalf("Expected 0 messages, got %d", len(agent.Messages))
		}
	})

	t.Run("nil messages returns error", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{
				ai.NewUserTextMessage("Existing message"),
			},
		}

		err := agent.ReplaceMessagesWith(nil)
		if err == nil {
			t.Error("ReplaceMessagesWith(nil) expected error, got nil")
		}

		expectedMsg := "messages cannot be nil"
		if err.Error() != expectedMsg {
			t.Errorf("ReplaceMessagesWith(nil) error = %q, want %q", err.Error(), expectedMsg)
		}

		// Verify existing messages are not modified when error occurs
		if len(agent.Messages) != 1 {
			t.Errorf("Messages were modified despite error, length = %d, want 1", len(agent.Messages))
		}
	})

	t.Run("replace preserves message references", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		msg1 := ai.NewUserTextMessage("Message 1")
		msg2 := ai.NewModelTextMessage("Message 2")
		newMessages := []*ai.Message{msg1, msg2}

		err := agent.ReplaceMessagesWith(newMessages)
		if err != nil {
			t.Errorf("ReplaceMessagesWith() unexpected error: %v", err)
		}

		// Verify the same slice is used (reference equality)
		if &agent.Messages[0] != &newMessages[0] {
			t.Error("Messages slice was copied instead of being assigned directly")
		}
	})
}

// ============================================================================
// Tests for Agent.ReplaceMessagesWithSystemMessages
// ============================================================================

func TestAgentReplaceMessagesWithSystemMessages(t *testing.T) {
	t.Run("replace empty messages with system messages", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		systemMessages := []string{
			"You are a helpful assistant",
			"Always be polite and concise",
		}

		err := agent.ReplaceMessagesWithSystemMessages(systemMessages)
		if err != nil {
			t.Errorf("ReplaceMessagesWithSystemMessages() unexpected error: %v", err)
		}

		if len(agent.Messages) != 2 {
			t.Fatalf("Expected 2 messages, got %d", len(agent.Messages))
		}

		// Verify all messages are system messages
		for i, msg := range agent.Messages {
			if msg.Role != ai.RoleSystem {
				t.Errorf("Message %d role = %v, want %v", i, msg.Role, ai.RoleSystem)
			}
			if msg.Content[0].Text != systemMessages[i] {
				t.Errorf("Message %d content = %q, want %q", i, msg.Content[0].Text, systemMessages[i])
			}
		}
	})

	t.Run("replace existing mixed messages with system messages", func(t *testing.T) {
		oldMessages := []*ai.Message{
			ai.NewUserTextMessage("User message"),
			ai.NewModelTextMessage("Model message"),
			ai.NewSystemTextMessage("Old system message"),
		}

		agent := &Agent{
			Messages: oldMessages,
		}

		systemMessages := []string{
			"System instruction 1",
			"System instruction 2",
		}

		err := agent.ReplaceMessagesWithSystemMessages(systemMessages)
		if err != nil {
			t.Errorf("ReplaceMessagesWithSystemMessages() unexpected error: %v", err)
		}

		if len(agent.Messages) != 2 {
			t.Fatalf("Expected 2 messages, got %d", len(agent.Messages))
		}

		// Verify all messages are system messages with correct content
		for i, msg := range agent.Messages {
			if msg.Role != ai.RoleSystem {
				t.Errorf("Message %d role = %v, want %v", i, msg.Role, ai.RoleSystem)
			}
			if msg.Content[0].Text != systemMessages[i] {
				t.Errorf("Message %d content = %q, want %q", i, msg.Content[0].Text, systemMessages[i])
			}
		}

		// Verify old messages are completely replaced
		for _, oldMsg := range oldMessages {
			found := false
			for _, currentMsg := range agent.Messages {
				if currentMsg == oldMsg {
					found = true
					break
				}
			}
			if found {
				t.Error("Old message still present after replacement")
			}
		}
	})

	t.Run("trim whitespace from system messages", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		systemMessages := []string{
			"  Message with leading spaces",
			"Message with trailing spaces  ",
			"  Message with both  ",
		}

		err := agent.ReplaceMessagesWithSystemMessages(systemMessages)
		if err != nil {
			t.Errorf("ReplaceMessagesWithSystemMessages() unexpected error: %v", err)
		}

		expectedMessages := []string{
			"Message with leading spaces",
			"Message with trailing spaces",
			"Message with both",
		}

		for i, msg := range agent.Messages {
			if msg.Content[0].Text != expectedMessages[i] {
				t.Errorf("Message %d content = %q, want %q", i, msg.Content[0].Text, expectedMessages[i])
			}
		}
	})

	t.Run("replace with empty slice", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{
				ai.NewUserTextMessage("Message 1"),
				ai.NewUserTextMessage("Message 2"),
			},
		}

		systemMessages := []string{}

		err := agent.ReplaceMessagesWithSystemMessages(systemMessages)
		if err != nil {
			t.Errorf("ReplaceMessagesWithSystemMessages() unexpected error: %v", err)
		}

		if len(agent.Messages) != 0 {
			t.Fatalf("Expected 0 messages, got %d", len(agent.Messages))
		}
	})

	t.Run("nil system messages returns error", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{
				ai.NewUserTextMessage("Existing message"),
			},
		}

		err := agent.ReplaceMessagesWithSystemMessages(nil)
		if err == nil {
			t.Error("ReplaceMessagesWithSystemMessages(nil) expected error, got nil")
		}

		expectedMsg := "systemMessages cannot be nil"
		if err.Error() != expectedMsg {
			t.Errorf("ReplaceMessagesWithSystemMessages(nil) error = %q, want %q", err.Error(), expectedMsg)
		}

		// Verify existing messages are not modified when error occurs
		if len(agent.Messages) != 1 {
			t.Errorf("Messages were modified despite error, length = %d, want 1", len(agent.Messages))
		}
	})

	t.Run("single system message", func(t *testing.T) {
		agent := &Agent{
			Messages: []*ai.Message{},
		}

		systemMessages := []string{"You are a helpful assistant"}

		err := agent.ReplaceMessagesWithSystemMessages(systemMessages)
		if err != nil {
			t.Errorf("ReplaceMessagesWithSystemMessages() unexpected error: %v", err)
		}

		if len(agent.Messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(agent.Messages))
		}

		if agent.Messages[0].Role != ai.RoleSystem {
			t.Errorf("Message role = %v, want %v", agent.Messages[0].Role, ai.RoleSystem)
		}

		if agent.Messages[0].Content[0].Text != "You are a helpful assistant" {
			t.Errorf("Message content = %q, want %q", agent.Messages[0].Content[0].Text, "You are a helpful assistant")
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
			ctx:                ctx,
			chatFlowWithMemory: nil,
		}

		_, err := agent.AskWithMemory("Hello")
		if err == nil {
			t.Error("AskWithMemory() expected error when chatFlow is nil, got nil")
		}

		expectedMsg := "chat flow is not initialized"
		if err.Error() != expectedMsg {
			t.Errorf("AskWithMemory() error message = %q, want %q", err.Error(), expectedMsg)
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
			ctx:                      ctx,
			chatStreamFlowWithMemory: nil,
		}

		callback := func(chunk ChatResponse) error {
			return nil
		}

		_, err := agent.AskStreamWithMemory("Hello", callback)
		if err == nil {
			t.Error("AskStreamWithMemory() expected error when chatStreamFlow is nil, got nil")
		}

		expectedMsg := "chat stream flow is not initialized"
		if err.Error() != expectedMsg {
			t.Errorf("AskStreamWithMemory() error message = %q, want %q", err.Error(), expectedMsg)
		}
	})
}

// ============================================================================
// Tests for Agent.Ask (without memory)
// ============================================================================

func TestAgentAskWithoutMemory(t *testing.T) {
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
// Tests for Agent.AskStream (without memory)
// ============================================================================

func TestAgentAskStreamWithoutMemory(t *testing.T) {
	t.Run("chatStreamFlow not initialized", func(t *testing.T) {
		ctx := context.Background()
		agent := &Agent{
			ctx:            ctx,
			chatStreamFlow: nil,
		}

		callback := func(chunk ChatResponse) error {
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
