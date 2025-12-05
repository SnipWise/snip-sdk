package main

import (
	"context"
	"log"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/snip"
)

func main() {
	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	// Create an agent with both chat flows enabled and HTTP server configuration
	agent, err := snip.NewAgent(ctx,
		snip.AgentConfig{
			Name:               "HTTP Agent",
			SystemInstructions: "You are a helpful assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		snip.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		snip.EnableChatFlowWithMemory(),
		snip.EnableChatStreamFlowWithMemory(),
		snip.EnableServer(snip.ConfigHTTP{
			Address:            "0.0.0.0:9100",
			ChatFlowPath:       "/api/chat",
			ChatStreamFlowPath: "/api/chat-stream",
			//ShutdownPath:       "-", // Disable shutdown endpoint
			ShutdownPath:       "/server/shutdown", // Disable shutdown endpoint

		}),
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
