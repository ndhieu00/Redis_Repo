package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
)

func cmdSINTER(args []string) []byte {
	if len(args) == 0 {
		return []byte(fmt.Sprintf(constant.ErrWrongArgCount, "SINTER"))
	}

	smallestKey := args[0]
	for i := 1; i < len(args); i++ {
		if _, exists := setStore[args[i]]; !exists {
			return resp.Encode([]any{})
		}
		if len(setStore[args[i]]) < len(setStore[smallestKey]) {
			smallestKey = args[i]
		}
	}

	result := make([]any, 0)

	// Check each member of the smallest set against all other sets
	for member := range setStore[smallestKey] {
		validMember := true
		for _, key := range args {
			if key == smallestKey {
				continue
			}

			if setStore[key].IsMember(member) == 0 { // Member not found
				validMember = false
				break
			}
		}

		if validMember {
			result = append(result, member)
		}
	}

	return resp.Encode(result)
}
