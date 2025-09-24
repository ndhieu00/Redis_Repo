package resp

import (
	"errors"
	"fmt"
	"log"
)

// DecodingError represents an error that occurred during decoding
type DecodingError struct {
	Position int
	Data     []byte
	Err      error
}

func (e *DecodingError) Error() string {
	return fmt.Sprintf("decoding error at position %d: %v (data: %q)", e.Position, e.Err, e.Data)
}

// DecodeResult represents the result of decoding a RESP value
type DecodeResult struct {
	Value  any
	Length int // Total length of data consumed
}

// readInteger decodes an integer from RESP format
// Example: :-5\r\n => -5
func readInteger(data []byte) (*DecodeResult, error) {
	if len(data) < 3 {
		return nil, &DecodingError{Position: 0, Data: data, Err: errors.New("insufficient data for integer")}
	}

	pos := 1
	var sign int64 = 1

	if pos < len(data) {
		switch data[pos] {
		case '-':
			sign = -1
			pos++
		case '+':
			pos++
		}
	}

	number, length, err := extractNumber(data[pos:])
	if err != nil {
		return nil, &DecodingError{Position: pos, Data: data, Err: err}
	}

	return &DecodeResult{
		Value:  number * sign,
		Length: pos + length,
	}, nil
}

// readError decodes an error from RESP format
// Example: -Example Error\r\n => Example Error
func readError(data []byte) (*DecodeResult, error) {
	if len(data) < 3 {
		return nil, &DecodingError{Position: 0, Data: data, Err: errors.New("insufficient data for error")}
	}

	pos := 1
	for pos < len(data) && data[pos] != CarriageReturnByte {
		pos++
	}

	if pos >= len(data) || pos+1 >= len(data) || data[pos+1] != LineFeedByte {
		return nil, &DecodingError{Position: pos, Data: data, Err: errors.New("missing CRLF terminator")}
	}

	return &DecodeResult{
		Value:  string(data[1:pos]),
		Length: pos + 2,
	}, nil
}

// readSimpleString decodes a simple string from RESP format
// Example: +Hello world\r\n => Hello world
func readSimpleString(data []byte) (*DecodeResult, error) {
	if len(data) < 3 {
		return nil, &DecodingError{Position: 0, Data: data, Err: errors.New("insufficient data for simple string")}
	}

	pos := 1
	for pos < len(data) && data[pos] != CarriageReturnByte {
		pos++
	}

	if pos >= len(data) || pos+1 >= len(data) || data[pos+1] != LineFeedByte {
		return nil, &DecodingError{Position: pos, Data: data, Err: errors.New("missing CRLF terminator")}
	}

	return &DecodeResult{
		Value:  string(data[1:pos]),
		Length: pos + 2,
	}, nil
}

// readBulkString decodes a bulk string from RESP format
// Example: $9\r\nhello\r\n\r\n => hello\r\n
func readBulkString(data []byte) (*DecodeResult, error) {
	if len(data) < 5 {
		return nil, &DecodingError{Position: 0, Data: data, Err: errors.New("insufficient data for bulk string")}
	}

	pos := 1
	length, lengthConsumed, err := extractNumber(data[pos:])
	if err != nil {
		return nil, &DecodingError{Position: pos, Data: data, Err: err}
	}
	pos += lengthConsumed

	// Handle nil bulk string
	if length == -1 {
		return &DecodeResult{
			Value:  nil,
			Length: pos,
		}, nil
	}

	// Check if we have enough data for the string
	if pos+int(length)+2 > len(data) {
		return nil, &DecodingError{Position: pos, Data: data, Err: errors.New("insufficient data for bulk string content")}
	}

	// Verify CRLF terminator
	if data[pos+int(length)] != CarriageReturnByte || data[pos+int(length)+1] != LineFeedByte {
		return nil, &DecodingError{Position: pos + int(length), Data: data, Err: errors.New("missing CRLF terminator")}
	}

	return &DecodeResult{
		Value:  string(data[pos : pos+int(length)]),
		Length: pos + int(length) + 2,
	}, nil
}

// readArray decodes an array from RESP format
// Example: *3\r\n$5\r\nhello\r\n$5\r\nworld\r\n:+25\r\n => ['hello', 'world', 25]
func readArray(data []byte) (*DecodeResult, error) {
	if len(data) < 4 {
		return nil, &DecodingError{Position: 0, Data: data, Err: errors.New("insufficient data for array")}
	}

	pos := 1
	length, lengthConsumed, err := extractNumber(data[pos:])
	if err != nil {
		return nil, &DecodingError{Position: pos, Data: data, Err: err}
	}
	pos += lengthConsumed

	// Handle nil array
	if length == -1 {
		return &DecodeResult{
			Value:  nil,
			Length: pos,
		}, nil
	}

	arrResult := make([]any, length)
	for i := range arrResult {
		result, err := decode(data[pos:])
		if err != nil {
			return nil, &DecodingError{Position: pos, Data: data, Err: fmt.Errorf("failed to decode array element %d: %w, current result arrResult: %s", i, err, arrResult)}
		}
		arrResult[i] = result.Value
		pos += result.Length
	}

	return &DecodeResult{
		Value:  arrResult,
		Length: pos,
	}, nil
}

// extractNumber extracts a number from RESP format and total length consumed
// Example: 5\r\n => (5, 3), -1\r\n => (-1, 4)
func extractNumber(data []byte) (int64, int, error) {
	if len(data) == 0 {
		return 0, 0, errors.New("empty data")
	}

	pos := 0
	var sign int64 = 1

	// Handle negative numbers
	if pos < len(data) && data[pos] == '-' {
		sign = -1
		pos++
	}

	result := int64(0)
	hasDigit := false

	for pos < len(data) && data[pos] != CarriageReturnByte {
		if data[pos] < '0' || data[pos] > '9' {
			return 0, 0, fmt.Errorf("invalid character in number: %c", data[pos])
		}
		result = result*10 + int64(data[pos]-'0')
		hasDigit = true
		pos++
	}

	if !hasDigit {
		return 0, 0, errors.New("no digits found in number")
	}

	if pos >= len(data) || pos+1 >= len(data) || data[pos+1] != LineFeedByte {
		return 0, 0, errors.New("missing CRLF terminator")
	}

	return result * sign, pos + 2, nil
}

// decode decodes a single RESP value from the given data
func decode(data []byte) (*DecodeResult, error) {
	if len(data) == 0 {
		return nil, &DecodingError{Position: 0, Data: data, Err: errors.New("empty data")}
	}

	sign := data[0]
	switch sign {
	case IntegerType.Sign:
		return readInteger(data)
	case SimpleStringType.Sign:
		return readSimpleString(data)
	case ErrorType.Sign:
		return readError(data)
	case BulkStringType.Sign:
		return readBulkString(data)
	case ArrayType.Sign:
		return readArray(data)
	default:
		// Log minimal info for debugging
		log.Printf("RESP decode: unsupported type '%c' at position 0", sign)
		// Find the end of the unsupported type for better error reporting
		pos := 0
		for pos < len(data) && data[pos] != CarriageReturnByte {
			pos++
		}
		return nil, &DecodingError{
			Position: 0,
			Data:     data,
			Err:      fmt.Errorf("unsupported RESP type: %c", sign),
		}
	}
}

// Decode decodes RESP data and returns the decoded value
// Returns one of these types:
// - int64: for integers (e.g., :42\r\n)
// - string: for simple strings, bulk strings, and errors (e.g., +hello\r\n, $5\r\nhello\r\n, -ERR message\r\n)
// - []any: for arrays (e.g., *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n)
// - nil: for nil bulk strings or nil arrays (e.g., $-1\r\n, *-1\r\n)
func Decode(data []byte) (any, error) {
	result, err := decode(data)
	if err != nil {
		return nil, err
	}
	return result.Value, nil
}
