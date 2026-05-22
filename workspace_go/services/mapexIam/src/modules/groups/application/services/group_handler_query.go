package services

import (
	ctx "context"
	"fmt"

	"mapexIam/src/modules/groups/application/dtos"
	"mapexIam/src/modules/groups/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildGroupListFilters assembles the Mongo filter for GetGroups: org
// filter + name/enabled filters + the optional MemberId narrowing (looks
// up GroupMember rows for the user and filters groups by id $in).
//
// Returns (filters, true) when MemberId narrows the result to zero matches
// so the caller can short-circuit with an empty paginated response.
func (s *GroupService) buildGroupListFilters(c ctx.Context, rc *reqCtx.RequestContext, query *dtos.GroupQueryDto) (model.Map, bool) {
	orgFilter, _ := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc, Query: query})
	filters := model.Map{}
	for k, v := range orgFilter {
		filters[k] = v
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.MemberId != nil && *query.MemberId != "" {
		groupIds, ok := s.resolveMemberGroupIds(c, *query.MemberId)
		if !ok {
			return filters, true
		}
		filters["_id"] = model.Map{"$in": groupIds}
	}
	return filters, false
}

// resolveMemberGroupIds returns (ids, true) when the user has at least one
// group; (_, false) signals an empty result and the caller short-circuits.
func (s *GroupService) resolveMemberGroupIds(c ctx.Context, memberId string) ([]model.ObjectId, bool) {
	memberships, err := s.deps.GroupMemberRepo.FindByUserId(c, memberId)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Group] Failed to query memberships for user %s: %v", memberId, err))
		return nil, false
	}
	if len(memberships) == 0 {
		return nil, false
	}
	groupIds := make([]model.ObjectId, len(memberships))
	for i, m := range memberships {
		groupIds[i] = m.GroupID
	}
	return groupIds, true
}

// emptyGroupListResult builds the paginated response for the
// "user is not a member of any group" short-circuit path.
func (s *GroupService) emptyGroupListResult(query *dtos.GroupQueryDto) *model.PaginatedResult[dtos.GroupResponse] {
	return &model.PaginatedResult[dtos.GroupResponse]{
		Items: []dtos.GroupResponse{},
		Pagination: model.Pagination{
			Page:       int64(query.GetPage()),
			PerPage:    int64(query.GetPerPage()),
			TotalItems: 0,
			TotalPages: 0,
		},
	}
}

// buildGroupListPagination derives the Mongo pagination opts from the
// query DTO using the shared per-page defaults.
func (s *GroupService) buildGroupListPagination(query *dtos.GroupQueryDto) *model.PaginationOpts {
	return &model.PaginationOpts{
		Page:    int64(query.GetPage()),
		PerPage: int64(query.GetPerPage()),
	}
}

// buildGroupListProjection forwards the projection helper so the
// orchestration in GetGroups stays free of orgfilter knowledge.
func (s *GroupService) buildGroupListProjection(query *dtos.GroupQueryDto) model.Map {
	return orgfilter.BuildProjection(query.Projection)
}

// mapGroupListResultWithCounts converts the entity list into DTOs and
// enriches each one with the live member count (single batch lookup).
func (s *GroupService) mapGroupListResultWithCounts(c ctx.Context, result *model.PaginatedResult[entities.Group]) *model.PaginatedResult[dtos.GroupResponse] {
	groupIds := make([]string, len(result.Items))
	for i, entity := range result.Items {
		groupIds[i] = entity.ID.Hex()
	}
	memberCounts, err := s.deps.GroupMemberRepo.CountByGroupIds(c, groupIds)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Group] Failed to get member counts: %v", err))
		memberCounts = map[string]int64{}
	}
	dtoItems := make([]dtos.GroupResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, _ := mapper.EntityToDto[entities.Group, dtos.GroupResponse](&entity)
		groupIdHex := entity.ID.Hex()
		if count, ok := memberCounts[groupIdHex]; ok {
			dto.MembersCount = &count
		} else {
			zero := int64(0)
			dto.MembersCount = &zero
		}
		dtoItems[i] = *dto
	}
	return &model.PaginatedResult[dtos.GroupResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}
}
