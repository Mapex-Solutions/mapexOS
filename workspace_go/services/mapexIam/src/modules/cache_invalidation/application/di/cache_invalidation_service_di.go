package di

import (
	authRepos "mapexIam/src/modules/auth/domain/repositories"
	authCacheRepos "mapexIam/src/modules/authorization_cache/domain/repositories"
	groupRepos "mapexIam/src/modules/groups/domain/repositories"
	membershipPorts "mapexIam/src/modules/memberships/application/ports"

	"go.uber.org/dig"
)

// CacheInvalidationServiceDependenciesInjection defines the dependency injection container
// for the CacheInvalidationService. All fields are driven ports (interfaces) following
// Hexagonal Architecture — no concrete infrastructure types.
type CacheInvalidationServiceDependenciesInjection struct {
	dig.In

	// Driven ports
	AuthCacheRepo     authCacheRepos.AuthCacheRepository
	CoverageCacheRepo authRepos.CoverageCacheRepository
	MembershipService membershipPorts.MembershipServicePort
	GroupMemberRepo   groupRepos.GroupMemberRepository
}
