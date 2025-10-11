package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
)

func cmdGET(args []string) []byte {
	if len(args) != 1 {
		return []byte(fmt.Sprintf(constant.ErrWrongArgCount, "GET"))
	}
	key := args[0]
	if key == "" {
		return []byte(constant.ErrEmptyKey)
	}

	vObject := dict.Get(key)
	if vObject == nil {
		return []byte(constant.RespNil)
	}

	return resp.Encode(vObject.Value)
}
