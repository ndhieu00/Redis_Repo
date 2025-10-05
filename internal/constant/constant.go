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
	ErrWrongArgCount = "-ERR wrong number of arguments for '%s' command\r\n"
	ErrEmptyKey      = "-ERR empty key\r\n"
	ErrInvalidTime   = "-ERR invalid time\r\n"
)

// Active Cleanup
const (
	ActiveCleanupFrequency                 = 100 // 100ms
	ActiveCleanupTimeLimit                 = 500 // 500ms
	ActiveCleanupSampleSize                = 20
	ActiveCleanupAcceptedExpiredProportion = 0.1 // The percentage of expired keys in the sample size is acceptable.
)
