package executor

import (
	"redis-repo/internal/core/resp"
)

func executeDel(args []string) []byte {
	count := 0
	for _, key := range args {
		if exist := dict.Delete(key); exist {
			count++
		}
	}
	return resp.Encode(count)
}
