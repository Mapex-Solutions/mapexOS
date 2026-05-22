package redis

import (
	"http_gateway/src/modules/datasources/application/ports"
)

// Compile-time check to ensure DataSourceCacheKeyBuilder implements the port.
var _ ports.DataSourceCacheKeyBuilderPort = (*DataSourceCacheKeyBuilder)(nil)

// NewDataSourceCacheKeyBuilder returns an implementation of
// ports.DataSourceCacheKeyBuilderPort backed by the Redis key layout.
func NewDataSourceCacheKeyBuilder() ports.DataSourceCacheKeyBuilderPort {
	return &DataSourceCacheKeyBuilder{}
}

// BuildKey constructs a Redis cache key for a DataSource entity.
//
// The key follows the pattern: DATA_SOURCE:{dataSourceId}
//
// Parameters:
//   - dataSourceId: The unique identifier (ObjectID as string) of the DataSource.
//
// Returns:
//   - string: A formatted cache key ready to be used with Redis operations.
//
// Example:
//
//	builder := redis.NewDataSourceCacheKeyBuilder()
//	cacheKey := builder.BuildKey("507f1f77bcf86cd799439011")
//	// Returns: "DATA_SOURCE:507f1f77bcf86cd799439011"
func (b *DataSourceCacheKeyBuilder) BuildKey(dataSourceId string) string {
	return dataSourceCacheKeyPrefix + ":" + dataSourceId
}
