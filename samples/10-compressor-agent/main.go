package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/snip"
)

func main() {
	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-128k-gguf:q4_k_m")

	// Create a compressor agent
	compressor, err := snip.NewCompressorAgent(
		ctx,
		snip.AgentConfig{
			Name:               "TextCompressor",
			SystemInstructions: "",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		snip.ModelConfig{
			Temperature: 0.3, // Lower temperature for more consistent compression
		},
	)
	if err != nil {
		log.Fatalf("Error creating compressor agent: %v", err)
	}

	// Example 1: Compress a long text (non-streaming)
	fmt.Println("=== Example 1: Non-streaming Compression ===")
	fmt.Println()

	longText := `
	User: Hello, I'm working on a Go project and I need to implement a feature that handles user authentication.
	Agent: That's great! For user authentication in Go, you have several options. You could use JWT tokens, session-based authentication, or OAuth2. What kind of application are you building?
	User: It's a REST API for a mobile app.
	Agent: Perfect! For a REST API serving a mobile app, I'd recommend using JWT (JSON Web Tokens). Here's why: JWTs are stateless, which means you don't need to store session data on the server. They're also great for mobile apps because you can store them securely on the device. Would you like me to show you how to implement JWT authentication in Go?
	User: Yes, please!
	Agent: Great! I'll show you a basic implementation. First, you'll need to install the jwt-go package. You can do this by running: go get github.com/golang-jwt/jwt/v5. Then, you'll need to create a middleware function that validates the JWT token on each request. You'll also need functions to generate tokens when users log in and to validate them on protected routes.
	User: That sounds good. What about password hashing?
	Agent: Excellent question! You should never store passwords in plain text. Go's crypto/bcrypt package is perfect for this. It automatically handles salting and uses a strong hashing algorithm. You can use bcrypt.GenerateFromPassword() to hash passwords when users register, and bcrypt.CompareHashAndPassword() to verify passwords during login.
	`

	compressed, err := compressor.CompressText(longText)
	if err != nil {
		log.Fatalf("Error compressing text: %v", err)
	}

	fmt.Println("Original text length:", len(longText), "characters")
	fmt.Println("Compressed text length:", len(compressed.Text), "characters")
	fmt.Println("Compression ratio:", fmt.Sprintf("%.1f%%", (1.0-float64(len(compressed.Text))/float64(len(longText)))*100))
	fmt.Println()
	fmt.Println("Compressed text:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(compressed.Text)
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()

	// Example 2: Compress text with streaming (real-time feedback)
	fmt.Println("=== Example 2: Streaming Compression ===")
	fmt.Println()

	anotherLongText := `
	User: I'm trying to optimize my database queries in my Go application. They're running really slow.
	Agent: I can help with that! First, let's identify the bottleneck. Are you using an ORM like GORM, or writing raw SQL queries?
	User: I'm using GORM.
	Agent: GORM is great, but it can sometimes generate inefficient queries. Here are some optimization strategies: 1) Use Preload() for eager loading instead of making multiple queries. 2) Add indexes to frequently queried columns. 3) Use Select() to only fetch the columns you need instead of SELECT *. 4) Consider using raw SQL for complex queries. 5) Enable query logging to see what SQL is being generated.
	User: How do I enable query logging?
	Agent: In GORM v2, you can enable logging like this: db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}). This will show you all the SQL queries being executed, which helps identify slow or inefficient queries.
	User: Thanks! What about connection pooling?
	Agent: Great question! Connection pooling is crucial for performance. Go's database/sql package handles this automatically, but you should configure it properly. Use SetMaxOpenConns() to limit maximum connections (usually 25-100), SetMaxIdleConns() for idle connections (usually 10-25), and SetConnMaxLifetime() to prevent stale connections (usually 5-10 minutes).
	`

	fmt.Println("Compressing with streaming (watch it happen in real-time):")
	fmt.Println(strings.Repeat("-", 80))

	streamCompressed, err := compressor.CompressTextStream(anotherLongText, func(chunk snip.ChatResponse) error {
		fmt.Print(chunk.Text)
		return nil
	})
	if err != nil {
		log.Fatalf("Error compressing text with streaming: %v", err)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()
	fmt.Println("Original text length:", len(anotherLongText), "characters")
	fmt.Println("Compressed text length:", len(streamCompressed.Text), "characters")
	fmt.Println("Compression ratio:", fmt.Sprintf("%.1f%%", (1.0-float64(len(streamCompressed.Text))/float64(len(anotherLongText)))*100))
	fmt.Println()
	fmt.Println("✅ Compression complete!")
}
