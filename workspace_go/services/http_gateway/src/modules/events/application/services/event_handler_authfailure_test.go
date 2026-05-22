package services

import (
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"http_gateway/src/bootstrap"
	dsDto "http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/events/application/di"
	"http_gateway/src/modules/events/application/ports"

	dsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/datasources"
	eventsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/events"
	eventsDto "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// Compile-time check — the inline mock must satisfy the EventBusPort.
var _ ports.EventBusPort = (*mockEventBus)(nil)

// mockEventBus captures Publish/PublishCore invocations from the events
// module so unit tests can assert what landed on NATS without spinning up a
// broker. Concurrency-safe because PublishAuthFailure publishes from a
// goroutine.
type mockEventBus struct {
	mu sync.Mutex
	wg sync.WaitGroup

	publishCalls     int
	lastPublishCfg   natsModel.PublishConfig
	publishErr       error

	publishCoreCalls   int
	lastPublishCoreCfg natsModel.PublishCoreConfig
	publishCoreErr     error

	flushCalls int
	flushErr   error
}

func (m *mockEventBus) Publish(cfg natsModel.PublishConfig) error {
	defer m.wg.Done()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishCalls++
	m.lastPublishCfg = cfg
	return m.publishErr
}

func (m *mockEventBus) PublishCore(cfg natsModel.PublishCoreConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishCoreCalls++
	m.lastPublishCoreCfg = cfg
	return m.publishCoreErr
}

func (m *mockEventBus) FlushConnection() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flushCalls++
	return m.flushErr
}

// expectOnePublish registers the goroutine fire-and-forget invocation count
// so the test can synchronously wait for it.
func (m *mockEventBus) expectOnePublish() {
	m.wg.Add(1)
}

// waitOrFail blocks up to timeout for the Publish goroutine to land.
func (m *mockEventBus) waitOrFail(t *testing.T, timeout time.Duration) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatalf("Publish: timed out after %s", timeout)
	}
}

// newAuthFailureService composes an EventService with the mock bus and a
// minimal Metrics struct. PublishAuthFailure does NOT touch the metrics, but
// the struct must be non-nil because the broader service references it.
func newAuthFailureService(bus *mockEventBus) *EventService {
	return &EventService{
		deps: di.EventServiceDependenciesInjection{
			NatsBus: bus,
			Metrics: &bootstrap.HttpGatewayMetrics{
				HeartbeatsTotal: prometheus.NewCounterVec(
					prometheus.CounterOpts{Name: "test_hb_total"},
					[]string{"status"},
				),
				HeartbeatDuration: prometheus.NewHistogram(
					prometheus.HistogramOpts{Name: "test_hb_duration"},
				),
			},
		},
	}
}

// orgPtr is a helper to allocate an ObjectID pointer used by DataSourceResponse.
func orgPtrAuth(t *testing.T, hex string) *model.ObjectId {
	t.Helper()
	o, err := model.ToObjectID(hex)
	if err != nil {
		t.Fatalf("ToObjectID(%s): %v", hex, err)
	}
	return &o
}

// stringPtrAuth returns a pointer to s — DataSourceResponse fields are *string.
func stringPtrAuth(s string) *string { return &s }

// TestPublishAuthFailure_BuildsRawEventDTO asserts the §3-refactor contract:
// the public method delegates to buildAuthFailurePayload + publishAuthFailureFireAndForget,
// the resulting RawEventDTO carries Success=false + the documented fields,
// and the publish lands on EventsRawSubject.
func TestPublishAuthFailure_BuildsRawEventDTO(t *testing.T) {
	bus := &mockEventBus{}
	bus.expectOnePublish()
	svc := newAuthFailureService(bus)

	ds := &dsContract.DataSourceResponse{
		ID:          orgPtrAuth(t, "68f5bbce1aef22967c3ebb30"),
		OrgId:       orgPtrAuth(t, "68f5bbce1aef22967c3ebb31"),
		PathKey:     stringPtrAuth("000001"),
		Name:        stringPtrAuth("My DataSource"),
		Description: stringPtrAuth("test source"),
	}

	svc.PublishAuthFailure(ds, map[string]any{"k": "v"}, "tracker-uuid-1", "Unauthorized - Invalid API Key")
	bus.waitOrFail(t, 1*time.Second)

	bus.mu.Lock()
	defer bus.mu.Unlock()

	if bus.publishCalls != 1 {
		t.Fatalf("Publish: want 1 call, got %d", bus.publishCalls)
	}
	if bus.lastPublishCfg.Subject != eventsContract.SubjectEventsRaw {
		t.Errorf("Subject: want %q, got %q", eventsContract.SubjectEventsRaw, bus.lastPublishCfg.Subject)
	}

	payload, ok := bus.lastPublishCfg.Data.(eventsDto.RawEventDTO)
	if !ok {
		t.Fatalf("Data: want eventsDto.RawEventDTO, got %T", bus.lastPublishCfg.Data)
	}
	if payload.Success {
		t.Errorf("Success: want false (auth failure), got true")
	}
	if payload.Error != "Unauthorized - Invalid API Key" {
		t.Errorf("Error: want %q, got %q", "Unauthorized - Invalid API Key", payload.Error)
	}
	if payload.EventTrackerId != "tracker-uuid-1" {
		t.Errorf("EventTrackerId: want %q, got %q", "tracker-uuid-1", payload.EventTrackerId)
	}
	if payload.OrgId != ds.OrgId.Hex() {
		t.Errorf("OrgId: want %q, got %q", ds.OrgId.Hex(), payload.OrgId)
	}
	if payload.PathKey != "000001" {
		t.Errorf("PathKey: want %q, got %q", "000001", payload.PathKey)
	}
	if payload.ThreadId != ds.ID.Hex() {
		t.Errorf("ThreadId: want %q, got %q", ds.ID.Hex(), payload.ThreadId)
	}
	if payload.Source != "http_gateway" {
		t.Errorf("Source: want %q, got %q", "http_gateway", payload.Source)
	}
}

// TestPublishAuthFailure_NilDataSource_NoPanic asserts the defensive branches
// inside buildAuthFailurePayload — when dataSource is nil the helper must
// still produce a valid (mostly-empty) RawEventDTO without panicking. The
// auth audit trail is best-effort and must not crash the request flow.
func TestPublishAuthFailure_NilDataSource_NoPanic(t *testing.T) {
	bus := &mockEventBus{}
	bus.expectOnePublish()
	svc := newAuthFailureService(bus)

	svc.PublishAuthFailure(nil, nil, "tracker-uuid-2", "missing Authorization header")
	bus.waitOrFail(t, 1*time.Second)

	bus.mu.Lock()
	defer bus.mu.Unlock()

	if bus.publishCalls != 1 {
		t.Fatalf("Publish: want 1 call, got %d", bus.publishCalls)
	}
	payload, ok := bus.lastPublishCfg.Data.(eventsDto.RawEventDTO)
	if !ok {
		t.Fatalf("Data: want eventsDto.RawEventDTO, got %T", bus.lastPublishCfg.Data)
	}
	if payload.OrgId != "" || payload.PathKey != "" || payload.ThreadId != "" || payload.Name != "" {
		t.Errorf("nil dataSource fields should be empty strings, got %+v", payload)
	}
	if payload.Error != "missing Authorization header" {
		t.Errorf("Error: want %q, got %q", "missing Authorization header", payload.Error)
	}
}

// silenceUnused keeps the dsDto import alive when other tests in this file
// don't reference it directly.
var _ = (*dsDto.DataSourceResponse)(nil)
