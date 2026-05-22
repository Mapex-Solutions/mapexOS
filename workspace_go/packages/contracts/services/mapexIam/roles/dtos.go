package dtos

import (
	"errors"
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type RoleId struct {
	RoleId string `params:"roleId" validate:"required,mongoid"`
}

type RoleCreate struct {
	Name        string   `json:"name" validate:"required,min=3,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Permissions []string `json:"permissions" validate:"required,min=1,dive,min=1"`
	IsSystem    bool     `json:"isSystem"`
	IsTemplate  bool     `json:"isTemplate"`
	Scope       string   `json:"scope" validate:"required,oneof=global local"` // "global" | "local"

	// Multi-tenant hierarchical fields (populated automatically by service from RequestContext)
	OrgID   *model.ObjectId `json:"orgId,omitempty" validate:"omitempty"` // null = system role or MAPEX exclusive
	PathKey *string         `json:"pathKey,omitempty" validate:"omitempty"` // For hierarchical queries
}

// Transform validates business rules after basic validation.
// If isSystem = false and orgId != null, scope is required.
func (dto *RoleCreate) Transform() error {
	if len(dto.Permissions) == 0 {
		return errors.New("at least one permission is required")
	}
	// If organization role, scope must be valid
	if !dto.IsSystem && dto.OrgID != nil {
		if dto.Scope != "global" && dto.Scope != "local" {
			return errors.New("scope must be 'global' or 'local' for organization roles")
		}
	}
	return nil
}

type RoleUpdate struct {
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=500"`
	Permissions *[]string `json:"permissions,omitempty" validate:"omitempty,min=1,dive,min=1"`
}

// RoleQuery represents query parameters for listing roles.
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
//   - Name: filter by role name (partial match)
//   - IsSystem: filter system roles
//   - Scope: filter by scope (global, local)
//   - Permission: filter roles that have this specific permission
//
// Organization filtering is handled automatically via RequestContext:
//   - No manual orgId/customerId/pathKey needed
//   - Context-aware filtering via X-Org-Context header
//   - Hierarchical queries via includeChildren parameter
type RoleQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Name       *string `query:"name" validate:"omitempty,max=100"`
	IsSystem   *bool   `query:"isSystem" validate:"omitempty"`
	IsTemplate *bool   `query:"isTemplate" validate:"omitempty"`
	Scope      *string `query:"scope" validate:"omitempty,oneof=global local"`
	Permission *string `query:"permission" validate:"omitempty"`
}

type RoleResponse struct {
	ID          *common.ObjectID `json:"id,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Permissions *[]string        `json:"permissions,omitempty"`
	IsSystem   *bool            `json:"isSystem,omitempty"`
	IsTemplate *bool            `json:"isTemplate,omitempty"`

	// Multi-tenant hierarchical fields
	OrgId   *common.ObjectID `json:"orgId,omitempty"`
	PathKey *string          `json:"pathKey,omitempty"`
	Scope      *string          `json:"scope,omitempty"`

	Created *time.Time `json:"created,omitempty"`
	Updated *time.Time `json:"updated,omitempty"`
}

func (r *RoleResponse) SetCreated(t *time.Time) { r.Created = t }
func (r *RoleResponse) SetUpdated(t *time.Time) { r.Updated = t }
