package constants

import "time"

// Counter cache configuration
const (
	// CounterCacheTTL is the cache duration for counter values (6 hours)
	CounterCacheTTL = 6 * time.Hour
)
