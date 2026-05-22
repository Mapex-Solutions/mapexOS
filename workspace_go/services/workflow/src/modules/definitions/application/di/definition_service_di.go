package di

import (
	"workflow/src/bootstrap"
	"workflow/src/modules/definitions/application/ports"
	"workflow/src/modules/definitions/domain/repositories"
	pluginPorts "workflow/src/modules/plugins/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// DefinitionServiceDependenciesInjection aggregates all dependencies required by DefinitionService.
//
// Following Hexagonal Architecture principles:
//   - Same-domain dependencies: Repositories (DefinitionRepository)
//   - Infrastructure ports: DefinitionStoragePort (MinIO)
//   - Messaging: NatsBus for FANOUT cache invalidation
//   - Metrics: WorkflowMetrics for Prometheus instrumentation
type DefinitionServiceDependenciesInjection struct {
	dig.In

	// DefinitionRepo provides persistence operations for WorkflowDefinition entities (MongoDB)
	DefinitionRepo repositories.DefinitionRepository

	// NatsBus provides FANOUT publishing for cache invalidation
	NatsBus natsModel.Fanout `name:"core"`

	// DefinitionStoragePort provides definition storage operations (MinIO L2)
	DefinitionStoragePort ports.DefinitionStoragePort

	// PluginLoader provides cached access to plugin manifests for status computation
	PluginLoader pluginPorts.PluginLoaderPort

	// Metrics provides service-specific Prometheus metrics for instrumentation
	Metrics *bootstrap.WorkflowMetrics
}
