package constants

// CacheInvalidationEventType defines the types of cache invalidation events
type CacheInvalidationEventType string

const (
	// Role events
	EventTypeRolePermissionsChanged CacheInvalidationEventType = "role.permissions.changed"
	EventTypeRoleDeleted            CacheInvalidationEventType = "role.deleted"

	// Organization events
	EventTypeOrgAccessPolicyChanged CacheInvalidationEventType = "organization.access_policy.changed"
	EventTypeOrgHierarchyChanged    CacheInvalidationEventType = "organization.hierarchy.changed"

	// Membership events
	EventTypeMembershipChanged CacheInvalidationEventType = "membership.changed"
	EventTypeMembershipDeleted CacheInvalidationEventType = "membership.deleted"

	// Group events
	EventTypeGroupChanged CacheInvalidationEventType = "group.changed"
	EventTypeGroupDeleted CacheInvalidationEventType = "group.deleted"
)
