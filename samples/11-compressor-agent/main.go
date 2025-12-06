package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/compressor"
)

func main() {
	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	// Create a compressor agent
	compressor, err := compressor.NewCompressorAgent(
		ctx,
		agents.AgentConfig{
			Name:               "MessageCompressor",
			SystemInstructions: "",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.3, // Lower temperature for more consistent compression
		},
	)
	if err != nil {
		log.Fatalf("Error creating compressor agent: %v", err)
	}

	// Example: Compress conversation messages
	fmt.Println("=== Example: Compressing Conversation Messages ===")
	fmt.Println()

	// Create a sample conversation with ai.Message objects
	messages := []*ai.Message{
		ai.NewUserTextMessage("Hello, I'm working on a Go project and I need to implement a feature that handles user authentication."),
		ai.NewModelTextMessage("That's great! For user authentication in Go, you have several options. You could use JWT tokens, session-based authentication, or OAuth2. What kind of application are you building?"),
		ai.NewUserTextMessage("It's a REST API for a mobile app."),
		ai.NewModelTextMessage("Perfect! For a REST API serving a mobile app, I'd recommend using JWT (JSON Web Tokens). Here's why: JWTs are stateless, which means you don't need to store session data on the server. They're also great for mobile apps because you can store them securely on the device. Would you like me to show you how to implement JWT authentication in Go?"),
		ai.NewUserTextMessage("Yes, please!"),
		ai.NewModelTextMessage("Great! I'll show you a basic implementation. First, you'll need to install the jwt-go package. You can do this by running: go get github.com/golang-jwt/jwt/v5. Then, you'll need to create a middleware function that validates the JWT token on each request. You'll also need functions to generate tokens when users log in and to validate them on protected routes."),
		ai.NewUserTextMessage("That sounds good. What about password hashing?"),
		ai.NewModelTextMessage("Excellent question! You should never store passwords in plain text. Go's crypto/bcrypt package is perfect for this. It automatically handles salting and uses a strong hashing algorithm. You can use bcrypt.GenerateFromPassword() to hash passwords when users register, and bcrypt.CompareHashAndPassword() to verify passwords during login."),
	}

	// Calculate original size
	originalSize := 0
	for _, msg := range messages {
		for _, content := range msg.Content {
			originalSize += len(content.Text)
		}
	}

	fmt.Println("Original conversation:")
	fmt.Println(strings.Repeat("-", 80))
	for i, msg := range messages {
		text := msg.Content[0].Text
		if len(text) > 100 {
			fmt.Printf("%d. [%s]: %s...\n", i+1, msg.Role, text[:100])
		} else {
			fmt.Printf("%d. [%s]: %s\n", i+1, msg.Role, text)
		}
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()

	// Example 1: Non-streaming compression
	fmt.Println("=== Non-streaming Compression ===")
	fmt.Println()

	compressed, err := compressor.CompressMessages(messages)
	if err != nil {
		log.Fatalf("Error compressing messages: %v", err)
	}

	fmt.Println("Compressed conversation:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(compressed.Text)
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()
	fmt.Println("Original size:", originalSize, "characters")
	fmt.Println("Compressed size:", len(compressed.Text), "characters")
	fmt.Println("Compression ratio:", fmt.Sprintf("%.1f%%", (1.0-float64(len(compressed.Text))/float64(originalSize))*100))
	fmt.Println()

	// Example 2: Streaming compression
	fmt.Println("=== Streaming Compression ===")
	fmt.Println()
	fmt.Println("Watch the compression happen in real-time:")
	fmt.Println(strings.Repeat("-", 80))

	streamCompressed, err := compressor.CompressMessagesStream(messages, func(chunk agents.ChatResponse) error {
		fmt.Print(chunk.Text)
		return nil
	})
	if err != nil {
		log.Fatalf("Error compressing messages with streaming: %v", err)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()
	fmt.Println("Original size:", originalSize, "characters")
	fmt.Println("Compressed size:", len(streamCompressed.Text), "characters")
	fmt.Println("Compression ratio:", fmt.Sprintf("%.1f%%", (1.0-float64(len(streamCompressed.Text))/float64(originalSize))*100))
	fmt.Println()
	fmt.Println("✅ Message compression complete!")
}
