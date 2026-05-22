package dtos

import (
	"errors"
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type GroupId struct {
	GroupId string `params:"groupId" validate:"required,mongoid"`
}

type GroupCreate struct {
	Name        string  `json:"name" validate:"required,min=3,max=150"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Enabled     bool    `json:"enabled"`

	// Roles for group membership (ALWAYS REQUIRED)
	// When a group is created, a Membership is automatically created with these roles.
	// This ensures groups ALWAYS have associated permissions.
	// Users cannot create groups without roles - this is enforced at API level.
	RoleIds []string `json:"roleIds" validate:"required,min=1,dive,mongoid"`

	// Multi-tenant fields (populated automatically by service from RequestContext)
	OrgID   *model.ObjectId `json:"orgId,omitempty" validate:"omitempty"`
	PathKey *string         `json:"pathKey,omitempty" validate:"omitempty"`
}

// Transform validates business rules after basic validation.
// - OrgID is required (populated from RequestContext by service)
// - RoleIds must have at least 1 role (groups MUST have permissions)
func (dto *GroupCreate) Transform() error {
	if dto.OrgID == nil {
		return errors.New("orgId is required for group creation")
	}
	if len(dto.RoleIds) == 0 {
		return errors.New("at least one role is required when creating a group")
	}
	return nil
}

type GroupUpdate struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=150"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Enabled     *bool   `json:"enabled,omitempty"`
}

// GroupQuery represents query parameters for listing groups.
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
//   - Name: filter by group name (partial match)
//   - Enabled: filter by enabled status
//   - MemberId: filter groups by member (user ID)
//
// Organization filtering is handled automatically via RequestContext:
//   - No manual orgId/customerId/pathKey needed
//   - Context-aware filtering via X-Org-Context header
//   - Hierarchical queries via includeChildren parameter
type GroupQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Name    *string `query:"name" validate:"omitempty,max=150"`
	Enabled *bool   `query:"enabled" validate:"omitempty"`
	MemberId *string `query:"memberId" validate:"omitempty,mongoid"`
}

type GroupResponse struct {
	ID          *common.ObjectID `json:"id,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`

	// MembersCount - total count of members in this group
	// Calculated from group_members junction collection for scalability
	MembersCount *int64 `json:"membersCount,omitempty"`

	// RoleIds - roles associated with this group via Membership
	// Only populated in GetGroupById, NOT in list operations (performance)
	// Frontend resolves role names via cache
	RoleIds []string `json:"roleIds,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	OrgID   *common.ObjectID `json:"orgId,omitempty"`
	PathKey *string          `json:"pathKey,omitempty"`
	Created *time.Time       `json:"created,omitempty"`
	Updated *time.Time       `json:"updated,omitempty"`
}

// GroupMemberResponse represents a member of a group
// Used for paginated member listing
type GroupMemberResponse struct {
	ID        *common.ObjectID `json:"id,omitempty"`
	UserID    *common.ObjectID `json:"userId,omitempty"`
	GroupID   *common.ObjectID `json:"groupId,omitempty"`
	OrgID     *common.ObjectID `json:"orgId,omitempty"`
	AddedAt   *time.Time       `json:"addedAt,omitempty"`
	AddedBy   *common.ObjectID `json:"addedBy,omitempty"`
	Created   *time.Time       `json:"created,omitempty"`

	// Denormalized user info (populated by service)
	UserEmail     *string `json:"userEmail,omitempty"`
	UserFirstName *string `json:"userFirstName,omitempty"`
	UserLastName  *string `json:"userLastName,omitempty"`
}

// GroupMembersQuery represents query parameters for listing group members
type GroupMembersQuery struct {
	query.BaseQueryDTO
}

// GroupMemberAddDto represents the request body for adding a member to a group
type GroupMemberAddDto struct {
	UserID string `json:"userId" validate:"required,mongoid"`
}

// GroupMemberIdDto represents path parameters for member operations
type GroupMemberIdDto struct {
	GroupId string `params:"groupId" validate:"required,mongoid"`
	UserId  string `params:"userId" validate:"required,mongoid"`
}

func (g *GroupMemberResponse) SetCreated(t *time.Time) { g.Created = t }

func (g *GroupResponse) SetCreated(t *time.Time) { g.Created = t }
func (g *GroupResponse) SetUpdated(t *time.Time) { g.Updated = t }
