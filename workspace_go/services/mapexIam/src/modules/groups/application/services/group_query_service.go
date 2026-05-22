package services

import (
	"context"

	"mapexIam/src/modules/groups/application/di"
	"mapexIam/src/modules/groups/application/ports"
)

// Compile-time check to ensure GroupQueryService implements GroupQueryServicePort
var _ ports.GroupQueryServicePort = (*GroupQueryService)(nil)

// NewGroupQueryService creates and returns a new instance of GroupQueryService.
//
// Parameters:
//   - deps: Dependency injection container providing repositories
//
// Returns:
//   - ports.GroupQueryServicePort: The service as an interface (Hexagonal Architecture)
func NewGroupQueryService(deps di.GroupQueryServiceDependenciesInjection) ports.GroupQueryServicePort {
	return &GroupQueryService{deps: deps}
}

// GetGroupBasicInfo retrieves basic group info (ID, Name, Description).
// Maps from the Group entity to a lightweight GroupBasicInfo DTO.
func (s *GroupQueryService) GetGroupBasicInfo(ctx context.Context, groupId string) (*ports.GroupBasicInfo, error) {
	return s.getGroupBasicInfo(ctx, groupId)
}

// GetAllUserGroupIds retrieves all group IDs a user belongs to across all organizations.
// Queries the GroupMember junction table by userId.
func (s *GroupQueryService) GetAllUserGroupIds(ctx context.Context, userId string) ([]string, error) {
	return s.getAllUserGroupIds(ctx, userId)
}

// GetUserIdsByGroupIds retrieves all user IDs that are members of any of the given groups.
// Uses GroupMember junction table bulk query for efficiency.
func (s *GroupQueryService) GetUserIdsByGroupIds(ctx context.Context, groupIds []string) ([]string, error) {
	return s.getUserIdsByGroupIds(ctx, groupIds)
}

// CountGroupsByUserIds returns a map of userId -> groupCount for batch enrichment.
// Queries GroupMember junction table for each user (individual queries).
func (s *GroupQueryService) CountGroupsByUserIds(ctx context.Context, userIds []string) (map[string]int, error) {
	return s.countGroupsByUserIds(ctx, userIds)
}
