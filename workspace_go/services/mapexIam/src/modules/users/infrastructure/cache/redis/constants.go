package redis

// Counter cache Redis key configuration
const (
	// CounterCacheKeyPrefix is the Redis key prefix for user counter cache
	CounterCacheKeyPrefix = "counter:users:"
)

// BuildCounterCacheKey builds the full cache key for a given orgId
func BuildCounterCacheKey(orgId string) string {
	return CounterCacheKeyPrefix + orgId
}
