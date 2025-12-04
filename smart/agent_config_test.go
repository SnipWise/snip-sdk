package smart

import (
	"testing"
)

// ============================================================================
// Tests for AgentConfig.Validate
// ============================================================================

func TestAgentConfigValidate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := AgentConfig{
			Name:               "test-agent",
			SystemInstructions: "You are helpful",
			ModelID:            "test-model",
			EngineURL:          "http://localhost:8080",
		}

		err := config.Validate()
		if err != nil {
			t.Errorf("Validate() unexpected error for valid config: %v", err)
		}
	})

	t.Run("missing name", func(t *testing.T) {
		config := AgentConfig{
			Name:               "",
			SystemInstructions: "You are helpful",
			ModelID:            "test-model",
			EngineURL:          "http://localhost:8080",
		}

		err := config.Validate()
		if err == nil {
			t.Error("Validate() expected error for missing name, got nil")
		}
		expectedMsg := "agent name is required"
		if err.Error() != expectedMsg {
			t.Errorf("Validate() error message = %q, want %q", err.Error(), expectedMsg)
		}
	})

	t.Run("missing modelID", func(t *testing.T) {
		config := AgentConfig{
			Name:               "test-agent",
			SystemInstructions: "You are helpful",
			ModelID:            "",
			EngineURL:          "http://localhost:8080",
		}

		err := config.Validate()
		if err == nil {
			t.Error("Validate() expected error for missing modelID, got nil")
		}
		expectedMsg := "model ID is required"
		if err.Error() != expectedMsg {
			t.Errorf("Validate() error message = %q, want %q", err.Error(), expectedMsg)
		}
	})

	t.Run("missing engineURL", func(t *testing.T) {
		config := AgentConfig{
			Name:               "test-agent",
			SystemInstructions: "You are helpful",
			ModelID:            "test-model",
			EngineURL:          "",
		}

		err := config.Validate()
		if err == nil {
			t.Error("Validate() expected error for missing engineURL, got nil")
		}
		expectedMsg := "engine URL is required"
		if err.Error() != expectedMsg {
			t.Errorf("Validate() error message = %q, want %q", err.Error(), expectedMsg)
		}
	})

	t.Run("all fields missing", func(t *testing.T) {
		config := AgentConfig{}

		err := config.Validate()
		if err == nil {
			t.Error("Validate() expected error for all missing fields, got nil")
		}
		// Should fail on the first check (name)
		expectedMsg := "agent name is required"
		if err.Error() != expectedMsg {
			t.Errorf("Validate() error message = %q, want %q", err.Error(), expectedMsg)
		}
	})

	t.Run("empty system instructions is allowed", func(t *testing.T) {
		config := AgentConfig{
			Name:               "test-agent",
			SystemInstructions: "",
			ModelID:            "test-model",
			EngineURL:          "http://localhost:8080",
		}

		err := config.Validate()
		if err != nil {
			t.Errorf("Validate() unexpected error for empty SystemInstructions: %v", err)
		}
	})

	t.Run("whitespace only fields", func(t *testing.T) {
		tests := []struct {
			name   string
			config AgentConfig
		}{
			{
				"whitespace name",
				AgentConfig{
					Name:      "   ",
					ModelID:   "test-model",
					EngineURL: "http://localhost:8080",
				},
			},
			{
				"whitespace modelID",
				AgentConfig{
					Name:      "test-agent",
					ModelID:   "   ",
					EngineURL: "http://localhost:8080",
				},
			},
			{
				"whitespace engineURL",
				AgentConfig{
					Name:      "test-agent",
					ModelID:   "test-model",
					EngineURL: "   ",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Note: Current implementation doesn't trim whitespace
				// Whitespace-only strings are considered valid
				// This test documents current behavior
				err := tt.config.Validate()
				if err != nil {
					// If validation should reject whitespace, update the implementation
					t.Logf("Validation rejected whitespace (current behavior): %v", err)
				} else {
					t.Logf("Validation accepted whitespace (current behavior)")
				}
			})
		}
	})
}
