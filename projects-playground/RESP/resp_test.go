package resp

import (
	"errors"
	"fmt"
	"testing"
)

func TestEncodeInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"positive integer", 42, ":42\r\n"},
		{"negative integer", -42, ":-42\r\n"},
		{"zero", 0, ":0\r\n"},
		{"int64 max", int64(9223372036854775807), ":9223372036854775807\r\n"},
		{"int64 min", int64(-9223372036854775808), ":-9223372036854775808\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
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
	}{
		{"valid string", "hello", "+hello\r\n"},
		{"empty string", "", "+\r\n"},
		{"string with spaces", "hello world", "+hello world\r\n"},
		{"string with CRLF", "hello\r\nworld", "$-1\r\n"},
		{"invalid type", 42, "$-1\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeSimpleString(tt.input)
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
	}{
		{"valid string", "hello", "$5\r\nhello\r\n"},
		{"empty string", "", "$0\r\n\r\n"},
		{"string with CRLF", "hello\r\nworld", "$12\r\nhello\r\nworld\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
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
	}{
		{"valid error", errors.New("test error"), "-test error\r\n"},
		{"empty error", errors.New(""), "-\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
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
	}{
		{"empty array", []any{}, "*0\r\n"},
		{"array with strings", []any{"hello", "world"}, "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"},
		{"array with mixed types", []any{"hello", 42, errors.New("error")}, "*3\r\n$5\r\nhello\r\n:42\r\n-error\r\n"},
		{"nested array", []any{[]any{"nested"}}, "*1\r\n*1\r\n$6\r\nnested\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
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
			result := Encode(tt.input)
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
		{"array one string", "*2\r\n$5\r\nhello\r\n$1\r\nh\r\n", []any{"hello", "h"}, false},
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
		name     string
		input    any
		expected any
	}{
		{"integer", (42), int64(42)},
		{"string", "hello world", "hello world"},
		{"error", errors.New("test error"), "test error"},
		{"array", []any{"hello", 42, errors.New("error")}, []any{"hello", int64(42), "error"}},
		{"empty array", []any{}, []any{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Encode(tt.input)

			decoded, err := Decode(encoded)
			if err != nil {
				t.Errorf("Decoding failed: %v", err)
				return
			}

			if tt.name == "array" || tt.name == "empty array" {
				expectedArr := tt.expected.([]any)
				decodedArr, ok := decoded.([]any)
				if !ok {
					t.Errorf("Expected []any, got %T", decoded)
					return
				}
				if len(expectedArr) != len(decodedArr) {
					t.Errorf("Array length mismatch: expected %d, got %d", len(expectedArr), len(decodedArr))
					return
				}
				for i, expected := range expectedArr {
					// Check types first
					expectedType := fmt.Sprintf("%T", expected)
					decodedType := fmt.Sprintf("%T", decodedArr[i])
					if expectedType != decodedType {
						t.Errorf("Array element type mismatch at index %d: expected %s, got %s", i, expectedType, decodedType)
						continue
					}
					// If types match, compare values
					expectedStr := fmt.Sprintf("%v", expected)
					decodedStr := fmt.Sprintf("%v", decodedArr[i])
					if expectedStr != decodedStr {
						t.Errorf("Array element value mismatch at index %d: expected %v, got %v", i, expectedStr, decodedStr)
					}
				}
			} else {
				// Check types first
				expectedType := fmt.Sprintf("%T", tt.expected)
				decodedType := fmt.Sprintf("%T", decoded)
				if expectedType != decodedType {
					t.Errorf("Type mismatch: expected %s, got %s", expectedType, decodedType)
					return
				}
				// If types match, compare values
				expectedStr := fmt.Sprintf("%v", tt.expected)
				decodedStr := fmt.Sprintf("%v", decoded)
				if expectedStr != decodedStr {
					t.Errorf("Value mismatch: expected %v, got %v", expectedStr, decodedStr)
				}
			}
		})
	}
}

func TestEncodeErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"unsupported type", map[string]string{"key": "value"}, "$-1\r\n"},
		{"nil input", nil, "$-1\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}
