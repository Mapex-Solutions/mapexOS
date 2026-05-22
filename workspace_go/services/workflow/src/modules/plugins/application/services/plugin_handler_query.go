package services

import (
	"workflow/src/modules/plugins/application/dtos"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildPluginListFilters assembles the multi-tenant visibility filter:
// org-local manifests, plus optional template manifests from the ancestor
// hierarchy when the query opts in. Field predicates (name, enabled,
// category, isTemplate exclusion) are merged on top.
func (s *PluginService) buildPluginListFilters(requestContext *reqCtx.RequestContext, query *dtos.PluginQueryDTO) model.Map {
	orConditions := s.buildPluginVisibilityConditions(requestContext, query)
	filters := model.Map{}
	if len(orConditions) > 0 {
		filters["$or"] = orConditions
	}
	if query.IsTemplate != nil && !*query.IsTemplate {
		filters["isTemplate"] = false
	}
	if query.Enabled != nil {
		filters["enabled"] = *query.Enabled
	}
	if query.Category != nil && *query.Category != "" {
		filters["category"] = *query.Category
	}
	if query.Name != nil && *query.Name != "" {
		filters["name"] = model.Map{"$regex": *query.Name, "$options": "i"}
	}
	return filters
}

// buildPluginVisibilityConditions collects the dynamic $or branches that
// drive multi-tenant visibility: my-org plus, when the caller asks, the
// template-ancestor branch. An empty result means a global view (system).
func (s *PluginService) buildPluginVisibilityConditions(requestContext *reqCtx.RequestContext, query *dtos.PluginQueryDTO) []model.Map {
	conditions := []model.Map{}
	orgFilter, _ := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{
		ReqContext: requestContext,
	})
	if len(orgFilter) > 0 {
		conditions = append(conditions, orgFilter)
	}
	if query.IsTemplate != nil && *query.IsTemplate {
		if templateFilter, err := orgfilter.BuildTemplateAncestorFilter(requestContext); err == nil && len(templateFilter) > 0 {
			conditions = append(conditions, templateFilter)
		}
	}
	return conditions
}

// buildPluginPagination converts the int-typed query DTO fields into the
// model's int64 PaginationOpts. Helper exists so the orchestration in
// service.go does not have to know the DTO/model field-type mismatch.
func (s *PluginService) buildPluginPagination(query *dtos.PluginQueryDTO) *model.PaginationOpts {
	return &model.PaginationOpts{
		Page:    int64(query.GetPage()),
		PerPage: int64(query.GetPerPage()),
	}
}
