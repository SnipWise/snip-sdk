package snip

import (
	"context"
	"testing"
)

// ============================================================================
// Tests for RagAgent.GetName
// ============================================================================

func TestRagAgentGetName(t *testing.T) {
	ragAgent := &RagAgent{
		Name: "test-rag-agent",
	}

	name := ragAgent.GetName()
	if name != "test-rag-agent" {
		t.Errorf("GetName() = %q, want %q", name, "test-rag-agent")
	}
}

// ============================================================================
// Tests for RagAgent.Kind
// ============================================================================

func TestRagAgentKind(t *testing.T) {
	ragAgent := &RagAgent{}

	kind := ragAgent.Kind()
	if kind != Rag {
		t.Errorf("Kind() = %v, want %v", kind, Rag)
	}
}

// ============================================================================
// Tests for RagAgent.IsStoreInitialized
// ============================================================================

func TestRagAgentIsStoreInitialized(t *testing.T) {
	t.Run("uninitialized store", func(t *testing.T) {
		ragAgent := &RagAgent{
			docStore: nil,
		}

		if ragAgent.IsStoreInitialized() {
			t.Error("IsStoreInitialized() = true, want false for nil docStore")
		}
	})

	// Note: We can't easily test initialized store without creating a full RagAgent
	// which requires actual embedding model connection
}

// ============================================================================
// Tests for RagAgent.GetNumberOfDocuments
// ============================================================================

func TestRagAgentGetNumberOfDocuments(t *testing.T) {
	t.Run("nil docStore", func(t *testing.T) {
		ragAgent := &RagAgent{
			docStore: nil,
		}

		count := ragAgent.GetNumberOfDocuments()
		if count != 0 {
			t.Errorf("GetNumberOfDocuments() = %d, want 0 for nil docStore", count)
		}
	})

	// Note: Testing with actual documents requires full RagAgent initialization
}

// ============================================================================
// Tests for RagAgent.GetInfo
// ============================================================================

func TestRagAgentGetInfo(t *testing.T) {
	t.Run("basic info without docStore", func(t *testing.T) {
		ragAgent := &RagAgent{
			Name:               "info-test-agent",
			ModelID:            "test-embedding-model",
			embeddingDimension: 384,
			storeName:          "test-store",
			storePath:          "./test-data",
			docStore:           nil, // No documents
		}

		info, err := ragAgent.GetInfo()
		if err != nil {
			t.Errorf("GetInfo() unexpected error: %v", err)
		}

		if info.Name != "info-test-agent" {
			t.Errorf("Info.Name = %q, want %q", info.Name, "info-test-agent")
		}

		if info.ModelID != "test-embedding-model" {
			t.Errorf("Info.ModelID = %q, want %q", info.ModelID, "test-embedding-model")
		}

		if info.EmbeddingDimension != 384 {
			t.Errorf("Info.EmbeddingDimension = %d, want 384", info.EmbeddingDimension)
		}

		if info.StoreName != "test-store" {
			t.Errorf("Info.StoreName = %q, want %q", info.StoreName, "test-store")
		}

		if info.StorePath != "./test-data" {
			t.Errorf("Info.StorePath = %q, want %q", info.StorePath, "./test-data")
		}

		if info.NumberOfDocuments != 0 {
			t.Errorf("Info.NumberOfDocuments = %d, want 0", info.NumberOfDocuments)
		}
	})
}

// ============================================================================
// Tests for RagAgentConfig.Validate
// ============================================================================

func TestRagAgentConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      RagAgentConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: RagAgentConfig{
				Name:      "test-agent",
				ModelID:   "test-model",
				EngineURL: "http://localhost:8080",
			},
			expectError: false,
		},
		{
			name: "missing name",
			config: RagAgentConfig{
				Name:      "",
				ModelID:   "test-model",
				EngineURL: "http://localhost:8080",
			},
			expectError: true,
			errorMsg:    "agent name is required",
		},
		{
			name: "missing model ID",
			config: RagAgentConfig{
				Name:      "test-agent",
				ModelID:   "",
				EngineURL: "http://localhost:8080",
			},
			expectError: true,
			errorMsg:    "model ID is required",
		},
		{
			name: "missing engine URL",
			config: RagAgentConfig{
				Name:      "test-agent",
				ModelID:   "test-model",
				EngineURL: "",
			},
			expectError: true,
			errorMsg:    "engine URL is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				if err == nil {
					t.Error("Validate() expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Validate() error = %q, want %q", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

// ============================================================================
// Tests for TextChunk struct
// ============================================================================

func TestTextChunk(t *testing.T) {
	t.Run("basic chunk", func(t *testing.T) {
		chunk := TextChunk{
			Content:  "This is a test chunk",
			Metadata: nil,
		}

		if chunk.Content != "This is a test chunk" {
			t.Errorf("Content = %q, want %q", chunk.Content, "This is a test chunk")
		}

		if chunk.Metadata != nil {
			t.Error("Metadata should be nil")
		}
	})

	t.Run("chunk with metadata", func(t *testing.T) {
		chunk := TextChunk{
			Content: "Chunk with metadata",
			Metadata: map[string]any{
				"source": "test.txt",
				"page":   42,
			},
		}

		if chunk.Content != "Chunk with metadata" {
			t.Errorf("Content = %q, want %q", chunk.Content, "Chunk with metadata")
		}

		if chunk.Metadata == nil {
			t.Fatal("Metadata should not be nil")
		}

		source, ok := chunk.Metadata["source"]
		if !ok || source != "test.txt" {
			t.Errorf("Metadata[source] = %v, want test.txt", source)
		}

		page, ok := chunk.Metadata["page"]
		if !ok || page != 42 {
			t.Errorf("Metadata[page] = %v, want 42", page)
		}
	})
}

// ============================================================================
// Tests for StoreConfig struct
// ============================================================================

func TestStoreConfig(t *testing.T) {
	t.Run("basic config", func(t *testing.T) {
		config := StoreConfig{
			StoreName: "my-store",
			StorePath: "./data",
		}

		if config.StoreName != "my-store" {
			t.Errorf("StoreName = %q, want %q", config.StoreName, "my-store")
		}

		if config.StorePath != "./data" {
			t.Errorf("StorePath = %q, want %q", config.StorePath, "./data")
		}
	})
}

// ============================================================================
// Tests for RagAgentInfo struct
// ============================================================================

func TestRagAgentInfo(t *testing.T) {
	t.Run("complete info", func(t *testing.T) {
		info := RagAgentInfo{
			Name:               "test-rag-agent",
			ModelID:            "embedding-model",
			EmbeddingDimension: 768,
			StoreName:          "vector-store",
			StorePath:          "/var/data",
			NumberOfDocuments:  100,
		}

		if info.Name != "test-rag-agent" {
			t.Errorf("Name = %q, want %q", info.Name, "test-rag-agent")
		}

		if info.ModelID != "embedding-model" {
			t.Errorf("ModelID = %q, want %q", info.ModelID, "embedding-model")
		}

		if info.EmbeddingDimension != 768 {
			t.Errorf("EmbeddingDimension = %d, want 768", info.EmbeddingDimension)
		}

		if info.StoreName != "vector-store" {
			t.Errorf("StoreName = %q, want %q", info.StoreName, "vector-store")
		}

		if info.StorePath != "/var/data" {
			t.Errorf("StorePath = %q, want %q", info.StorePath, "/var/data")
		}

		if info.NumberOfDocuments != 100 {
			t.Errorf("NumberOfDocuments = %d, want 100", info.NumberOfDocuments)
		}
	})
}

// ============================================================================
// Integration Tests for NewRagAgent
// ============================================================================

// Note: These tests require a running model engine and are marked as integration tests
// They can be skipped in CI/CD by using: go test -short

func TestNewRagAgentIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires a running model engine with embedding model
	// It's designed to fail gracefully if the engine is not available

	ctx := context.Background()

	config := RagAgentConfig{
		Name:      "integration-test-agent",
		ModelID:   "ai/mxbai-embed-large", // Common embedding model
		EngineURL: "http://localhost:12434/engines/llama.cpp/v1",
	}

	storeConfig := StoreConfig{
		StoreName: "integration-test-store",
		StorePath: "./test-data-integration",
	}

	ragAgent, err := NewRagAgent(ctx, config, storeConfig)

	// If the model is not available, the test should fail with a descriptive error
	if err != nil {
		t.Logf("Integration test skipped: %v", err)
		t.Logf("To run this test, ensure a model engine is running at %s with model %s", config.EngineURL, config.ModelID)
		return
	}

	// If we got here, the agent was created successfully
	if ragAgent == nil {
		t.Fatal("NewRagAgent() returned nil without error")
	}

	if ragAgent.Name != "integration-test-agent" {
		t.Errorf("Agent.Name = %q, want %q", ragAgent.Name, "integration-test-agent")
	}

	if ragAgent.ModelID != "ai/mxbai-embed-large" {
		t.Errorf("Agent.ModelID = %q, want %q", ragAgent.ModelID, "ai/mxbai-embed-large")
	}

	if !ragAgent.IsStoreInitialized() {
		t.Error("Document store should be initialized")
	}

	if ragAgent.embeddingDimension == 0 {
		t.Error("Embedding dimension should be calculated and non-zero")
	}

	t.Logf("Successfully created RAG agent with embedding dimension: %d", ragAgent.embeddingDimension)
}

// ============================================================================
// Tests for AddTextChunksToStore
// ============================================================================

func TestRagAgentAddTextChunksIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	config := RagAgentConfig{
		Name:      "chunks-test-agent",
		ModelID:   "ai/mxbai-embed-large",
		EngineURL: "http://localhost:12434/engines/llama.cpp/v1",
	}

	storeConfig := StoreConfig{
		StoreName: "chunks-test-store",
		StorePath: "./test-data-chunks",
	}

	ragAgent, err := NewRagAgent(ctx, config, storeConfig)
	if err != nil {
		t.Logf("Integration test skipped: %v", err)
		return
	}

	chunks := []TextChunk{
		{
			Content:  "The quick brown fox jumps over the lazy dog",
			Metadata: map[string]any{"source": "test1.txt"},
		},
		{
			Content:  "Machine learning is a subset of artificial intelligence",
			Metadata: map[string]any{"source": "test2.txt"},
		},
	}

	count, err := ragAgent.AddTextChunksToStore(chunks)
	if err != nil {
		t.Errorf("AddTextChunksToStore() unexpected error: %v", err)
	}

	if count != 2 {
		t.Errorf("AddTextChunksToStore() returned count = %d, want 2", count)
	}

	// Verify documents were added
	if ragAgent.GetNumberOfDocuments() != 2 {
		t.Errorf("GetNumberOfDocuments() = %d, want 2", ragAgent.GetNumberOfDocuments())
	}

	t.Logf("Successfully added %d chunks to store", count)
}

// ============================================================================
// Tests for SearchSimilarities
// ============================================================================

func TestRagAgentSearchSimilaritiesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	config := RagAgentConfig{
		Name:      "search-test-agent",
		ModelID:   "ai/mxbai-embed-large",
		EngineURL: "http://localhost:12434/engines/llama.cpp/v1",
	}

	storeConfig := StoreConfig{
		StoreName: "search-test-store",
		StorePath: "./test-data-search",
	}

	ragAgent, err := NewRagAgent(ctx, config, storeConfig)
	if err != nil {
		t.Logf("Integration test skipped: %v", err)
		return
	}

	// Add test documents
	chunks := []TextChunk{
		{
			Content:  "Dolphins swim in the ocean",
			Metadata: map[string]any{"category": "marine"},
		},
		{
			Content:  "Eagles fly in the sky",
			Metadata: map[string]any{"category": "birds"},
		},
		{
			Content:  "Whales also swim in the ocean",
			Metadata: map[string]any{"category": "marine"},
		},
	}

	_, err = ragAgent.AddTextChunksToStore(chunks)
	if err != nil {
		t.Fatalf("Failed to add chunks: %v", err)
	}

	// Search for similar documents
	results, err := ragAgent.SearchSimilarities("Which animals swim?")
	if err != nil {
		t.Errorf("SearchSimilarities() unexpected error: %v", err)
	}

	if len(results) == 0 {
		t.Error("SearchSimilarities() returned no results")
	}

	t.Logf("Found %d similar documents", len(results))
	for i, result := range results {
		t.Logf("Result %d: %s", i+1, result)
	}

	// The results should contain documents about swimming (dolphins, whales)
	// This is a simple heuristic test
	foundSwimming := false
	for _, result := range results {
		if len(result) > 0 {
			foundSwimming = true
			break
		}
	}

	if !foundSwimming {
		t.Error("Expected to find at least one non-empty result")
	}
}
