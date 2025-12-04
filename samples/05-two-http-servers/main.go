package main

import (
	"context"
	"log"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/smart"
)

func main() {
	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	// Create an agentOne with both chat flows enabled and HTTP server configuration
	agentOne := smart.NewAgent(ctx,
		smart.AgentConfig{
			Name:               "HTTP_Agent_1",
			SystemInstructions: "You are a helpful assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		smart.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		smart.EnableChatFlowWithMemory(),
		smart.EnableChatStreamFlowWithMemory(),
		smart.EnableServer(smart.ConfigHTTP{
			Address:            "0.0.0.0:9100",
			ChatFlowPath:       "/api/chat",
			ChatStreamFlowPath: "/api/chat-stream",
			//ShutdownPath:       "-", // Disable shutdown endpoint
			ShutdownPath:       "/server/shutdown", // Disable shutdown endpoint

		}),
	)

	agentTwo := smart.NewAgent(ctx,
		smart.AgentConfig{
			Name:               "HTTP_Agent_2",
			SystemInstructions: "You are a helpful assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		smart.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		smart.EnableChatFlowWithMemory(),
		smart.EnableChatStreamFlowWithMemory(),
		smart.EnableServer(smart.ConfigHTTP{
			Address:            "0.0.0.0:9200",
			ChatFlowPath:       "/api/chat",
			ChatStreamFlowPath: "/api/chat-stream",
			//ShutdownPath:       "-", // Disable shutdown endpoint
			ShutdownPath:       "/server/shutdown", // Disable shutdown endpoint

		}),
	)


	// Start agentOne in a goroutine
	go func() {
		log.Println("Starting Agent One server on port 9100...")
		if err := agentOne.Serve(); err != nil {
			log.Fatalf("Server One error: %v", err)
		}
	}()

	// Start agentTwo in the main thread (handles Ctrl+C automatically)
	log.Println("Starting Agent Two server on port 9200...")
	if err := agentTwo.Serve(); err != nil {
		log.Fatalf("Server Two error: %v", err)
	}
}
