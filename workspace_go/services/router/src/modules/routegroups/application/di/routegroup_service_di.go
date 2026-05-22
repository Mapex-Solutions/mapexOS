package di

import (
	"router/src/bootstrap"
	"router/src/modules/routegroups/domain/repositories"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	"go.uber.org/dig"
)

// RouteGroupServiceDependenciesInjection aggregates all dependencies required
// by the RouteGroupService.
//
// This struct follows the Dependency Injection pattern using uber/dig, enabling
// automatic dependency resolution and loose coupling between layers.
//
// Dependencies:
//   - RouteGroupRepo: Repository for route group persistence operations
//   - CacheRepo: Repository for caching route group data
//   - AppCache: Service-private Redis cache for counter and general caching
//   - Metrics: Prometheus metrics for observability
//
// The dig.In tag enables automatic dependency injection by the dig container.
type RouteGroupServiceDependenciesInjection struct {
	dig.In
	RouteGroupRepo repositories.RouteGroupRepository
	CacheRepo      repositories.CacheRepository
	AppCache       common.AppCache
	Metrics        *bootstrap.RouterMetrics
}
