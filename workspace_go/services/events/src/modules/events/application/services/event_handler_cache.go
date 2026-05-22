package services

import (
	"context"
	"encoding/json"
	"fmt"

	templateContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assettemplates"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// parseTemplateInvalidatePayload decodes one FANOUT payload from the
// mapexos.fanout.template.invalidate stream. Returns ok=false on invalid
// JSON or missing templateId — the kit's Subscribe callback Acks regardless,
// so we just drop and log the bad message.
func (s *EventService) parseTemplateInvalidatePayload(msg *natsModel.Message) (templateContract.TemplateInvalidatePayload, bool) {
	var event templateContract.TemplateInvalidatePayload
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] TemplateInvalidate invalid JSON: %v", err))
		return event, false
	}
	if event.TemplateId == "" {
		logger.Warn("[SERVICE:Event] TemplateInvalidate missing templateId, discarding")
		return event, false
	}
	return event, true
}

// applyTemplateInvalidate clears the local TieredCache (L0+L1) entry for the
// {orgId}/{templateId} pair so the next read falls back to L2 (MinIO) for
// the source-of-truth template.
func (s *EventService) applyTemplateInvalidate(event templateContract.TemplateInvalidatePayload) {
	cacheKey := event.OrgId + "/" + event.TemplateId
	if err := s.deps.TemplateCache.Invalidate(context.Background(), cacheKey); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] Template cache invalidate failed for %s: %v", cacheKey, err))
		return
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Event] Template cache invalidated key=%s", cacheKey))
}
