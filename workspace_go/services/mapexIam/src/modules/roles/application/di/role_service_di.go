package di

import (
	"mapexIam/src/modules/roles/application/ports"
	"mapexIam/src/modules/roles/domain/repositories"
	orgPorts "mapexIam/src/modules/organizations/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// RoleServiceDependenciesInjection defines the dependency injection container for RoleService.
// This struct aggregates all dependencies required by RoleService using dig.In.
//
// Architecture Pattern: Dependency Injection with Uber Dig
//   - dig.In: Instructs Dig to inject all fields automatically
//   - Provides clean constructor signature (single parameter instead of multiple)
//   - Scalable: Adding new dependencies doesn't change constructor signature
//
// DDD/Hexagonal Architecture:
//   - Same domain: Uses repository directly
//   - Other domains: Uses service port (NOT repository) to maintain bounded context
//
// Cache Invalidation:
//   - Uses NATS event bus to publish domain events
//   - Centralized consumer handles cache invalidation logic
//
// Used by:
//   - RoleService (as single dependency container)
type RoleServiceDependenciesInjection struct {
	dig.In

	// Same domain repository
	Repo repositories.RoleRepository

	// Other domain services (using ports for DDD bounded context)
	OrgService orgPorts.OrganizationServicePort

	// Event bus for domain event publishing (uses Publisher interface for Hexagonal Architecture)
	NatsBus natsModel.Publisher
}

// RoleServiceDI is used to inject RoleService into other services.
// This struct is used by consumers of RoleService (e.g., other modules).
type RoleServiceDI struct {
	dig.In

	RoleService ports.RoleServicePort
}
