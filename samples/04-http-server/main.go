package main

import (
	"context"
	"log"

	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/chat"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

func main() {
	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	// Create an agent with both chat flows enabled and HTTP server configuration
	agent, err := chat.NewChatAgent(ctx,
		agents.AgentConfig{
			Name:               "HTTP Agent",
			SystemInstructions: "You are a helpful assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		chat.EnableChatFlowWithMemory(),
		chat.EnableChatStreamFlowWithMemory(),
		chat.EnableServer(chat.ConfigHTTP{
			Address:            "0.0.0.0:9100",
			ChatFlowPath:       "/api/chat",
			ChatStreamFlowPath: "/api/chat-stream",
			//ShutdownPath:       "-", // Disable shutdown endpoint
			ShutdownPath: "/server/shutdown", // Disable shutdown endpoint

		}),
		chat.WithLogLevel(logger.LevelDebug), // chat.WithVerbose(true)
	)
	if err != nil {
		log.Fatalf("Error creating agent: %v", err)
	}

	log.Printf("Agent '%s' created successfully", agent.GetName())
	log.Println("Available endpoints:")
	log.Println("  POST http://0.0.0.0:9100/api/chat")
	log.Println("  POST http://0.0.0.0:9100/api/chat-stream")
	log.Println("  POST http://0.0.0.0:9100/shutdown (stop the server)")
	log.Println("")

	// Start the HTTP server (handles Ctrl+C automatically)
	if err := agent.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
