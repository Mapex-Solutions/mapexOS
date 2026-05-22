package services

import (
	ctx "context"
	"fmt"

	events "mapexIam/src/modules/cache_invalidation/application/events"
	"mapexIam/src/modules/organizations/application/di"
	"mapexIam/src/modules/organizations/application/dtos"
	"mapexIam/src/modules/organizations/application/ports"
	"mapexIam/src/modules/organizations/domain/entities"

	common "github.com/Mapex-Solutions/MapexOS/contracts/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
	"github.com/Mapex-Solutions/mapexGoKit/utils/pathkey"
)

// Compile-time check to ensure OrganizationService implements OrganizationServicePort
var _ ports.OrganizationServicePort = (*OrganizationService)(nil)

// New creates and returns a new instance of OrganizationService.
func New(deps di.OrganizationServiceDependenciesInjection) ports.OrganizationServicePort {
	return &OrganizationService{deps: deps}
}

// CreateOrganization orchestrates the creation flow:
// resolve parent (404 on miss) -> map DTO to entity -> assign code, pathKey,
// depth, childCount -> persist -> set customerId once we have the new id ->
// bump parent.childCount -> publish hierarchy + created NATS events ->
// return the response DTO.
func (s *OrganizationService) CreateOrganization(c ctx.Context, dto *dtos.CreateOrganizationDto) (*dtos.OrganizationResponse, error) {
	parent, err := s.loadParentOrgIfPresent(c, dto.ParentOrgID)
	if err != nil {
		return nil, err
	}

	orgEntity, _ := mapper.DtoToEntity[dtos.CreateOrganizationDto, entities.Organization](dto)
	s.assignHierarchyFields(orgEntity, dto, parent)

	createdOrg, err := s.deps.Repo.Create(c, orgEntity)
	if err != nil {
		return nil, err
	}

	createdOrg = s.persistCustomerId(c, createdOrg, dto.Type, parent)
	s.incrementParentChildCount(c, parent)

	if parent != nil {
		s.publishOrgHierarchyChangedEvent(createdOrg.ID.Hex(), parent, "created")
	}
	s.publishOrgCreatedEvent(createdOrg)

	resp, _ := mapper.EntityToDto[entities.Organization, dtos.OrganizationResponse](createdOrg)
	return resp, nil
}

// GetOrganizationById fetches a single organization by id. Returns 404
// when the id is unknown.
func (s *OrganizationService) GetOrganizationById(c ctx.Context, organizationId *string) (*dtos.OrganizationResponse, error) {
	org, err := s.deps.Repo.FindById(c, organizationId)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Organization not found"}}
	}
	resp, _ := mapper.EntityToDto[entities.Organization, dtos.OrganizationResponse](org)
	return resp, nil
}

// UpdateOrganizationById orchestrates a partial update: load the prior
// entity (404 on miss) -> apply the patch in Mongo -> when AccessPolicy.RolePolicy
// flipped, publish a NATS event so coverage caches drop affected principals.
func (s *OrganizationService) UpdateOrganizationById(c ctx.Context, organizationId *string, dto *dtos.UpdateOrganizationDto) (*dtos.OrganizationResponse, error) {
	oldOrg, err := s.deps.Repo.FindById(c, organizationId)
	if err != nil {
		return nil, err
	}
	if oldOrg == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Organization not found"}}
	}

	fields, _ := mapper.DtoToMap(dto)
	updated, _ := s.deps.Repo.FindByIdAndUpdate(c, organizationId, fields)
	if updated.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Organization not found"}}
	}

	if accessPolicyChanged(oldOrg, dto) {
		go s.publishOrgAccessPolicyChanged(updated, oldOrg, dto)
	}

	resp, _ := mapper.EntityToDto[entities.Organization, dtos.OrganizationResponse](updated)
	return resp, nil
}

// DeleteOrganizationById orchestrates removal: load the org first (so we
// can later announce hierarchy change) -> delete from repository (404 on
// miss) -> publish a hierarchy-changed NATS event so coverage caches drop
// recursive memberships rooted on this branch.
func (s *OrganizationService) DeleteOrganizationById(c ctx.Context, organizationId *string) (map[string]bool, error) {
	org, _ := s.deps.Repo.FindById(c, organizationId)
	if err := s.deps.Repo.DeleteById(c, organizationId); err != nil {
		if err.Error() == "document not found" {
			return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Organization not found"}}
		}
		return nil, err
	}

	if org != nil && org.ParentOrgID != nil {
		parentIdHex := org.ParentOrgID.Hex()
		if parentOrg, _ := s.deps.Repo.FindById(c, &parentIdHex); parentOrg != nil {
			s.publishOrgHierarchyChangedEvent(org.ID.Hex(), parentOrg, "deleted")
		}
	}
	return map[string]bool{"success": true}, nil
}

// GetOrganizations orchestrates the paginated list:
// build the org filter (rewriting "orgId" to "parentOrgId" for direct
// children semantics) -> apply optional type/parentOrgId/enabled/depth/name
// filters -> resolve pagination + projection -> delegate to repository ->
// map entities to response DTOs.
func (s *OrganizationService) GetOrganizations(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.OrganizationQueryDto) (*model.PaginatedResult[dtos.OrganizationResponse], error) {
	filters, err := s.buildOrgListFilters(requestContext, query)
	if err != nil {
		return nil, err
	}

	pagination := &model.PaginationOpts{Page: int64(query.GetPage()), PerPage: int64(query.GetPerPage())}
	projection := orgfilter.BuildProjection(query.Projection)

	result, err := s.deps.Repo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]dtos.OrganizationResponse, len(result.Items))
	for i, entity := range result.Items {
		dto, _ := mapper.EntityToDto[entities.Organization, dtos.OrganizationResponse](&entity)
		dtoItems[i] = *dto
	}
	return &model.PaginatedResult[dtos.OrganizationResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}, nil
}

// GetOrganizationsTree orchestrates the cursor-paginated tree query:
// when scoped to a single org, lock the result to that org's pathKey
// range (recursive scope = self + descendants); otherwise return the
// global tree -> resolve cursor opts -> minimal projection -> delegate
// to repository -> shape into TreeItem + cursor info.
func (s *OrganizationService) GetOrganizationsTree(c ctx.Context, orgId *string, query *dtos.TreeQueryDto) (*dtos.TreeResponseDto, error) {
	filters, err := s.buildTreeFilters(c, orgId)
	if err != nil {
		return nil, err
	}

	cursorOpts := &model.CursorOpts{
		Cursor:    query.Cursor,
		Direction: model.CursorDirection(query.Direction),
		Limit:     query.Limit,
		SortAsc:   *query.SortAsc,
	}
	projection := model.Map{"_id": 1, "name": 1, "type": 1}

	result, err := s.deps.Repo.FindWithCursor(c, filters, cursorOpts, projection)
	if err != nil {
		return nil, err
	}

	treeItems := make([]dtos.TreeItemDto, len(result.Items))
	for i, org := range result.Items {
		treeItems[i] = dtos.TreeItemDto{ID: org.ID, Name: org.Name, Type: org.Type}
	}
	return &dtos.TreeResponseDto{
		Items: treeItems,
		Cursor: common.CursorInfo{
			Next:        result.NextCursor,
			Previous:    result.PrevCursor,
			HasNext:     result.HasNext,
			HasPrevious: result.HasPrevious,
		},
	}, nil
}

// GetChildOrganizationsByPathKey returns every descendant of the given
// pathKey via Mongo range query. Used by the coverage cache builder; no
// pagination since the consumer needs the full subtree.
func (s *OrganizationService) GetChildOrganizationsByPathKey(c ctx.Context, parentPathKey string) ([]entities.Organization, error) {
	filters := model.Map{
		"pathKey": model.Map{
			"$gt": parentPathKey,
			"$lt": parentPathKey + "~",
		},
	}
	pagination := &model.PaginationOpts{Page: 1, PerPage: 10000}
	projection := model.Map{
		"_id":          1,
		"name":         1,
		"type":         1,
		"pathKey":      1,
		"accessPolicy": 1,
	}

	result, err := s.deps.Repo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Organization] Failed to get children for pathKey=%s", parentPathKey))
		return nil, err
	}
	logger.Info(fmt.Sprintf("[SERVICE:Organization] Found %d children for pathKey=%s", len(result.Items), parentPathKey))
	return result.Items, nil
}

// accessPolicyChanged reports whether the update DTO carries a different
// AccessPolicy.RolePolicy from the existing entity. RolePolicy controls
// recursive role inheritance — a flip invalidates downstream caches.
func accessPolicyChanged(oldOrg *entities.Organization, dto *dtos.UpdateOrganizationDto) bool {
	return dto.AccessPolicy != nil &&
		dto.AccessPolicy.RolePolicy != "" &&
		dto.AccessPolicy.RolePolicy != oldOrg.AccessPolicy.RolePolicy
}

// publishOrgAccessPolicyChanged emits the NATS cache-invalidation event
// after a successful AccessPolicy flip. Runs in its own goroutine so the
// HTTP response latency is not bound to NATS publish.
func (s *OrganizationService) publishOrgAccessPolicyChanged(updated, oldOrg *entities.Organization, dto *dtos.UpdateOrganizationDto) {
	logger.Info(fmt.Sprintf("[SERVICE:Organization] AccessPolicy changed for org=%s - publishing cache invalidation event", updated.ID.Hex()))

	pathKeyStr := updated.PathKey
	event := events.NewOrgAccessPolicyChangedEvent(
		updated.ID.Hex(),
		oldOrg.AccessPolicy.RolePolicy,
		dto.AccessPolicy.RolePolicy,
		pathKeyStr,
		"",
	)

	subject := fmt.Sprintf("mapexos.cache.invalidation.organization.%s.access_policy.changed", updated.ID.Hex())
	if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{Subject: subject, Data: event}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Organization] Failed to publish OrgAccessPolicyChangedEvent for org=%s", updated.ID.Hex()))
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:Organization] Published OrgAccessPolicyChangedEvent for org=%s to subject=%s", updated.ID.Hex(), subject))
}

// buildTreeFilters returns the pathKey range filter restricting the tree
// to one org's recursive scope, or an empty filter for ROOT users.
func (s *OrganizationService) buildTreeFilters(c ctx.Context, orgId *string) (model.Map, error) {
	filters := model.Map{}
	if orgId == nil || *orgId == "" {
		return filters, nil
	}
	org, err := s.deps.Repo.FindById(c, orgId)
	if err != nil || org == nil {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.FORBIDDEN,
			Errors: []string{"Invalid organization context"},
		}
	}
	filters["pathKey"] = model.Map{
		"$gte": org.PathKey,
		"$lt":  pathkey.CalculateNextSiblingPathKey(org.PathKey),
	}
	return filters, nil
}
