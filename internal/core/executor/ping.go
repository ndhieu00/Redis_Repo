package executor

import "redis-repo/internal/core/resp"

// ExecutePing handles the PING command
func executePing() []byte {
	return resp.EncodeSimpleString("PONG")
}
