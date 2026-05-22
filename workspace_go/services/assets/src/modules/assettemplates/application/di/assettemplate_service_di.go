package di

import (
	"assets/src/bootstrap"
	assetsPorts "assets/src/modules/assets/application/ports"
	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/domain/repositories"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// AssetTemplateServiceDependenciesInjection aggregates all dependencies required
// by the AssetTemplateService.
//
// This struct follows the Dependency Injection pattern using uber/dig, enabling
// automatic dependency resolution and loose coupling between layers.
//
// Dependencies:
//   - AssetTemplateRepo: Repository for asset template persistence operations
//   - NatsBus: NATS Fanout interface for publishing cache invalidation events
//   - TemplateStoragePort: Port for template script storage (abstracts MinIO)
//   - TieredCache: L1 (disk) + L2 (MinIO) cache for scripts
//
// The dig.In tag enables automatic dependency injection by the dig container.
type AssetTemplateServiceDependenciesInjection struct {
	dig.In
	AssetTemplateRepo repositories.AssetTemplateRepository

	// AppCache provides service-private cache (Redis DB 0)
	// Used for: counter cache with 6h TTL
	AppCache common.AppCache

	// NatsBus provides FANOUT publishing for cache invalidation
	// Uses interface for testability (can be mocked in unit tests)
	NatsBus natsModel.Fanout `name:"core"`

	// TemplateStoragePort provides template script storage operations (abstracts MinIO)
	// Follows Hexagonal Architecture: application depends on port, not MinIO implementation
	// Key format: templates/{templateId}.json
	// Consumed by: JS-Executor via TieredCache
	TemplateStoragePort ports.TemplateStoragePort

	// TieredCache provides L1 (disk) + L2 (MinIO) caching for scripts
	// L0 (RAM) disabled - scripts are large and rarely change
	// Used by: GetAssetTemplateById to fetch scripts with caching
	TieredCache common.TieredCache `name:"templates"`

	// CacheKeyBuilder produces Redis keys for the application cache
	// (counter). Abstracted as a port so the application layer never
	// imports infrastructure/cache/redis directly.
	CacheKeyBuilder ports.CacheKeyBuilderPort

	// Metrics provides service-specific Prometheus metrics for instrumentation
	Metrics *bootstrap.AssetsMetrics

	// L2WritesPublisher publishes a durable retry hint to the L2 writes
	// stream when the synchronous MinIO write fails. Reuses the
	// platform-shared port defined in the assets module — the publisher
	// concern is the same; only the subject/MsgId prefix varies per
	// module.
	L2WritesPublisher assetsPorts.L2WritesPublisherPort
}
