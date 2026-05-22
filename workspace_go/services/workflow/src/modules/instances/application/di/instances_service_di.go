package di

import (
	"workflow/src/modules/instances/application/ports"
	"workflow/src/modules/instances/domain/repositories"
	runtimePorts "workflow/src/modules/runtime/application/ports"

	"go.uber.org/dig"
)

// InstancesServiceDependenciesInjection aggregates all dependencies required by InstancesService.
//
// Following Hexagonal Architecture principles:
//   - Same-domain: InstanceRepository (MongoDB) for CRUD queries
//   - Infrastructure: InstanceLoader (TieredCache + MongoDB) for cached reads
type InstancesServiceDependenciesInjection struct {
	dig.In

	// InstanceRepo provides persistence operations for WorkflowInstance entities (MongoDB)
	InstanceRepo repositories.InstanceRepository

	// InstanceLoader provides cached instance reads (TieredCache L0/L1 → MongoDB fallback)
	InstanceLoader ports.InstanceLoaderPort

	// RuntimeService provides workflow execution operations
	RuntimeService runtimePorts.RuntimeServicePort
}
