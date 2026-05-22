package redis

import (
	"assets/src/modules/assets/application/ports"
)

// Compile-time check: the adapter fulfils the application-layer port.
var _ ports.CacheKeyBuilderPort = (*CacheKeyBuilderAdapter)(nil)

// NewCacheKeyBuilderAdapter returns the adapter as the application port
// interface — the DI container resolves it by interface, not by concrete
// type (matches the pattern used by other adapters in this module).
func NewCacheKeyBuilderAdapter() ports.CacheKeyBuilderPort {
	return &CacheKeyBuilderAdapter{}
}

// BuildCounterCacheKey builds the Redis key for the per-org asset counter.
// Key format: counter:assets:{orgId}.
func (a *CacheKeyBuilderAdapter) BuildCounterCacheKey(orgId string) string {
	return CounterCacheKeyPrefix + orgId
}
