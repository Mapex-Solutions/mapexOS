package dtos

// ListResolveRequest is the service-to-service request to resolve a
// classification slug into an org-scoped list id, creating the list when it does
// not exist yet. The asset-template install flow sends one per classification
// level (manufacturer, category, model). Resolution matches on (Type, Slug,
// OrgId) — never a global/system list — so a locally-installed vendor never
// leaks into another org.
type ListResolveRequest struct {
	Type              string  `json:"type" validate:"required"`
	Slug              string  `json:"slug" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	ParentId          *string `json:"parentId,omitempty" validate:"omitempty,mongoid"`
	OrgId             string  `json:"orgId" validate:"required,mongoid"`
	PathKey           *string `json:"pathKey,omitempty"`
	ShareWithChildren bool    `json:"shareWithChildren"`
}

// ListResolveResponse carries the resolved (found-or-created) list id.
type ListResolveResponse struct {
	Id string `json:"id"`
}
