package ports

// TriggerCacheKeyBuilderPort is the driven port for producing Redis cache keys
// used by the triggers module.
//
// Implementations live under `infrastructure/cache/redis/` and encode the
// actual key format. The application layer depends ONLY on this interface,
// never on the concrete key strings.
//
// Methods:
//   - TriggerKey: Redis key for a single trigger cache entry (by triggerId).
//   - CounterKey: Redis key for the per-org trigger counter cache (by orgId).
type TriggerCacheKeyBuilderPort interface {
	// TriggerKey returns the cache key for a trigger identified by triggerId.
	TriggerKey(triggerId string) string

	// CounterKey returns the cache key for the per-org trigger counter.
	CounterKey(orgId string) string
}
