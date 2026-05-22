package nats

import (
	"context"
	"fmt"
	"time"

	"assets/src/modules/healthmonitor/application/ports"
	"assets/src/modules/healthmonitor/domain/entities"

	hmContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/google/uuid"
)

// Compile-time check
var _ ports.AlertPublisherPort = (*alertPublisher)(nil)

// NewAlertPublisher creates a publisher that dual-publishes asset connectivity
// transitions:
//
//  1. ALWAYS to mapexos.events.asset_status_save (EVENTS-ASSET-STATUS stream)
//     for persistence to the asset_status_history ClickHouse table.
//  2. CONDITIONALLY to mapexos.route.execute (ROUTE-GROUPS stream) only when
//     the alert carries non-empty RouteGroupIds.
//
// Each publish is independent — failure of one is logged but does not
// short-circuit the other. No synchronous retries here; redelivery is
// handled by the respective consumer retry policies.
func NewAlertPublisher(publisher natsModel.Publisher) ports.AlertPublisherPort {
	return &alertPublisher{publisher: publisher}
}

// PublishOffline publishes a sensor.offline transition (dual publish).
func (p *alertPublisher) PublishOffline(ctx context.Context, event entities.AlertEvent) error {
	return p.publishAlert(event)
}

// PublishOnline publishes a sensor.online transition (dual publish).
func (p *alertPublisher) PublishOnline(ctx context.Context, event entities.AlertEvent) error {
	return p.publishAlert(event)
}

// publishAlert runs the dual-publish sequence for one state transition.
//
// Shared identity across both payloads:
//   - eventId:  UUID generated once per transition. Placed inside the
//     StandardizedPayload nested event (route) AND at the top level of the
//     persistence payload (as a ClickHouse column). Lets downstream joins
//     correlate a persisted row with the routed event.
//   - created:  ISO-8601 RFC3339Nano timestamp generated once; used as the
//     StandardizedPayload.created field AND as the persistence-row timestamp.
//
// Persistence ALWAYS fires first. Route fires next only if the event carries
// non-empty RouteGroupIds. Errors are logged; we return nil so the scanner
// records "alert sent" even if one path failed — the other likely succeeded
// and NATS redelivery handles the rest.
func (p *alertPublisher) publishAlert(event entities.AlertEvent) error {
	eventId := uuid.New().String()
	created := time.Now().UTC().Format(time.RFC3339Nano)

	// (1) Persistence — ALWAYS.
	persistencePayload := buildPersistencePayload(event, eventId, created)
	if err := p.publisher.Publish(natsModel.PublishConfig{
		Subject: hmContract.AssetStatusSaveSubject,
		Data:    persistencePayload,
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:HealthAlert] Persistence publish failed: type=%s assetUUID=%s", event.Type, event.AssetUUID))
	} else {
		logger.Debug(fmt.Sprintf("[INFRA:HealthAlert] Persistence published: type=%s assetUUID=%s eventId=%s", event.Type, event.AssetUUID, eventId))
	}

	// (2) Router — CONDITIONAL on configured route groups.
	if len(event.RouteGroupIds) > 0 {
		routePayload := buildRoutePayload(event, eventId, created)
		if err := p.publisher.Publish(natsModel.PublishConfig{
			Subject: hmContract.RouterSubject,
			Data:    routePayload,
		}); err != nil {
			logger.Error(err, fmt.Sprintf("[INFRA:HealthAlert] Route publish failed: type=%s assetUUID=%s routes=%d", event.Type, event.AssetUUID, len(event.RouteGroupIds)))
		} else {
			logger.Debug(fmt.Sprintf("[INFRA:HealthAlert] Route published: type=%s assetUUID=%s eventId=%s routes=%d", event.Type, event.AssetUUID, eventId, len(event.RouteGroupIds)))
		}
	}

	return nil
}

// buildRoutePayload assembles the mapexos.route.execute payload. The top
// level carries routing metadata (orgId, assetUUID, pathKey, eventSource,
// fresh eventTrackerId) and the nested event conforms to StandardizedPayload.
func buildRoutePayload(event entities.AlertEvent, eventId string, created string) map[string]interface{} {
	return map[string]interface{}{
		"orgId":          event.OrgId,
		"assetUUID":      event.AssetUUID,
		"pathKey":        event.PathKey,
		"eventSource":    hmContract.EventSourceHealthStatus,
		"eventTrackerId": uuid.New().String(),
		"event":          buildStandardizedEvent(event, eventId, created),
	}
}

// buildStandardizedEvent builds the nested event object per the
// StandardizedPayload contract: { eventType, eventId, data, metadata, created }.
// `data` is guaranteed non-empty (ObjectAndNotBeEmpty in Zod) — assetUUID and
// assetName are always present; lastSeenAt/thresholdMinutes/missCount are
// included only when set.
func buildStandardizedEvent(event entities.AlertEvent, eventId string, created string) map[string]interface{} {
	data := map[string]interface{}{
		"assetUUID": event.AssetUUID,
		"assetName": event.AssetName,
	}
	if event.LastSeenAt != nil {
		data["lastSeenAt"] = event.LastSeenAt.UTC().Format(time.RFC3339Nano)
	}
	if event.ThresholdMinutes > 0 {
		data["thresholdMinutes"] = event.ThresholdMinutes
	}
	if event.MissCount > 0 {
		data["missCount"] = event.MissCount
	}

	return map[string]interface{}{
		"eventType": event.Type,
		"eventId":   eventId,
		"data":      data,
		"metadata":  map[string]interface{}{"source": hmContract.EventSource},
		"created":   created,
	}
}

// buildPersistencePayload produces the FLAT row shape for
// mapexos.events.asset_status_save. Fields map 1:1 to asset_status_history
// columns so the events-MS consumer can bulk-insert without unwrapping.
func buildPersistencePayload(event entities.AlertEvent, eventId string, created string) map[string]interface{} {
	payload := map[string]interface{}{
		"orgId":     event.OrgId,
		"pathKey":   event.PathKey,
		"assetUUID": event.AssetUUID,
		"assetName": event.AssetName,
		"eventId":   eventId,
		"eventType": event.Type,
		"created":   created,
	}
	if event.LastSeenAt != nil {
		payload["lastSeenAt"] = event.LastSeenAt.UTC().Format(time.RFC3339Nano)
	}
	if event.ThresholdMinutes > 0 {
		payload["thresholdMinutes"] = event.ThresholdMinutes
	}
	if event.MissCount > 0 {
		payload["missCount"] = event.MissCount
	}
	return payload
}
