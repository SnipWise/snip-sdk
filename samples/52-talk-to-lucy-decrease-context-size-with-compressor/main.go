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
	"github.com/snipwise/snip-sdk/snip/compressor"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/lucy-gguf:q4_k_m")
	compressorModelId := env.GetEnvOrDefault("COMPRESSOR_MODEL", "hf.co/menlo/jan-nano-128k-gguf:q4_k_m")

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

	// Create the main chat agent
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

	// Create a compressor agent for context compression
	compressor, err := compressor.NewCompressorAgent(
		ctx,
		agents.AgentConfig{
			Name:               "MessageCompressor",
			SystemInstructions: "",
			ModelID:            compressorModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.3, // Lower temperature for more consistent compression
		},
	)
	if err != nil {
		fmt.Printf("Error creating compressor agent: %v\n", err)
		return
	}

	// First question
	fmt.Println("=== First Question ===")
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
	fmt.Printf("📏 Context Size after adding knowledge: %d\n", agent0.GetCurrentContextSize())

	// Second question - with knowledge base
	fmt.Println("\n=== Second Question (with knowledge base) ===")
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
	fmt.Printf("📏 Context Size: %d\n", agent0.GetCurrentContextSize())

	// Third question
	fmt.Println("\n=== Third Question ===")
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

	// Display current context state
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("📏 Current Context Size:", agent0.GetCurrentContextSize())
	fmt.Printf("📝 Number of messages: %d\n", len(agent0.GetMessages()))
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

	// NOW: Compress the conversation to reduce context size
	fmt.Println("\n🗜️ === Compressing conversation with CompressorAgent ===")

	messages := agent0.GetMessages()
	fmt.Printf("Original messages count: %d\n", len(messages))
	fmt.Printf("Original context size: %d\n", agent0.GetCurrentContextSize())

	// Compress the conversation
	fmt.Println("\nCompressing conversation (streaming)...")
	fmt.Println(strings.Repeat("-", 80))

	compressed, err := compressor.CompressMessagesStream(messages, func(chunk agents.ChatResponse) error {
		fmt.Print(chunk.Text)
		return nil
	})
	if err != nil {
		fmt.Printf("\n❌ Error compressing messages: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("\n✅ Compression complete!\n")
	fmt.Printf("Compressed size: %d characters\n", len(compressed.Text))

	// Replace the agent's messages with compressed context
	fmt.Println("\n--- Replacing messages with compressed context ---")

	// Create new messages: original system instructions + compressed conversation
	newMessages := []*ai.Message{
		ai.NewSystemTextMessage(systemInstructions),
		ai.NewSystemTextMessage("Previous conversation summary:\n" + compressed.Text),
	}

	err = agent0.ReplaceMessagesWith(newMessages)
	if err != nil {
		fmt.Printf("Error replacing messages: %v\n", err)
		return
	}

	fmt.Println("\n--- After compression and replacement ---")
	fmt.Println(strings.Repeat("+", 50))
	fmt.Printf("📏 New Context Size: %d\n", agent0.GetCurrentContextSize())
	fmt.Printf("📝 Number of messages: %d\n", len(agent0.GetMessages()))
	for i, msg := range agent0.GetMessages() {
		// Get first 100 characters for compressed message
		messageText := ""
		if len(msg.Content) > 0 && len(msg.Content[0].Text) > 0 {
			if len(msg.Content[0].Text) > 100 {
				messageText = msg.Content[0].Text[:100] + "..."
			} else {
				messageText = msg.Content[0].Text
			}
		}
		fmt.Printf("Message %d: Role=%s, Content=%q\n", i, msg.Role, messageText)
	}
	fmt.Println(strings.Repeat("+", 50))

	// Test a new question with the compressed context
	fmt.Println("\n=== Testing with compressed context ===")
	answer, err = agent0.AskStreamWithMemory("Can you remind me what we discussed about Hawaiian pizza?",
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
	fmt.Printf("📏 Final Context Size: %d\n", agent0.GetCurrentContextSize())

	fmt.Println("\n✨ Demo completed successfully!")
	fmt.Println("The compressor agent successfully reduced the context size while preserving important information.")

}
