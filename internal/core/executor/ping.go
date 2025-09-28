package executor

import (
	"errors"
	"redis-repo/internal/core/resp"
)

// ExecutePing handles the PING command
func executePing(args []string) []byte {
	switch len(args) {
	case 0:
		return resp.EncodeSimpleString("PONG")
	case 1:
		return resp.Encode(args[0])
	default:
		return resp.Encode(errors.New("wrong number of arguments for 'ping' command"))
	}
}
