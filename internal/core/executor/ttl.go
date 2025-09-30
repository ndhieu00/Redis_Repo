package executor

import (
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
	"time"
)

func executeTTL(args []string) []byte {
	if len(args) != 1 {
		return resp.Encode("ERR wrong number of arguments for command")
	}
	key := args[0]

	vObj := dictStore.Get(key)
	if vObj == nil {
		return constant.TtlKeyNotExist
	}

	expiryTime, exist := dictStore.GetExpiryTime(key)
	now := uint64(time.Now().UnixMilli())

	if !exist {
		return constant.TtlKeyExistNoExpire
	}

	if expiryTime < now {
		dictStore.Delete(key)
		return constant.TtlKeyNotExist
	}

	remainMs := expiryTime - now
	return resp.Encode(remainMs / 1000)
}
