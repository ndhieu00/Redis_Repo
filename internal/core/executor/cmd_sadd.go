package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
	"redis-repo/internal/data_structure"
)

func cmdSADD(args []string) []byte {
	if len(args) < 2 {
		return []byte(fmt.Sprintf(constant.ErrWrongArgCount, "SADD"))
	}
	keySet := args[0]
	members := args[1:]

	set, exists := setStore[keySet]
	if !exists {
		setStore[keySet] = data_structure.NewSet(members)
		return resp.Encode(len(members))
	}

	added := set.Add(members)

	return resp.Encode(added)
}
