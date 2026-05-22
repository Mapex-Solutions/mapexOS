package di

import (
	"assets/src/bootstrap"
	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/domain/repositories"
	assettemplatePort "assets/src/modules/assettemplates/application/ports"
	healthPorts "assets/src/modules/healthmonitor/application/ports"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// AssetServiceDependenciesInjection aggregates all dependencies required by AssetService.
// Following Hexagonal Architecture principles:
//   - Same-domain dependencies: Repositories (AssetRepository)
//   - Cross-domain dependencies: AssetTemplateRepository (for fetching scripts)
//   - Infrastructure ports: AppCache, AssetStoragePort, RouteGroupPort
//   - Messaging: NatsBus for FANOUT cache invalidation
//
// Note: SharedCache (DB 5) was removed. Asset read model is now published to MinIO (L2)
// and consumed via TieredCache by other services. Cache invalidation uses NATS FANOUT.
//
// This struct uses dig.In to enable automatic dependency injection by uber/dig.
type AssetServiceDependenciesInjection struct {
	dig.In

	// AssetRepo provides persistence operations for Asset entities (MongoDB - source of truth)
	AssetRepo repositories.AssetRepository

	// AppCache provides service-private cache (Redis DB 0)
	// Used for: per-org asset counter cache.
	AppCache common.AppCache

	// AssetTemplateRepo provides access to asset templates (for scripts)
	AssetTemplateRepo assettemplatePort.AssetTemplateRepository

	// NatsBus provides FANOUT publishing for cache invalidation
	// Uses interface for testability (can be mocked in unit tests)
	NatsBus natsModel.Fanout `name:"core"`

	// RouteGroupPort provides RouteGroup lookup operations (abstracts Router service).
	// Follows Hexagonal Architecture: application depends on port, not HTTP implementation.
	// Used for: name enrichment on asset responses AND HealthMonitor router-kind
	// validation at Create/Update (validateHealthMonitorConfig).
	RouteGroupPort ports.RouteGroupPort

	// AssetStoragePort provides asset read model storage operations (abstracts MinIO).
	// Follows Hexagonal Architecture: application depends on port, not MinIO implementation.
	// Key format: {orgId}/{assetUUID}.json.
	// Consumed by: Router, JS-Executor, Events, mapex-mqtt-broker plugin via TieredCache.
	// The read model carries PasswordHash + CurrentCert so the broker plugin
	// decides MQTT CONNECTs locally (no auth callout, no parallel bucket).
	AssetStoragePort ports.AssetStoragePort

	// CacheKeyBuilder produces Redis keys for the application cache
	// (per-org counter). Abstracted as a port so the application layer
	// never imports infrastructure/cache/redis directly.
	CacheKeyBuilder ports.CacheKeyBuilderPort

	// Metrics provides service-specific Prometheus metrics for instrumentation
	Metrics *bootstrap.AssetsMetrics

	// HealthRepo provides real-time health status data from Redis
	// Used for API enrichment (GetById, List) with healthStatus and lastSeenAt
	HealthRepo healthPorts.HealthRepository

	// HealthLifecycle clears Redis health state when a user toggles
	// HealthMonitor.Enabled from true to false, or when an asset is deleted.
	HealthLifecycle healthPorts.HealthLifecyclePort

	// L2WritesPublisher publishes a durable retry hint to the L2 writes
	// stream when the synchronous MinIO write fails. The in-module
	// fallback consumer drains the stream and reconciles against
	// current Mongo state, re-emitting FANOUT on success.
	L2WritesPublisher ports.L2WritesPublisherPort
}
