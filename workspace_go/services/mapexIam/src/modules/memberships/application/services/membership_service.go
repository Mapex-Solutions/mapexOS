package services

import (
	ctx "context"
	"fmt"

	events "mapexIam/src/modules/cache_invalidation/application/events"
	"mapexIam/src/modules/memberships/application/di"
	"mapexIam/src/modules/memberships/application/dtos"
	"mapexIam/src/modules/memberships/application/ports"
	"mapexIam/src/modules/memberships/domain/entities"

	contractsCacheInvalidation "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/cache_invalidation"
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// Compile-time check to ensure MembershipService implements MembershipServicePort
var _ ports.MembershipServicePort = (*MembershipService)(nil)

// New creates and returns a new instance of MembershipService.
func New(deps di.MembershipServiceDependenciesInjection) ports.MembershipServicePort {
	return &MembershipService{deps: deps}
}

// CreateMembership orchestrates membership creation:
// idempotent existing-row check (same assignee+org returns the existing
// row) -> resolve organization for multi-tenant fields (404 on miss) ->
// translate string ids into ObjectIDs -> persist -> publish a NATS
// cache-invalidation event for user memberships.
func (s *MembershipService) CreateMembership(c ctx.Context, dto *dtos.CreateMembershipDto) (*dtos.MembershipResponse, error) {
	if existing, _ := s.deps.Repo.FindByAssigneeAndOrg(c, dto.AssigneeType, dto.AssigneeID, dto.OrgID); existing != nil {
		logger.Info(fmt.Sprintf("[SERVICE:Membership] Membership already exists for assignee=%s org=%s - returning existing", dto.AssigneeID, dto.OrgID))
		resp, _ := mapper.EntityToDto[entities.Membership, dtos.MembershipResponse](existing)
		return resp, nil
	}

	membershipEntity, _ := mapper.DtoToEntity[dtos.CreateMembershipDto, entities.Membership](dto)
	if err := s.bindMembershipTenantFields(c, membershipEntity, dto); err != nil {
		return nil, err
	}
	if err := convertMembershipObjectIds(membershipEntity, dto); err != nil {
		return nil, err
	}

	created, err := s.deps.Repo.Create(c, membershipEntity)
	if err != nil {
		return nil, err
	}

	if dto.AssigneeType == "user" {
		go s.publishMembershipChanged(created.ID.Hex(), dto.AssigneeID, dto.OrgID)
	}

	resp, _ := mapper.EntityToDto[entities.Membership, dtos.MembershipResponse](created)
	return resp, nil
}

// GetMembershipById fetches a single membership by id. Returns 404 when
// the id is unknown.
func (s *MembershipService) GetMembershipById(c ctx.Context, membershipId *string) (*dtos.MembershipResponse, error) {
	entity, err := s.deps.Repo.FindById(c, membershipId)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Membership not found"}}
	}
	resp, _ := mapper.EntityToDto[entities.Membership, dtos.MembershipResponse](entity)
	return resp, nil
}

// UpdateMembershipById orchestrates a partial update: build the $set map
// with optional roleIds string->ObjectID translation -> apply patch ->
// publish a NATS cache-invalidation event for user memberships.
func (s *MembershipService) UpdateMembershipById(c ctx.Context, membershipId *string, dto *dtos.UpdateMembershipDto) (*dtos.MembershipResponse, error) {
	fields, _ := mapper.DtoToMap(dto)
	if dto.RoleIds != nil && len(*dto.RoleIds) > 0 {
		roleObjectIds, err := convertStringsToObjectIds(*dto.RoleIds, "Invalid role ID format")
		if err != nil {
			return nil, err
		}
		fields["roleIds"] = roleObjectIds
	}

	updated, _ := s.deps.Repo.FindByIdAndUpdate(c, membershipId, fields)
	if updated.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Membership not found"}}
	}

	if updated.AssigneeType == "user" {
		go s.publishMembershipChanged(updated.ID.Hex(), updated.AssigneeID.Hex(), updated.OrgID.Hex())
	}

	resp, _ := mapper.EntityToDto[entities.Membership, dtos.MembershipResponse](updated)
	return resp, nil
}

// DeleteMembershipById orchestrates removal: load (404 on miss) so we
// can announce the deletion -> delete -> publish a NATS event for user
// memberships so coverage caches drop the principal-org link.
func (s *MembershipService) DeleteMembershipById(c ctx.Context, membershipId *string) (map[string]bool, error) {
	membership, err := s.deps.Repo.FindById(c, membershipId)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Membership not found"}}
	}

	if delErr := s.deps.Repo.DeleteById(c, membershipId); delErr != nil {
		if delErr.Error() == "document not found" {
			return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Membership not found"}}
		}
		return nil, delErr
	}

	if membership.AssigneeType == "user" {
		go s.publishMembershipDeleted(*membershipId, membership.AssigneeID.Hex(), membership.OrgID.Hex())
	}
	return map[string]bool{"success": true}, nil
}

// GetMemberships orchestrates the paginated list (HTTP path):
// build the org filter from RequestContext + apply optional
// assigneeType / assigneeId / userId / roleId / scope / enabled
// filters -> resolve pagination + projection -> delegate to repository
// -> map entities to response DTOs.
func (s *MembershipService) GetMemberships(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.MembershipQueryDto) (*model.PaginatedResult[dtos.MembershipResponse], error) {
	filters, err := s.buildMembershipListFilters(requestContext, query)
	if err != nil {
		return nil, err
	}
	pagination := &model.PaginationOpts{Page: int64(query.GetPage()), PerPage: int64(query.GetPerPage())}
	projection := orgfilter.BuildProjection(query.Projection)

	result, err := s.deps.Repo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]dtos.MembershipResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, _ := mapper.EntityToDto[entities.Membership, dtos.MembershipResponse](&entity)
		dtoItems[i] = *dto
	}
	return &model.PaginatedResult[dtos.MembershipResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}, nil
}

// GetUserCoverage returns the customers/orgs a user can reach through
// their memberships. Cache-aside via GetOrSetEx: try Redis -> on miss
// build by iterating user memberships, deduping customerIds, and
// resolving each via OrgService -> filter by type=customer -> shape the
// response DTO. 24h TTL.
func (s *MembershipService) GetUserCoverage(c ctx.Context, userId string) (*dtos.MeCoverageResponse, error) {
	cacheKey := fmt.Sprintf("user:coverage:%s", userId)
	cacheTTL := int(24 * 60 * 60)

	var response dtos.MeCoverageResponse
	_, err := s.deps.Cache.GetOrSetEx(common.GetOrSetParams{
		Ctx:      c,
		CacheKey: cacheKey,
		CacheTTL: cacheTTL,
		Dest:     &response,
		Callback: func() (interface{}, error) {
			return s.buildUserCoverage(c, userId)
		},
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetAllUserMemberships returns every membership reachable by a user —
// direct ones (assigneeType="user") plus memberships of any group the
// user belongs to. Group lookup uses GroupQueryService to keep the
// cross-module surface as a service port.
func (s *MembershipService) GetAllUserMemberships(c ctx.Context, userId string) ([]*entities.Membership, error) {
	directMemberships, err := s.deps.Repo.FindByUserId(c, &userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get direct memberships: %w", err)
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Membership] Found %d direct memberships for user=%s", len(directMemberships), userId))

	groupIds, err := s.deps.GroupQueryService.GetAllUserGroupIds(c, userId)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Membership] Failed to get group IDs for user=%s: %v", userId, err))
		groupIds = nil
	}

	if len(groupIds) == 0 {
		logger.Debug(fmt.Sprintf("[SERVICE:Membership] User=%s is not in any groups, returning %d direct memberships", userId, len(directMemberships)))
		return directMemberships, nil
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Membership] User=%s belongs to %d groups: %v", userId, len(groupIds), groupIds))

	groupMemberships, err := s.deps.Repo.FindByGroupIds(c, groupIds)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Membership] Failed to get group memberships: %v", err))
		return directMemberships, nil
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Membership] Found %d group memberships for user=%s", len(groupMemberships), userId))

	allMemberships := make([]*entities.Membership, 0, len(directMemberships)+len(groupMemberships))
	allMemberships = append(allMemberships, directMemberships...)
	allMemberships = append(allMemberships, groupMemberships...)
	logger.Info(fmt.Sprintf("[SERVICE:Membership] Total memberships for user=%s: %d (direct: %d, via groups: %d)",
		userId, len(allMemberships), len(directMemberships), len(groupMemberships)))
	return allMemberships, nil
}

// GetAllMemberships returns every membership matching the filters using
// cursor pagination internally — the caller receives the full set, not
// a single page. Used by cache builders / consumers that need exhaustive
// reads. Skips coverage filtering (internal-only).
func (s *MembershipService) GetAllMemberships(c ctx.Context, query *dtos.MembershipQueryDto) ([]*entities.Membership, error) {
	filters, err := s.buildAllMembershipsFilters(query)
	if err != nil {
		return nil, err
	}
	projection := orgfilter.BuildProjection(query.Projection)

	allMemberships := make([]*entities.Membership, 0)
	page := int64(1)
	batchSize := int64(300)

	for {
		result, err := s.deps.Repo.FindWithFilters(c, filters, &model.PaginationOpts{Page: page, PerPage: batchSize}, projection)
		if err != nil {
			return nil, err
		}
		for _, item := range result.Items {
			membershipCopy := item
			allMemberships = append(allMemberships, &membershipCopy)
		}
		if int64(len(result.Items)) < batchSize {
			break
		}
		page++
	}
	return allMemberships, nil
}

// GetDirectUserMemberships returns only direct (assigneeType="user")
// memberships for a user. Thin wrapper over the repository — no
// orchestration to surface.
func (s *MembershipService) GetDirectUserMemberships(c ctx.Context, userId string) ([]*entities.Membership, error) {
	return s.deps.Repo.FindByUserId(c, &userId)
}

// GetMembershipsByGroupIds returns memberships where assigneeType="group"
// and assigneeId is in the given list. Thin wrapper over the repository.
func (s *MembershipService) GetMembershipsByGroupIds(c ctx.Context, groupIds []string) ([]*entities.Membership, error) {
	return s.deps.Repo.FindByGroupIds(c, groupIds)
}

// GetAssigneeIdsByOrgIds collects unique assignee ids matching the org
// list and assignee type via batched cursor pagination (300/batch).
// Used by coverage cache builders.
func (s *MembershipService) GetAssigneeIdsByOrgIds(c ctx.Context, orgIds []model.ObjectId, assigneeType string) ([]string, error) {
	filters := model.Map{
		"orgId":        model.Map{"$in": orgIds},
		"assigneeType": assigneeType,
	}

	assigneeIdsMap := make(map[string]bool)
	page := int64(1)
	batchSize := int64(300)

	for {
		result, err := s.deps.Repo.FindWithFilters(c, filters, &model.PaginationOpts{Page: page, PerPage: batchSize}, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch memberships for assigneeType=%s: %w", assigneeType, err)
		}
		for _, membership := range result.Items {
			if !membership.AssigneeID.IsZero() {
				assigneeIdsMap[membership.AssigneeID.Hex()] = true
			}
		}
		if int64(len(result.Items)) < batchSize {
			break
		}
		page++
	}

	assigneeIds := make([]string, 0, len(assigneeIdsMap))
	for id := range assigneeIdsMap {
		assigneeIds = append(assigneeIds, id)
	}
	return assigneeIds, nil
}

// publishMembershipChanged emits the NATS cache-invalidation event after
// a user-membership create/update. Best-effort: log on failure; never
// blocks the caller because it runs in a goroutine.
func (s *MembershipService) publishMembershipChanged(membershipId, userId, orgId string) {
	logger.Info(fmt.Sprintf("[SERVICE:Membership] Membership changed for user=%s org=%s - publishing cache invalidation event", userId, orgId))
	event := events.NewMembershipChangedEvent(membershipId, userId, orgId, "", "", "")
	subject := fmt.Sprintf(contractsCacheInvalidation.MembershipChangedSubjectFormat, membershipId)
	if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{Subject: subject, Data: event}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Membership] Failed to publish MembershipChangedEvent for user=%s", userId))
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:Membership] Published MembershipChangedEvent for user=%s to subject=%s", userId, subject))
}

// publishMembershipDeleted emits the NATS cache-invalidation event after
// a user-membership delete. Best-effort.
func (s *MembershipService) publishMembershipDeleted(membershipId, userId, orgId string) {
	logger.Info(fmt.Sprintf("[SERVICE:Membership] Membership deleted for user=%s org=%s - publishing cache invalidation event", userId, orgId))
	event := events.NewMembershipDeletedEvent(membershipId, userId, orgId, "")
	subject := fmt.Sprintf(contractsCacheInvalidation.MembershipDeletedSubjectFormat, membershipId)
	if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{Subject: subject, Data: event}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Membership] Failed to publish MembershipDeletedEvent for user=%s", userId))
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:Membership] Published MembershipDeletedEvent for user=%s to subject=%s", userId, subject))
}
