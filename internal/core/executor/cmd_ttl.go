package executor

import (
	"fmt"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
	"time"
)

func cmdTTL(args []string) []byte {
	if len(args) != 1 {
		return []byte(fmt.Sprintf(constant.ErrWrongArgCount, "TTL"))
	}
	key := args[0]
	if key == "" {
		return []byte(constant.ErrEmptyKey)
	}

	vObj := dict.Get(key)
	if vObj == nil {
		return []byte(constant.TtlKeyNotExist)
	}

	expiryTime, exist := dict.GetExpiryTime(key)
	now := uint64(time.Now().UnixMilli())

	if !exist {
		return []byte(constant.TtlKeyExistNoExpire)
	}

	if expiryTime < now {
		dict.Delete(key)
		return []byte(constant.TtlKeyNotExist)
	}

	remainMs := expiryTime - now
	return resp.Encode(remainMs / 1000)
}
