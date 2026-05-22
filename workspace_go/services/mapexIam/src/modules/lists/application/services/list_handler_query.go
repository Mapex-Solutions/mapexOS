package services

import (
	"mapexIam/src/modules/lists/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildListListFilters assembles the Mongo filter map for the lists list.
// Visibility = ANY of (org-scope, system, ancestor) — encoded via $or —
// minus the explicit isSystem=false / isTemplate=false exclusions, plus
// the optional type / name / enabled / parentId filters.
func (s *ListService) buildListListFilters(rc *reqCtx.RequestContext, query *dtos.ListQueryDTO) model.Map {
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
	if query.Type != nil && *query.Type != "" {
		filters["type"] = *query.Type
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	if query.ParentId != nil && *query.ParentId != "" {
		if parentIdObj, err := model.ToObjectID(*query.ParentId); err == nil {
			filters["parentId"] = parentIdObj
		}
	}
	return filters
}

// buildListProjection wraps orgfilter.BuildProjection so the service
// signature stays focused on its inputs.
func (s *ListService) buildListProjection(query *dtos.ListQueryDTO) model.Map {
	return orgfilter.BuildProjection(query.Projection)
}
