package services

import (
	ctx "context"
	"fmt"

	"mapexIam/src/modules/memberships/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildMembershipListFilters assembles the Mongo filter map for the
// HTTP-facing list path. The orgfilter helper applies coverage from
// RequestContext; module-specific filters layer on top.
func (s *MembershipService) buildMembershipListFilters(rc *reqCtx.RequestContext, query *dtos.MembershipQueryDto) (model.Map, error) {
	filters := model.Map{}
	orgFilter, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc, Query: query})
	if err != nil {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to build org filter: %v", err)},
		}
	}
	for k, v := range orgFilter {
		filters[k] = v
	}
	return applyMembershipQueryFilters(filters, query)
}

// buildAllMembershipsFilters mirrors the list-mode filter builder but
// skips coverage filtering (internal-only callers like cache builders).
// It supports an extra OrgID filter that the public path doesn't expose.
func (s *MembershipService) buildAllMembershipsFilters(query *dtos.MembershipQueryDto) (model.Map, error) {
	filters, err := applyMembershipQueryFilters(model.Map{}, query)
	if err != nil {
		return nil, err
	}
	if query.OrgID != nil && *query.OrgID != "" {
		orgObjectID, err := model.ToObjectID(*query.OrgID)
		if err != nil {
			return nil, &customErrors.ServerCustomError{
				Code:   httpStatus.BAD_REQUEST,
				Errors: []string{"Invalid orgId format"},
			}
		}
		filters["orgId"] = orgObjectID
	}
	return filters, nil
}

// applyMembershipQueryFilters layers the module-specific filters
// (assigneeType, assigneeId, userId, roleId, scope, enabled) onto an
// existing filter map. Returns a 400 for any malformed string id.
func applyMembershipQueryFilters(filters model.Map, query *dtos.MembershipQueryDto) (model.Map, error) {
	if query.AssigneeType != nil && *query.AssigneeType != "" {
		filters["assigneeType"] = *query.AssigneeType
	}
	if query.AssigneeID != nil && *query.AssigneeID != "" {
		assigneeObjectID, err := model.ToObjectID(*query.AssigneeID)
		if err != nil {
			return nil, &customErrors.ServerCustomError{
				Code:   httpStatus.BAD_REQUEST,
				Errors: []string{"Invalid assigneeId format"},
			}
		}
		filters["assigneeId"] = assigneeObjectID
	}
	if query.UserID != nil && *query.UserID != "" {
		userObjectID, err := model.ToObjectID(*query.UserID)
		if err != nil {
			return nil, &customErrors.ServerCustomError{
				Code:   httpStatus.BAD_REQUEST,
				Errors: []string{"Invalid userId format"},
			}
		}
		filters["assigneeType"] = "user"
		filters["assigneeId"] = userObjectID
	}
	if query.RoleID != nil && *query.RoleID != "" {
		roleObjectID, err := model.ToObjectID(*query.RoleID)
		if err != nil {
			return nil, &customErrors.ServerCustomError{
				Code:   httpStatus.BAD_REQUEST,
				Errors: []string{"Invalid roleId format"},
			}
		}
		filters["roleIds"] = roleObjectID
	}
	if query.Scope != nil && *query.Scope != "" {
		filters["scope"] = *query.Scope
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	return filters, nil
}

// buildUserCoverage assembles the MeCoverageResponse for a user by
// walking their memberships, deduping customerIds, resolving each via
// OrgService, and filtering down to type=customer entries.
func (s *MembershipService) buildUserCoverage(c ctx.Context, userId string) (*dtos.MeCoverageResponse, error) {
	memberships, err := s.deps.Repo.FindByUserId(c, &userId)
	if err != nil {
		return nil, err
	}

	customerIDsMap := make(map[string]bool)
	for _, membership := range memberships {
		if membership != nil && membership.CustomerID != nil && !membership.CustomerID.IsZero() {
			customerIDsMap[membership.CustomerID.Hex()] = true
		}
	}

	customerIDs := make([]string, 0, len(customerIDsMap))
	for id := range customerIDsMap {
		customerIDs = append(customerIDs, id)
	}

	if len(customerIDs) == 0 {
		userObjectID, _ := model.ToObjectID(userId)
		return &dtos.MeCoverageResponse{
			UserID:    &userObjectID,
			Customers: []*dtos.CustomerCoverage{},
		}, nil
	}

	customers := make([]*dtos.CustomerCoverage, 0)
	for _, customerID := range customerIDs {
		org, err := s.deps.OrgService.GetOrganizationById(c, &customerID)
		if err != nil || org == nil {
			continue
		}
		if org.Type != nil && *org.Type == "customer" && org.ID != nil {
			customerObjectID, _ := model.ToObjectID(*org.ID)
			customers = append(customers, &dtos.CustomerCoverage{
				CustomerID:   &customerObjectID,
				CustomerName: org.Name,
				PathKey:      org.PathKey,
			})
		}
	}

	userObjectID, _ := model.ToObjectID(userId)
	return &dtos.MeCoverageResponse{
		UserID:    &userObjectID,
		Customers: customers,
	}, nil
}
