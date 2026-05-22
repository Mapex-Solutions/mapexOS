package services

import (
	"mapexIam/src/modules/roles/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildRoleListFilters assembles the Mongo filter map for the role list.
// Roles are visible when ANY of (org-scope, system, ancestor) match —
// encoded via $or — minus the explicit isSystem=false / isTemplate=false
// exclusions, plus the optional name / scope / permission filters.
func (s *RoleService) buildRoleListFilters(rc *reqCtx.RequestContext, query *dtos.RoleQueryDto) model.Map {
	orConditions := []model.Map{}
	if orgFilter, _ := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc, Query: query}); len(orgFilter) > 0 {
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
	if query.Scope != nil && *query.Scope != "" {
		filters["scope"] = *query.Scope
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.Permission != nil && *query.Permission != "" {
		filters["permissions"] = *query.Permission
	}
	return filters
}

// buildRoleProjection wraps orgfilter.BuildProjection so the service
// signature stays focused on its inputs.
func (s *RoleService) buildRoleProjection(query *dtos.RoleQueryDto) model.Map {
	return orgfilter.BuildProjection(query.Projection)
}
