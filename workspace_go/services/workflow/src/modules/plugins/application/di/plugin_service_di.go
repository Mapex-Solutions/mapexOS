package di

import (
	"workflow/src/bootstrap"
	"workflow/src/modules/plugins/application/ports"
	"workflow/src/modules/plugins/domain/repositories"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// PluginServiceDependenciesInjection aggregates all dependencies required by PluginService.
//
// Following Hexagonal Architecture principles:
//   - Same-domain dependencies: Repositories (PluginManifestRepository)
//   - Application ports: PluginLoaderPort (TieredCache L0→L1→MongoDB)
//   - Messaging: NatsBus for FANOUT cache invalidation across pods
//   - Metrics: WorkflowMetrics for Prometheus instrumentation
type PluginServiceDependenciesInjection struct {
	dig.In

	// PluginRepo provides persistence operations for PluginManifest entities (MongoDB)
	PluginRepo repositories.PluginManifestRepository

	// PluginLoader provides cached access to plugin manifests (TieredCache L0→L1→MongoDB)
	PluginLoader ports.PluginLoaderPort

	// NatsBus provides FANOUT publishing for cache invalidation across pods
	NatsBus natsModel.Fanout `name:"core"`

	// Metrics provides service-specific Prometheus metrics for instrumentation
	Metrics *bootstrap.WorkflowMetrics
}
