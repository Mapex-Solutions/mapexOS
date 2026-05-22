package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	"mapexIam/src/modules/cache_invalidation/application/events"
	"mapexIam/src/modules/cache_invalidation/domain/constants"

	membershipPorts "mapexIam/src/modules/memberships/application/ports"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// handleCacheInvalidationEvent routes events to appropriate handlers
func (s *CacheInvalidationService) handleCacheInvalidationEvent(msg []byte) error {
	var baseEvent events.BaseCacheInvalidationEvent
	if err := json.Unmarshal(msg, &baseEvent); err != nil {
		// ACK malformed messages to prevent infinite redelivery (e.g. old double-encoded messages)
		logger.Warn(fmt.Sprintf("[SERVICE:CacheInvalidation] Skipping unparseable message (ACK to prevent redelivery): %v", err))
		return nil
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Received event: %s", baseEvent.EventType))

	switch baseEvent.EventType {
	case constants.EventTypeRolePermissionsChanged:
		return s.handleRolePermissionsChanged(msg)
	case constants.EventTypeRoleDeleted:
		return s.handleRoleDeleted(msg)
	case constants.EventTypeOrgAccessPolicyChanged:
		return s.handleOrgAccessPolicyChanged(msg)
	case constants.EventTypeMembershipChanged:
		return s.handleMembershipChanged(msg)
	case constants.EventTypeMembershipDeleted:
		return s.handleMembershipDeleted(msg)
	case constants.EventTypeGroupChanged:
		return s.handleGroupChanged(msg)
	case constants.EventTypeGroupDeleted:
		return s.handleGroupDeleted(msg)
	case constants.EventTypeOrgHierarchyChanged:
		return s.handleOrgHierarchyChanged(msg)
	default:
		logger.Warn(fmt.Sprintf("[SERVICE:CacheInvalidation] Unknown event type: %s", baseEvent.EventType))
		return nil
	}
}

/**
 * resolveMembershipsToUserOrgPairs resolves a list of memberships to unique user-org pairs.
 * For user memberships, adds the assignee directly.
 * For group memberships, resolves all group members via groupMemberRepo.
 */
func (s *CacheInvalidationService) resolveMembershipsToUserOrgPairs(memberships []*membershipPorts.Membership) map[UserOrgPair]bool {
	pairs := make(map[UserOrgPair]bool)

	for _, membership := range memberships {
		orgID := membership.OrgID.Hex()

		if membership.AssigneeType == "user" {
			pairs[UserOrgPair{UserID: membership.AssigneeID.Hex(), OrgID: orgID}] = true

		} else if membership.AssigneeType == "group" {
			groupMembers, err := s.deps.GroupMemberRepo.FindByGroupId(ctx.Background(), membership.AssigneeID.Hex())
			if err != nil {
				logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to resolve members for group=%s", membership.AssigneeID.Hex()))
				continue
			}
			for _, member := range groupMembers {
				pairs[UserOrgPair{UserID: member.UserID.Hex(), OrgID: orgID}] = true
			}
		}
	}

	return pairs
}

// invalidateAuthForPairs invalidates the authorization cache for each user-org pair.
func (s *CacheInvalidationService) invalidateAuthForPairs(pairs map[UserOrgPair]bool) {
	for pair := range pairs {
		if err := s.deps.AuthCacheRepo.InvalidateUserAuth(ctx.Background(), pair.UserID, pair.OrgID); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate auth cache for user=%s org=%s", pair.UserID, pair.OrgID))
		}
	}
}

/**
 * collectUserIdsFromMemberships resolves a list of memberships to unique user IDs.
 * For user memberships, adds the assignee directly.
 * For group memberships, resolves members via groupMemberRepo.
 */
func (s *CacheInvalidationService) collectUserIdsFromMemberships(memberships []*membershipPorts.Membership, dest map[string]bool) {
	for _, membership := range memberships {
		if membership.AssigneeType == "user" {
			dest[membership.AssigneeID.Hex()] = true
		} else if membership.AssigneeType == "group" {
			groupMembers, err := s.deps.GroupMemberRepo.FindByGroupId(ctx.Background(), membership.AssigneeID.Hex())
			if err != nil {
				logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to get members for group=%s", membership.AssigneeID.Hex()))
				continue
			}
			for _, member := range groupMembers {
				dest[member.UserID.Hex()] = true
			}
		}
	}
}

/**
 * invalidateAuthAndCoverageForGroupMembers fetches all members of a group
 * and invalidates both auth and coverage cache for each user.
 */
func (s *CacheInvalidationService) invalidateAuthAndCoverageForGroupMembers(groupID, orgID string) (int, error) {
	members, err := s.deps.GroupMemberRepo.FindByGroupId(ctx.Background(), groupID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, member := range members {
		userID := member.UserID.Hex()

		if err := s.deps.AuthCacheRepo.InvalidateUserAuth(ctx.Background(), userID, orgID); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate auth cache for user=%s org=%s", userID, orgID))
			continue
		}

		if err := s.deps.CoverageCacheRepo.InvalidateCache(ctx.Background(), userID); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate coverage cache for user=%s", userID))
		}

		count++
	}
	return count, nil
}
