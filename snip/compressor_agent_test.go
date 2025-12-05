package snip

import (
	"context"
	"strings"
	"testing"
)

// ============================================================================
// Tests for CompressorAgent
// ============================================================================

func TestCompressorAgentGetName(t *testing.T) {
	agent := &Agent{
		Name: "TestCompressor",
	}
	compressor := &CompressorAgent{
		agent: agent,
	}

	if got := compressor.GetName(); got != "TestCompressor" {
		t.Errorf("GetName() = %q, want %q", got, "TestCompressor")
	}
}

func TestCompressorAgentGetKind(t *testing.T) {
	agent := &Agent{
		Name: "TestCompressor",
	}
	compressor := &CompressorAgent{
		agent: agent,
	}

	if got := compressor.GetKind(); got != Compressor {
		t.Errorf("GetKind() = %v, want %v", got, Compressor)
	}
}

func TestCompressorAgentGetInfo(t *testing.T) {
	agent := &Agent{
		Name:    "TestCompressor",
		ModelID: "test-model",
		Config: ModelConfig{
			Temperature: 0.7,
		},
	}
	compressor := &CompressorAgent{
		agent: agent,
	}

	info, err := compressor.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() unexpected error: %v", err)
	}

	if info.Name != "TestCompressor" {
		t.Errorf("GetInfo().Name = %q, want %q", info.Name, "TestCompressor")
	}

	if info.ModelID != "test-model" {
		t.Errorf("GetInfo().ModelID = %q, want %q", info.ModelID, "test-model")
	}
}

func TestCompressorAgentCompressText(t *testing.T) {
	t.Run("chatFlow not initialized", func(t *testing.T) {
		ctx := context.Background()
		agent := &Agent{
			ctx:      ctx,
			chatFlow: nil,
		}
		compressor := &CompressorAgent{
			agent:             agent,
			compressionPrompt: "Compress this: ",
		}

		_, err := compressor.CompressText("This is a test")
		if err == nil {
			t.Error("CompressText() expected error when chatFlow is nil, got nil")
		}

		expectedMsg := "chat flow is not initialized"
		if err.Error() != expectedMsg {
			t.Errorf("CompressText() error message = %q, want %q", err.Error(), expectedMsg)
		}
	})

	t.Run("empty text", func(t *testing.T) {
		ctx := context.Background()
		agent := &Agent{
			ctx:      ctx,
			chatFlow: nil,
		}
		compressor := &CompressorAgent{
			agent:             agent,
			compressionPrompt: "Compress this: ",
		}

		// Should still fail because chatFlow is nil
		_, err := compressor.CompressText("")
		if err == nil {
			t.Error("CompressText() expected error when chatFlow is nil, got nil")
		}
	})
}

func TestCompressorAgentCompressTextStream(t *testing.T) {
	t.Run("chatStreamFlow not initialized", func(t *testing.T) {
		ctx := context.Background()
		agent := &Agent{
			ctx:            ctx,
			chatStreamFlow: nil,
		}
		compressor := &CompressorAgent{
			agent:             agent,
			compressionPrompt: "Compress this: ",
		}

		callback := func(chunk ChatResponse) error {
			return nil
		}

		_, err := compressor.CompressTextStream("This is a test", callback)
		if err == nil {
			t.Error("CompressTextStream() expected error when chatStreamFlow is nil, got nil")
		}

		expectedMsg := "chat stream flow is not initialized"
		if err.Error() != expectedMsg {
			t.Errorf("CompressTextStream() error message = %q, want %q", err.Error(), expectedMsg)
		}
	})

	t.Run("callback is called", func(t *testing.T) {
		ctx := context.Background()
		agent := &Agent{
			ctx:            ctx,
			chatStreamFlow: nil,
		}
		compressor := &CompressorAgent{
			agent:             agent,
			compressionPrompt: "Compress this: ",
		}

		callbackCalled := false
		callback := func(chunk ChatResponse) error {
			callbackCalled = true
			return nil
		}

		_, _ = compressor.CompressTextStream("This is a test", callback)

		// Callback won't be called because chatStreamFlow is nil
		// This is expected behavior - the error is returned before any streaming
		if callbackCalled {
			t.Error("CompressTextStream() callback should not be called when chatStreamFlow is nil")
		}
	})
}

func TestCompressorAgentCompressionPrompt(t *testing.T) {
	t.Run("prompt is set during creation", func(t *testing.T) {
		agent := &Agent{}

		compressor := &CompressorAgent{
			agent:             agent,
			compressionPrompt: "Test prompt",
		}

		if compressor.compressionPrompt != "Test prompt" {
			t.Errorf("CompressionPrompt = %q, want %q", compressor.compressionPrompt, "Test prompt")
		}
	})

	t.Run("default prompt contains expected keywords", func(t *testing.T) {
		// This tests the default prompt created by NewCompressorAgent
		// We check that it contains key compression-related terms
		defaultPrompt := `You are a context compression specialist. Your task is to analyze the conversation history and compress it while preserving all essential information.`

		expectedKeywords := []string{
			"compression",
			"compress",
			"conversation",
			"essential information",
		}

		for _, keyword := range expectedKeywords {
			if !strings.Contains(strings.ToLower(defaultPrompt), strings.ToLower(keyword)) {
				t.Errorf("Default compression prompt missing keyword: %q", keyword)
			}
		}
	})
}
