package executor

import (
	"redis-repo/internal/constant"
	"redis-repo/internal/data_structure"
	"time"
)

var dict *data_structure.Dict
var setStore map[string]data_structure.Set

func init() {
	dict = data_structure.NewDict()
	setStore = make(map[string]data_structure.Set)
}

// Clean some expired keys, follows Redis's solution
func CleanupExpiredKeys() {
	deleted, total := 0, 0
	startTime := time.Now().UnixMilli()

	dict.IterateExpiredKeys(func(key string, expiryTime uint64) bool {
		if dict.HasExpired(key) {
			dict.Delete(key)
			deleted++
		}
		total++

		// Check batches using a sample size, and stop the cleanup once the ratio of expired keys is within the acceptable range
		if total == constant.ActiveCleanupSampleSize {
			if float64(deleted/total) < constant.ActiveCleanupAcceptedExpiredProportion {
				return false // Stop iteration
			}

			// Reset variables to continue clean up
			total = 0
			deleted = 0
		}

		// Ensure the time for active clean up does not take a lot
		now := time.Now().UnixMilli()
		if now-startTime > constant.ActiveCleanupTimeLimit {
			return false // Stop iteration
		}

		return true // Continue iteration
	})
}
