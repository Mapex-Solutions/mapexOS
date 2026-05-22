package dtos

import (
	"errors"
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

type MembershipId struct {
	MembershipId string `params:"membershipId" validate:"required,mongoid"`
}

type MembershipCreate struct {
	AssigneeType string   `json:"assigneeType" validate:"required,oneof=user group"`
	AssigneeID   string   `json:"assigneeId" validate:"required,mongoid"`
	OrgID        string   `json:"orgId" validate:"required,mongoid"`
	RoleIds      []string `json:"roleIds" validate:"required,min=1,dive,mongoid"`
	Scope        string   `json:"scope" validate:"required,oneof=local recursive"`
	Enabled      bool     `json:"enabled"`
}

// Transform validates business rules after basic validation.
func (dto *MembershipCreate) Transform() error {
	if len(dto.RoleIds) == 0 {
		return errors.New("at least one role is required")
	}
	return nil
}

type MembershipUpdate struct {
	RoleIds *[]string `json:"roleIds,omitempty" validate:"omitempty,min=1,dive,mongoid"`
	Scope   *string   `json:"scope,omitempty" validate:"omitempty,oneof=local recursive"`
	Enabled *bool     `json:"enabled,omitempty"`
}

// MembershipQuery represents query parameters for listing memberships.
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
//   - AssigneeType: filter by assignee type (user, group)
//   - AssigneeID: filter by specific assignee
//   - UserID: filter by user (when assigneeType=user)
//   - RoleID: filter memberships with this role
//   - Scope: filter by scope (local, recursive)
//   - Enabled: filter by enabled status
//
// Organization filtering is handled automatically via RequestContext:
//   - No manual orgId/customerId/orgPathKey needed
//   - Context-aware filtering via X-Org-Context header
//   - Hierarchical queries via includeChildren parameter
type MembershipQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	AssigneeType *string `query:"assigneeType" validate:"omitempty,oneof=user group"`
	AssigneeID   *string `query:"assigneeId" validate:"omitempty,mongoid"`
	UserID       *string `query:"userId" validate:"omitempty,mongoid"`
	RoleID       *string `query:"roleId" validate:"omitempty,mongoid"`
	Scope        *string `query:"scope" validate:"omitempty,oneof=local recursive"`
	Enabled      *bool   `query:"enabled" validate:"omitempty"`

	// Internal-only filters (not exposed in HTTP routes, used for cache building)
	OrgID *string `query:"orgId" validate:"omitempty,mongoid"`
}

type MembershipResponse struct {
	ID           *common.ObjectID `json:"id,omitempty"`
	AssigneeType *string          `json:"assigneeType,omitempty"`
	AssigneeID   *common.ObjectID `json:"assigneeId,omitempty"`
	OrgID        *common.ObjectID   `json:"orgId,omitempty"`
	OrgPathKey   *string            `json:"orgPathKey,omitempty"`
	CustomerID   *common.ObjectID   `json:"customerId,omitempty"`
	RoleIds      *[]common.ObjectID `json:"roleIds,omitempty"`
	Scope        *string            `json:"scope,omitempty"`
	Enabled      *bool            `json:"enabled,omitempty"`
	Created      *time.Time       `json:"created,omitempty"`
	Updated      *time.Time       `json:"updated,omitempty"`
}

func (m *MembershipResponse) SetCreated(t *time.Time) { m.Created = t }
func (m *MembershipResponse) SetUpdated(t *time.Time) { m.Updated = t }

// CustomerCoverage represents a customer that the user has access to
type CustomerCoverage struct {
	CustomerID   *common.ObjectID `json:"customerId"`
	CustomerName *string          `json:"customerName"`
	PathKey      *string          `json:"pathKey"`
}

// MeCoverageResponse contains the list of customers/organizations a user has access to
type MeCoverageResponse struct {
	UserID    *common.ObjectID    `json:"userId"`
	Customers []*CustomerCoverage `json:"customers"`
}
