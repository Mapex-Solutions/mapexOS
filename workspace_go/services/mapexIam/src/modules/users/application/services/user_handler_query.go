package services

import (
	ctx "context"
	"fmt"

	constants "mapexIam/src/modules/users/application/constants"
	"mapexIam/src/modules/users/application/dtos"
	"mapexIam/src/modules/users/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildUserListFilters assembles the Mongo filter for GetUsers. Multi-tenant
// filtering happens here: users do not have an orgId field of their own, so
// the org filter is resolved through MembershipService (direct + group paths).
//
// Returns (filters, true, nil) when the org filter resolves to zero accessible
// users — the caller short-circuits with an empty paginated response.
func (s *UserService) buildUserListFilters(c ctx.Context, rc *reqCtx.RequestContext, query *dtos.UserQueryDto) (model.Map, bool, error) {
	filters := model.Map{}
	if len(rc.ScopedOrgIds) > 0 {
		userObjectIDs, ok, err := s.resolveAccessibleUserIds(c, rc.ScopedOrgIds)
		if err != nil {
			return nil, false, err
		}
		if !ok {
			return filters, true, nil
		}
		filters["_id"] = model.Map{"$in": userObjectIDs}
	}
	s.applyUserListFieldFilters(filters, query)
	return filters, false, nil
}

// resolveAccessibleUserIds walks the two membership paths (direct user and
// via group) to compute the user ids the caller can list, then converts the
// final set to ObjectIDs for the user query. Returns ok=false when the set
// is empty so the caller can short-circuit.
func (s *UserService) resolveAccessibleUserIds(c ctx.Context, scopedOrgIds []string) ([]model.ObjectId, bool, error) {
	orgObjectIDs := make([]model.ObjectId, 0, len(scopedOrgIds))
	for _, orgId := range scopedOrgIds {
		if orgObjectID, err := model.ToObjectID(orgId); err == nil {
			orgObjectIDs = append(orgObjectIDs, orgObjectID)
		}
	}
	userIdsMap := make(map[string]bool)
	directUserIds, err := s.deps.MembershipService.GetAssigneeIdsByOrgIds(c, orgObjectIDs, "user")
	if err != nil {
		return nil, false, &customErrors.ServerCustomError{
			Code:   status.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to fetch user memberships for organization filtering: %v", err)},
		}
	}
	for _, userId := range directUserIds {
		userIdsMap[userId] = true
	}
	groupIds, err := s.deps.MembershipService.GetAssigneeIdsByOrgIds(c, orgObjectIDs, "group")
	if err == nil && len(groupIds) > 0 {
		groupUserIds, gErr := s.deps.GroupQueryService.GetUserIdsByGroupIds(c, groupIds)
		if gErr == nil {
			for _, userId := range groupUserIds {
				userIdsMap[userId] = true
			}
		}
	}
	if len(userIdsMap) == 0 {
		return nil, false, nil
	}
	userObjectIDs := make([]model.ObjectId, 0, len(userIdsMap))
	for userId := range userIdsMap {
		if userObjectID, err := model.ToObjectID(userId); err == nil {
			userObjectIDs = append(userObjectIDs, userObjectID)
		}
	}
	return userObjectIDs, true, nil
}

// applyUserListFieldFilters layers per-field filters (email, firstName,
// lastName, enabled) onto the filters map built from the org scoping.
func (s *UserService) applyUserListFieldFilters(filters model.Map, query *dtos.UserQueryDto) {
	if query.Email != nil && *query.Email != "" {
		filters["email"] = model.Map{"$regex": *query.Email, "$options": "i"}
	}
	if query.FirstName != nil && *query.FirstName != "" {
		filters["firstName"] = model.Map{"$regex": *query.FirstName, "$options": "i"}
	}
	if query.LastName != nil && *query.LastName != "" {
		filters["lastName"] = model.Map{"$regex": *query.LastName, "$options": "i"}
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
}

// emptyUserListResult builds the paginated response for the
// "no accessible users" short-circuit path.
func (s *UserService) emptyUserListResult(query *dtos.UserQueryDto) *model.PaginatedResult[dtos.UserResponse] {
	return &model.PaginatedResult[dtos.UserResponse]{
		Items: []dtos.UserResponse{},
		Pagination: model.Pagination{
			Page:       int64(query.GetPage()),
			PerPage:    int64(query.GetPerPage()),
			TotalItems: 0,
			TotalPages: 0,
		},
	}
}

// buildUserListPagination derives the Mongo pagination opts from the
// query DTO using the shared per-page defaults.
func (s *UserService) buildUserListPagination(query *dtos.UserQueryDto) *model.PaginationOpts {
	return &model.PaginationOpts{
		Page:    int64(query.GetPage()),
		PerPage: int64(query.GetPerPage()),
	}
}

// buildUserListProjection forwards the projection helper so the
// orchestration in GetUsers stays free of orgfilter knowledge.
func (s *UserService) buildUserListProjection(query *dtos.UserQueryDto) model.Map {
	return orgfilter.BuildProjection(query.Projection)
}

// mapUserListResultWithGroupCounts converts the entity list into DTOs and
// enriches each one with its groupsCount via the GroupQueryService batch.
func (s *UserService) mapUserListResultWithGroupCounts(c ctx.Context, result *model.PaginatedResult[entities.User]) *model.PaginatedResult[dtos.UserResponse] {
	dtoItems := make([]dtos.UserResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, _ := mapper.EntityToDto[entities.User, dtos.UserResponse](&entity)
		dtoItems[i] = *dto
	}
	if len(dtoItems) == 0 {
		return &model.PaginatedResult[dtos.UserResponse]{Items: dtoItems, Pagination: result.Pagination}
	}
	userIds := make([]string, len(dtoItems))
	for i, dto := range dtoItems {
		if dto.ID != nil {
			userIds[i] = dto.ID.Hex()
		}
	}
	groupCountMap, err := s.deps.GroupQueryService.CountGroupsByUserIds(c, userIds)
	if err == nil {
		for i := range dtoItems {
			if dtoItems[i].ID != nil {
				userId := dtoItems[i].ID.Hex()
				count := groupCountMap[userId]
				dtoItems[i].GroupsCount = &count
			}
		}
	}
	return &model.PaginatedResult[dtos.UserResponse]{Items: dtoItems, Pagination: result.Pagination}
}

// userOrgIdFromContext extracts the orgId string from the request context
// (or "" when absent) for the per-org counter cache key.
func (s *UserService) userOrgIdFromContext(rc *reqCtx.RequestContext) string {
	if rc.OrgContext != nil {
		return *rc.OrgContext
	}
	return ""
}

// tryCachedUserCount returns (count, true) on a Redis hit; (0, false) when
// the caller must fall back to the repository.
func (s *UserService) tryCachedUserCount(c ctx.Context, cacheKey string) (int64, bool) {
	var count int64
	if err := s.deps.AppCache.Get(c, cacheKey, &count); err == nil {
		return count, true
	}
	return 0, false
}

// countUsersFromRepo runs the org-scoped CountDocuments query.
func (s *UserService) countUsersFromRepo(c ctx.Context, rc *reqCtx.RequestContext) (int64, error) {
	logger.Debug(fmt.Sprintf("[SERVICE:User] Counter cache miss for orgId=%s", s.userOrgIdFromContext(rc)))
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc})
	if err != nil {
		return 0, err
	}
	return s.deps.Repo.CountDocuments(c, orgFilter)
}

// cacheUserCount stores a fresh count under the per-org cache key.
// Best-effort; failures are intentionally ignored.
func (s *UserService) cacheUserCount(c ctx.Context, cacheKey string, count int64) {
	_ = s.deps.AppCache.SetEx(c, cacheKey, count, constants.CounterCacheTTL)
}
