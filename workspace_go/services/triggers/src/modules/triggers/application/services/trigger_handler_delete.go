package services

import (
	ctx "context"

	"triggers/src/modules/triggers/domain/entities"
)

// invalidateTriggerCounterCache drops the per-org counter cache after a
// successful delete using the deleted entity's OrgId (the request context
// may have moved on by the time we reach the post-delete step).
func (s *TriggerService) invalidateTriggerCounterCache(c ctx.Context, trigger *entities.Trigger) {
	if trigger == nil || trigger.OrgID == nil {
		return
	}
	counterKey := s.deps.CacheKeyBuilder.CounterKey(trigger.OrgID.Hex())
	_ = s.deps.AppCache.Del(c, counterKey)
}
