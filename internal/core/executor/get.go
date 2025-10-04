package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
)

func executeGet(args []string) []byte {
	if len(args) != 1 {
		return resp.Encode(fmt.Sprintf(constant.ErrWrongArgCount, "GET"))
	}
	key := args[0]
	if key == "" {
		return resp.Encode(constant.ErrEmptyKey)
	}

	vObject := dictStore.Get(key)
	if vObject == nil {
		return []byte(constant.RespNil)
	}

	return resp.Encode(vObject.Value)
}
