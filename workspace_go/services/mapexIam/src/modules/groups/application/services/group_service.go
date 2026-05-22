package services

import (
	ctx "context"
	"fmt"
	"time"

	"mapexIam/src/modules/groups/application/dtos"
	"mapexIam/src/modules/groups/application/di"
	"mapexIam/src/modules/groups/application/ports"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// Compile-time check to ensure GroupService implements GroupServicePort.
var _ ports.GroupServicePort = (*GroupService)(nil)

// New creates and returns a new instance of GroupService.
func New(deps di.GroupServiceDependenciesInjection) ports.GroupServicePort {
	return &GroupService{deps: deps}
}

// CreateGroup orchestrates group creation: validate org context, stamp
// multi-tenant fields, persist the group, attach a Membership for its
// permissions, invalidate the per-org counter cache, return the response DTO.
// Membership failure rolls the new group back so the system stays consistent.
func (s *GroupService) CreateGroup(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateGroupDto) (*dtos.GroupResponse, error) {
	if err := s.applyGroupCreateScope(requestContext, dto); err != nil {
		return nil, err
	}
	createdGroup, err := s.persistNewGroup(c, dto)
	if err != nil {
		return nil, err
	}
	s.invalidateGroupCounterCache(c, requestContext)
	if err := s.attachGroupMembership(c, dto, createdGroup); err != nil {
		return nil, err
	}
	return s.buildGroupResponse(createdGroup), nil
}

// GetGroupById fetches a group, enriches it with the live member count and
// the role list pulled from its associated Membership, and returns the DTO.
func (s *GroupService) GetGroupById(c ctx.Context, groupId *string) (*dtos.GroupResponse, error) {
	groupEntityData, err := s.fetchGroupOr404(c, groupId)
	if err != nil {
		return nil, err
	}
	response := s.buildGroupResponse(groupEntityData)
	s.enrichGroupMemberCount(c, response, *groupId)
	s.enrichGroupRoleIds(c, response, groupEntityData, groupId)
	return response, nil
}

// UpdateGroupById applies a partial update to a group and fans out a
// fire-and-forget cache-invalidation event so consumers refresh their cached
// principals. 404 when the id is unknown.
func (s *GroupService) UpdateGroupById(c ctx.Context, groupId *string, dto *dtos.UpdateGroupDto) (*dtos.GroupResponse, error) {
	if _, err := s.fetchGroupOr404(c, groupId); err != nil {
		return nil, err
	}
	fields := s.buildGroupUpdateFields(dto)
	updated, err := s.applyGroupUpdate(c, groupId, fields)
	if err != nil {
		return nil, err
	}
	go s.publishGroupChanged(updated.ID.Hex(), s.orgIdHex(updated.OrgID))
	return s.buildGroupResponse(updated), nil
}

// DeleteGroupById removes a group, but only when it has no active members.
// On success it tears down the group's Membership rows, drops the per-org
// counter cache, and fans out a deletion event for cache invalidation.
func (s *GroupService) DeleteGroupById(c ctx.Context, groupId *string) (map[string]bool, error) {
	group, err := s.fetchGroupOr404(c, groupId)
	if err != nil {
		return nil, err
	}
	if err := s.ensureGroupIsEmpty(c, *groupId); err != nil {
		return nil, err
	}
	s.deleteGroupMemberships(c, group, groupId)
	if err := s.deleteGroupFromRepo(c, groupId); err != nil {
		return nil, err
	}
	s.invalidateGroupCounterCacheForOrg(c, group)
	go s.publishGroupDeleted(*groupId, s.orgIdHex(group.OrgID))
	return map[string]bool{"success": true}, nil
}

// GetGroups returns the paginated, filtered group list scoped to the
// caller's org context, optionally narrowed by member id, with each result
// enriched with its member count.
func (s *GroupService) GetGroups(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.GroupQueryDto) (*model.PaginatedResult[dtos.GroupResponse], error) {
	filters, empty := s.buildGroupListFilters(c, requestContext, query)
	if empty {
		return s.emptyGroupListResult(query), nil
	}
	pagination := s.buildGroupListPagination(query)
	projection := s.buildGroupListProjection(query)
	result, err := s.deps.Repo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		return nil, err
	}
	return s.mapGroupListResultWithCounts(c, result), nil
}

// CountGroups returns the per-org group count with cache-aside semantics
// (6h TTL); the counter cache is invalidated on every Create/Delete.
func (s *GroupService) CountGroups(c ctx.Context, requestContext *reqCtx.RequestContext) (int64, error) {
	orgId := s.orgIdFromContext(requestContext)
	cacheKey := s.deps.CounterCache.BuildKey(orgId)
	if count, ok := s.tryCachedGroupCount(c, cacheKey); ok {
		return count, nil
	}
	count, err := s.countGroupsFromRepo(c, requestContext)
	if err != nil {
		return 0, err
	}
	s.cacheGroupCount(c, cacheKey, count)
	return count, nil
}

// AddMemberToGroup adds a user to a group via the GroupMember junction
// table. Idempotent: returns nil when the user is already a member.
func (s *GroupService) AddMemberToGroup(c ctx.Context, groupId string, userId string) error {
	group, groupObjectID, userObjectID, err := s.resolveAddMemberInputs(c, groupId, userId)
	if err != nil {
		return err
	}
	if exists, _ := s.deps.GroupMemberRepo.FindByGroupAndUser(c, groupId, userId); exists != nil {
		return nil
	}
	if err := s.persistGroupMember(c, group, groupObjectID, userObjectID); err != nil {
		return err
	}
	go s.publishGroupChanged(groupId, s.orgIdHex(group.OrgID))
	return nil
}

// RemoveMemberFromGroup removes a user from a group's junction table.
// Idempotent: returns nil when the user is not currently a member.
func (s *GroupService) RemoveMemberFromGroup(c ctx.Context, groupId string, userId string) error {
	group, err := s.fetchGroupOr404(c, &groupId)
	if err != nil {
		return err
	}
	if exists, _ := s.deps.GroupMemberRepo.FindByGroupAndUser(c, groupId, userId); exists == nil {
		return nil
	}
	if err := s.deps.GroupMemberRepo.DeleteByGroupAndUser(c, groupId, userId); err != nil {
		return &customErrors.ServerCustomError{
			Code:   status.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to remove user from group %s: %v", groupId, err)},
		}
	}
	go s.publishGroupChanged(groupId, s.orgIdHex(group.OrgID))
	return nil
}

// GetUserGroupsInOrg returns the IDs of every group a user belongs to in
// the given org. Used by access-resolution to merge group permissions.
func (s *GroupService) GetUserGroupsInOrg(c ctx.Context, userId, orgId string) ([]string, error) {
	members, err := s.deps.GroupMemberRepo.FindByUserIdAndOrgId(c, userId, orgId)
	if err != nil {
		return nil, &customErrors.ServerCustomError{
			Code:   status.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to query user groups: %v", err)},
		}
	}
	groupIds := make([]string, len(members))
	for i, member := range members {
		groupIds[i] = member.GroupID.Hex()
	}
	return groupIds, nil
}

// GetGroupMembers returns the paginated member list of one group, enriched
// with the user-detail fields (email, first/last name) so the UI does not
// need a follow-up users lookup.
func (s *GroupService) GetGroupMembers(c ctx.Context, groupId string, query *dtos.GroupMembersQueryDto) (*model.PaginatedResult[dtos.GroupMemberResponse], error) {
	if _, err := s.fetchGroupOr404(c, &groupId); err != nil {
		return nil, err
	}
	page, perPage := s.cappedPaginationOpts(query)
	members, totalItems, err := s.deps.GroupMemberRepo.FindByGroupIdPaginated(c, groupId, page, perPage)
	if err != nil {
		return nil, &customErrors.ServerCustomError{
			Code:   status.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to get group members: %v", err)},
		}
	}
	dtoItems := s.enrichGroupMemberDetails(c, members)
	return s.buildGroupMembersResult(dtoItems, page, perPage, totalItems), nil
}

// silence imports kept here so the file compiles even when a public method
// happens to drop one — keeps build noise out of unrelated diffs.
var _ = time.Time{}
