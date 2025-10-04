package constant

// RESP Protocol Response Constants
const (
	RespOk  = "+OK\r\n"
	RespNil = "$-1\r\n"
)

// TTL Response Constants
const (
	TtlKeyNotExist      = ":-2\r\n"
	TtlKeyExistNoExpire = ":-1\r\n"
)

// Error Messages
const (
	ErrWrongArgCount = "ERR wrong number of arguments for '%s' command"
	ErrEmptyKey      = "ERR empty key"
	ErrInvalidTime   = "ERR invalid time"
)
