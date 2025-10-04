package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
	"time"
)

func executeTTL(args []string) []byte {
	if len(args) != 1 {
		return resp.Encode(fmt.Sprintf(constant.ErrWrongArgCount, "TTL"))
	}
	key := args[0]
	if key == "" {
		return resp.Encode(constant.ErrEmptyKey)
	}

	vObj := dictStore.Get(key)
	if vObj == nil {
		return []byte(constant.TtlKeyNotExist)
	}

	expiryTime, exist := dictStore.GetExpiryTime(key)
	now := uint64(time.Now().UnixMilli())

	if !exist {
		return []byte(constant.TtlKeyExistNoExpire)
	}

	if expiryTime < now {
		dictStore.Delete(key)
		return []byte(constant.TtlKeyNotExist)
	}

	remainMs := expiryTime - now
	return resp.Encode(remainMs / 1000)
}
