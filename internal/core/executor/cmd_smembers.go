package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
)

func cmdSMEMBERS(args []string) []byte {
	if len(args) != 1 {
		return []byte(fmt.Sprintf(constant.ErrWrongArgCount, "SMEMBERS"))
	}

	keySet := args[0]
	set, exists := setStore[keySet]
	if !exists {
		return resp.Encode([]any{}) // Return empty array
	}

	ans := make([]any, 0, len(set))
	for member := range set {
		ans = append(ans, member)
	}

	return resp.Encode(ans)
}
