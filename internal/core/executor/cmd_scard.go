package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
)

func cmdSCARD(args []string) []byte {
	if len(args) != 1 {
		return []byte(fmt.Sprintf(constant.ErrWrongArgCount, "SCARD"))
	}

	keySet := args[0]
	set, exists := setStore[keySet]
	if !exists {
		return resp.Encode(0) // Return 0 for non-existing set
	}

	return resp.Encode(len(set))
}
