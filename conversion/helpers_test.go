package conversion

import (
	"testing"
)

// ============================================================================
// Tests for StringToInt
// ============================================================================

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"positive integer", "42", 42},
		{"negative integer", "-10", -10},
		{"zero", "0", 0},
		{"large number", "999999", 999999},
		{"invalid string", "abc", 0},
		{"empty string", "", 0},
		{"float string", "3.14", 0},
		{"mixed string", "123abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToInt(tt.input)
			if result != tt.expected {
				t.Errorf("StringToInt(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Tests for StringToFloat
// ============================================================================

func TestStringToFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"positive float", "3.14", 3.14},
		{"negative float", "-2.5", -2.5},
		{"zero", "0", 0.0},
		{"integer as float", "42", 42.0},
		{"scientific notation", "1.5e2", 150.0},
		{"invalid string", "abc", 0.0},
		{"empty string", "", 0.0},
		{"mixed string", "3.14abc", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToFloat(tt.input)
			if result != tt.expected {
				t.Errorf("StringToFloat(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Tests for StringToBool
// ============================================================================

func TestStringToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"true lowercase", "true", true},
		{"true uppercase", "TRUE", true},
		{"true mixed case", "True", true},
		{"1 as true", "1", true},
		{"false lowercase", "false", false},
		{"false uppercase", "FALSE", false},
		{"false mixed case", "False", false},
		{"0 as false", "0", false},
		{"invalid string", "abc", false},
		{"empty string", "", false},
		{"yes", "yes", false}, // strconv.ParseBool doesn't accept "yes"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToBool(tt.input)
			if result != tt.expected {
				t.Errorf("StringToBool(%q) = %t, want %t", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Tests for StringToIntErr
// ============================================================================

func TestStringToIntErr(t *testing.T) {
	t.Run("valid conversion", func(t *testing.T) {
		result, err := StringToIntErr("42")
		if err != nil {
			t.Errorf("StringToIntErr(\"42\") unexpected error: %v", err)
		}
		if result != 42 {
			t.Errorf("StringToIntErr(\"42\") = %d, want 42", result)
		}
	})

	t.Run("negative number", func(t *testing.T) {
		result, err := StringToIntErr("-10")
		if err != nil {
			t.Errorf("StringToIntErr(\"-10\") unexpected error: %v", err)
		}
		if result != -10 {
			t.Errorf("StringToIntErr(\"-10\") = %d, want -10", result)
		}
	})

	t.Run("invalid conversion", func(t *testing.T) {
		_, err := StringToIntErr("abc")
		if err == nil {
			t.Error("StringToIntErr(\"abc\") expected error, got nil")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := StringToIntErr("")
		if err == nil {
			t.Error("StringToIntErr(\"\") expected error, got nil")
		}
	})

	t.Run("float string", func(t *testing.T) {
		_, err := StringToIntErr("3.14")
		if err == nil {
			t.Error("StringToIntErr(\"3.14\") expected error, got nil")
		}
	})
}

// ============================================================================
// Tests for StringToFloatErr
// ============================================================================

func TestStringToFloatErr(t *testing.T) {
	t.Run("valid conversion", func(t *testing.T) {
		result, err := StringToFloatErr("3.14")
		if err != nil {
			t.Errorf("StringToFloatErr(\"3.14\") unexpected error: %v", err)
		}
		if result != 3.14 {
			t.Errorf("StringToFloatErr(\"3.14\") = %f, want 3.14", result)
		}
	})

	t.Run("negative float", func(t *testing.T) {
		result, err := StringToFloatErr("-2.5")
		if err != nil {
			t.Errorf("StringToFloatErr(\"-2.5\") unexpected error: %v", err)
		}
		if result != -2.5 {
			t.Errorf("StringToFloatErr(\"-2.5\") = %f, want -2.5", result)
		}
	})

	t.Run("integer as float", func(t *testing.T) {
		result, err := StringToFloatErr("42")
		if err != nil {
			t.Errorf("StringToFloatErr(\"42\") unexpected error: %v", err)
		}
		if result != 42.0 {
			t.Errorf("StringToFloatErr(\"42\") = %f, want 42.0", result)
		}
	})

	t.Run("invalid conversion", func(t *testing.T) {
		_, err := StringToFloatErr("abc")
		if err == nil {
			t.Error("StringToFloatErr(\"abc\") expected error, got nil")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := StringToFloatErr("")
		if err == nil {
			t.Error("StringToFloatErr(\"\") expected error, got nil")
		}
	})
}

// ============================================================================
// Tests for StringToBoolErr
// ============================================================================

func TestStringToBoolErr(t *testing.T) {
	t.Run("true lowercase", func(t *testing.T) {
		result, err := StringToBoolErr("true")
		if err != nil {
			t.Errorf("StringToBoolErr(\"true\") unexpected error: %v", err)
		}
		if result != true {
			t.Errorf("StringToBoolErr(\"true\") = %t, want true", result)
		}
	})

	t.Run("false lowercase", func(t *testing.T) {
		result, err := StringToBoolErr("false")
		if err != nil {
			t.Errorf("StringToBoolErr(\"false\") unexpected error: %v", err)
		}
		if result != false {
			t.Errorf("StringToBoolErr(\"false\") = %t, want false", result)
		}
	})

	t.Run("1 as true", func(t *testing.T) {
		result, err := StringToBoolErr("1")
		if err != nil {
			t.Errorf("StringToBoolErr(\"1\") unexpected error: %v", err)
		}
		if result != true {
			t.Errorf("StringToBoolErr(\"1\") = %t, want true", result)
		}
	})

	t.Run("0 as false", func(t *testing.T) {
		result, err := StringToBoolErr("0")
		if err != nil {
			t.Errorf("StringToBoolErr(\"0\") unexpected error: %v", err)
		}
		if result != false {
			t.Errorf("StringToBoolErr(\"0\") = %t, want false", result)
		}
	})

	t.Run("invalid conversion", func(t *testing.T) {
		_, err := StringToBoolErr("abc")
		if err == nil {
			t.Error("StringToBoolErr(\"abc\") expected error, got nil")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := StringToBoolErr("")
		if err == nil {
			t.Error("StringToBoolErr(\"\") expected error, got nil")
		}
	})
}

// ============================================================================
// Tests for StringToIntOrDefault
// ============================================================================

func TestStringToIntOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultValue int
		expected     int
	}{
		{"valid conversion", "42", 100, 42},
		{"negative number", "-10", 100, -10},
		{"zero", "0", 100, 0},
		{"invalid string uses default", "abc", 100, 100},
		{"empty string uses default", "", 50, 50},
		{"float string uses default", "3.14", 75, 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToIntOrDefault(tt.input, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("StringToIntOrDefault(%q, %d) = %d, want %d", tt.input, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Tests for StringToFloatOrDefault
// ============================================================================

func TestStringToFloatOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultValue float64
		expected     float64
	}{
		{"valid conversion", "3.14", 100.0, 3.14},
		{"negative float", "-2.5", 100.0, -2.5},
		{"zero", "0", 100.0, 0.0},
		{"integer as float", "42", 100.0, 42.0},
		{"invalid string uses default", "abc", 100.0, 100.0},
		{"empty string uses default", "", 50.0, 50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToFloatOrDefault(tt.input, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("StringToFloatOrDefault(%q, %f) = %f, want %f", tt.input, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Tests for StringToBoolOrDefault
// ============================================================================

func TestStringToBoolOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultValue bool
		expected     bool
	}{
		{"true lowercase", "true", false, true},
		{"false lowercase", "false", true, false},
		{"1 as true", "1", false, true},
		{"0 as false", "0", true, false},
		{"invalid string uses default true", "abc", true, true},
		{"invalid string uses default false", "xyz", false, false},
		{"empty string uses default", "", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToBoolOrDefault(tt.input, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("StringToBoolOrDefault(%q, %t) = %t, want %t", tt.input, tt.defaultValue, result, tt.expected)
			}
		})
	}
}
