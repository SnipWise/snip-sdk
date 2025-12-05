package snip

import (
	"context"
	"testing"
)

// ============================================================================
// Integration Tests for GetModelsList
// ============================================================================

func TestGetModelsListIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	modelRunnerEndpoint := "http://localhost:12434/engines/llama.cpp/v1"

	models, err := GetModelsList(ctx, modelRunnerEndpoint)

	if err != nil {
		t.Logf("Integration test skipped: %v", err)
		t.Logf("To run this test, ensure a model engine is running at %s", modelRunnerEndpoint)
		return
	}

	if models == nil {
		t.Error("GetModelsList() returned nil models slice")
	}

	t.Logf("Found %d models", len(models))
	for _, model := range models {
		t.Logf("  - %s", model)
	}

	// We can't assert specific models since it depends on what's loaded
	// but we can verify the function works correctly
}

// ============================================================================
// Integration Tests for IsModelAvailable
// ============================================================================

func TestIsModelAvailableIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	modelRunnerEndpoint := "http://localhost:12434/engines/llama.cpp/v1"

	t.Run("check for non-existent model", func(t *testing.T) {
		// This model should not exist
		modelID := "non-existent-model-12345"

		available := IsModelAvailable(ctx, modelRunnerEndpoint, modelID)

		if available {
			t.Errorf("IsModelAvailable() = true for non-existent model %q", modelID)
		}

		t.Logf("Correctly identified that model %q is not available", modelID)
	})

	t.Run("check for potentially available model", func(t *testing.T) {
		// Try to get the list of models first
		models, err := GetModelsList(ctx, modelRunnerEndpoint)
		if err != nil {
			t.Logf("Could not get models list: %v", err)
			return
		}

		if len(models) == 0 {
			t.Log("No models available on the engine")
			return
		}

		// Check if the first model is available
		modelID := models[0]
		available := IsModelAvailable(ctx, modelRunnerEndpoint, modelID)

		if !available {
			t.Errorf("IsModelAvailable() = false for existing model %q", modelID)
		} else {
			t.Logf("Correctly identified that model %q is available", modelID)
		}
	})
}

// ============================================================================
// Unit Tests for function behavior (without actual API calls)
// ============================================================================

func TestGetModelsListWithInvalidEndpoint(t *testing.T) {
	ctx := context.Background()

	// Test with invalid endpoint
	invalidEndpoint := "http://invalid-host-that-does-not-exist:99999/v1"

	models, err := GetModelsList(ctx, invalidEndpoint)

	if err == nil {
		t.Error("GetModelsList() expected error for invalid endpoint, got nil")
	}

	if models == nil {
		// This is acceptable - returning nil on error
	} else if len(models) != 0 {
		t.Errorf("GetModelsList() should return empty slice on error, got %d models", len(models))
	}
}

func TestIsModelAvailableWithInvalidEndpoint(t *testing.T) {
	ctx := context.Background()

	// Test with invalid endpoint
	invalidEndpoint := "http://invalid-host-that-does-not-exist:99999/v1"
	modelID := "any-model"

	available := IsModelAvailable(ctx, invalidEndpoint, modelID)

	if available {
		t.Error("IsModelAvailable() should return false for invalid endpoint")
	}
}
