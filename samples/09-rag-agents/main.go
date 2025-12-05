package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/smart"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")

	firstEmbeddingModelId := env.GetEnvOrDefault("EMBEDDING_MODEL_1", "ai/embeddinggemma")
	secondEmbeddingModelId := env.GetEnvOrDefault("EMBEDDING_MODEL_2", "ai/granite-embedding-multilingual")
	thirdEmbeddingModelId := env.GetEnvOrDefault("EMBEDDING_MODEL_3", "ai/mxbai-embed-large")

	// Stack errors during RAG agents creation
	var errs []error

	ragAgent01, err := smart.NewRagAgent(ctx,
		smart.RagAgentConfig{
			Name:      "RAG_Agent_1",
			ModelID:   firstEmbeddingModelId,
			EngineURL: engineURL,
		},
		smart.StoreConfig{
			StoreName: "RAG_Store_1",
			StorePath: "./data",
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("RAG_Agent_1: %w", err))
	}

	ragAgent02, err := smart.NewRagAgent(ctx,
		smart.RagAgentConfig{
			Name:      "RAG_Agent_2",
			ModelID:   secondEmbeddingModelId,
			EngineURL: engineURL,
		},
		smart.StoreConfig{
			StoreName: "RAG_Store_2",
			StorePath: "./data",
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("RAG_Agent_2: %w", err))
	}

	ragAgent03, err := smart.NewRagAgent(ctx,
		smart.RagAgentConfig{
			Name:      "RAG_Agent_3",
			ModelID:   thirdEmbeddingModelId,
			EngineURL: engineURL,
		},
		smart.StoreConfig{
			StoreName: "RAG_Store_3",
			StorePath: "./data",
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("RAG_Agent_3: %w", err))
	}

	// Check if any errors were stacked
	if len(errs) > 0 {
		fmt.Printf("Error creating RAG agents: %v\n", errors.Join(errs...))
		return
	}

	fmt.Println(ragAgent01.GetInfo())
	fmt.Println(ragAgent02.GetInfo())
	fmt.Println(ragAgent03.GetInfo())

	if ragAgent01.GetNumberOfDocuments() == 0 {
		txtChunks := []string{
			"Squirrels run in the forest",
			"Birds fly in the sky",
			"Frogs swim in the pond",
			"Fishes swim in the sea",
			"Lions roar in the savannah",
			"Eagles soar above the mountains",
			"Dolphins leap out of the ocean",
			"Bears fish in the river",
		}
		var chunks []smart.TextChunk
		for _, txt := range txtChunks {
			chunk := smart.TextChunk{
				Content:  txt,
				Metadata: map[string]any{"source": "example"},
			}
			chunks = append(chunks, chunk)
		}

		_, err = ragAgent01.AddTextChunksToStore(chunks)
		if err != nil {
			fmt.Printf("Error adding text chunks to RAG_Agent_1: %v\n", err)
			return
		}

		_, err = ragAgent02.AddTextChunksToStore(chunks)
		if err != nil {
			fmt.Printf("Error adding text chunks to RAG_Agent_2: %v\n", err)
			return
		}

		_, err = ragAgent03.AddTextChunksToStore(chunks)
		if err != nil {
			fmt.Printf("Error adding text chunks to RAG_Agent_3: %v\n", err)
			return
		}

		fmt.Println("Text chunks added to all RAG agents successfully.")
	} else {
		fmt.Println("RAG agents already have documents in their stores.")
	}

	query := "Which animals swim?"
	similarities1, err := ragAgent01.SearchSimilarities(query)
	if err != nil {
		fmt.Printf("Error searching similarities with RAG_Agent_1: %v\n", err)
		return
	}
	fmt.Printf("RAG_Agent_1 similarities for query '%s':\n", query)
	for _, sim := range similarities1 {
		fmt.Println("  -", sim)
	}

	similarities2, err := ragAgent02.SearchSimilarities(query)
	if err != nil {
		fmt.Printf("Error searching similarities with RAG_Agent_2: %v\n", err)
		return
	}
	fmt.Printf("RAG_Agent_2 similarities for query '%s':\n", query)
	for _, sim := range similarities2 {
		fmt.Println("  -", sim)
	}

	similarities3, err := ragAgent03.SearchSimilarities(query)
	if err != nil {
		fmt.Printf("Error searching similarities with RAG_Agent_3: %v\n", err)
		return
	}
	fmt.Printf("RAG_Agent_3 similarities for query '%s':\n", query)
	for _, sim := range similarities3 {
		fmt.Println("  -", sim)
	}

}
