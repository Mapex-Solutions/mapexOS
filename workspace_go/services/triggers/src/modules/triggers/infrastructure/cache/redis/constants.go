// Package redis holds Redis-specific constants for the triggers module's
// application cache.
//
// Scope: infrastructure/cache/redis — keys and key-prefix helpers that are
// technology-specific and intentionally not part of the domain or
// cross-service contract.
//
// Roles covered:
//   - Per-trigger cache (single trigger lookup by ID)
//   - Per-org counter cache for trigger totals
package redis

const (
	// TriggerCacheKeyPrefix defines the prefix for trigger cache keys.
	// Format: TRIGGER:{triggerId}
	TriggerCacheKeyPrefix = "TRIGGER"

	// CounterCacheKeyPrefix is the Redis key prefix for the trigger counter cache.
	// Format: counter:triggers:{orgId}
	CounterCacheKeyPrefix = "counter:triggers:"
)

// BuildTriggerCacheKey builds the Redis key for a single trigger cache entry.
// Key format: TRIGGER:{triggerId}
func BuildTriggerCacheKey(triggerId string) string {
	return TriggerCacheKeyPrefix + ":" + triggerId
}

// BuildCounterCacheKey builds the Redis key for the per-org trigger counter.
// Key format: counter:triggers:{orgId}
func BuildCounterCacheKey(orgId string) string {
	return CounterCacheKeyPrefix + orgId
}
