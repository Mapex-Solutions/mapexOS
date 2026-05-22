package services

import (
	ctx "context"
	"fmt"
	"time"

	events "mapexIam/src/modules/cache_invalidation/application/events"
	constants "mapexIam/src/modules/groups/application/constants"
	"mapexIam/src/modules/groups/application/dtos"
	"mapexIam/src/modules/groups/domain/entities"
	membershipDtos "mapexIam/src/modules/memberships/application/dtos"

	cacheInvalidation "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/cache_invalidation"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// applyGroupCreateScope validates org context, stamps OrgID/PathKey on the
// DTO, and rejects requests missing the org binding (groups always belong
// to an org).
func (s *GroupService) applyGroupCreateScope(rc *reqCtx.RequestContext, dto *dtos.CreateGroupDto) error {
	if err := orgfilter.ValidateOrgContextForNonSystem(rc); err != nil {
		return &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{err.Error()}}
	}
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		if orgObjectId, err := model.ToObjectID(*rc.OrgContext); err == nil {
			dto.OrgID = &orgObjectId
		}
	}
	if rc.OrgContextData != nil && rc.OrgContextData.PathKey != "" {
		pathKey := rc.OrgContextData.PathKey
		dto.PathKey = &pathKey
	}
	return nil
}

// persistNewGroup maps the DTO to a domain entity, sets default scope,
// and writes it to the repository. Returns BAD_REQUEST when OrgID is missing.
func (s *GroupService) persistNewGroup(c ctx.Context, dto *dtos.CreateGroupDto) (*entities.Group, error) {
	groupEntity, _ := mapper.DtoToEntity[dtos.CreateGroupDto, entities.Group](dto)
	if dto.OrgID == nil {
		return nil, &customErrors.ServerCustomError{
			Code:   status.BAD_REQUEST,
			Errors: []string{"OrgID is required for group creation"},
		}
	}
	groupEntity.Scope = "local"
	return s.deps.Repo.Create(c, groupEntity)
}

// attachGroupMembership creates the Membership row that carries the group's
// permissions (Group -> Membership(assigneeType=group) -> Roles). Failure
// rolls the new group back so the system stays consistent.
func (s *GroupService) attachGroupMembership(c ctx.Context, dto *dtos.CreateGroupDto, createdGroup *entities.Group) error {
	if len(dto.RoleIds) == 0 {
		return nil
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Group] Creating membership for group: GroupID=%s, OrgID=%s, Roles=%v",
		createdGroup.ID.Hex(), dto.OrgID.Hex(), dto.RoleIds))
	membershipDto := &membershipDtos.CreateMembershipDto{
		AssigneeType: "group",
		AssigneeID:   createdGroup.ID.Hex(),
		OrgID:        dto.OrgID.Hex(),
		RoleIds:      dto.RoleIds,
		Scope:        createdGroup.Scope,
		Enabled:      true,
	}
	if _, err := s.deps.MembershipService.CreateMembership(c, membershipDto); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Group] Failed to create membership for group %s - rolling back group creation", createdGroup.ID.Hex()))
		groupId := createdGroup.ID.Hex()
		_ = s.deps.Repo.DeleteById(c, &groupId)
		return &customErrors.ServerCustomError{
			Code:   status.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to create membership for group: %v", err)},
		}
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Group] Membership created successfully for group %s", createdGroup.ID.Hex()))
	return nil
}

// invalidateGroupCounterCache drops the per-org counter cache after a
// successful create so the next CountGroups hits the repository.
func (s *GroupService) invalidateGroupCounterCache(c ctx.Context, rc *reqCtx.RequestContext) {
	if rc.OrgContext != nil {
		_ = s.deps.AppCache.Del(c, s.deps.CounterCache.BuildKey(*rc.OrgContext))
	}
}

// invalidateGroupCounterCacheForOrg mirrors the create-side helper but uses
// the deleted entity's OrgId because the request context may have moved on.
func (s *GroupService) invalidateGroupCounterCacheForOrg(c ctx.Context, group *entities.Group) {
	if group.OrgID != nil {
		_ = s.deps.AppCache.Del(c, s.deps.CounterCache.BuildKey(group.OrgID.Hex()))
	}
}

// fetchGroupOr404 loads a group by id and returns NOT_FOUND when the repo
// returns nil. Used by every read/update/delete orchestration.
func (s *GroupService) fetchGroupOr404(c ctx.Context, groupId *string) (*entities.Group, error) {
	group, err := s.deps.Repo.FindById(c, groupId)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Group not found"}}
	}
	return group, nil
}

// buildGroupResponse converts a Group entity to its response DTO.
func (s *GroupService) buildGroupResponse(group *entities.Group) *dtos.GroupResponse {
	resp, _ := mapper.EntityToDto[entities.Group, dtos.GroupResponse](group)
	return resp
}

// enrichGroupMemberCount populates GroupResponse.MembersCount from the
// junction table; on lookup failure the count is logged and left at 0.
func (s *GroupService) enrichGroupMemberCount(c ctx.Context, response *dtos.GroupResponse, groupId string) {
	memberCount, err := s.deps.GroupMemberRepo.CountByGroupId(c, groupId)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Group] Failed to get member count for group %s: %v", groupId, err))
		memberCount = 0
	}
	response.MembersCount = &memberCount
}

// enrichGroupRoleIds populates GroupResponse.RoleIds by reading the group's
// associated Membership (Group -> Membership(assigneeType=group) -> RoleIds).
func (s *GroupService) enrichGroupRoleIds(c ctx.Context, response *dtos.GroupResponse, groupEntity *entities.Group, groupId *string) {
	if groupEntity.OrgID == nil {
		return
	}
	orgId := groupEntity.OrgID.Hex()
	q := &membershipDtos.MembershipQueryDto{
		AssigneeType: ptrString("group"),
		AssigneeID:   groupId,
		OrgID:        &orgId,
	}
	memberships, err := s.deps.MembershipService.GetAllMemberships(c, q)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Group] Failed to get membership for group %s: %v", *groupId, err))
		return
	}
	if len(memberships) == 0 || len(memberships[0].RoleIds) == 0 {
		return
	}
	roleIds := make([]string, len(memberships[0].RoleIds))
	for i, roleId := range memberships[0].RoleIds {
		roleIds[i] = roleId.Hex()
	}
	response.RoleIds = roleIds
}

// buildGroupUpdateFields composes the $set map for an update, stamping the
// updated timestamp at the same time.
func (s *GroupService) buildGroupUpdateFields(dto *dtos.UpdateGroupDto) map[string]interface{} {
	fields, _ := mapper.DtoToMap(dto)
	fields["updated"] = time.Now()
	return fields
}

// applyGroupUpdate runs the partial update against the repository and
// translates a missing document into the canonical 404 contract error.
func (s *GroupService) applyGroupUpdate(c ctx.Context, groupId *string, fields map[string]interface{}) (*entities.Group, error) {
	group, _ := s.deps.Repo.FindByIdAndUpdate(c, groupId, fields)
	if group.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Group not found"}}
	}
	return group, nil
}

// ensureGroupIsEmpty refuses deletion when the group still has members.
// Members must be removed via RemoveMemberFromGroup first.
func (s *GroupService) ensureGroupIsEmpty(c ctx.Context, groupId string) error {
	memberCount, err := s.deps.GroupMemberRepo.CountByGroupId(c, groupId)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Group] Failed to count group members: %v", err))
		memberCount = 0
	}
	if memberCount > 0 {
		return &customErrors.ServerCustomError{
			Code:   status.BAD_REQUEST,
			Errors: []string{fmt.Sprintf("Cannot delete group with %d members. Remove all members first or use force delete.", memberCount)},
		}
	}
	return nil
}

// deleteGroupMemberships tears down every Membership row associated with
// the group (assigneeType=group). Best-effort — failures are logged.
func (s *GroupService) deleteGroupMemberships(c ctx.Context, group *entities.Group, groupId *string) {
	if group.OrgID == nil {
		return
	}
	orgId := group.OrgID.Hex()
	q := &membershipDtos.MembershipQueryDto{
		AssigneeType: ptrString("group"),
		AssigneeID:   groupId,
		OrgID:        &orgId,
	}
	memberships, err := s.deps.MembershipService.GetAllMemberships(c, q)
	if err != nil || len(memberships) == 0 {
		return
	}
	for _, membership := range memberships {
		membershipId := membership.ID.Hex()
		if _, delErr := s.deps.MembershipService.DeleteMembershipById(c, &membershipId); delErr != nil {
			logger.Warn(fmt.Sprintf("[SERVICE:Group] Failed to delete membership %s for group %s: %v", membershipId, *groupId, delErr))
			continue
		}
		logger.Debug(fmt.Sprintf("[SERVICE:Group] Deleted membership %s for group %s", membershipId, *groupId))
	}
}

// deleteGroupFromRepo removes the document and translates the driver's
// "document not found" string into the canonical 404 contract error.
func (s *GroupService) deleteGroupFromRepo(c ctx.Context, groupId *string) error {
	if err := s.deps.Repo.DeleteById(c, groupId); err != nil {
		if err.Error() == "document not found" {
			return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Group not found"}}
		}
		return err
	}
	return nil
}

// publishGroupChanged emits the cache-invalidation event when the group
// itself or its membership changed. Best-effort and async.
func (s *GroupService) publishGroupChanged(groupId, orgId string) {
	logger.Info(fmt.Sprintf("[SERVICE:Group] Group changed group=%s org=%s - publishing cache invalidation event", groupId, orgId))
	event := events.NewGroupChangedEvent(groupId, orgId, "")
	subject := fmt.Sprintf(cacheInvalidation.GroupChangedSubjectFormat, groupId)
	if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{Subject: subject, Data: event}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Group] Failed to publish GroupChangedEvent for group=%s", groupId))
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:Group] Published GroupChangedEvent for group=%s to subject=%s", groupId, subject))
}

// publishGroupDeleted emits the cache-invalidation event after the group is
// removed so consumers drop every cached principal that referenced it.
func (s *GroupService) publishGroupDeleted(groupId, orgId string) {
	logger.Info(fmt.Sprintf("[SERVICE:Group] Group deleted group=%s org=%s - publishing cache invalidation event", groupId, orgId))
	event := events.NewGroupDeletedEvent(groupId, orgId, "")
	subject := fmt.Sprintf(cacheInvalidation.GroupDeletedSubjectFormat, groupId)
	if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{Subject: subject, Data: event}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Group] Failed to publish GroupDeletedEvent for group=%s", groupId))
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:Group] Published GroupDeletedEvent for group=%s to subject=%s", groupId, subject))
}

// orgIdHex unwraps the optional OrgID into the empty string when absent so
// log lines and event payloads carry a stable type.
func (s *GroupService) orgIdHex(orgID *model.ObjectId) string {
	if orgID == nil {
		return ""
	}
	return orgID.Hex()
}

// orgIdFromContext extracts the orgId string from the request context (or
// "" when absent) for cache-key building in CountGroups.
func (s *GroupService) orgIdFromContext(rc *reqCtx.RequestContext) string {
	if rc.OrgContext != nil {
		return *rc.OrgContext
	}
	return ""
}

// tryCachedGroupCount returns (count, true) on a Redis hit; (0, false) when
// the caller must fall back to the repository.
func (s *GroupService) tryCachedGroupCount(c ctx.Context, cacheKey string) (int64, bool) {
	var count int64
	if err := s.deps.AppCache.Get(c, cacheKey, &count); err == nil {
		return count, true
	}
	return 0, false
}

// countGroupsFromRepo runs the org-scoped CountDocuments query.
func (s *GroupService) countGroupsFromRepo(c ctx.Context, rc *reqCtx.RequestContext) (int64, error) {
	logger.Debug(fmt.Sprintf("[SERVICE:Group] Counter cache miss for orgId=%s", s.orgIdFromContext(rc)))
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc})
	if err != nil {
		return 0, err
	}
	return s.deps.Repo.CountDocuments(c, orgFilter)
}

// cacheGroupCount stores a fresh count under the per-org cache key.
// Best-effort; failures are intentionally ignored.
func (s *GroupService) cacheGroupCount(c ctx.Context, cacheKey string, count int64) {
	_ = s.deps.AppCache.SetEx(c, cacheKey, count, constants.CounterCacheTTL)
}

// ptrString returns a pointer to its string argument — used to build
// pointer-typed query DTOs without temporary variables.
func ptrString(s string) *string {
	return &s
}
