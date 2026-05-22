package di

import (
	"http_gateway/src/bootstrap"
	"http_gateway/src/modules/datasources/application/ports"
	"http_gateway/src/modules/datasources/domain/repositories"

	"go.uber.org/dig"
)

// DataSourceServiceDependenciesInjection aggregates all dependencies required
// by the DataSourceService.
//
// This struct follows the Dependency Injection pattern using uber/dig, enabling
// automatic dependency resolution and loose coupling between layers.
//
// Dependencies:
//   - DataSourceRepo: Repository for data source persistence operations
//   - CacheRepo: Repository for caching data source lookups (Redis)
//   - CacheKeyBuilder: Port that builds cache keys for DataSource entities
//     (keeps the Redis key layout out of the application layer)
//   - Metrics: Service-specific Prometheus metrics for instrumentation
//
// The dig.In tag enables automatic dependency injection by the dig container.
type DataSourceServiceDependenciesInjection struct {
	dig.In
	DataSourceRepo  repositories.DataSourceRepository
	CacheRepo       repositories.CacheRepository
	CacheKeyBuilder ports.DataSourceCacheKeyBuilderPort
	Metrics         *bootstrap.HttpGatewayMetrics
}
