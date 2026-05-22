package dtos

import (
	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

type UserId struct {
	UserId string `params:"userId" validate:"required,mongoid"`
}

type AuthProvider struct {
	Type       string                 `json:"type" validate:"required,oneof=internal google github microsoft keycloak"`
	ExternalID *string                `json:"externalId,omitempty" validate:"omitempty,min=1"`
	Metadata   map[string]interface{} `json:"metadata" validate:"omitempty"` // if optional: `omitempty`
}

type UserCreate struct {
	Email                   string  `json:"email" validate:"required,email,max=254"`
	Password                *string `json:"password,omitempty" validate:"omitempty,min=8,max=72"`
	ChangePasswordNextLogin bool    `json:"changePasswordNextLogin" validate:"-"`
	FirstName               string  `json:"firstName" validate:"required,min=2,max=100"`
	LastName                string  `json:"lastName" validate:"required,min=2,max=100"`
	Phone                   *string `json:"phone,omitempty" validate:"omitempty,e164"`
	JobTitle                *string `json:"jobTitle,omitempty" validate:"omitempty,max=120"`
	Enabled                 bool    `json:"enabled" validate:"-"`
	Avatar                  *string `json:"avatar,omitempty" validate:"omitempty,url"`
	StartTour               bool    `json:"startTour" validate:"-"`
	// AuthProvider removed from create DTO - V1 always uses internal auth.
	// Next version: auth provider will be determined by the customer's Organization.AuthConfig
}

// AuthProviderUpdate removed from contracts - V1 always uses internal auth.
// Next version: will be restored when Organization.AuthConfig drives user auth

type UserUpdate struct {
	Email                   *string `json:"email,omitempty" validate:"omitempty,email,max=254"`
	Password                *string `json:"password,omitempty" validate:"omitempty,min=8,max=72"`
	ChangePasswordNextLogin *bool   `json:"changePasswordNextLogin,omitempty" validate:"omitempty"`

	FirstName *string `json:"firstName,omitempty" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"lastName,omitempty" validate:"omitempty,min=2,max=100"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,e164"`
	JobTitle  *string `json:"jobTitle,omitempty" validate:"omitempty,max=120"`
	Enabled   *bool   `json:"enabled,omitempty" validate:"omitempty"`
	Avatar    *string `json:"avatar,omitempty" validate:"omitempty,url"`
	StartTour *bool   `json:"startTour,omitempty" validate:"omitempty"`
	// AuthProvider removed from update DTO - V1 always uses internal auth.
	// Next version: auth provider will be determined by the customer's Organization.AuthConfig
}

// UserGroupInfo represents a group the user belongs to (for detail view)
type UserGroupInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// UserMembershipInfo represents an organization membership (for detail view)
type UserMembershipInfo struct {
	OrgID     string   `json:"orgId"`
	OrgName   string   `json:"orgName"`
	OrgType   string   `json:"orgType"` // vendor, partner, customer
	Scope     string   `json:"scope"`   // local, recursive
	RoleNames []string `json:"roleNames"`
	Via       string   `json:"via"` // "direct" or "Group: {groupName}"
}

type UserResponse struct {
	ID                      *common.ObjectID `json:"id,omitempty"`
	Email                   *string          `json:"email,omitempty"`
	ChangePasswordNextLogin *bool            `json:"changePasswordNextLogin,omitempty"`
	AuthProvider            *AuthProvider    `json:"authProvider,omitempty"`
	FirstName               *string          `json:"firstName,omitempty"`
	LastName                *string          `json:"lastName,omitempty"`
	Phone                   *string          `json:"phone,omitempty"`
	JobTitle                *string          `json:"jobTitle,omitempty"`
	Enabled                 *bool            `json:"enabled,omitempty"`
	Avatar                  *string          `json:"avatar,omitempty"`
	StartTour               *bool            `json:"startTour,omitempty"`
	Created                 *common.NullTime `json:"created,omitempty"`
	Updated                 *common.NullTime `json:"updated,omitempty"`

	// Enriched fields (populated by service layer)
	// GroupsCount: number of groups user belongs to (for listing)
	// Groups: detailed group info (for detail view - GetUserById)
	// Memberships: organization access info (for detail view - GetUserById)
	GroupsCount *int                  `json:"groupsCount,omitempty"`
	Groups      []UserGroupInfo       `json:"groups,omitempty"`
	Memberships []UserMembershipInfo  `json:"memberships,omitempty"`
}

func (u *UserResponse) SetCreated(t *common.NullTime) { u.Created = t }
func (u *UserResponse) SetUpdated(t *common.NullTime) { u.Updated = t }

// UserQuery represents query parameters for listing users.
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
//   - Email: filter by email (partial match)
//   - FirstName: filter by first name (partial match)
//   - LastName: filter by last name (partial match)
//   - Enabled: filter by enabled status
//
// Note: Users are typically global entities, but org filtering still applies
// based on memberships and RequestContext for multi-tenant isolation.
type UserQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Email     *string `query:"email" validate:"omitempty,max=254"`
	FirstName *string `query:"firstName" validate:"omitempty,max=100"`
	LastName  *string `query:"lastName" validate:"omitempty,max=100"`
	Enabled   *bool   `query:"enabled" validate:"omitempty"`
}
