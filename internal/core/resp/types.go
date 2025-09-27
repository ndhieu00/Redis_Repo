package resp

// RESP protocol constant
const (
	CarriageReturnByte = byte('\r')
	LineFeedByte       = byte('\n')
	CRLFString         = "\r\n"
)

// CRLFBytes represents the CRLF sequence as bytes
var CRLFBytes = []byte{CarriageReturnByte, LineFeedByte}

// Pre-encoded RESP nil value
var RespNil = []byte("$-1\r\n")

// RESP data types
type DataType struct {
	Name string
	Sign byte
}

var (
	IntegerType = DataType{
		Name: "integer",
		Sign: ':',
	}

	SimpleStringType = DataType{
		Name: "simple_string",
		Sign: '+',
	}

	BulkStringType = DataType{
		Name: "bulk_string",
		Sign: '$',
	}

	ErrorType = DataType{
		Name: "error",
		Sign: '-',
	}

	ArrayType = DataType{
		Name: "array",
		Sign: '*',
	}
)
