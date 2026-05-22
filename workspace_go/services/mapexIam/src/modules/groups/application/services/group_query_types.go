package services

import (
	"mapexIam/src/modules/groups/application/di"
)

// GroupQueryService provides read-only group queries for cross-domain consumption.
// This service exists to break circular dependencies between GroupService and its consumers.
//
// Dependency graph (no cycles):
//   - GroupQueryService → GroupRepo, GroupMemberRepo (same domain only)
//   - UserService → GroupQueryService (via port)
//   - MembershipService → GroupQueryService (via port)
//
// Architecture Pattern: Dependency Injection
//   - Uses GroupQueryServiceDependenciesInjection struct to aggregate all dependencies
type GroupQueryService struct {
	deps di.GroupQueryServiceDependenciesInjection
}
