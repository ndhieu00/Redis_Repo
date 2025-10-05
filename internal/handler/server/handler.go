package server

import "redis-repo/internal/core/executor"

// HandleSystemCleanup handles system-level cleanup operations
func HandleSystemCleanup() {
	executor.CleanupExpiredKeys()
}
