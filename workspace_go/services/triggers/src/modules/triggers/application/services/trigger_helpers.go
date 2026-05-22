package services

import (
	"context"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// mapperResponseOpts defines options for Entity → Response DTO conversion.
// Shared across all trigger handler files (ObjectId → string coercion).
var mapperResponseOpts = mapper.MapperOptions{ObjectIdToString: true}

// invalidateTriggerCache invalidates the trigger cache entry for the given triggerId.
// Called after updates and deletes to avoid serving stale data.
func (s *TriggerService) invalidateTriggerCache(ctx context.Context, triggerId string) {
	cacheKey := s.deps.CacheKeyBuilder.TriggerKey(triggerId)
	if err := s.deps.CacheRepository.Del(ctx, cacheKey); err != nil {
		logger.Warn("[SERVICE:Trigger] Failed to invalidate trigger cache: " + err.Error())
	}
}
