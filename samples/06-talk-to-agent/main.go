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
	time.Sleep(2 * time.Second)

	// Example 0: Get agent information
	log.Println("📝 Example 0: Get remote agent information")
	log.Println("---")

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

	//Example 1: Ask a question using non-streaming mode
	log.Println("📝 Example 1: Non-streaming question")
	log.Println("Question: What is Go programming language?")
	log.Println("---")

	response, err := remoteAgent.Ask("What is Go programming language?")
	if err != nil {
		log.Fatalf("❌ Error asking question: %v", err)
	}

	fmt.Println(response)
	log.Println("")
	log.Println("✅ Non-streaming request completed")
	log.Println("")

	// Example 2: Ask a question using streaming mode
	log.Println("📝 Example 2: Streaming question")
	log.Println("Question: Explain what are goroutines in 2 sentences")
	log.Println("---")

	fullResponse, err := remoteAgent.AskStream(
		"Explain what are goroutines in 2 sentences",
		func(chunk snip.ChatResponse) error {
			// Print the chunk text if present
			if chunk.Text != "" {
				fmt.Print(chunk.Text)
			}
			// Check if this is the final chunk with metadata
			if chunk.FinishReason != "" {
				fmt.Printf("\n[Streaming completed - FinishReason: %s]", chunk.FinishReason)
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalf("❌ Error asking streaming question: %v", err)
	}

	log.Println("")
	log.Println("")
	log.Println("✅ Streaming request completed")
	log.Printf("📊 Total response length: %d characters", len(fullResponse.Text))
	log.Println("")

	// Example 3: Multiple questions to test memory
	log.Println("📝 Example 3: Follow-up question (testing memory)")
	log.Println("Question: What are the benefits of using it?")
	log.Println("---")

	response3, err := remoteAgent.Ask("What are the benefits of using it?")
	if err != nil {
		log.Fatalf("❌ Error asking follow-up question: %v", err)
	}

	fmt.Println(response3)
	log.Println("")
	log.Println("✅ Follow-up question completed")

	log.Println("📏 GetCurrentContextSize", remoteAgent.GetCurrentContextSize())

	// info, err = remoteAgent.GetInfo()
	// if err != nil {
	// 	log.Fatalf("❌ Error getting agent info: %v", err)
	// }

	// log.Printf("Agent Name: %s", info.Name)
	// log.Printf("Model ID: %s", info.ModelID)
	// log.Printf("Temperature: %.2f", info.Config.Temperature)
	// log.Printf("TopP: %.2f", info.Config.TopP)
	// log.Println("")
	// log.Println("✅ Agent information retrieved")
	// log.Println("")

}
