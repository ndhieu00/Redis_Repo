package executor

import (
	"redis-repo/internal/constant"
	"redis-repo/internal/data_structure"
	"time"
)

var dict *data_structure.Dict

func init() {
	dict = data_structure.NewDict()
}

// Clean some expired keys, follows Redis's solution
func CleanupExpiredKeys() {
	deleted, total := 0, 0
	startTime := time.Now().UnixMilli()

	for key := range dict.GetExpiredDictStore() {
		if dict.HasExpired(key) {
			dict.Delete(key)
			deleted++
		}
		total++

		// Check batches using a sample size, and stop the cleanup once the ratio of expired keys is within the acceptable range
		if total == constant.ActiveCleanupSampleSize {
			if float64(deleted/total) < constant.ActiveCleanupAcceptedExpiredProportion {
				break
			}

			// Reset variables to continue clean up
			total = 0
			deleted = 0
		}

		// Ensure the time for active clean up does not take a lot
		now := time.Now().UnixMilli()
		if now-startTime > constant.ActiveCleanupTimeLimit {
			break
		}
	}
}
