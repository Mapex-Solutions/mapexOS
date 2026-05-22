package services

import (
	ctx "context"
	"fmt"
	"time"

	"router/src/modules/routegroups/application/constants"
	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/domain/entities"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// applyRouteGroupScope enforces multi-tenant scope rules on a Create DTO:
// - System resources clear org context.
// - Templates require Vendor/Customer org context.
// - Local resources require a non-empty org context.
// On non-system resources it also stamps OrgID + PathKey from the request.
func (s *RouteGroupService) applyRouteGroupScope(rc *reqCtx.RequestContext, dto *dtos.RouteGroupCreateDTO) error {
	if dto.IsSystem != nil && *dto.IsSystem {
		dto.OrgID = nil
		dto.PathKey = nil
		return nil
	}
	if dto.IsTemplate != nil && *dto.IsTemplate {
		if err := orgfilter.ValidateTemplateCreation(rc.OrgContextData.PathKey); err != nil {
			return &customErrors.ServerCustomError{Code: status.FORBIDDEN, Errors: []string{err.Error()}}
		}
	} else if err := orgfilter.ValidateOrgContextForNonSystem(rc); err != nil {
		return &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{err.Error()}}
	}
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		orgObjectId, _ := model.ToObjectID(*rc.OrgContext)
		dto.OrgID = &orgObjectId
	}
	if rc.OrgContextData != nil && rc.OrgContextData.PathKey != "" {
		pathKey := rc.OrgContextData.PathKey
		dto.PathKey = &pathKey
	}
	return nil
}

// persistNewRouteGroup maps the Create DTO to a domain entity and writes it
// to the repository. Returns BAD_REQUEST when the DTO cannot be mapped.
func (s *RouteGroupService) persistNewRouteGroup(c ctx.Context, dto *dtos.RouteGroupCreateDTO) (*entities.RouteGroup, error) {
	listEntity, err := mapper.DtoToEntity[dtos.RouteGroupCreateDTO, entities.RouteGroup](dto)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"failed to map route group data"}}
	}
	return s.deps.RouteGroupRepo.Create(c, listEntity)
}

// applyRouteGroupUpdate runs the partial update against the repository.
// Returns NOT_FOUND when no document matched the id.
func (s *RouteGroupService) applyRouteGroupUpdate(c ctx.Context, listId *string, fieldsToUpdate map[string]any) (*entities.RouteGroup, error) {
	listEntityData, err := s.deps.RouteGroupRepo.FindByIdAndUpdate(c, listId, fieldsToUpdate)
	if err != nil {
		return nil, err
	}
	if listEntityData.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"route group not found"}}
	}
	return listEntityData, nil
}

// applyRouteGroupUpdateFromDto bundles DTO->map conversion and the partial
// update so UpdateRouteGroupById can show its orchestration in 5-15 lines.
func (s *RouteGroupService) applyRouteGroupUpdateFromDto(c ctx.Context, listId *string, dto *dtos.RouteGroupUpdateDTO) (*entities.RouteGroup, error) {
	fieldsToUpdate, err := mapper.DtoToMap(dto)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"failed to map route group update data"}}
	}
	return s.applyRouteGroupUpdate(c, listId, fieldsToUpdate)
}

// runRouteGroupCreatePipeline executes the persist -> warm cache -> map
// response steps that follow scope validation in CreateRouteGroup. Kept as
// a single helper so the public method body stays in the orchestration
// budget and the metric outcome handling lives in one place.
func (s *RouteGroupService) runRouteGroupCreatePipeline(c ctx.Context, rc *reqCtx.RequestContext, dto *dtos.RouteGroupCreateDTO) (*dtos.RouteGroupResponse, error) {
	entity, err := s.persistNewRouteGroup(c, dto)
	if err != nil {
		return nil, err
	}
	s.warmCachesAfterCreate(c, rc, entity)
	return s.mapRouteGroupResponse(entity)
}

// warmCachesAfterCreate stores the new entity under its per-id key and drops
// the per-org counter cache so the next CountRouteGroups hits the repo.
func (s *RouteGroupService) warmCachesAfterCreate(c ctx.Context, rc *reqCtx.RequestContext, entity *entities.RouteGroup) {
	s.cacheRouteGroupEntity(c, entity.ID.Hex(), entity)
	s.invalidateCounterCacheForCreate(c, rc)
}

// fetchRouteGroupForDelete loads the RouteGroup so DeleteRouteGroupById can
// access its OrgId for counter-cache invalidation. 404 on miss.
func (s *RouteGroupService) fetchRouteGroupForDelete(c ctx.Context, listId *string) (*entities.RouteGroup, error) {
	routeGroup, err := s.deps.RouteGroupRepo.FindById(c, listId)
	if err != nil || routeGroup == nil {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"route group not found"}}
	}
	return routeGroup, nil
}

// deleteRouteGroupFromRepo removes the document and translates the driver's
// "document not found" string into the canonical 404 contract error.
func (s *RouteGroupService) deleteRouteGroupFromRepo(c ctx.Context, listId *string) error {
	if err := s.deps.RouteGroupRepo.DeleteById(c, listId); err != nil {
		if err.Error() == "document not found" {
			return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"route group not found"}}
		}
		return err
	}
	return nil
}

// loadCachedRouteGroup runs the cache-aside read against Redis and falls
// back to the repository on miss. cacheMiss is true when the callback ran.
func (s *RouteGroupService) loadCachedRouteGroup(c ctx.Context, routeGroupId *string) (*entities.RouteGroup, bool, error) {
	cacheKey := BuildRouteGroupCacheKey(*routeGroupId)
	var entityData entities.RouteGroup
	cacheMiss := false
	_, err := s.deps.CacheRepo.GetOrSetEx(common.GetOrSetParams{
		Ctx:      c,
		CacheKey: cacheKey,
		CacheTTL: int(constants.RouteGroupCacheTTL.Seconds()),
		Dest:     &entityData,
		Callback: func() (interface{}, error) {
			cacheMiss = true
			return s.deps.RouteGroupRepo.FindById(c, routeGroupId)
		},
	})
	return &entityData, cacheMiss, err
}

// cacheRouteGroupEntity warms the per-id Redis cache after a write.
func (s *RouteGroupService) cacheRouteGroupEntity(c ctx.Context, id string, entity *entities.RouteGroup) {
	cacheKey := BuildRouteGroupCacheKey(id)
	s.deps.CacheRepo.SetEx(c, cacheKey, entity, constants.RouteGroupCacheTTL)
}

// deleteRouteGroupCache evicts the per-id Redis entry after a delete.
func (s *RouteGroupService) deleteRouteGroupCache(c ctx.Context, listId *string) {
	cacheKey := BuildRouteGroupCacheKey(*listId)
	s.deps.CacheRepo.Del(c, cacheKey)
}

// invalidateCounterCacheForCreate drops the per-org counter cache after a
// successful create so the next CountRouteGroups call hits the repository.
func (s *RouteGroupService) invalidateCounterCacheForCreate(c ctx.Context, rc *reqCtx.RequestContext) {
	if rc.OrgContext != nil {
		counterKey := constants.BuildCounterCacheKey(*rc.OrgContext)
		_ = s.deps.AppCache.Del(c, counterKey)
	}
}

// invalidateCounterCacheForDelete mirrors the create-side invalidation but
// uses the deleted entity's OrgId since the request context may have moved on.
func (s *RouteGroupService) invalidateCounterCacheForDelete(c ctx.Context, routeGroup *entities.RouteGroup) {
	if routeGroup.OrgId != nil && !routeGroup.OrgId.IsZero() {
		counterKey := constants.BuildCounterCacheKey(routeGroup.OrgId.Hex())
		_ = s.deps.AppCache.Del(c, counterKey)
	}
}

// mapRouteGroupResponse converts a RouteGroup entity into its response DTO,
// translating mapper failures into the canonical 400 contract error.
func (s *RouteGroupService) mapRouteGroupResponse(entity *entities.RouteGroup) (*dtos.RouteGroupResponse, error) {
	listResponse, err := mapper.EntityToDto[entities.RouteGroup, dtos.RouteGroupResponse](entity)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"failed to map route group response"}}
	}
	return listResponse, nil
}

// buildRouteGroupListFilters assembles the Mongo filter for GetRouteGroups:
// org filter ($or with system and ancestor templates) + the optional
// per-field filters (name, version, enabled, kinds).
func (s *RouteGroupService) buildRouteGroupListFilters(rc *reqCtx.RequestContext, query *dtos.RouteGroupQueryDTO) model.Map {
	orConditions := []model.Map{}
	orgFilter, _ := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc, Query: query})
	if len(orgFilter) > 0 {
		orConditions = append(orConditions, orgFilter)
	}
	if query.IsSystem == nil || *query.IsSystem {
		orConditions = append(orConditions, model.Map{"isSystem": true})
	}
	if query.IsTemplate != nil && *query.IsTemplate {
		if templateFilter, err := orgfilter.BuildTemplateAncestorFilter(rc); err == nil && len(templateFilter) > 0 {
			orConditions = append(orConditions, templateFilter)
		}
	}
	filters := model.Map{}
	if len(orConditions) > 0 {
		filters["$or"] = orConditions
	}
	if query.IsSystem != nil && !*query.IsSystem {
		filters["isSystem"] = false
	}
	if query.IsTemplate != nil && !*query.IsTemplate {
		filters["isTemplate"] = false
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	if query.Version != nil && *query.Version != "" {
		filters["version"] = *query.Version
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if len(query.Kinds) > 0 {
		filters["routers.0"] = model.Map{"$exists": true}
		filters["routers"] = model.Map{
			"$not": model.Map{
				"$elemMatch": model.Map{"kind": model.Map{"$nin": query.Kinds}},
			},
		}
	}
	return filters
}

// buildRouteGroupListPagination derives the Mongo pagination opts from the
// query DTO using the shared per-page defaults.
func (s *RouteGroupService) buildRouteGroupListPagination(query *dtos.RouteGroupQueryDTO) *model.PaginationOpts {
	return &model.PaginationOpts{
		Page:    int64(query.GetPage()),
		PerPage: int64(query.GetPerPage()),
	}
}

// buildRouteGroupListProjection forwards the projection helper so the
// orchestration in GetRouteGroups stays free of orgfilter knowledge.
func (s *RouteGroupService) buildRouteGroupListProjection(query *dtos.RouteGroupQueryDTO) model.Map {
	return orgfilter.BuildProjection(query.Projection)
}

// mapRouteGroupListResult converts each entity in a paginated query result to
// its response DTO and re-wraps with the original pagination metadata.
func (s *RouteGroupService) mapRouteGroupListResult(result *model.PaginatedResult[entities.RouteGroup]) (*model.PaginatedResult[dtos.RouteGroupResponse], error) {
	dtoItems := make([]dtos.RouteGroupResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, err := mapper.EntityToDto[entities.RouteGroup, dtos.RouteGroupResponse](&entity)
		if err != nil {
			return nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"failed to map route group response"}}
		}
		dtoItems[i] = *dto
	}
	return &model.PaginatedResult[dtos.RouteGroupResponse]{Items: dtoItems, Pagination: result.Pagination}, nil
}

// buildCounterCacheKey forwards the constants helper so the orchestration in
// CountRouteGroups doesn't reach into the constants package directly.
func (s *RouteGroupService) buildCounterCacheKey(orgId string) string {
	return constants.BuildCounterCacheKey(orgId)
}

// tryCachedRouteGroupCount returns (count, true) on a Redis hit and (0, false)
// when the caller must fall back to the repository.
func (s *RouteGroupService) tryCachedRouteGroupCount(c ctx.Context, cacheKey string) (int64, bool) {
	var count int64
	if err := s.deps.AppCache.Get(c, cacheKey, &count); err == nil {
		return count, true
	}
	return 0, false
}

// countRouteGroupsFromRepo runs the org-scoped CountDocuments query.
func (s *RouteGroupService) countRouteGroupsFromRepo(c ctx.Context, rc *reqCtx.RequestContext) (int64, error) {
	logger.Debug(fmt.Sprintf("[SERVICE:RouteGroup] Counter cache miss for orgId=%s", orgIdFromContext(rc)))
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc})
	if err != nil {
		return 0, err
	}
	return s.deps.RouteGroupRepo.CountDocuments(c, orgFilter)
}

// cacheRouteGroupCount stores a fresh count under the per-org cache key with
// the standard counter TTL. Best-effort; failures are intentionally ignored.
func (s *RouteGroupService) cacheRouteGroupCount(c ctx.Context, cacheKey string, count int64) {
	_ = s.deps.AppCache.SetEx(c, cacheKey, count, constants.CounterCacheTTL)
}

// orgIdFromContext extracts the OrgContext id (or "") for log lines; kept as
// a tiny private to avoid duplicating the nil-check pattern.
func orgIdFromContext(rc *reqCtx.RequestContext) string {
	if rc.OrgContext != nil {
		return *rc.OrgContext
	}
	return ""
}

// recordRouteGroupOp emits the count + duration metrics for one CRUD attempt
// so every exit path of every public method stays observability-consistent.
func (s *RouteGroupService) recordRouteGroupOp(op, outcome string, start time.Time) {
	s.deps.Metrics.RouteGroupOperations.WithLabelValues(op, outcome).Inc()
	s.deps.Metrics.RouteGroupOperationDuration.WithLabelValues(op).Observe(time.Since(start).Seconds())
}

// recordRouteGroupError records the error metric for the given operation and
// returns the same error so call sites can short-circuit in a single line:
// `return nil, s.recordRouteGroupError("create", start, err)`.
func (s *RouteGroupService) recordRouteGroupError(op string, start time.Time, err error) error {
	s.recordRouteGroupOp(op, "error", start)
	return err
}

// recordRouteGroupCacheOutcome bumps the per-call hit/miss counter so the
// orchestration in GetRouteGroupEntityById stays a single named-step call.
func (s *RouteGroupService) recordRouteGroupCacheOutcome(miss bool) {
	if miss {
		s.deps.Metrics.RouteGroupCacheTotal.WithLabelValues("miss").Inc()
		return
	}
	s.deps.Metrics.RouteGroupCacheTotal.WithLabelValues("hit").Inc()
}

// recordRouteGroupListSuccess emits the per-page result-count histogram on top
// of the standard list-success metrics so the orchestration in GetRouteGroups
// stays a single named-step call.
func (s *RouteGroupService) recordRouteGroupListSuccess(start time.Time, dtoResult *model.PaginatedResult[dtos.RouteGroupResponse]) {
	s.recordRouteGroupOp("list", "success", start)
	s.deps.Metrics.RouteGroupListResultsCount.Observe(float64(len(dtoResult.Items)))
}
