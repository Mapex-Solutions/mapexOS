package di

import (
	"events/src/modules/retention/application/ports"
	"events/src/modules/retention/domain/repositories"

	"go.uber.org/dig"
)

// RetentionServiceDependenciesInjection aggregates all dependencies required
// by the RetentionService.
//
// This struct follows the Dependency Injection pattern using uber/dig, enabling
// automatic dependency resolution and loose coupling between layers.
//
// Dependencies:
//   - RetentionRepo:   Repository for retention policy persistence operations.
//   - CacheRepo:       Repository for caching retention policy data.
//   - ClickHouseConn:  Port-scoped wrapper over the ClickHouse connection used
//     to apply TTL changes (ALTER TABLE ... MODIFY TTL) when a retention
//     policy is updated. Keeps the concrete driver type out of the
//     application layer.
//
// The dig.In tag enables automatic dependency injection by the dig container.
type RetentionServiceDependenciesInjection struct {
	dig.In
	RetentionRepo  repositories.RetentionRepository
	CacheRepo      repositories.CacheRepository
	ClickHouseConn ports.ClickHouseConnPort
}
