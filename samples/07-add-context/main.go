package main

import (
	"context"
	"fmt"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/smart"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	agent0 := smart.NewAgent(ctx,
		"Local Agent",
		"You are a helpful assistant.",
		chatModelId,
		engineURL,
		smart.Config{
			Temperature: 0.5,
			TopP:        0.9,
		},
		smart.EnableChatStreamFlowWithMemory(),
	)

	_ = agent0.AddContextToMessages("Philippe Charrière is a French Solutions Architect at Docker.")

	_, err := agent0.AskStream("Who is Philippe Charrière?",
		func(chunk string) error {
			fmt.Print(chunk)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}
}