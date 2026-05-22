package ports

// DataSourceCacheKeyBuilderPort defines the contract for building cache keys
// for DataSource entities.
//
// This port hides the cache-technology-specific prefix (Redis key layout) from
// the application layer, preserving the Dependency Inversion Principle: the
// application expresses "I need a cache key for this dataSource" and the
// infrastructure chooses the concrete key format.
//
// Implementations:
//   - infrastructure/cache/redis.DataSourceCacheKeyBuilder
type DataSourceCacheKeyBuilderPort interface {
	// BuildKey returns the cache key used to store/read the DataSource whose
	// unique identifier (ObjectID as string) is dataSourceId.
	BuildKey(dataSourceId string) string
}
