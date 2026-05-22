package di

import (
	"mapexIam/src/modules/groups/application/ports"
	"mapexIam/src/modules/groups/domain/repositories"
	membershipPorts "mapexIam/src/modules/memberships/application/ports"
	orgPorts "mapexIam/src/modules/organizations/application/ports"
	userPorts "mapexIam/src/modules/users/application/ports"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// GroupServiceDependenciesInjection defines the dependency injection container for GroupService.
// This struct aggregates all dependencies required by GroupService using dig.In.
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
//   - GroupService (as single dependency container)
type GroupServiceDependenciesInjection struct {
	dig.In

	// Same domain repositories
	Repo            repositories.GroupRepository
	GroupMemberRepo repositories.GroupMemberRepository

	// Event bus for domain event publishing (uses Publisher interface for Hexagonal Architecture)
	NatsBus natsModel.Publisher

	// Other domain services (using ports for DDD bounded context)
	OrgService        orgPorts.OrganizationServicePort
	MembershipService membershipPorts.MembershipServicePort
	UserService       userPorts.UserServicePort

	// AppCache provides service-private cache (Redis DB 0)
	AppCache common.AppCache

	// CounterCache exposes counter cache key construction via an application
	// port, keeping infrastructure details out of the application layer.
	CounterCache ports.CounterCachePort
}

// GroupServiceDI is used to inject GroupService into other services.
// This struct is used by consumers of GroupService (e.g., other modules).
type GroupServiceDI struct {
	dig.In

	GroupService ports.GroupServicePort
}
