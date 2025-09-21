package resp

import (
	"errors"
	"testing"
)

func TestEncodeInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		hasError bool
	}{
		{"positive integer", 42, ":42\r\n", false},
		{"negative integer", -42, ":-42\r\n", false},
		{"zero", 0, ":0\r\n", false},
		{"int64 max", int64(9223372036854775807), ":9223372036854775807\r\n", false},
		{"int64 min", int64(-9223372036854775808), ":-9223372036854775808\r\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestEncodeSimpleString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		hasError bool
	}{
		{"valid string", "hello", "+hello\r\n", false},
		{"empty string", "", "+\r\n", false},
		{"string with spaces", "hello world", "+hello world\r\n", false},
		{"string with CRLF", "hello\r\nworld", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EncodeSimpleString(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestEncodeBulkString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		hasError bool
	}{
		{"valid string", "hello", "$5\r\nhello\r\n", false},
		{"empty string", "", "$0\r\n\r\n", false},
		{"string with CRLF", "hello\r\nworld", "$12\r\nhello\r\nworld\r\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestEncodeError(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		hasError bool
	}{
		{"valid error", errors.New("test error"), "-test error\r\n", false},
		{"empty error", errors.New(""), "-\r\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestEncodeArray(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		hasError bool
	}{
		{"empty array", []any{}, "*0\r\n", false},
		{"array with strings", []any{"hello", "world"}, "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", false},
		{"array with mixed types", []any{"hello", 42, errors.New("error")}, "*3\r\n$5\r\nhello\r\n:42\r\n-error\r\n", false},
		{"nested array", []any{[]any{"nested"}}, "*1\r\n*1\r\n$6\r\nnested\r\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestAutoDetectType(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"string becomes bulk string", "hello", "$5\r\nhello\r\n"},
		{"error becomes error", errors.New("test"), "-test\r\n"},
		{"array becomes array", []any{"test"}, "*1\r\n$4\r\ntest\r\n"},
		{"integer becomes integer", 42, ":42\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestDecodeInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		hasError bool
	}{
		{"positive integer", ":42\r\n", 42, false},
		{"negative integer", ":-42\r\n", -42, false},
		{"zero", ":0\r\n", 0, false},
		{"large number", ":9223372036854775807\r\n", 9223372036854775807, false},
		{"invalid format", ":abc\r\n", 0, true},
		{"missing CRLF", ":42", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decode([]byte(tt.input))
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %d, got %v", tt.expected, result)
			}
		})
	}
}

func TestDecodeSimpleString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"valid string", "+hello\r\n", "hello", false},
		{"empty string", "+\r\n", "", false},
		{"string with spaces", "+hello world\r\n", "hello world", false},
		{"missing CRLF", "+hello", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decode([]byte(tt.input))
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDecodeBulkString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
		hasError bool
	}{
		{"valid string", "$5\r\nhello\r\n", "hello", false},
		{"empty string", "$0\r\n\r\n", "", false},
		{"nil bulk string", "$-1\r\n", nil, false},
		{"string with CRLF", "$12\r\nhello\r\nworld\r\n", "hello\r\nworld", false},
		{"invalid length", "$abc\r\nhello\r\n", nil, true},
		{"insufficient data", "$10\r\nhello\r\n", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decode([]byte(tt.input))
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDecodeError(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"valid error", "-test error\r\n", "test error", false},
		{"empty error", "-\r\n", "", false},
		{"missing CRLF", "-test error", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decode([]byte(tt.input))
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDecodeArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []any
		hasError bool
	}{
		{"empty array", "*0\r\n", []any{}, false},
		{"array with strings", "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", []any{"hello", "world"}, false},
		{"nil array", "*-1\r\n", nil, false},
		{"invalid length", "*abc\r\n", nil, true},
		{"insufficient elements", "*2\r\n$5\r\nhello\r\n", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decode([]byte(tt.input))
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Compare arrays
			if tt.expected == nil && result != nil {
				t.Errorf("Expected nil, got %v", result)
			} else if tt.expected != nil {
				arr, ok := result.([]any)
				if !ok {
					t.Errorf("Expected []any, got %T", result)
					return
				}
				if len(arr) != len(tt.expected) {
					t.Errorf("Expected length %d, got %d", len(tt.expected), len(arr))
					return
				}
				for i, v := range arr {
					if v != tt.expected[i] {
						t.Errorf("Expected [%d] = %v, got %v", i, tt.expected[i], v)
					}
				}
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{"integer", int64(42)},
		{"string", "hello world"},
		{"error", errors.New("test error")},
		{"array", []any{"hello", 42, errors.New("error")}},
		{"empty array", []any{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded, err := Encode(tt.input)
			if err != nil {
				t.Errorf("Encoding failed: %v", err)
				return
			}

			// Decode
			decoded, err := Decode(encoded)
			if err != nil {
				t.Errorf("Decoding failed: %v", err)
				return
			}

			// Compare
			if tt.name == "array" {
				// Special handling for arrays
				expectedArr := tt.input.([]any)
				decodedArr, ok := decoded.([]any)
				if !ok {
					t.Errorf("Expected []any, got %T", decoded)
					return
				}
				if len(expectedArr) != len(decodedArr) {
					t.Errorf("Array length mismatch: expected %d, got %d", len(expectedArr), len(decodedArr))
					return
				}
				// Note: This is a simplified comparison - in practice you'd want more sophisticated comparison
			} else {
				if decoded != tt.input {
					t.Errorf("Round trip failed: expected %v, got %v", tt.input, decoded)
				}
			}
		})
	}
}
