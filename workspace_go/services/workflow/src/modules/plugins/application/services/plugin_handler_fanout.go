package services

import (
	"context"
	"encoding/json"
	"fmt"

	"workflow/src/modules/plugins/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/* PRIVATE METHODS — FANOUT INVALIDATION HANDLER */

// parseFanoutPayload decodes the FANOUT message body into a PluginInvalidatePayload.
// Returns ok=false on malformed input. Invalid payloads are logged and dropped (no retry),
// preserving the original module-level behavior.
func (s *PluginService) parseFanoutPayload(msg *natsModel.Message) (ports.PluginInvalidatePayload, bool) {
	var payload ports.PluginInvalidatePayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Plugins] Failed to unmarshal FANOUT payload: %v", err))
		return payload, false
	}
	return payload, true
}

// invalidatePluginCache invalidates both the targeted pluginId entry and the full
// "all enabled plugins" list in the local TieredCache (L0/L1). Failures are logged
// but never propagate — best-effort, matching the original handler semantics.
func (s *PluginService) invalidatePluginCache(ctx context.Context, payload ports.PluginInvalidatePayload) {
	if err := s.deps.PluginLoader.Invalidate(ctx, payload.PluginID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Plugins] FANOUT invalidate failed for %s: %v", payload.PluginID, err))
	}
	if err := s.deps.PluginLoader.InvalidateAll(ctx); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Plugins] FANOUT invalidateAll failed: %v", err))
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Plugins] FANOUT received: invalidated %s (action: %s)", payload.PluginID, payload.Action))
}
