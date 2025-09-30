package executor

import (
	"errors"
	"fmt"
	"log"
	"redis-repo/internal/constant"
	"redis-repo/internal/core/resp"
	"strconv"
	"strings"
	"time"
)

// Support SET key value [EX seconds|PX milliseconds|EXAT timestamp|PXAT milliseconds-timestamp]
func executeSet(args []string) []byte {
	if len(args) != 2 && len(args) != 4 {
		return resp.Encode(errors.New("ERR wrong number of arguments for 'SET'"))
	}

	var expiryTimeMs uint64
	var err error
	// Has expiry time
	if len(args) == 4 {
		typeExp, timeStr := args[2], args[3]

		switch strings.ToUpper(typeExp) {
		case "EX": // TimeStr in seconds
			expiryTimeMs, err = expiryTimeMsFromEX(timeStr)
		case "PX": // TimeStr in milliseconds
			expiryTimeMs, err = expiryTimeMsFromPX(timeStr)
		case "EXAT": // TimeStr in timestamp
			expiryTimeMs, err = expiryTimeMsFromEXAT(timeStr)
		case "PXAT": // TimeStr in milliseconds-timestamp
			expiryTimeMs, err = expiryTimeMsFromPXAT(timeStr)
		default:
			return resp.Encode(errors.New("invalid type of expiry time"))
		}

		if err != nil {
			log.Println(err)
			return resp.Encode(errors.New("invalid time"))
		}
	}

	dictStore.Set(args[0], args[1], expiryTimeMs)

	return constant.RespOk
}

func expiryTimeMsFromEX(timeStr string) (uint64, error) {
	ttlSec, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return 0, err
	}
	if ttlSec <= 0 {
		return 0, fmt.Errorf("expiryTimeMsFromEX: invalid ttl %q, must be >= 0", timeStr)
	}

	expiryTimeMs := uint64(time.Now().Unix()+ttlSec) * 1000 // to milliseconds

	return expiryTimeMs, nil
}

func expiryTimeMsFromPX(timeStr string) (uint64, error) {
	ttlMs, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return 0, err
	}
	if ttlMs <= 0 {
		return 0, fmt.Errorf("expiryTimeMsFromPX: invalid ttl %q, must be >= 0", timeStr)
	}

	expiryTimeMs := uint64(time.Now().UnixMilli() + ttlMs)

	return expiryTimeMs, nil
}

func expiryTimeMsFromEXAT(timeStr string) (uint64, error) {
	expiryTimeSec, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return 0, err
	}
	if expiryTimeSec <= time.Now().Unix() {
		return 0, fmt.Errorf("expiryTimeMsFromEXAT: invalid timestamp %q, must be >= now", timeStr)
	}

	expiryTimeMs := uint64(expiryTimeSec) * 1000

	return expiryTimeMs, nil
}

func expiryTimeMsFromPXAT(timeStr string) (uint64, error) {
	expiryTimeMs, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return 0, err
	}
	if expiryTimeMs <= time.Now().UnixMilli() {
		return 0, fmt.Errorf("expiryTimeMsFromEXAT: invalid timestamp (milliseconds) %q, must be >= now", timeStr)
	}

	return uint64(expiryTimeMs), nil
}
