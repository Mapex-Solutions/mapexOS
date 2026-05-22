package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	"mapexIam/src/modules/cache_invalidation/application/events"

	membershipDtos "mapexIam/src/modules/memberships/application/dtos"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// handleOrgAccessPolicyChanged invalidates cache when organization access policy changes
func (s *CacheInvalidationService) handleOrgAccessPolicyChanged(data []byte) error {
	var event events.OrgAccessPolicyChangedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Processing OrgAccessPolicyChanged for org=%s", event.OrganizationID))

	orgID := event.OrganizationID
	allMemberships, err := s.deps.MembershipService.GetAllMemberships(ctx.Background(), &membershipDtos.MembershipQueryDto{
		OrgID: &orgID,
	})
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to get memberships for org=%s", event.OrganizationID))
		return err
	}

	pairs := s.resolveMembershipsToUserOrgPairs(allMemberships)
	s.invalidateAuthForPairs(pairs)

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Invalidated cache for %d user-org pairs in org=%s", len(pairs), event.OrganizationID))
	return nil
}

/**
 * handleOrgHierarchyChanged invalidates coverage cache when org hierarchy changes
 * (new org created or org deleted). Finds users with recursive memberships on
 * ancestor orgs and invalidates their coverage cache so it gets rebuilt with the new hierarchy.
 */
func (s *CacheInvalidationService) handleOrgHierarchyChanged(data []byte) error {
	var event events.OrgHierarchyChangedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Processing OrgHierarchyChanged action=%s org=%s ancestors=%d",
		event.Action, event.OrganizationID, len(event.AncestorOrgIds)))

	// Collect all unique user IDs whose coverage cache needs invalidation
	uniqueUserIds := make(map[string]bool)
	recursiveScope := "recursive"

	// For each ancestor org, find memberships with scope="recursive"
	// These users' coverage caches include the ancestor's descendants via hierarchy expansion,
	// so adding/removing a child org makes their cache stale
	for _, ancestorOrgId := range event.AncestorOrgIds {
		orgId := ancestorOrgId
		memberships, err := s.deps.MembershipService.GetAllMemberships(ctx.Background(), &membershipDtos.MembershipQueryDto{
			OrgID: &orgId,
			Scope: &recursiveScope,
		})
		if err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to get recursive memberships for ancestor org=%s", ancestorOrgId))
			continue
		}
		s.collectUserIdsFromMemberships(memberships, uniqueUserIds)
	}

	// For delete: also invalidate users with direct memberships on the deleted org itself
	if event.Action == "deleted" {
		orgId := event.OrganizationID
		directMemberships, err := s.deps.MembershipService.GetAllMemberships(ctx.Background(), &membershipDtos.MembershipQueryDto{
			OrgID: &orgId,
		})
		if err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to get direct memberships for deleted org=%s", event.OrganizationID))
		} else {
			s.collectUserIdsFromMemberships(directMemberships, uniqueUserIds)
		}
	}

	// Invalidate coverage cache for all affected users
	invalidatedCount := 0
	for userID := range uniqueUserIds {
		if err := s.deps.CoverageCacheRepo.InvalidateCache(ctx.Background(), userID); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate coverage cache for user=%s", userID))
			continue
		}
		invalidatedCount++
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Invalidated coverage cache for %d users due to org hierarchy change (org=%s action=%s)",
		invalidatedCount, event.OrganizationID, event.Action))
	return nil
}
