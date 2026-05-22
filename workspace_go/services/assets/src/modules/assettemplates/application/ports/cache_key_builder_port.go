package ports

// CacheKeyBuilderPort abstracts Redis key construction for the
// assettemplates module's application cache. It exposes ONLY the
// key-formatting helper the application layer needs, without leaking
// Redis-specific constants (key prefixes) into the application package.
//
// The application layer depends on this interface, not on any
// infrastructure/cache/redis symbol. The concrete implementation lives in
// infrastructure/cache/redis/cache_key_builder_adapter.go.
type CacheKeyBuilderPort interface {
	// BuildCounterCacheKey returns the Redis key used to cache the
	// per-org asset template counter (total templates scoped to orgId).
	//
	// Parameters:
	//   - orgId: Organization ID (hex string). May be empty for
	//     unscoped/global counts.
	//
	// Returns:
	//   - The full Redis key string.
	BuildCounterCacheKey(orgId string) string
}
