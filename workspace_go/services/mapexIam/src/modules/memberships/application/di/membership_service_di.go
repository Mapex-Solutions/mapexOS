package di

import (
	groupQueryPorts "mapexIam/src/modules/groups/application/ports"
	"mapexIam/src/modules/memberships/application/ports"
	"mapexIam/src/modules/memberships/domain/repositories"
	orgPorts "mapexIam/src/modules/organizations/application/ports"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// MembershipServiceDependenciesInjection defines the dependency injection container for MembershipService.
// This struct aggregates all dependencies required by MembershipService using dig.In.
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
// Used by:
//   - MembershipService (as single dependency container)
type MembershipServiceDependenciesInjection struct {
	dig.In

	// Same domain repository and services
	Repo  repositories.MembershipRepository
	Cache common.CacheGetOrSetEx

	// Event bus for domain event publishing (uses Publisher interface for Hexagonal Architecture)
	NatsBus natsModel.Publisher

	// Other domain services (using ports for DDD bounded context)
	OrgService orgPorts.OrganizationServicePort

	// Cross-domain query service for group membership resolution (no circular dependency)
	// Used by GetAllUserMemberships to include memberships inherited from groups
	GroupQueryService groupQueryPorts.GroupQueryServicePort
}

// MembershipServiceDI is used to inject MembershipService into other services.
// This struct is used by consumers of MembershipService (e.g., other modules).
type MembershipServiceDI struct {
	dig.In

	MembershipService ports.MembershipServicePort
}
