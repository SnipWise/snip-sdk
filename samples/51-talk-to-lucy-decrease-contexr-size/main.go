package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/files"

	"github.com/snipwise/snip-sdk/snip"
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

	agent0, err := snip.NewAgent(ctx,
		snip.AgentConfig{
			Name:               "Bob_Agentic_Agent",
			SystemInstructions: systemInstructions,
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		snip.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		snip.EnableChatStreamFlowWithMemory(),
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	answer, err := agent0.AskStream("What is the best pizza of the world?",
		func(chunk snip.ChatResponse) error {
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

	answer, err = agent0.AskStream("Who invented Hawaiian pizza?",
		func(chunk snip.ChatResponse) error {
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

	answer, err = agent0.AskStream("What is Hawaiian pizza?",
		func(chunk snip.ChatResponse) error {
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
	answer, err = agent0.AskStream("Tell me about pizza.",
		func(chunk snip.ChatResponse) error {
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
