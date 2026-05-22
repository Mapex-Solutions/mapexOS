package dtos

import (
	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type ListId struct {
	ListId string `params:"listId" validate:"required,mongoid"`
}

type ListCreate struct {
	Type    string `json:"type" validate:"required,max=100"`
	Name    string `json:"name" validate:"required,max=254"`
	Value   string `json:"value" validate:"required,max=254"`
	Enabled bool   `json:"enabled" validate:"required"`

	// Hierarchical reference - links to parent list item (e.g., manufacturer for asset types)
	ParentId *model.ObjectId `json:"parentId,omitempty" validate:"omitempty,mongoid"`

	// Flexible metadata for type-specific fields
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Template visibility flags
	IsSystem   bool `json:"isSystem"`
	IsTemplate bool `json:"isTemplate"`

	// Multi-tenant fields (populated automatically by service from RequestContext)
	OrgID   *model.ObjectId `json:"orgId,omitempty" validate:"omitempty"`
	PathKey *string         `json:"pathKey,omitempty" validate:"omitempty"`
}

type ListUpdate struct {
	Name     *string                `json:"name,omitempty" validate:"omitempty,max=254"`
	Value    *string                `json:"value,omitempty" validate:"omitempty,max=254"`
	Enabled  *bool                  `json:"enabled,omitempty" validate:"omitempty"`
	ParentId *model.ObjectId        `json:"parentId,omitempty" validate:"omitempty,mongoid"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ListQuery represents query parameters for listing lists.
// Embeds BaseQueryDTO for standard pagination, sorting, and hierarchy support.
//
// Standard fields (from BaseQueryDTO):
//   - Projection: comma-separated fields to return
//   - Page: page number (default: 1)
//   - PerPage: items per page (default: 20)
//   - Sort: sort order (default: "created:desc")
//   - IncludeChildren: include child orgs hierarchically (default: false)
//
// Module-specific filters:
//   - Type: filter by list type (e.g., asset_manufacturer, asset_model, asset_category, or any custom type)
//   - IsSystem: filter system lists
//   - IsTemplate: filter template lists
//   - Name: filter by name (partial match)
//   - Enabled: filter by enabled status (true/false)
//   - ParentId: filter by parent list item ID (for hierarchical lists)
//
// Organization filtering is handled automatically via RequestContext:
//   - No manual orgId needed
//   - Context-aware filtering via X-Org-Context header
//   - Hierarchical queries via includeChildren parameter
type ListQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Type       *string `query:"type" validate:"omitempty"`
	IsSystem   *bool   `query:"isSystem" validate:"omitempty"`
	IsTemplate *bool   `query:"isTemplate" validate:"omitempty"`
	Name       *string `query:"name" validate:"omitempty,max=254"`
	Enabled    *bool   `query:"enabled" validate:"omitempty"`
	ParentId   *string `query:"parentId" validate:"omitempty,mongoid"`
}

type ListResponse struct {
	ID       *common.ObjectID       `json:"id,omitempty"`
	Type     *string                `json:"type,omitempty"`
	Name     *string                `json:"name,omitempty"`
	Value    *string                `json:"value,omitempty"`
	Enabled  *bool                  `json:"enabled,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Hierarchical data - parent info (populated from parentId lookup)
	ParentId   *model.ObjectId `json:"parentId,omitempty"`
	ParentName *string         `json:"parentName,omitempty"`
	ParentType *string         `json:"parentType,omitempty"`

	// Template visibility flags
	IsSystem   *bool `json:"isSystem,omitempty"`
	IsTemplate *bool `json:"isTemplate,omitempty"`

	// Multi-tenant fields
	OrgId   *model.ObjectId `json:"orgId,omitempty"`
	PathKey *string         `json:"pathKey,omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

func (u *ListResponse) SetCreated(t *common.NullTime) { u.Created = t }
func (u *ListResponse) SetUpdated(t *common.NullTime) { u.Updated = t }
