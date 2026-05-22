package di

import (
	
	"workflow/src/bootstrap"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	enginePorts "workflow/src/modules/engine/application/ports"
	instancePorts "workflow/src/modules/instances/application/ports"
	pluginPorts "workflow/src/modules/plugins/application/ports"
	"workflow/src/modules/runtime/application/ports"
	"workflow/src/modules/runtime/domain/repositories"

	"go.uber.org/dig"
)

// RuntimeServiceDependenciesInjection aggregates all dependencies required by RuntimeService.
//
// Following Hexagonal Architecture principles:
//   - Domain: ExecutionStateRepo (NATS KV) for per-step state checkpoint
//   - Application: RuntimePublisher (NATS JetStream) for state events + resume/callback publishing
//   - Application: DefinitionLoader (TieredCache + MongoDB) for cached definition access
//   - Application: InstanceLoader (TieredCache + MongoDB) for cached instance config access
//   - Cross-module: ConditionEvaluatorPort, ValueResolverPort (from engine module)
type RuntimeServiceDependenciesInjection struct {
	dig.In

	// ExecutionStateRepo provides hot state persistence for workflow executions (NATS KV).
	ExecutionStateRepo repositories.ExecutionStateRepository

	// RuntimePublisher provides JetStream publishing for state events, timers, and callbacks.
	RuntimePublisher ports.RuntimePublisherPort

	// ConditionEvaluator provides condition group evaluation (from engine module)
	ConditionEvaluator enginePorts.ConditionEvaluatorPort

	// ValueResolver provides field value resolution (from engine module)
	ValueResolver enginePorts.ValueResolverPort

	// DefinitionLoader provides cached access to WorkflowDefinition entities.
	DefinitionLoader ports.DefinitionLoaderPort

	// InstanceLoader provides cached access to WorkflowInstance config entities.
	InstanceLoader instancePorts.InstanceLoaderPort

	// VaultService provides credential decryption via vault MS internal API
	VaultService ports.VaultPort

	// PluginRepo provides plugin manifest lookup for plugin execution (from plugins module)
	PluginRepo pluginPorts.PluginManifestRepository

	// ShutdownManager provides the IsShuttingDown() flag for graceful walker drain
	ShutdownManager *shutdown.ShutdownManager

	// Metrics provides Prometheus counters/histograms for runtime observability
	Metrics *bootstrap.WorkflowMetrics
}
