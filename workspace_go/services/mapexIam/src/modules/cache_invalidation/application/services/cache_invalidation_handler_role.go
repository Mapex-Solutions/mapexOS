package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	"mapexIam/src/modules/cache_invalidation/application/events"

	membershipDtos "mapexIam/src/modules/memberships/application/dtos"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// handleRolePermissionsChanged invalidates cache when role permissions change
func (s *CacheInvalidationService) handleRolePermissionsChanged(data []byte) error {
	var event events.RolePermissionsChangedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Processing RolePermissionsChanged for role=%s", event.RoleID))

	if err := s.deps.AuthCacheRepo.InvalidateRole(ctx.Background(), event.RoleID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate role cache for role=%s", event.RoleID))
	}

	roleID := event.RoleID
	allMemberships, err := s.deps.MembershipService.GetAllMemberships(ctx.Background(), &membershipDtos.MembershipQueryDto{
		RoleID: &roleID,
	})
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to get memberships for role=%s", event.RoleID))
		return err
	}

	pairs := s.resolveMembershipsToUserOrgPairs(allMemberships)
	s.invalidateAuthForPairs(pairs)

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Invalidated cache for %d user-org pairs affected by role=%s", len(pairs), event.RoleID))
	return nil
}

// handleRoleDeleted invalidates cache when a role is deleted
func (s *CacheInvalidationService) handleRoleDeleted(data []byte) error {
	var event events.RoleDeletedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Processing RoleDeleted for role=%s", event.RoleID))

	if err := s.deps.AuthCacheRepo.InvalidateRole(ctx.Background(), event.RoleID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate role cache for role=%s", event.RoleID))
	}

	// Memberships still exist at this point (role is deleted but memberships aren't cleaned up yet)
	roleID := event.RoleID
	allMemberships, err := s.deps.MembershipService.GetAllMemberships(ctx.Background(), &membershipDtos.MembershipQueryDto{
		RoleID: &roleID,
	})
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to get memberships for deleted role=%s", event.RoleID))
		return err
	}

	pairs := s.resolveMembershipsToUserOrgPairs(allMemberships)
	s.invalidateAuthForPairs(pairs)

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Invalidated role cache + %d user-org pairs for deleted role=%s", len(pairs), event.RoleID))
	return nil
}
