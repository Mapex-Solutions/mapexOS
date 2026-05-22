package dtos

import (
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

type OrganizationId struct {
	OrganizationId string `params:"organizationId" validate:"required,mongoid"`
}

type Address struct {
	City    string `json:"city,omitempty" validate:"omitempty,max=100"`
	State   string `json:"state,omitempty" validate:"omitempty,max=100"`
	Country string `json:"country,omitempty" validate:"omitempty,max=100"`
	ZipCode string `json:"zipCode,omitempty" validate:"omitempty,max=20"`
}

type AuthConfig struct {
	ProviderType     string                 `json:"providerType" validate:"required,oneof=keycloak internal"`
	IssuerURL        *string                `json:"issuerUrl,omitempty"`
	ClientID         *string                `json:"clientId,omitempty"`
	JWTClaimMappings map[string]string      `json:"jwtClaimMappings,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type AccessPolicy struct {
	RolePolicy   string `json:"rolePolicy" validate:"required,oneof=strict merge"`
	DefaultScope string `json:"defaultScope" validate:"required,oneof=local recursive"`

	// REMOVED: AllowDirectPermissions
	// V1 uses pure role-based architecture (no direct permissions)
}

type OrganizationCreate struct {
	Name         string       `json:"name" validate:"required,min=3,max=150"`
	Type         string       `json:"type" validate:"required,oneof=vendor customer site building floor zone"`
	ParentOrgID  *string      `json:"parentOrgId,omitempty" validate:"omitempty,mongoid"`
	Enabled      bool         `json:"enabled"`
	Address      *Address     `json:"address,omitempty" validate:"omitempty"`
	Phone        string       `json:"phone,omitempty" validate:"omitempty,e164"`
	AuthConfig   AuthConfig   `json:"authConfig" validate:"required"`
	AccessPolicy AccessPolicy `json:"accessPolicy" validate:"required"`
}

// Transform performs any necessary transformations on the DTO.
func (dto *OrganizationCreate) Transform() error {
	return nil
}

// GetType implements mapexGoKit orghierarchy.OrganizationCreateContract.
func (dto *OrganizationCreate) GetType() string {
	return dto.Type
}

// GetParentOrgID implements mapexGoKit orghierarchy.OrganizationCreateContract.
func (dto *OrganizationCreate) GetParentOrgID() *string {
	return dto.ParentOrgID
}

type OrganizationUpdate struct {
	Name         *string       `json:"name,omitempty" validate:"omitempty,min=3,max=150"`
	Type         *string       `json:"type,omitempty" validate:"omitempty,oneof=vendor customer site building floor zone"`
	ParentOrgID  *string       `json:"parentOrgId,omitempty" validate:"omitempty,mongoid"`
	Enabled      *bool         `json:"enabled,omitempty"`
	Address      *Address      `json:"address,omitempty"`
	Phone        *string       `json:"phone,omitempty" validate:"omitempty,e164"`
	AuthConfig   *AuthConfig   `json:"authConfig,omitempty"`
	AccessPolicy *AccessPolicy `json:"accessPolicy,omitempty"`
}

// Transform performs any necessary transformations on the DTO.
func (dto *OrganizationUpdate) Transform() error {
	return nil
}

// OrganizationQuery represents query parameters for listing organizations.
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
//   - Type: filter by organization type (vendor, customer, site, building, floor, zone)
//   - ParentOrgID: filter by parent organization
//   - Name: filter by name (partial match)
//   - Enabled: filter by enabled status
//   - Depth: filter by depth in hierarchy
//
// Organization filtering is handled automatically via RequestContext:
//   - No manual customerId/pathKey needed
//   - Context-aware filtering via X-Org-Context header
//   - Hierarchical queries via includeChildren parameter
type OrganizationQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Type        *string `query:"type" validate:"omitempty,oneof=vendor customer site building floor zone"`
	ParentOrgID *string `query:"parentOrgId" validate:"omitempty,mongoid"`
	Name        *string `query:"name" validate:"omitempty,max=150"`
	Enabled     *bool   `query:"enabled" validate:"omitempty"`
	Depth       *int    `query:"depth" validate:"omitempty,min=0,max=10"`
}

type OrganizationResponse struct {
	ID           *common.ObjectID  `json:"id,omitempty"`
	Name         *string           `json:"name,omitempty"`
	Type         *string           `json:"type,omitempty"`
	ParentOrgID  *common.ObjectID  `json:"parentOrgId,omitempty"`
	Code         *string           `json:"code,omitempty"`
	PathKey      *string           `json:"pathKey,omitempty"`
	Depth        *int              `json:"depth,omitempty"`
	CustomerID   *common.ObjectID  `json:"customerId,omitempty"`
	ChildCount   *int              `json:"childCount,omitempty"`
	Enabled      *bool             `json:"enabled,omitempty"`
	Address      *Address          `json:"address,omitempty"`
	Phone        *string           `json:"phone,omitempty"`
	AuthConfig   *AuthConfig       `json:"authConfig,omitempty"`
	AccessPolicy *AccessPolicy     `json:"accessPolicy,omitempty"`
	Created      *time.Time        `json:"created,omitempty"`
	Updated      *time.Time        `json:"updated,omitempty"`

}

func (u *OrganizationResponse) SetCreated(t *time.Time) { u.Created = t }
func (u *OrganizationResponse) SetUpdated(t *time.Time) { u.Updated = t }

// TreeQuery contains query parameters for hierarchical tree navigation.
// This endpoint uses cursor-based pagination for efficient navigation
// through organization hierarchies in UI components (selects, tree views).
type TreeQuery struct {
	// Cursor is the ID to start from (hex string of ObjectID).
	// Empty string means start from the beginning.
	Cursor string `query:"cursor" validate:"omitempty,hexadecimal,len=24"`

	// Direction specifies pagination direction: "next" or "previous".
	Direction string `query:"direction" validate:"omitempty,oneof=next previous" default:"next"`

	// Limit is the maximum number of items to return.
	Limit int64 `query:"limit" validate:"omitempty,min=1,max=300" default:"50"`

	// SortAsc controls sort direction. true = ascending, false = descending.
	SortAsc *bool `query:"sortAsc" validate:"omitempty" default:"true"`
}

// TreeItem represents minimal organization data for tree navigation.
// Contains only essential fields needed for UI components to display
// hierarchical structures (organization pickers, tree views).
type TreeItem struct {
	ID   common.ObjectID `json:"id"`
	Name string          `json:"name"`
	Type string          `json:"type"` // vendor, customer, site, building, floor, zone
}

// TreeResponse contains paginated tree items with cursor navigation metadata.
type TreeResponse struct {
	Items  []TreeItem        `json:"items"`
	Cursor common.CursorInfo `json:"cursor"`
}
