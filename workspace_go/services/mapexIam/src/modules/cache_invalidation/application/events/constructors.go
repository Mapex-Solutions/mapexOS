package events

import (
	"mapexIam/src/modules/cache_invalidation/domain/constants"
	"time"
)

// NewRolePermissionsChangedEvent creates a new RolePermissionsChangedEvent
func NewRolePermissionsChangedEvent(roleID string, oldPermissions, newPermissions []string, actorID string) *RolePermissionsChangedEvent {
	return &RolePermissionsChangedEvent{
		BaseCacheInvalidationEvent: BaseCacheInvalidationEvent{
			EventType: constants.EventTypeRolePermissionsChanged,
			Timestamp: time.Now(),
			ActorID:   actorID,
		},
		RoleID:         roleID,
		OldPermissions: oldPermissions,
		NewPermissions: newPermissions,
	}
}

// NewRoleDeletedEvent creates a new RoleDeletedEvent
func NewRoleDeletedEvent(roleID string, actorID string) *RoleDeletedEvent {
	return &RoleDeletedEvent{
		BaseCacheInvalidationEvent: BaseCacheInvalidationEvent{
			EventType: constants.EventTypeRoleDeleted,
			Timestamp: time.Now(),
			ActorID:   actorID,
		},
		RoleID: roleID,
	}
}

// NewOrgAccessPolicyChangedEvent creates a new OrgAccessPolicyChangedEvent
func NewOrgAccessPolicyChangedEvent(orgID, oldPolicy, newPolicy, orgKey, actorID string) *OrgAccessPolicyChangedEvent {
	return &OrgAccessPolicyChangedEvent{
		BaseCacheInvalidationEvent: BaseCacheInvalidationEvent{
			EventType: constants.EventTypeOrgAccessPolicyChanged,
			Timestamp: time.Now(),
			ActorID:   actorID,
		},
		OrganizationID:  orgID,
		OldRolePolicy:   oldPolicy,
		NewRolePolicy:   newPolicy,
		OrganizationKey: orgKey,
	}
}

// NewMembershipChangedEvent creates a new MembershipChangedEvent
func NewMembershipChangedEvent(membershipID, userID, orgID, oldRoleID, newRoleID, actorID string) *MembershipChangedEvent {
	return &MembershipChangedEvent{
		BaseCacheInvalidationEvent: BaseCacheInvalidationEvent{
			EventType: constants.EventTypeMembershipChanged,
			Timestamp: time.Now(),
			ActorID:   actorID,
		},
		MembershipID:   membershipID,
		UserID:         userID,
		OrganizationID: orgID,
		OldRoleID:      oldRoleID,
		NewRoleID:      newRoleID,
	}
}

// NewMembershipDeletedEvent creates a new MembershipDeletedEvent
func NewMembershipDeletedEvent(membershipID, userID, orgID, actorID string) *MembershipDeletedEvent {
	return &MembershipDeletedEvent{
		BaseCacheInvalidationEvent: BaseCacheInvalidationEvent{
			EventType: constants.EventTypeMembershipDeleted,
			Timestamp: time.Now(),
			ActorID:   actorID,
		},
		MembershipID:   membershipID,
		UserID:         userID,
		OrganizationID: orgID,
	}
}

// NewGroupChangedEvent creates a new GroupChangedEvent
func NewGroupChangedEvent(groupID, orgID, actorID string) *GroupChangedEvent {
	return &GroupChangedEvent{
		BaseCacheInvalidationEvent: BaseCacheInvalidationEvent{
			EventType: constants.EventTypeGroupChanged,
			Timestamp: time.Now(),
			ActorID:   actorID,
		},
		GroupID:        groupID,
		OrganizationID: orgID,
	}
}

// NewGroupDeletedEvent creates a new GroupDeletedEvent
func NewGroupDeletedEvent(groupID, orgID, actorID string) *GroupDeletedEvent {
	return &GroupDeletedEvent{
		BaseCacheInvalidationEvent: BaseCacheInvalidationEvent{
			EventType: constants.EventTypeGroupDeleted,
			Timestamp: time.Now(),
			ActorID:   actorID,
		},
		GroupID:        groupID,
		OrganizationID: orgID,
	}
}

// NewOrgHierarchyChangedEvent creates a new OrgHierarchyChangedEvent
func NewOrgHierarchyChangedEvent(orgID string, ancestorOrgIds []string, action string, actorID string) *OrgHierarchyChangedEvent {
	return &OrgHierarchyChangedEvent{
		BaseCacheInvalidationEvent: BaseCacheInvalidationEvent{
			EventType: constants.EventTypeOrgHierarchyChanged,
			Timestamp: time.Now(),
			ActorID:   actorID,
		},
		OrganizationID: orgID,
		AncestorOrgIds: ancestorOrgIds,
		Action:         action,
	}
}
