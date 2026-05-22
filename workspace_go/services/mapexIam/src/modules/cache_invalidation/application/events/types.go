package events

import (
	"mapexIam/src/modules/cache_invalidation/domain/constants"
	"time"
)

// BaseCacheInvalidationEvent contains common fields for all cache invalidation events
type BaseCacheInvalidationEvent struct {
	EventType constants.CacheInvalidationEventType `json:"eventType"`
	Timestamp time.Time                            `json:"timestamp"`
	ActorID   string                               `json:"actorId,omitempty"` // User who triggered the change
}

// RolePermissionsChangedEvent is published when a role's permissions are updated
type RolePermissionsChangedEvent struct {
	BaseCacheInvalidationEvent
	RoleID         string   `json:"roleId"`
	OldPermissions []string `json:"oldPermissions"`
	NewPermissions []string `json:"newPermissions"`
}

// RoleDeletedEvent is published when a role is deleted
type RoleDeletedEvent struct {
	BaseCacheInvalidationEvent
	RoleID string `json:"roleId"`
}

// OrgAccessPolicyChangedEvent is published when an organization's access policy changes
type OrgAccessPolicyChangedEvent struct {
	BaseCacheInvalidationEvent
	OrganizationID  string `json:"organizationId"`
	OldRolePolicy   string `json:"oldRolePolicy"`
	NewRolePolicy   string `json:"newRolePolicy"`
	OrganizationKey string `json:"organizationKey,omitempty"` // PathKey for hierarchical queries
}

// OrgHierarchyChangedEvent is published when an organization is created or deleted,
// affecting the coverage cache of users with recursive memberships on ancestor orgs.
type OrgHierarchyChangedEvent struct {
	BaseCacheInvalidationEvent
	OrganizationID string   `json:"organizationId"` // The created/deleted org
	AncestorOrgIds []string `json:"ancestorOrgIds"` // All ancestor org IDs (for recursive membership lookup)
	Action         string   `json:"action"`          // "created" | "deleted"
}

// MembershipChangedEvent is published when a membership is created or updated
type MembershipChangedEvent struct {
	BaseCacheInvalidationEvent
	MembershipID   string `json:"membershipId"`
	UserID         string `json:"userId"`
	OrganizationID string `json:"organizationId"`
	OldRoleID      string `json:"oldRoleId,omitempty"`
	NewRoleID      string `json:"newRoleId"`
}

// MembershipDeletedEvent is published when a membership is deleted
type MembershipDeletedEvent struct {
	BaseCacheInvalidationEvent
	MembershipID   string `json:"membershipId"`
	UserID         string `json:"userId"`
	OrganizationID string `json:"organizationId"`
}

// GroupChangedEvent is published when a group is created or updated
type GroupChangedEvent struct {
	BaseCacheInvalidationEvent
	GroupID        string `json:"groupId"`
	OrganizationID string `json:"organizationId"`
}

// GroupDeletedEvent is published when a group is deleted
type GroupDeletedEvent struct {
	BaseCacheInvalidationEvent
	GroupID        string `json:"groupId"`
	OrganizationID string `json:"organizationId"`
}
