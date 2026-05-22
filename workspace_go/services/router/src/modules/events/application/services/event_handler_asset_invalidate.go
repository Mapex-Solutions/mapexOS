package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	eventTypes "router/src/modules/events/application/types"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// parseAssetInvalidateEvent parses and validates a FANOUT cache invalidation
// payload. Returns the parsed event or nil when the message must be discarded
// (invalid JSON or missing assetUUID).
func (s *EventService) parseAssetInvalidateEvent(data []byte) *eventTypes.AssetInvalidateEvent {
	var event eventTypes.AssetInvalidateEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil
	}
	if event.AssetUUID == "" {
		return nil
	}
	return &event
}

// processAssetInvalidateMessage invalidates local cache (L0+L1) when asset data changes.
// FANOUT consumer - each service instance receives this message.
// Key format: {orgId}/{assetUUID}
func (s *EventService) processAssetInvalidateMessage(idx int, msg *natsModel.Message) {
	var event eventTypes.AssetInvalidateEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] AssetInvalidate invalid JSON: %v", err))
		return
	}

	if event.OrgId == "" {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] AssetInvalidate missing required field 'orgId': %+v", event))
		return
	}

	if event.AssetUUID == "" {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] AssetInvalidate missing required field 'assetUUID': %+v", event))
		return
	}

	cacheKey := event.OrgId + "/" + event.AssetUUID

	if err := s.deps.AssetCache.Invalidate(ctx.Background(), cacheKey); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] Cache invalidate failed for %s: %v", cacheKey, err))
		s.deps.Metrics.CacheInvalidationsTotal.WithLabelValues("error").Inc()
	} else {
		s.deps.Metrics.CacheInvalidationsTotal.WithLabelValues("success").Inc()
	}
}
