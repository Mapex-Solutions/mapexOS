package services

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"assets/src/bootstrap"
	assetPorts "assets/src/modules/assets/application/ports"
	"assets/src/modules/healthmonitor/application/di"
	"assets/src/modules/healthmonitor/application/ports"
	"assets/src/modules/healthmonitor/domain/entities"

	hmContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Compile-time checks — the inline fakes must implement the real ports.
var (
	_ ports.HealthRepository    = (*mockHealthRepo)(nil)
	_ ports.AlertPublisherPort  = (*mockAlertPublisher)(nil)
	_ assetPorts.AssetRepository = (*mockAssetRepo)(nil)
)

/*
 * MOCKS — inline private types co-located with their consumer tests.
 *
 * Each mock implements its real port. Test bodies seed return values via
 * exported struct fields and assert behavior via call counters / captured
 * arguments. No third-party mocking lib (testify/gomock) is used here —
 * matches the convention in healthmonitor_handler_enrichment_test.go.
 */

// mockHealthRepo records call counts per method and lets each test stub the
// return value for the bools that drive the heartbeat state machine.
type mockHealthRepo struct {
	updateLastSeenCalls   int
	resetMissCounterCalls int
	registerOrgCalls      int
	removeAlertedCalls    int
	markKnownOnlineCalls  int
	isKnownOnlineCalls    int
	setLastConnectAtCalls int
	getLastConnectAtCalls int
	markAlertedCalls      int

	updateLastSeenErr     error
	removeAlertedReturn   bool // first return — true means "won the SREM"
	removeAlertedErr      error
	isAlertedReturn       bool // first return — true means "asset is already offline"
	isAlertedErr          error
	isKnownOnlineReturn   bool
	getLastConnectAtValue *time.Time // nil → "never connected" branch in presence handler
	getLastConnectAtErr   error
}

func (m *mockHealthRepo) UpdateLastSeen(_ context.Context, _ string, _ string, _ time.Time) error {
	m.updateLastSeenCalls++
	return m.updateLastSeenErr
}

func (m *mockHealthRepo) ResetMissCounter(_ context.Context, _ string, _ string) error {
	m.resetMissCounterCalls++
	return nil
}

func (m *mockHealthRepo) IsAlerted(_ context.Context, _ string, _ string) (bool, error) {
	return m.isAlertedReturn, m.isAlertedErr
}

func (m *mockHealthRepo) RemoveAlerted(_ context.Context, _ string, _ string) (bool, error) {
	m.removeAlertedCalls++
	return m.removeAlertedReturn, m.removeAlertedErr
}

func (m *mockHealthRepo) RegisterOrg(_ context.Context, _ string) error {
	m.registerOrgCalls++
	return nil
}

func (m *mockHealthRepo) IsKnownOnline(_ context.Context, _ string, _ string) (bool, error) {
	m.isKnownOnlineCalls++
	return m.isKnownOnlineReturn, nil
}

func (m *mockHealthRepo) MarkKnownOnline(_ context.Context, _ string, _ string) error {
	m.markKnownOnlineCalls++
	return nil
}

func (m *mockHealthRepo) FindStale(_ context.Context, _ string, _ time.Time, _ int64, _ int64) ([]string, error) {
	return nil, nil
}

func (m *mockHealthRepo) IncrementMiss(_ context.Context, _ string, _ string) (int64, error) {
	return 0, nil
}

func (m *mockHealthRepo) MarkAlerted(_ context.Context, _ string, _ string) error {
	m.markAlertedCalls++
	return nil
}

func (m *mockHealthRepo) GetActiveOrgs(_ context.Context) ([]string, error) {
	return nil, nil
}

func (m *mockHealthRepo) GetLastSeen(_ context.Context, _ string, _ string) (*time.Time, error) {
	return nil, nil
}

func (m *mockHealthRepo) GetLastSeenBatch(_ context.Context, _ string, _ []string) (map[string]*time.Time, error) {
	return nil, nil
}

func (m *mockHealthRepo) IsAlertedBatch(_ context.Context, _ string, _ []string) (map[string]bool, error) {
	return nil, nil
}

func (m *mockHealthRepo) RemoveAsset(_ context.Context, _ string, _ string) error {
	return nil
}

func (m *mockHealthRepo) SetLastConnectAt(_ context.Context, _ string, _ string, _ time.Time) error {
	m.setLastConnectAtCalls++
	return nil
}

func (m *mockHealthRepo) GetLastConnectAt(_ context.Context, _ string, _ string) (*time.Time, error) {
	m.getLastConnectAtCalls++
	return m.getLastConnectAtValue, m.getLastConnectAtErr
}

// mockAssetRepo implements assetPorts.AssetRepository. Only the two methods
// reached by the heartbeat path (FindByAssetUUID + UpdateHealthStatusWithChangedAt)
// have meaningful behavior; the rest panic if reached so a test that drifts
// through the wrong code path fails loudly.
type mockAssetRepo struct {
	asset                          *assetPorts.Asset
	findErr                        error
	findByAssetUUIDCalls           int
	updateHealthStatusCalls        int
	updateHealthStatusLastStatus   string
	updateHealthStatusErr          error
}

func (m *mockAssetRepo) Create(_ context.Context, _ *assetPorts.Asset) (*assetPorts.Asset, error) {
	panic("mockAssetRepo.Create not expected on heartbeat path")
}

func (m *mockAssetRepo) FindById(_ context.Context, _ *string) (*assetPorts.Asset, error) {
	panic("mockAssetRepo.FindById not expected on heartbeat path")
}

func (m *mockAssetRepo) FindByAssetUUID(_ context.Context, _ *string) (*assetPorts.Asset, error) {
	m.findByAssetUUIDCalls++
	return m.asset, m.findErr
}

func (m *mockAssetRepo) FindByMqttUsername(_ context.Context, _ string) (*assetPorts.Asset, error) {
	panic("mockAssetRepo.FindByMqttUsername not expected on heartbeat path")
}

func (m *mockAssetRepo) FindByIdAndUpdate(_ context.Context, _ *string, _ map[string]any) (*assetPorts.Asset, error) {
	panic("mockAssetRepo.FindByIdAndUpdate not expected on heartbeat path")
}

func (m *mockAssetRepo) DeleteById(_ context.Context, _ *string) error {
	panic("mockAssetRepo.DeleteById not expected on heartbeat path")
}

func (m *mockAssetRepo) FindWithFilters(
	_ context.Context, _ model.Map, _ *model.PaginationOpts, _ model.Map,
) (*model.PaginatedResult[assetPorts.Asset], error) {
	panic("mockAssetRepo.FindWithFilters not expected on heartbeat path")
}

func (m *mockAssetRepo) FindWithFiltersAndTemplate(
	_ context.Context, _ model.Map, _ model.Map, _ *model.PaginationOpts, _ model.Map,
) (*model.PaginatedResult[assetPorts.AssetWithTemplate], error) {
	panic("mockAssetRepo.FindWithFiltersAndTemplate not expected on heartbeat path")
}

func (m *mockAssetRepo) CountDocuments(_ context.Context, _ model.Map) (int64, error) {
	panic("mockAssetRepo.CountDocuments not expected on heartbeat path")
}

func (m *mockAssetRepo) UpdateHealthStatusWithChangedAt(_ context.Context, _ *string, status string, _ time.Time) error {
	m.updateHealthStatusCalls++
	m.updateHealthStatusLastStatus = status
	return m.updateHealthStatusErr
}

// mockAlertPublisher captures the events handed to PublishOnline / PublishOffline
// so tests can assert RouteGroupIds (the load-bearing field for monitor-only
// vs monitor+route mode).
type mockAlertPublisher struct {
	publishOnlineCalls  int
	publishOfflineCalls int
	lastOnlineEvent     *entities.AlertEvent
	lastOfflineEvent    *entities.AlertEvent
}

func (m *mockAlertPublisher) PublishOnline(_ context.Context, event entities.AlertEvent) error {
	m.publishOnlineCalls++
	captured := event
	m.lastOnlineEvent = &captured
	return nil
}

func (m *mockAlertPublisher) PublishOffline(_ context.Context, event entities.AlertEvent) error {
	m.publishOfflineCalls++
	captured := event
	m.lastOfflineEvent = &captured
	return nil
}

/*
 * FIXTURES
 */

// newTestMetrics returns a minimal AssetsMetrics wired to a throwaway prometheus
// registry — only the counters touched by the heartbeat AND presence paths need
// to be present so both test files share the same fixture without duplicate
// registry-name clashes.
func newTestMetrics() *bootstrap.AssetsMetrics {
	reg := metrics.NewRegistry("assets_hm_heartbeat_test")
	return &bootstrap.AssetsMetrics{
		Registry: reg,
		HealthHeartbeatsReceived: reg.NewCounter(metrics.CounterOpts{
			Subsystem: "health", Name: "heartbeats_received_total", Help: "test",
		}),
		HealthOnlineTransitions: reg.NewCounter(metrics.CounterOpts{
			Subsystem: "health", Name: "online_transitions_total", Help: "test",
		}),
		HealthAlertsPublished: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "health", Name: "alerts_published_total", Help: "test",
		}, []string{"type"}),
		HealthRedisErrors: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "health", Name: "redis_errors_total", Help: "test",
		}, []string{"operation"}),
		HealthPresenceReceived: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "health", Name: "presence_received_total", Help: "test",
		}, []string{"action"}),
		HealthPresenceFiltered: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "health", Name: "presence_filtered_total", Help: "test",
		}, []string{"reason"}),
		HealthPresenceProcessed: reg.NewCounterVec(metrics.CounterOpts{
			Subsystem: "health", Name: "presence_processed_total", Help: "test",
		}, []string{"outcome"}),
		HealthPresenceHandlerDuration: reg.NewHistogram(metrics.HistogramOpts{
			Subsystem: "health", Name: "presence_handler_duration_seconds", Help: "test",
			Buckets: prometheus.DefBuckets,
		}),
	}
}

// Make sure prometheus import is referenced — silences an unused-import error
// when the test file is built standalone.
var _ = prometheus.DefBuckets

// newHeartbeatService composes a HealthMonitorService with the three mocked
// ports plus a real (throwaway) AssetsMetrics struct. Returns the concrete
// type so tests can call the public HandleHeartbeat entry directly.
func newHeartbeatService(
	healthRepo *mockHealthRepo,
	assetRepo *mockAssetRepo,
	publisher *mockAlertPublisher,
) *HealthMonitorService {
	return &HealthMonitorService{
		deps: di.HealthMonitorServiceDI{
			HealthRepo:     healthRepo,
			AlertPublisher: publisher,
			AssetRepo:      assetRepo,
			Metrics:        newTestMetrics(),
		},
	}
}

// newHeartbeatMessage marshals a HeartbeatEvent into a *natsModel.Message ready
// for HandleHeartbeat. Uses NewTestMessage with a non-nil but empty callbacks
// struct so Ack/Nack/Reject become no-ops without dereferencing the (nil)
// jetstream rawMsg. Tests that need to assert ack-vs-nack supply real callbacks.
func newHeartbeatMessage(t *testing.T, hb hmContract.HeartbeatEvent) *natsModel.Message {
	t.Helper()
	payload, err := json.Marshal(hb)
	if err != nil {
		t.Fatalf("marshal heartbeat: %v", err)
	}
	return natsModel.NewTestMessage(payload, 0, &natsModel.TestMessageCallbacks{})
}

// monitoredAsset builds an asset with the minimum fields the heartbeat path
// reads: OrgID, AssetUUID, Name, PathKey, and HealthMonitor.
func monitoredAsset(orgHex string, uuid string, hm *assetPorts.HealthMonitorConfig) *assetPorts.Asset {
	orgID, _ := model.ToObjectID(orgHex)
	return &assetPorts.Asset{
		ID:            model.ObjectId(primitive.NewObjectID()),
		OrgID:         orgID,
		AssetUUID:     uuid,
		Name:          "test-asset",
		PathKey:       "P0001",
		HealthMonitor: hm,
	}
}

/*
 * TESTS
 */

// TestHandleHeartbeat_DisabledAsset_NoSideEffects — gate behavior introduced
// in T6: when HealthMonitor.Enabled=false, the handler must drop the message
// after Acking it, with zero Redis writes and zero NATS publishes.
func TestHandleHeartbeat_DisabledAsset_NoSideEffects(t *testing.T) {
	healthRepo := &mockHealthRepo{}
	publisher := &mockAlertPublisher{}
	assetRepo := &mockAssetRepo{
		asset: monitoredAsset(
			"68f5bbce1aef22967c3ebb30",
			"uuid-1",
			&assetPorts.HealthMonitorConfig{Enabled: false},
		),
	}

	svc := newHeartbeatService(healthRepo, assetRepo, publisher)
	msg := newHeartbeatMessage(t, hmContract.HeartbeatEvent{
		OrgId:     "org-1",
		AssetUUID: "uuid-1",
		Timestamp: time.Now().Unix(),
	})

	svc.HandleHeartbeat(msg)

	if assetRepo.findByAssetUUIDCalls != 1 {
		t.Errorf("FindByAssetUUID: want 1 call, got %d", assetRepo.findByAssetUUIDCalls)
	}
	if healthRepo.updateLastSeenCalls != 0 {
		t.Errorf("UpdateLastSeen MUST NOT be called when disabled, got %d calls", healthRepo.updateLastSeenCalls)
	}
	if healthRepo.removeAlertedCalls != 0 {
		t.Errorf("RemoveAlerted MUST NOT be called when disabled, got %d calls", healthRepo.removeAlertedCalls)
	}
	if assetRepo.updateHealthStatusCalls != 0 {
		t.Errorf("UpdateHealthStatusWithChangedAt MUST NOT be called when disabled, got %d calls", assetRepo.updateHealthStatusCalls)
	}
	if publisher.publishOnlineCalls != 0 {
		t.Errorf("PublishOnline MUST NOT be called when disabled, got %d calls", publisher.publishOnlineCalls)
	}
	if publisher.publishOfflineCalls != 0 {
		t.Errorf("PublishOffline MUST NOT be called when disabled, got %d calls", publisher.publishOfflineCalls)
	}
}

// TestHandleHeartbeat_MonitorOnly_PublisherReceivesEmptyArrays — when the
// asset is monitored but no OnlineRouteGroupIds are configured, the offline→online
// transition must still publish to the persistence/router boundary, with an
// empty (or nil) RouteGroupIds slice. Router gating is the consumer's concern.
func TestHandleHeartbeat_MonitorOnly_PublisherReceivesEmptyArrays(t *testing.T) {
	healthRepo := &mockHealthRepo{
		removeAlertedReturn: true, // won the SREM → transition fires
	}
	publisher := &mockAlertPublisher{}
	assetRepo := &mockAssetRepo{
		asset: monitoredAsset(
			"68f5bbce1aef22967c3ebb30",
			"uuid-1",
			&assetPorts.HealthMonitorConfig{
				Enabled:              true,
				ThresholdMinutes:     10,
				OfflineRouteGroupIds: nil,
				OnlineRouteGroupIds:  nil,
			},
		),
	}

	svc := newHeartbeatService(healthRepo, assetRepo, publisher)
	msg := newHeartbeatMessage(t, hmContract.HeartbeatEvent{
		OrgId:     "org-1",
		AssetUUID: "uuid-1",
		Timestamp: time.Now().Unix(),
	})

	svc.HandleHeartbeat(msg)

	if publisher.publishOnlineCalls != 1 {
		t.Fatalf("PublishOnline: want 1 call (transition), got %d", publisher.publishOnlineCalls)
	}
	if publisher.lastOnlineEvent == nil {
		t.Fatal("lastOnlineEvent: expected captured event, got nil")
	}
	if len(publisher.lastOnlineEvent.RouteGroupIds) != 0 {
		t.Errorf("RouteGroupIds: want empty/nil for monitor-only mode, got %v",
			publisher.lastOnlineEvent.RouteGroupIds)
	}
	if publisher.publishOfflineCalls != 0 {
		t.Errorf("PublishOffline MUST NOT be called on online transition, got %d calls",
			publisher.publishOfflineCalls)
	}
	if assetRepo.updateHealthStatusLastStatus != "online" {
		t.Errorf("UpdateHealthStatusWithChangedAt status: want %q, got %q",
			"online", assetRepo.updateHealthStatusLastStatus)
	}
}

// TestHandleHeartbeat_MonitorAndRoute_PublisherReceivesArrays — when the asset
// declares OnlineRouteGroupIds, the publisher MUST receive that exact slice so
// router fan-out can fire downstream.
func TestHandleHeartbeat_MonitorAndRoute_PublisherReceivesArrays(t *testing.T) {
	healthRepo := &mockHealthRepo{
		removeAlertedReturn: true, // won the SREM → transition fires
	}
	publisher := &mockAlertPublisher{}
	expectedRouteGroups := []string{"x"}
	assetRepo := &mockAssetRepo{
		asset: monitoredAsset(
			"68f5bbce1aef22967c3ebb30",
			"uuid-1",
			&assetPorts.HealthMonitorConfig{
				Enabled:             true,
				ThresholdMinutes:    10,
				OnlineRouteGroupIds: expectedRouteGroups,
			},
		),
	}

	svc := newHeartbeatService(healthRepo, assetRepo, publisher)
	msg := newHeartbeatMessage(t, hmContract.HeartbeatEvent{
		OrgId:     "org-1",
		AssetUUID: "uuid-1",
		Timestamp: time.Now().Unix(),
	})

	svc.HandleHeartbeat(msg)

	if publisher.publishOnlineCalls != 1 {
		t.Fatalf("PublishOnline: want 1 call (transition), got %d", publisher.publishOnlineCalls)
	}
	if publisher.lastOnlineEvent == nil {
		t.Fatal("lastOnlineEvent: expected captured event, got nil")
	}
	if !reflect.DeepEqual(publisher.lastOnlineEvent.RouteGroupIds, expectedRouteGroups) {
		t.Errorf("RouteGroupIds: want %v, got %v",
			expectedRouteGroups, publisher.lastOnlineEvent.RouteGroupIds)
	}
}

// TestHandleHeartbeat_OriginAgnostic verifies the heartbeat consumer
// invariant: it cannot tell apart heartbeats originating from js-executor (implicit)
// from heartbeats originating from http_gateway (POST /api/v1/heartbeat with
// body { assetUUID }). Both publishers target `mapexos.asset.heartbeat.{orgId}`
// with payload {orgId, assetUUID, ts} and the handler MUST produce identical
// Redis side-effects regardless of source.
//
// We exercise the invariant by feeding TWO messages with identical payloads
// (representing the same logical heartbeat coming from two different origins)
// and asserting Redis state-mutation calls are exactly N×.
func TestHandleHeartbeat_OriginAgnostic(t *testing.T) {
	healthRepo := &mockHealthRepo{}
	publisher := &mockAlertPublisher{}
	assetRepo := &mockAssetRepo{
		asset: monitoredAsset(
			"68f5bbce1aef22967c3ebb30",
			"uuid-origin-agnostic",
			&assetPorts.HealthMonitorConfig{
				Enabled:          true,
				ThresholdMinutes: 10,
			},
		),
	}

	svc := newHeartbeatService(healthRepo, assetRepo, publisher)

	hb := hmContract.HeartbeatEvent{
		OrgId:     "org-1",
		AssetUUID: "uuid-origin-agnostic",
		Timestamp: time.Now().Unix(),
	}

	// Origin 1: simulates implicit publish from js-executor.
	svc.HandleHeartbeat(newHeartbeatMessage(t, hb))
	// Origin 2: simulates MQTT republish (or HTTP gateway) — same payload shape.
	svc.HandleHeartbeat(newHeartbeatMessage(t, hb))

	if assetRepo.findByAssetUUIDCalls != 2 {
		t.Errorf("FindByAssetUUID: want 2 calls (one per heartbeat), got %d", assetRepo.findByAssetUUIDCalls)
	}
	if healthRepo.updateLastSeenCalls != 2 {
		t.Errorf("UpdateLastSeen: want 2 calls (origin-agnostic), got %d", healthRepo.updateLastSeenCalls)
	}
	if healthRepo.removeAlertedCalls != 2 {
		t.Errorf("RemoveAlerted: want 2 calls (origin-agnostic), got %d", healthRepo.removeAlertedCalls)
	}
}

// nackingMessage builds a heartbeat message whose Ack/Nack/Reject callbacks
// record which terminal action ran. Used by Bug #1 fix tests to assert the
// 3-state contract of loadAssetForHeartbeat (transient err → Nack, drop → Ack).
type ackTracker struct {
	ackCalled    int
	nackCalled   int
	rejectCalled int
	nackErr      error
}

func newTrackedHeartbeatMessage(t *testing.T, hb hmContract.HeartbeatEvent, tracker *ackTracker) *natsModel.Message {
	t.Helper()
	payload, err := json.Marshal(hb)
	if err != nil {
		t.Fatalf("marshal heartbeat: %v", err)
	}
	return natsModel.NewTestMessage(payload, 0, &natsModel.TestMessageCallbacks{
		OnAck:    func() error { tracker.ackCalled++; return nil },
		OnNack:   func(err error) error { tracker.nackCalled++; tracker.nackErr = err; return nil },
		OnReject: func(_ string) error { tracker.rejectCalled++; return nil },
	})
}

// TestHandleHeartbeat_TransientLookupError_Nacks asserts the Bug #1 fix:
// when AssetRepo.FindByAssetUUID returns a transient error, the handler MUST
// Nack the message so NATS redelivers — the message is NOT silently dropped.
// Zero Redis writes, zero AlertPublisher publishes.
func TestHandleHeartbeat_TransientLookupError_Nacks(t *testing.T) {
	healthRepo := &mockHealthRepo{}
	publisher := &mockAlertPublisher{}
	assetRepo := &mockAssetRepo{
		findErr: errSentinel("simulated mongo timeout"),
	}

	svc := newHeartbeatService(healthRepo, assetRepo, publisher)
	tracker := &ackTracker{}
	msg := newTrackedHeartbeatMessage(t, hmContract.HeartbeatEvent{
		OrgId:     "org-1",
		AssetUUID: "uuid-1",
		Timestamp: time.Now().Unix(),
	}, tracker)

	svc.HandleHeartbeat(msg)

	if tracker.nackCalled != 1 {
		t.Errorf("Nack: want 1 call (transient error), got %d", tracker.nackCalled)
	}
	if tracker.ackCalled != 0 {
		t.Errorf("Ack MUST NOT be called on transient error, got %d", tracker.ackCalled)
	}
	if tracker.nackErr == nil {
		t.Error("Nack: expected non-nil error to drive NATS retry, got nil")
	}
	if healthRepo.updateLastSeenCalls != 0 {
		t.Errorf("UpdateLastSeen MUST NOT be called on transient error, got %d", healthRepo.updateLastSeenCalls)
	}
	if publisher.publishOnlineCalls != 0 || publisher.publishOfflineCalls != 0 {
		t.Errorf("AlertPublisher MUST NOT be called on transient error, got online=%d offline=%d",
			publisher.publishOnlineCalls, publisher.publishOfflineCalls)
	}
}

// TestHandleHeartbeat_AssetNotFound_Acks asserts the Bug #1 fix:
// when AssetRepo.FindByAssetUUID returns (nil, nil) — no error, just absent —
// the handler MUST Ack the message (legitimate drop) and NOT touch any
// downstream side effects.
func TestHandleHeartbeat_AssetNotFound_Acks(t *testing.T) {
	healthRepo := &mockHealthRepo{}
	publisher := &mockAlertPublisher{}
	assetRepo := &mockAssetRepo{
		asset:   nil,
		findErr: nil,
	}

	svc := newHeartbeatService(healthRepo, assetRepo, publisher)
	tracker := &ackTracker{}
	msg := newTrackedHeartbeatMessage(t, hmContract.HeartbeatEvent{
		OrgId:     "org-1",
		AssetUUID: "uuid-not-found",
		Timestamp: time.Now().Unix(),
	}, tracker)

	svc.HandleHeartbeat(msg)

	if tracker.ackCalled != 1 {
		t.Errorf("Ack: want 1 call (legitimate drop), got %d", tracker.ackCalled)
	}
	if tracker.nackCalled != 0 {
		t.Errorf("Nack MUST NOT be called when asset is missing, got %d", tracker.nackCalled)
	}
	if healthRepo.updateLastSeenCalls != 0 {
		t.Errorf("UpdateLastSeen MUST NOT be called for absent asset, got %d", healthRepo.updateLastSeenCalls)
	}
	if publisher.publishOnlineCalls != 0 || publisher.publishOfflineCalls != 0 {
		t.Errorf("AlertPublisher MUST NOT be called for absent asset, got online=%d offline=%d",
			publisher.publishOnlineCalls, publisher.publishOfflineCalls)
	}
}

// errSentinel is a minimal error type so test bodies stay free of fmt.Errorf
// noise when seeding mock errors.
type errSentinel string

func (e errSentinel) Error() string { return string(e) }
