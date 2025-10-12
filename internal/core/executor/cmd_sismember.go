package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
)

func cmdSMISMEMBER(args []string) []byte {
	if len(args) < 2 {
		return []byte(fmt.Sprintf(constant.ErrWrongArgCount, "SMISMEMBER"))
	}

	keySet := args[0]
	members := args[1:]
	ans := make([]any, len(members))

	set, exists := setStore[keySet]
	if !exists {
		// Initialize all elements to 0 (member not found)
		for i := range ans {
			ans[i] = 0
		}
		return resp.Encode(ans)
	}

	for i, member := range members {
		ans[i] = set.IsMember(member)
	}

	return resp.Encode(ans)
}
