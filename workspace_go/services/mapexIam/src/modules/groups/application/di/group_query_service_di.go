package di

import (
	"mapexIam/src/modules/groups/domain/repositories"

	"go.uber.org/dig"
)

// GroupQueryServiceDependenciesInjection defines the dependency injection container for GroupQueryService.
// This service depends ONLY on same-domain repositories (no cross-domain dependencies).
// This design prevents circular dependencies:
//   - GroupService → UserService, MembershipService
//   - GroupQueryService → GroupRepo, GroupMemberRepo (ONLY)
//
// Used by:
//   - GroupQueryService (as single dependency container)
type GroupQueryServiceDependenciesInjection struct {
	dig.In

	// Same domain repositories (no cross-domain dependencies)
	Repo            repositories.GroupRepository
	GroupMemberRepo repositories.GroupMemberRepository
}
