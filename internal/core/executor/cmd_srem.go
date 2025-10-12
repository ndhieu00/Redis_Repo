package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
)

func cmdSREM(args []string) []byte {
	if len(args) < 2 {
		return []byte(fmt.Sprintf(constant.ErrWrongArgCount, "SREM"))
	}
	keySet := args[0]
	members := args[1:]

	set, exists := setStore[keySet]
	if !exists {
		return resp.Encode(0) // Nothing to remove
	}

	removed := set.Remove(members)
	return resp.Encode(removed)
}
