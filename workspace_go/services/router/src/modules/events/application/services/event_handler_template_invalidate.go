package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	eventTypes "router/src/modules/events/application/types"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// processTemplateInvalidateMessage parses a single FANOUT invalidation payload and
// removes the matching entry from the TemplateCache (L0+L1). Invalid payloads are
// discarded (Reject) — retries would not recover a malformed message.
// Key format: {orgId}/{templateId}
func (s *EventService) processTemplateInvalidateMessage(idx int, msg *natsModel.Message) {
	var event eventTypes.TemplateInvalidateEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] Template invalidate invalid JSON: %v", err))
		return
	}

	if event.TemplateId == "" {
		logger.Warn("[SERVICE:Event] Template invalidate missing templateId, discarding")
		return
	}

	cacheKey := event.OrgId + "/" + event.TemplateId

	if err := s.deps.TemplateCache.Invalidate(ctx.Background(), cacheKey); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] Template cache invalidate failed for %s: %v", cacheKey, err))
	} else {
		logger.Debug(fmt.Sprintf("[SERVICE:Event] Template cache invalidated key=%s", cacheKey))
	}
}
