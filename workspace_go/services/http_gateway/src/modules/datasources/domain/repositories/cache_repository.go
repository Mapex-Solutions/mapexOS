package repositories

import (
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
)

// CacheRepository combines all cache operations needed by DataSource domain.
//
// This interface extends both basic cache operations (Set, Get, Del) and
// the GetOrSet pattern with TTL support, providing a complete abstraction
// over the caching infrastructure.
//
// By using this interface instead of concrete implementations, the domain
// layer remains independent of specific cache technologies (Redis, Memcached, etc.),
// adhering to Dependency Inversion Principle and improving testability.
//
// Implementations:
//   - *redisModel.RedisClient implements this interface
type CacheRepository interface {
	common.Cache
	common.CacheGetOrSetEx
}
