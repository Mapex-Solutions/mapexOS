package services

import (
	"encoding/json"
	"fmt"

	"mapexIam/src/modules/cache_invalidation/application/events"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * handleGroupChanged invalidates cache when a group is created or updated.
 * This includes when users are added/removed from the group.
 */
func (s *CacheInvalidationService) handleGroupChanged(data []byte) error {
	var event events.GroupChangedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Processing GroupChanged for group=%s org=%s", event.GroupID, event.OrganizationID))

	invalidatedCount, err := s.invalidateAuthAndCoverageForGroupMembers(event.GroupID, event.OrganizationID)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to get members for group=%s", event.GroupID))
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Invalidated cache for %d members in group=%s", invalidatedCount, event.GroupID))
	return nil
}

// handleGroupDeleted invalidates cache when a group is deleted
func (s *CacheInvalidationService) handleGroupDeleted(data []byte) error {
	var event events.GroupDeletedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Processing GroupDeleted for group=%s org=%s", event.GroupID, event.OrganizationID))

	// Note: Members should still exist at this point (deleted after cache invalidation)
	invalidatedCount, err := s.invalidateAuthAndCoverageForGroupMembers(event.GroupID, event.OrganizationID)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to get members for group=%s", event.GroupID))
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Invalidated cache for %d members before group=%s deletion", invalidatedCount, event.GroupID))
	return nil
}
