package di

import (
	"go.uber.org/dig"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	authPorts "mapexIam/src/modules/auth/application/ports"
	membershipPorts "mapexIam/src/modules/memberships/application/ports"
	orgPorts "mapexIam/src/modules/organizations/application/ports"
	rolePorts "mapexIam/src/modules/roles/application/ports"
)

// AuthorizationCacheRepositoryDI defines the dependency injection container for AuthorizationCacheRepository.
// This struct aggregates all dependencies required by AuthorizationCacheRepository using dig.In.
//
// Architecture Pattern: Dependency Injection with Uber Dig
//   - dig.In: Instructs Dig to inject all fields automatically
//   - Provides clean constructor signature (single parameter instead of multiple)
//   - Scalable: Adding new dependencies doesn't change constructor signature
//
// DDD/Hexagonal Architecture:
//   - All cross-domain and infrastructure dependencies exposed via ports (interfaces)
//   - Never uses concrete drivers directly (no *redisLock.LockManager) — preserves
//     the DI pattern where all fields are port interfaces.
//
// Used by:
//   - AuthorizationCacheRepository (as single dependency container)
type AuthorizationCacheRepositoryDI struct {
	dig.In

	Cache             common.SharedCache
	LockMgr           authPorts.LockManagerPort
	MembershipService membershipPorts.MembershipServicePort
	RoleService       rolePorts.RoleServicePort
	OrgService        orgPorts.OrganizationServicePort
}
