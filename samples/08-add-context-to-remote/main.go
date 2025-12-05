package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/snip"
)

func main() {
	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	//chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:0.5B-F16")
	// Create an agentOne with both chat flows enabled and HTTP server configuration
	agentOne, err := snip.NewAgent(ctx,
		snip.AgentConfig{
			Name:               "AGENT_ONE",
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
			InformationPath:    "/api/information",
			//ShutdownPath:       "-", // Disable shutdown endpoint
			ShutdownPath: "/server/shutdown", // Enable shutdown endpoint

		}),
	)
	if err != nil {
		log.Fatalf("Error creating agent: %v", err)
	}

	// Create a remote agent that connects to the server
	remoteAgent := snip.NewRemoteAgent(
		"Remote Knowledge Agent",
		snip.ConfigHTTP{
			Address:            "0.0.0.0:9100",
			ChatFlowPath:       "/api/chat",
			ChatStreamFlowPath: "/api/chat-stream",
			InformationPath:    "/api/information",
		},
	)

	// Start agentOne in a goroutine
	go func() {
		log.Println("Starting Agent One server on port 9100...")
		if err := agentOne.Serve(); err != nil {
			log.Fatalf("Server One error: %v", err)
		}
	}()

	// Wait a moment for the server to start
	log.Println("⏳ Waiting for the server to start and the model to load...")
	// In production code, implement a proper health check or wait mechanism
	// Here we just use a simple sleep for demonstration purposes
	//
	time.Sleep(2 * time.Second)

	_ = remoteAgent.AddSystemMessage("Philippe Charrière is a French Solutions Architect at Docker.")

	info, err := remoteAgent.GetInfo()
	if err != nil {
		log.Fatalf("❌ Error getting agent info: %v", err)
	}

	log.Printf("Agent Name: %s", info.Name)
	log.Printf("Model ID: %s", info.ModelID)
	log.Printf("Temperature: %.2f", info.Config.Temperature)
	log.Printf("TopP: %.2f", info.Config.TopP)
	log.Println("")
	log.Println("✅ Agent information retrieved")
	log.Println("")

	// Example 2: Ask a question using streaming mode
	log.Println("📝 Streaming question")
	log.Println("Question: Who is Philippe Charrière?")
	log.Println("---")

	fullResponse, err := remoteAgent.AskStream(
		"Who is Philippe Charrière?",
		func(chunk snip.ChatResponse) error {
			fmt.Print(chunk.Text)
			return nil
		},
	)
	if err != nil {
		log.Fatalf("❌ Error asking streaming question: %v", err)
	}
	fmt.Println("")
	log.Println("---")
	log.Println("✅ Streaming request completed")
	log.Printf("📊 Total response length: %d characters", len(fullResponse.Text))
	log.Println("")

		log.Println("===")

	//fmt.Println(remoteAgent.GetMessages())
	for _, msg := range remoteAgent.GetMessages() {
		fmt.Printf("Role: %s, Content: %v\n", msg.Role, msg.Content[0].Text)
	}
}
