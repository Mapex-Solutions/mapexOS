package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	"mapexIam/src/modules/cache_invalidation/application/events"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// handleMembershipChanged invalidates cache when a membership is created or updated
func (s *CacheInvalidationService) handleMembershipChanged(data []byte) error {
	var event events.MembershipChangedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Processing MembershipChanged for user=%s org=%s", event.UserID, event.OrganizationID))

	if err := s.deps.AuthCacheRepo.InvalidateUserAuth(ctx.Background(), event.UserID, event.OrganizationID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate auth cache for user=%s org=%s", event.UserID, event.OrganizationID))
		return err
	}

	if err := s.deps.CoverageCacheRepo.InvalidateCache(ctx.Background(), event.UserID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate coverage cache for user=%s", event.UserID))
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Invalidated auth + coverage cache for user=%s org=%s", event.UserID, event.OrganizationID))
	return nil
}

// handleMembershipDeleted invalidates cache when a membership is deleted
func (s *CacheInvalidationService) handleMembershipDeleted(data []byte) error {
	var event events.MembershipDeletedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Processing MembershipDeleted for user=%s org=%s", event.UserID, event.OrganizationID))

	if err := s.deps.AuthCacheRepo.InvalidateUserAuth(ctx.Background(), event.UserID, event.OrganizationID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate auth cache for user=%s org=%s", event.UserID, event.OrganizationID))
		return err
	}

	if err := s.deps.CoverageCacheRepo.InvalidateCache(ctx.Background(), event.UserID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:CacheInvalidation] Failed to invalidate coverage cache for user=%s", event.UserID))
	}

	logger.Info(fmt.Sprintf("[SERVICE:CacheInvalidation] Invalidated auth + coverage cache for user=%s org=%s", event.UserID, event.OrganizationID))
	return nil
}
