package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/toolbox/files"

	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/chat"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/lucy-gguf:q4_k_m")
	//chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:latest")

	systemInstructions, err := files.ReadTextFile(env.GetEnvOrDefault("SYSTEM_INSTRUCTION_PATH", "./system-instructions.md"))
	if err != nil {
		fmt.Printf("Error reading system instructions: %v\n", err)
		return
	}
	knowledgeBase, err := files.ReadTextFile(env.GetEnvOrDefault("KNOWLEDGE_BASE_PATH", "./knowledge-base.md"))
	if err != nil {
		fmt.Printf("Error reading knowledge base: %v\n", err)
		return
	}

	agent0, err := chat.NewChatAgent(ctx,
		agents.AgentConfig{
			Name:               "Bob_Agentic_Agent",
			SystemInstructions: systemInstructions,
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		chat.EnableChatStreamFlowWithMemory(),
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	answer, err := agent0.AskStreamWithMemory("What is the best pizza of the world?",
		func(chunk agents.ChatResponse) error {
			fmt.Print(chunk.Text)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}

	fmt.Println("\n✋ FinishReason:", answer.FinishReason)
	if answer.IsFinishReasonLength() {
		fmt.Println("⚠️ The answer was cut off due to length limits.")
	}
	if answer.IsFinishReasonStop() {
		fmt.Println("✅ The answer was completed successfully.")
	}
	fmt.Println(strings.Repeat("*", 50))
	fmt.Println("📏 Current Context Size:", agent0.GetCurrentContextSize())
	fmt.Println(strings.Repeat("*", 50))


	fmt.Println("\n--- Adding knowledge base context ---")

	// Log context size before adding
	fmt.Printf("Knowledge base size: %d characters\n", len(knowledgeBase))
	fmt.Printf("Current messages in history: %d\n", len(agent0.GetMessages()))

	err = agent0.AddSystemMessage(knowledgeBase)
	if err != nil {
		fmt.Printf("Error adding system message: %v\n", err)
		return
	}

	fmt.Printf("Messages after adding context: %d\n", len(agent0.GetMessages()))

	answer, err = agent0.AskStreamWithMemory("Who invented Hawaiian pizza?",
		func(chunk agents.ChatResponse) error {
			fmt.Print(chunk.Text)
			return nil
		},
	)
	fmt.Println("\n✋ FinishReason:", answer.FinishReason)
	if answer.IsFinishReasonLength() {
		fmt.Println("⚠️ The answer was cut off due to length limits.")
	}
	if answer.IsFinishReasonStop() {
		fmt.Println("✅ The answer was completed successfully.")
	}

	if err != nil {
		fmt.Printf("\n❌ Error asking question: %v\n", err)
		fmt.Printf("Partial answer received: %q\n", answer)
		return
	}
	fmt.Println()
	fmt.Println("\n--- Add again knowledge base context ---")

	answer, err = agent0.AskStreamWithMemory("What is Hawaiian pizza?",
		func(chunk agents.ChatResponse) error {
			fmt.Print(chunk.Text)
			return nil
		},
	)
	fmt.Println("\n✋ FinishReason:", answer.FinishReason)
	if answer.IsFinishReasonLength() {
		fmt.Println("⚠️ The answer was cut off due to length limits.")
	}
	if answer.IsFinishReasonStop() {
		fmt.Println("✅ The answer was completed successfully.")
	}

	if err != nil {
		fmt.Printf("\n❌ Error asking question: %v\n", err)
		fmt.Printf("Partial answer received: %q\n", answer)
		return
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Current Context Size:", agent0.GetCurrentContextSize())
	for i, msg := range agent0.GetMessages() {
		// Get first 50 characters of the message content
		messageText := ""
		if len(msg.Content) > 0 && len(msg.Content[0].Text) > 0 {
			if len(msg.Content[0].Text) > 50 {
				messageText = msg.Content[0].Text[:50] + "..."
			} else {
				messageText = msg.Content[0].Text
			}
		}
		fmt.Printf("Message %d: Role=%s, Content=%q\n", i, msg.Role, messageText)
	}
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("\n--- Replacing messages with new system messages ---")

	// Create new messages: 2 system messages
	newMessages := []*ai.Message{
		ai.NewSystemTextMessage("You are a helpful assistant specialized in Italian cuisine."),
		ai.NewSystemTextMessage("You should always mention that pizza is best enjoyed fresh from a wood-fired oven."),
	}

	err = agent0.ReplaceMessagesWith(newMessages)
	if err != nil {
		fmt.Printf("Error replacing messages: %v\n", err)
		return
	}

	fmt.Println("\n--- After replacing messages ---")
	fmt.Println(strings.Repeat("+", 50))
	fmt.Println("Context Size:", agent0.GetCurrentContextSize())
	for i, msg := range agent0.GetMessages() {
		// Get first 50 characters of the message content
		messageText := ""
		if len(msg.Content) > 0 && len(msg.Content[0].Text) > 0 {
			if len(msg.Content[0].Text) > 50 {
				messageText = msg.Content[0].Text[:50] + "..."
			} else {
				messageText = msg.Content[0].Text
			}
		}
		fmt.Printf("Message %d: Role=%s, Content=%q\n", i, msg.Role, messageText)
	}
	fmt.Println(strings.Repeat("+", 50))

	// Test a new question with the replaced context
	fmt.Println("\n--- Asking a question with new context ---")
	answer, err = agent0.AskStreamWithMemory("Tell me about pizza.",
		func(chunk agents.ChatResponse) error {
			fmt.Print(chunk.Text)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("\n❌ Error asking question: %v\n", err)
		return
	}
	fmt.Println("\n✋ FinishReason:", answer.FinishReason)

	// Demonstrate ReplaceMessagesWithSystemMessages
	fmt.Println("\n--- Using ReplaceMessagesWithSystemMessages ---")

	systemMessages := []string{
		"You are a helpful assistant specialized in French cuisine.",
		"You should always emphasize the importance of using fresh, local ingredients.",
		"You are passionate about traditional cooking techniques.",
	}

	err = agent0.ReplaceMessagesWithSystemMessages(systemMessages)
	if err != nil {
		fmt.Printf("Error replacing messages with system messages: %v\n", err)
		return
	}

	fmt.Println("\n--- After ReplaceMessagesWithSystemMessages ---")
	fmt.Println(strings.Repeat("~", 50))
	fmt.Println("Context Size:", agent0.GetCurrentContextSize())
	fmt.Printf("Number of messages: %d\n", len(agent0.GetMessages()))
	for i, msg := range agent0.GetMessages() {
		messageText := ""
		if len(msg.Content) > 0 && len(msg.Content[0].Text) > 0 {
			messageText = msg.Content[0].Text
		}
		fmt.Printf("Message %d: Role=%s, Content=%q\n", i, msg.Role, messageText)
	}
	fmt.Println(strings.Repeat("~", 50))

	// Test with the new system messages
	fmt.Println("\n--- Asking about French cuisine ---")
	answer, err = agent0.AskStreamWithMemory("What is the best French dish?",
		func(chunk agents.ChatResponse) error {
			fmt.Print(chunk.Text)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("\n❌ Error asking question: %v\n", err)
		return
	}
	fmt.Println("\n✋ FinishReason:", answer.FinishReason)

}
