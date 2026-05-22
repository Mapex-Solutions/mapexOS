package dtos

import (
	"errors"

	membershipDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/memberships"
	userDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/users"
)

// MembershipData represents direct role assignment during user onboarding.
// OrgID comes from RequestContext.OrgContext (current org).
// Scope can be provided in request (optional) - falls back to Organization.AccessPolicy.DefaultScope.
type MembershipData struct {
	Roles []string `json:"roles" validate:"required,min=1,dive,min=1"`
	Scope *string  `json:"scope,omitempty" validate:"omitempty,oneof=local recursive"`
}

// ExistingGroupData represents using an existing group during onboarding.
// User will be added to the group's members[] array and inherit its roles.
// The group ALREADY has a Membership (assigneeType='group') with roles defined.
//
// No roles needed - they come from the group's existing Membership.
type ExistingGroupData struct {
	GroupID string `json:"groupId" validate:"required,mongoid"`
}

// NewGroupData represents creating a new group during onboarding.
// Onboarding service will delegate to Groups service to create:
//   - The Group entity
//   - The Membership for the group (assigneeType='group', with the specified roles)
//
// Then the user is added to the new group's members[] array.
//
// Notes:
//   - IsSystem is ALWAYS false (enforced by Groups service)
//   - IsTemplate defaults to false
//   - OrgID comes from RequestContext
//   - Scope comes from Organization.AccessPolicy.DefaultScope
type NewGroupData struct {
	Name        string   `json:"name" validate:"required,min=3,max=150"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	RoleIds     []string `json:"roleIds" validate:"required,min=1,dive,mongoid"`
}

// GroupAccessData represents group-based access during user onboarding.
// EXACTLY ONE of ExistingGroup or NewGroup must be provided.
//
// Use cases:
//   - ExistingGroup: User joins an existing group (inherits its roles)
//   - NewGroup: Create a new group with specific roles, then user joins it
type GroupAccessData struct {
	ExistingGroup *ExistingGroupData `json:"existingGroup,omitempty" validate:"omitempty"`
	NewGroup      *NewGroupData      `json:"newGroup,omitempty" validate:"omitempty"`
}

// Transform validates that exactly one of ExistingGroup or NewGroup is provided.
func (g *GroupAccessData) Transform() error {
	if g.ExistingGroup != nil && g.NewGroup != nil {
		return errors.New("only one of existingGroup or newGroup can be provided")
	}
	if g.ExistingGroup == nil && g.NewGroup == nil {
		return errors.New("either existingGroup or newGroup must be provided")
	}
	return nil
}

// CreateUserWithMemberships combines user creation with multiple memberships and group assignments
type CreateUserWithMemberships struct {
	// User data
	Email                   string  `json:"email" validate:"required,email,max=254"`
	Password                *string `json:"password,omitempty" validate:"omitempty,min=8,max=72"`
	ChangePasswordNextLogin bool    `json:"changePasswordNextLogin" validate:"-"`
	FirstName               string  `json:"firstName" validate:"required,min=2,max=100"`
	LastName                string  `json:"lastName" validate:"required,min=2,max=100"`
	Phone                   *string `json:"phone,omitempty" validate:"omitempty,e164"`
	JobTitle                *string `json:"jobTitle,omitempty" validate:"omitempty,max=120"`
	Enabled                 bool    `json:"enabled" validate:"-"`
	Avatar                  *string `json:"avatar,omitempty" validate:"omitempty,url"`
	// AuthProvider removed - V1 always uses internal auth.
	// Next version: auth provider will be determined by the customer's Organization.AuthConfig

	// Direct role memberships (optional - for simple cases or small clients)
	Memberships []MembershipData `json:"memberships,omitempty" validate:"omitempty,dive"`

	// Group access (recommended for large clients - cleaner RBAC)
	// Each entry represents either joining an existing group OR creating a new group.
	// EXACTLY ONE of existingGroup or newGroup must be provided per entry.
	Groups []GroupAccessData `json:"groups,omitempty" validate:"omitempty,dive"`
}

// Transform validates business rules after basic validation.
// At least one membership type (direct roles or groups) must be provided.
// Each group entry must have exactly one of existingGroup or newGroup.
func (dto *CreateUserWithMemberships) Transform() error {
	if len(dto.Memberships) == 0 && len(dto.Groups) == 0 {
		return errors.New("at least one membership (direct roles) or group must be provided")
	}

	// Validate each group entry
	for i := range dto.Groups {
		if err := dto.Groups[i].Transform(); err != nil {
			return err
		}
	}

	return nil
}

// UpdateUserWithAccessParams contains the path parameters for user update endpoint
type UpdateUserWithAccessParams struct {
	UserID string `params:"userId" validate:"required,mongoid"`
}

// UpdateUserWithAccess contains user data updates and access configuration changes.
// OrgID comes from RequestContext.OrgContext (current org).
// Scope can be provided per membership (optional) - falls back to Organization.AccessPolicy.DefaultScope.
//
// Backend logic (DIFF-BASED for groups):
//  1. Updates user data (only provided fields)
//  2. For groups: Calculate diff (add new, remove missing from current org)
//  3. For memberships: Replace direct memberships in current org with new ones
//  4. All within a MongoDB transaction
//
// Note: Empty groups/memberships arrays mean "remove all" for that type.
// To keep existing access unchanged, omit the field entirely (nil vs empty array).
type UpdateUserWithAccess struct {
	// User data (all optional for partial updates)
	FirstName               *string `json:"firstName,omitempty" validate:"omitempty,min=2,max=100"`
	LastName                *string `json:"lastName,omitempty" validate:"omitempty,min=2,max=100"`
	Phone                   *string `json:"phone,omitempty" validate:"omitempty,e164"`
	JobTitle                *string `json:"jobTitle,omitempty" validate:"omitempty,max=120"`
	Enabled                 *bool   `json:"enabled,omitempty" validate:"-"`
	Avatar                  *string `json:"avatar,omitempty" validate:"omitempty,url"`
	Password                *string `json:"password,omitempty" validate:"omitempty,min=8,max=72"`
	ChangePasswordNextLogin *bool   `json:"changePasswordNextLogin,omitempty" validate:"-"`

	// Direct role memberships - DECLARATIVE: array represents desired state
	// nil = don't change existing memberships
	// [] = remove all direct memberships
	// [items] = replace with these memberships (scope optional per item)
	Memberships *[]MembershipData `json:"memberships,omitempty" validate:"omitempty,dive"`

	// Group access - DECLARATIVE: array represents desired groups
	// nil = don't change group memberships
	// [] = remove user from all groups in current org
	// [items] = ensure user is member of exactly these groups (diff-based add/remove)
	Groups *[]GroupAccessData `json:"groups,omitempty" validate:"omitempty,dive"`
}

// Transform validates business rules after basic validation.
// For UPDATE: at least one field must be provided (user data OR access config).
// Each group entry must have exactly one of existingGroup or newGroup.
func (dto *UpdateUserWithAccess) Transform() error {
	// Check if any user data is provided
	hasUserData := dto.FirstName != nil || dto.LastName != nil || dto.Phone != nil ||
		dto.JobTitle != nil || dto.Enabled != nil || dto.Avatar != nil ||
		dto.Password != nil || dto.ChangePasswordNextLogin != nil

	// Check if any access config is provided
	hasAccessConfig := dto.Memberships != nil || dto.Groups != nil

	if !hasUserData && !hasAccessConfig {
		return errors.New("at least one field (user data or access config) must be provided")
	}

	// Validate each group entry if groups provided
	if dto.Groups != nil {
		for i := range *dto.Groups {
			if err := (*dto.Groups)[i].Transform(); err != nil {
				return err
			}
		}
	}

	return nil
}

// UserOnboardingResponse contains the created user and their memberships
type UserOnboardingResponse struct {
	User        *userDtos.UserResponse                `json:"user"`
	Memberships []*membershipDtos.MembershipResponse `json:"memberships"`
}
