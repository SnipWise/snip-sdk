package models

import (
	"testing"
)

// ============================================================================
// Tests for ModelConfig.ToOpenAIParams
// ============================================================================

func TestModelConfigToOpenAIParams(t *testing.T) {
	t.Run("all fields set", func(t *testing.T) {
		seed := int64(42)
		config := ModelConfig{
			Temperature:      0.7,
			TopP:             0.9,
			MaxTokens:        1000,
			FrequencyPenalty: 0.5,
			PresencePenalty:  0.3,
			Stop:             []string{"STOP", "END"},
			Seed:             &seed,
		}

		params := config.ToOpenAIParams()

		if params == nil {
			t.Fatal("ToOpenAIParams() returned nil")
		}

		// Check Temperature
		if !params.Temperature.Valid() {
			t.Error("Temperature not set in params")
		} else if params.Temperature.Value != 0.7 {
			t.Errorf("Temperature = %f, want 0.7", params.Temperature.Value)
		}

		// Check TopP
		if !params.TopP.Valid() {
			t.Error("TopP not set in params")
		} else if params.TopP.Value != 0.9 {
			t.Errorf("TopP = %f, want 0.9", params.TopP.Value)
		}

		// Check MaxTokens
		if !params.MaxTokens.Valid() {
			t.Error("MaxTokens not set in params")
		} else if params.MaxTokens.Value != 1000 {
			t.Errorf("MaxTokens = %d, want 1000", params.MaxTokens.Value)
		}

		// Check FrequencyPenalty
		if !params.FrequencyPenalty.Valid() {
			t.Error("FrequencyPenalty not set in params")
		} else if params.FrequencyPenalty.Value != 0.5 {
			t.Errorf("FrequencyPenalty = %f, want 0.5", params.FrequencyPenalty.Value)
		}

		// Check PresencePenalty
		if !params.PresencePenalty.Valid() {
			t.Error("PresencePenalty not set in params")
		} else if params.PresencePenalty.Value != 0.3 {
			t.Errorf("PresencePenalty = %f, want 0.3", params.PresencePenalty.Value)
		}

		// Check Seed
		if !params.Seed.Valid() {
			t.Error("Seed not set in params")
		} else if params.Seed.Value != 42 {
			t.Errorf("Seed = %d, want 42", params.Seed.Value)
		}
	})

	t.Run("zero values not set", func(t *testing.T) {
		config := ModelConfig{
			Temperature:      0,
			TopP:             0,
			MaxTokens:        0,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
			Seed:             nil,
		}

		params := config.ToOpenAIParams()

		if params == nil {
			t.Fatal("ToOpenAIParams() returned nil")
		}

		// Zero values should not be set
		if params.Temperature.Valid() {
			t.Error("Temperature should not be set for zero value")
		}
		if params.TopP.Valid() {
			t.Error("TopP should not be set for zero value")
		}
		if params.MaxTokens.Valid() {
			t.Error("MaxTokens should not be set for zero value")
		}
		if params.FrequencyPenalty.Valid() {
			t.Error("FrequencyPenalty should not be set for zero value")
		}
		if params.PresencePenalty.Valid() {
			t.Error("PresencePenalty should not be set for zero value")
		}
		if params.Seed.Valid() {
			t.Error("Seed should not be set when not provided")
		}
	})

	t.Run("partial fields set", func(t *testing.T) {
		config := ModelConfig{
			Temperature: 0.8,
			MaxTokens:   500,
			// Other fields are zero/nil
		}

		params := config.ToOpenAIParams()

		if params == nil {
			t.Fatal("ToOpenAIParams() returned nil")
		}

		// Check set fields
		if !params.Temperature.Valid() {
			t.Error("Temperature should be set")
		} else if params.Temperature.Value != 0.8 {
			t.Errorf("Temperature = %f, want 0.8", params.Temperature.Value)
		}

		if !params.MaxTokens.Valid() {
			t.Error("MaxTokens should be set")
		} else if params.MaxTokens.Value != 500 {
			t.Errorf("MaxTokens = %d, want 500", params.MaxTokens.Value)
		}

		// Check unset fields
		if params.TopP.Valid() {
			t.Error("TopP should not be set")
		}
		if params.FrequencyPenalty.Valid() {
			t.Error("FrequencyPenalty should not be set")
		}
		if params.PresencePenalty.Valid() {
			t.Error("PresencePenalty should not be set")
		}
	})

	t.Run("negative values", func(t *testing.T) {
		config := ModelConfig{
			FrequencyPenalty: -0.5,
			PresencePenalty:  -0.3,
		}

		params := config.ToOpenAIParams()

		if params == nil {
			t.Fatal("ToOpenAIParams() returned nil")
		}

		// Negative values should be set (they are valid for penalties)
		if !params.FrequencyPenalty.Valid() {
			t.Error("FrequencyPenalty should be set for negative value")
		} else if params.FrequencyPenalty.Value != -0.5 {
			t.Errorf("FrequencyPenalty = %f, want -0.5", params.FrequencyPenalty.Value)
		}

		if !params.PresencePenalty.Valid() {
			t.Error("PresencePenalty should be set for negative value")
		} else if params.PresencePenalty.Value != -0.3 {
			t.Errorf("PresencePenalty = %f, want -0.3", params.PresencePenalty.Value)
		}
	})

	t.Run("extreme values", func(t *testing.T) {
		seed := int64(9999999)
		config := ModelConfig{
			Temperature:      2.0,
			TopP:             1.0,
			MaxTokens:        4096,
			FrequencyPenalty: 2.0,
			PresencePenalty:  2.0,
			Seed:             &seed,
		}

		params := config.ToOpenAIParams()

		if params == nil {
			t.Fatal("ToOpenAIParams() returned nil")
		}

		// All extreme values should be set
		if !params.Temperature.Valid() || params.Temperature.Value != 2.0 {
			t.Error("Temperature not correctly set")
		}
		if !params.TopP.Valid() || params.TopP.Value != 1.0 {
			t.Error("TopP not correctly set")
		}
		if !params.MaxTokens.Valid() || params.MaxTokens.Value != 4096 {
			t.Error("MaxTokens not correctly set")
		}
		if !params.FrequencyPenalty.Valid() || params.FrequencyPenalty.Value != 2.0 {
			t.Error("FrequencyPenalty not correctly set")
		}
		if !params.PresencePenalty.Valid() || params.PresencePenalty.Value != 2.0 {
			t.Error("PresencePenalty not correctly set")
		}
		if !params.Seed.Valid() || params.Seed.Value != 9999999 {
			t.Error("Seed not correctly set")
		}
	})

	t.Run("seed with zero value", func(t *testing.T) {
		seed := int64(0)
		config := ModelConfig{
			Seed: &seed,
		}

		params := config.ToOpenAIParams()

		if params == nil {
			t.Fatal("ToOpenAIParams() returned nil")
		}

		// Seed should be set even if it's 0 (because it's a pointer)
		if !params.Seed.Valid() {
			t.Error("Seed should be set when pointer is provided, even if value is 0")
		} else if params.Seed.Value != 0 {
			t.Errorf("Seed = %d, want 0", params.Seed.Value)
		}
	})

	t.Run("stop sequences", func(t *testing.T) {
		config := ModelConfig{
			Stop: []string{"STOP", "END", "TERMINATE"},
		}

		params := config.ToOpenAIParams()

		if params == nil {
			t.Fatal("ToOpenAIParams() returned nil")
		}

		// Note: Current implementation doesn't handle Stop parameter
		// This test documents the current behavior
		// The Stop field is not converted to params
		t.Log("Stop parameter is currently not converted (known limitation)")
	})

	t.Run("empty stop sequences", func(t *testing.T) {
		config := ModelConfig{
			Stop: []string{},
		}

		params := config.ToOpenAIParams()

		if params == nil {
			t.Fatal("ToOpenAIParams() returned nil")
		}

		// Empty Stop array should not be set
		t.Log("Empty Stop parameter is not converted")
	})
}
