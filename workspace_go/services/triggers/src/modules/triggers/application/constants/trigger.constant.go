package constants

import "time"

// Cache TTL configuration for Trigger (application behavior — not tech-specific).
const (
	// TriggerCacheTTL defines the time-to-live for Trigger cache entries.
	TriggerCacheTTL = 60 * time.Minute

	// CounterCacheTTL is the cache duration for counter values (6 hours).
	CounterCacheTTL = 6 * time.Hour
)
