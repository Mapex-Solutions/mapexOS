package nats

import (
	"context"
	"errors"
	"testing"
	"time"

	"assets/src/modules/healthmonitor/domain/entities"

	hmContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// spyPublisher records every Publish call. errOn is used to fail publishes
// whose subject matches the map key (to simulate per-path failures).
type spyPublisher struct {
	calls []natsModel.PublishConfig
	errOn map[string]error
}

func (s *spyPublisher) Publish(config natsModel.PublishConfig) error {
	s.calls = append(s.calls, config)
	if err, ok := s.errOn[config.Subject]; ok {
		return err
	}
	return nil
}

func sampleEvent(includeOptional bool, routeGroupIds []string) entities.AlertEvent {
	lastSeen := time.Date(2026, 4, 21, 10, 14, 0, 0, time.UTC)
	e := entities.AlertEvent{
		Type:          hmContract.EventTypeOffline,
		OrgId:         "org-1",
		AssetUUID:     "asset-xyz",
		AssetName:     "Door Sensor",
		PathKey:       "org-1/floor-2",
		RouteGroupIds: routeGroupIds,
	}
	if includeOptional {
		e.LastSeenAt = &lastSeen
		e.ThresholdMinutes = 10
		e.MissCount = 3
	}
	return e
}

func TestPublishAlert_DualPublishWithRouteGroups(t *testing.T) {
	spy := &spyPublisher{}
	p := &alertPublisher{publisher: spy}
	event := sampleEvent(true, []string{"rg-1", "rg-2"})

	if err := p.publishAlert(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spy.calls) != 2 {
		t.Fatalf("expected 2 publishes, got %d", len(spy.calls))
	}

	// First publish: persistence
	if spy.calls[0].Subject != hmContract.AssetStatusSaveSubject {
		t.Fatalf("first publish subject = %q, want %q", spy.calls[0].Subject, hmContract.AssetStatusSaveSubject)
	}
	persistence, ok := spy.calls[0].Data.(map[string]interface{})
	if !ok {
		t.Fatal("persistence payload is not a map")
	}
	for _, field := range []string{"orgId", "pathKey", "assetUUID", "assetName", "eventId", "eventType", "created", "lastSeenAt", "thresholdMinutes", "missCount"} {
		if _, present := persistence[field]; !present {
			t.Errorf("persistence payload missing field %q", field)
		}
	}
	// Persistence must be FLAT — no nested "event" object.
	if _, hasEvent := persistence["event"]; hasEvent {
		t.Error("persistence payload must be flat, not nested under 'event'")
	}
	if persistence["eventType"] != hmContract.EventTypeOffline {
		t.Errorf("persistence eventType = %v, want %v", persistence["eventType"], hmContract.EventTypeOffline)
	}

	// Second publish: router
	if spy.calls[1].Subject != hmContract.RouterSubject {
		t.Fatalf("second publish subject = %q, want %q", spy.calls[1].Subject, hmContract.RouterSubject)
	}
	route, ok := spy.calls[1].Data.(map[string]interface{})
	if !ok {
		t.Fatal("route payload is not a map")
	}
	if route["eventSource"] != hmContract.EventSourceHealthStatus {
		t.Errorf("route eventSource = %v, want %q", route["eventSource"], hmContract.EventSourceHealthStatus)
	}
	if _, ok := route["eventTrackerId"].(string); !ok {
		t.Error("route payload missing string eventTrackerId")
	}
	nestedEvent, ok := route["event"].(map[string]interface{})
	if !ok {
		t.Fatal("route.event is not a map")
	}
	// StandardizedPayload shape check.
	for _, field := range []string{"eventType", "eventId", "data", "metadata", "created"} {
		if _, present := nestedEvent[field]; !present {
			t.Errorf("nested event missing %q", field)
		}
	}
	metadata, ok := nestedEvent["metadata"].(map[string]interface{})
	if !ok || metadata["source"] != hmContract.EventSource {
		t.Errorf("nested event metadata.source = %v, want %q", metadata["source"], hmContract.EventSource)
	}
	data, ok := nestedEvent["data"].(map[string]interface{})
	if !ok || len(data) == 0 {
		t.Error("nested event data must be non-empty map")
	}
	if data["lastSeenAt"] == nil {
		t.Error("nested event data.lastSeenAt should be present when event.LastSeenAt is set")
	}

	// Shared identity — eventId and created match across both payloads.
	if persistence["eventId"] != nestedEvent["eventId"] {
		t.Errorf("eventId mismatch: persistence=%v, route=%v", persistence["eventId"], nestedEvent["eventId"])
	}
	if persistence["created"] != nestedEvent["created"] {
		t.Errorf("created mismatch: persistence=%v, route=%v", persistence["created"], nestedEvent["created"])
	}
}

func TestPublishAlert_NoRouteGroupsSkipsRoutePublish(t *testing.T) {
	spy := &spyPublisher{}
	p := &alertPublisher{publisher: spy}
	event := sampleEvent(true, nil)

	if err := p.publishAlert(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spy.calls) != 1 {
		t.Fatalf("expected 1 publish (persistence only), got %d", len(spy.calls))
	}
	if spy.calls[0].Subject != hmContract.AssetStatusSaveSubject {
		t.Fatalf("only publish should be persistence, got subject=%q", spy.calls[0].Subject)
	}
}

func TestPublishAlert_OptionalFieldsOmittedWhenZero(t *testing.T) {
	spy := &spyPublisher{}
	p := &alertPublisher{publisher: spy}
	event := sampleEvent(false, []string{"rg-1"}) // no LastSeenAt / Threshold / MissCount

	if err := p.publishAlert(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spy.calls) != 2 {
		t.Fatalf("expected 2 publishes, got %d", len(spy.calls))
	}

	persistence := spy.calls[0].Data.(map[string]interface{})
	for _, field := range []string{"lastSeenAt", "thresholdMinutes", "missCount"} {
		if _, present := persistence[field]; present {
			t.Errorf("persistence payload must omit %q when zero-value", field)
		}
	}

	route := spy.calls[1].Data.(map[string]interface{})
	nestedData := route["event"].(map[string]interface{})["data"].(map[string]interface{})
	for _, field := range []string{"lastSeenAt", "thresholdMinutes", "missCount"} {
		if _, present := nestedData[field]; present {
			t.Errorf("nested event.data must omit %q when zero-value", field)
		}
	}
	// Data must still be non-empty (assetUUID + assetName).
	if len(nestedData) < 2 {
		t.Error("nested event.data must contain assetUUID + assetName at minimum")
	}
}

func TestPublishAlert_PersistenceFailureDoesNotBlockRoute(t *testing.T) {
	spy := &spyPublisher{
		errOn: map[string]error{
			hmContract.AssetStatusSaveSubject: errors.New("persistence stream down"),
		},
	}
	p := &alertPublisher{publisher: spy}
	event := sampleEvent(true, []string{"rg-1"})

	if err := p.publishAlert(event); err != nil {
		t.Fatalf("publishAlert should not propagate errors, got %v", err)
	}

	if len(spy.calls) != 2 {
		t.Fatalf("expected both publishes attempted, got %d", len(spy.calls))
	}
	if spy.calls[1].Subject != hmContract.RouterSubject {
		t.Errorf("route publish should still run after persistence failure, got subject=%q", spy.calls[1].Subject)
	}
}

func TestPublishAlert_EntryPointsUseDualPublish(t *testing.T) {
	spy := &spyPublisher{}
	p := &alertPublisher{publisher: spy}
	event := sampleEvent(false, []string{"rg-1"})

	if err := p.PublishOffline(context.Background(), event); err != nil {
		t.Fatalf("PublishOffline error: %v", err)
	}
	if err := p.PublishOnline(context.Background(), event); err != nil {
		t.Fatalf("PublishOnline error: %v", err)
	}

	// 2 publishes per call (persistence + route); 2 calls → 4 total.
	if len(spy.calls) != 4 {
		t.Fatalf("expected 4 publishes across PublishOffline + PublishOnline, got %d", len(spy.calls))
	}
}
