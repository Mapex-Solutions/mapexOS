package services

import (
	ctx "context"
	"fmt"

	"mapexIam/src/modules/organizations/application/dtos"
	"mapexIam/src/modules/organizations/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// loadParentOrgIfPresent fetches the parent organization (if specified)
// or returns nil for top-level orgs. Missing parent ids surface as 404.
func (s *OrganizationService) loadParentOrgIfPresent(c ctx.Context, parentOrgID *string) (*entities.Organization, error) {
	if parentOrgID == nil {
		return nil, nil
	}
	parent, err := s.deps.Repo.FindById(c, parentOrgID)
	if err != nil || parent == nil {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Parent organization not found"},
		}
	}
	return parent, nil
}

// assignHierarchyFields fills code, parentOrgId, pathKey, depth, and
// initial childCount on a freshly mapped entity.
func (s *OrganizationService) assignHierarchyFields(orgEntity *entities.Organization, dto *dtos.CreateOrganizationDto, parent *entities.Organization) {
	childCount := 0
	parentPathKey := ""
	if parent != nil {
		childCount = parent.ChildCount
		orgEntity.ParentOrgID = &parent.ID
		parentPathKey = parent.PathKey
	}
	orgEntity.Code = buildOrganizationCode(childCount, dto.Type)
	orgEntity.PathKey = buildPathKey(parentPathKey, orgEntity.Code)
	orgEntity.Depth = calculateDepth(parent)
	orgEntity.ChildCount = 0
}

// persistCustomerId computes the customerId based on org type/parent and
// persists it via a follow-up update — needed because customerId
// references the org's own _id which is only known after Create.
func (s *OrganizationService) persistCustomerId(c ctx.Context, createdOrg *entities.Organization, orgType string, parent *entities.Organization) *entities.Organization {
	customerID := determineCustomerID(orgType, createdOrg.ID, parent)
	if customerID == nil {
		return createdOrg
	}
	orgIDHex := createdOrg.ID.Hex()
	updated, _ := s.deps.Repo.FindByIdAndUpdate(c, &orgIDHex, map[string]any{"customerId": customerID})
	return updated
}

// incrementParentChildCount bumps the parent's childCount by one. No-op
// when there is no parent (top-level org).
func (s *OrganizationService) incrementParentChildCount(c ctx.Context, parent *entities.Organization) {
	if parent == nil {
		return
	}
	parentIDHex := parent.ID.Hex()
	s.deps.Repo.FindByIdAndUpdate(c, &parentIDHex, map[string]any{"childCount": parent.ChildCount + 1})
}

// buildOrgListFilters assembles the Mongo filter map for the org list,
// rewriting the orgfilter helper's "orgId" key to "parentOrgId" so the
// query returns direct children of the active scope.
func (s *OrganizationService) buildOrgListFilters(rc *reqCtx.RequestContext, query *dtos.OrganizationQueryDto) (model.Map, error) {
	filters := model.Map{}
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc, Query: query})
	if err != nil {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to build org filter: %v", err)},
		}
	}
	for k, v := range orgFilter {
		if k == "orgId" {
			filters["parentOrgId"] = v
		} else {
			filters[k] = v
		}
	}

	if query.Type != nil {
		filters["type"] = *query.Type
	}
	if query.ParentOrgID != nil {
		parentOrgObjectID, err := model.ToObjectID(*query.ParentOrgID)
		if err != nil {
			return nil, &customErrors.ServerCustomError{
				Code:   httpStatus.BAD_REQUEST,
				Errors: []string{"Invalid parentOrgId format"},
			}
		}
		filters["parentOrgId"] = parentOrgObjectID
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	if query.Depth != nil {
		filters["depth"] = *query.Depth
	}
	if query.Name != nil {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	return filters, nil
}
