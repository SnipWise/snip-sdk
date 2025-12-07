package tools

import (
	"encoding/json"
	"fmt"
)

// Transform converts a map[string]interface{} output to a typed struct using JSON marshaling
func Transform[T any](output any) (T, error) {
	var result T

	// Convert output to JSON bytes
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		return result, fmt.Errorf("failed to marshal output: %w", err)
	}

	// Unmarshal JSON bytes into the typed struct
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal to type: %w", err)
	}

	return result, nil
}
