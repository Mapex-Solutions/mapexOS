package redis

// Counter cache Redis key configuration
const (
	// CounterCacheKeyPrefix is the Redis key prefix for group counter cache
	CounterCacheKeyPrefix = "counter:groups:"
)

// BuildCounterCacheKey builds the full cache key for a given orgId
func BuildCounterCacheKey(orgId string) string {
	return CounterCacheKeyPrefix + orgId
}
