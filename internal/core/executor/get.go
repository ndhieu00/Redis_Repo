package executor

import (
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
)

func executeGet(args []string) []byte {
	if len(args) != 1 {
		return resp.Encode("ERR wrong number of arguments for command")
	}
	key := args[0]
	vObject := dictStore.Get(key)
	if vObject == nil {
		return constant.RespNil
	}

	return resp.Encode(vObject.Value)
}
