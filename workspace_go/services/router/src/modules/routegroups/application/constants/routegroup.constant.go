package constants

import "time"

// Cache TTL configuration for RouteGroup
const (
	// RouteGroupCacheTTL defines the time-to-live for RouteGroup cache entries
	RouteGroupCacheTTL = 60 * time.Minute

	// RouteGroupCacheKeyPrefix defines the prefix for RouteGroup cache keys
	// Format: ROUTE_GROUP:{routeGroupId}
	RouteGroupCacheKeyPrefix = "ROUTE_GROUP"
)

// Counter cache configuration
const (
	// CounterCacheKeyPrefix is the Redis key prefix for route group counter cache
	CounterCacheKeyPrefix = "counter:route_groups:"

	// CounterCacheTTL is the cache duration for counter values (6 hours)
	CounterCacheTTL = 6 * time.Hour
)

// BuildCounterCacheKey builds the full cache key for a given orgId
func BuildCounterCacheKey(orgId string) string {
	return CounterCacheKeyPrefix + orgId
}
