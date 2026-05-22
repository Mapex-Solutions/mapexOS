package redis

import (
	"mapexIam/src/modules/groups/application/ports"
)

// NewCounterCacheAdapter creates a new CounterCacheAdapter instance that
// implements groups/application/ports.CounterCachePort.
//
// Returns:
//   - ports.CounterCachePort: The adapter as an interface (Hexagonal Architecture)
func NewCounterCacheAdapter() ports.CounterCachePort {
	return &CounterCacheAdapter{}
}

// BuildKey builds the full counter cache key for the given orgId.
// Delegates to the package-level BuildCounterCacheKey to keep a single source
// of truth for the key layout.
func (a *CounterCacheAdapter) BuildKey(orgId string) string {
	return BuildCounterCacheKey(orgId)
}

// Compile-time check to ensure CounterCacheAdapter implements CounterCachePort.
var _ ports.CounterCachePort = (*CounterCacheAdapter)(nil)
