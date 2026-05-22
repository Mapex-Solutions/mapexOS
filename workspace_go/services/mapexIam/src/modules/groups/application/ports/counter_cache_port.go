package ports

// CounterCachePort is the driven port for building counter cache keys for groups.
// The application layer depends on this port to avoid coupling to the
// infrastructure/cache/redis package (Hexagonal Architecture).
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: CounterCachePort (this interface)
//   - Adapter: infrastructure/cache/redis.CounterCacheAdapter
//
// Only methods actually consumed by the application layer are exposed here.
type CounterCachePort interface {
	// BuildKey builds the full counter cache key for the given orgId.
	BuildKey(orgId string) string
}
