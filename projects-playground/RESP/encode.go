package resp

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

// convertToInt64 converts any integer type to int64
func convertToInt64(data any) int64 {
	switch v := data.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	default:
		return 0
	}
}

// encodeInteger encodes an int64 value to RESP format
func encodeInteger(intVal int64) ([]byte, error) {
	return []byte(fmt.Sprintf("%c%d%s", IntegerType.Sign, intVal, CRLFString)), nil
}

// encodeSimpleString encodes a string to RESP format
func encodeSimpleString(str string) ([]byte, error) {
	if strings.Contains(str, CRLFString) {
		return nil, fmt.Errorf("simple string cannot contain CRLF")
	}

	return []byte(fmt.Sprintf("%c%s%s", SimpleStringType.Sign, str, CRLFString)), nil
}

// encodeBulkString encodes a string to RESP format
func encodeBulkString(str string) ([]byte, error) {
	return []byte(fmt.Sprintf("%c%d%s%s%s",
		BulkStringType.Sign,
		len(str),
		CRLFString,
		str,
		CRLFString)), nil
}

// encodeError encodes an error to RESP format
func encodeError(err error) ([]byte, error) {
	return []byte(fmt.Sprintf("%c%s%s", ErrorType.Sign, err.Error(), CRLFString)), nil
}

// encodeArray encodes an array to RESP format
func encodeArray(arr []any) ([]byte, error) {
	var buf bytes.Buffer
	for _, item := range arr {
		encoded, err := encode(item)
		if err != nil {
			return nil, fmt.Errorf("failed to encode array element: %w", err)
		}
		buf.Write(encoded)
	}

	return []byte(fmt.Sprintf("%c%d%s%s",
		ArrayType.Sign,
		len(arr),
		CRLFString,
		buf.String())), nil
}

func encode(data any) ([]byte, error) {
	switch v := data.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return encodeInteger(convertToInt64(v))
	case string:
		return encodeBulkString(v)
	case error:
		return encodeError(v)
	case []any:
		return encodeArray(v)
	default:
		return nil, fmt.Errorf("unsupported type %T", data)
	}
}

// Encode encodes data to RESP format
//
// Accepts these types:
// - int, int8, int16, int32, int64, uint8, uint16, uint32, uint64: encoded as integer
// - string: encoded as bulk string (e.g., "hello" -> $5\r\nhello\r\n)
// - error: encoded as error (e.g., errors.New("msg") -> -msg\r\n)
// - []any: encoded as array (e.g., []any{"hello", 42} -> *2\r\n$5\r\nhello\r\n:42\r\n)
func Encode(data any) []byte {
	result, err := encode(data)
	if err != nil {
		log.Printf("error encoding data: %v", err)
		return RespNil
	}
	return result
}

// EncodeSimpleString encodes a string as a simple string (not bulk string)
func EncodeSimpleString(data any) []byte {
	str, ok := data.(string)
	if !ok {
		log.Printf("RESP encode: expected string, got %T", data)
		return RespNil
	}
	result, err := encodeSimpleString(str)
	if err != nil {
		log.Printf("error encoding data: %v", err)
		return RespNil
	}
	return result
}
