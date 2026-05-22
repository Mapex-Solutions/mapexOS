package di

import (
	groupQueryPorts "mapexIam/src/modules/groups/application/ports"
	orgPorts "mapexIam/src/modules/organizations/application/ports"
	rolePorts "mapexIam/src/modules/roles/application/ports"
	"mapexIam/src/modules/users/application/ports"
	"mapexIam/src/modules/users/domain/repositories"

	membershipPorts "mapexIam/src/modules/memberships/application/ports"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	"go.uber.org/dig"
)

// UserServiceDependenciesInjection defines the dependency injection container for UserService.
// This struct aggregates all dependencies required by UserService using dig.In.
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
//   - UserService (as single dependency container)
type UserServiceDependenciesInjection struct {
	dig.In

	// Same domain repository
	Repo repositories.UserRepository

	// Cross-domain service ports (Hexagonal Architecture)
	MembershipService membershipPorts.MembershipServicePort

	// Cross-domain query service for group resolution (no circular dependency)
	// Replaces direct access to GroupRepo, GroupMemberRepo, and MembershipRepo
	GroupQueryService groupQueryPorts.GroupQueryServicePort

	// Cross-domain service for getting organization details (name, type)
	OrgService orgPorts.OrganizationServicePort

	// Cross-domain service for getting role details (name)
	RoleService rolePorts.RoleServicePort

	// AppCache provides service-private cache (Redis DB 0)
	AppCache common.AppCache

	// CounterCache exposes counter cache key construction via an application
	// port, keeping infrastructure details out of the application layer.
	CounterCache ports.CounterCachePort
}

// UserServiceDI is used to inject UserService into other services.
// This struct is used by consumers of UserService (e.g., AuthService).
type UserServiceDI struct {
	dig.In

	UserService ports.UserServicePort
}
