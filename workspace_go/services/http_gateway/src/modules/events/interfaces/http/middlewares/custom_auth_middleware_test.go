package middlewares

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"http_gateway/src/bootstrap"
	dsDto "http_gateway/src/modules/datasources/application/dtos"
	dsPort "http_gateway/src/modules/datasources/application/ports"
	"http_gateway/src/modules/events/application/dtos"
	eventsPort "http_gateway/src/modules/events/application/ports"

	dsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/datasources"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
)

// Compile-time checks — minimal mocks must satisfy the real ports.
var (
	_ dsPort.DataSourceServicePort = (*mockDataSourceService)(nil)
	_ eventsPort.EventServicePort  = (*mockEventService)(nil)
)

// mockDataSourceService implements only GetDataSourceById; every other port
// method panics so a misrouted code path fails loudly during the test.
type mockDataSourceService struct {
	response *dsDto.DataSourceResponse
	err      error
	calls    int
}

func (m *mockDataSourceService) GetDataSourceById(_ context.Context, _ *string) (*dsDto.DataSourceResponse, error) {
	m.calls++
	return m.response, m.err
}

func (m *mockDataSourceService) GetDataSources(_ context.Context, _ *reqCtx.RequestContext, _ *dsDto.DataSourceQueryDTO) (*model.PaginatedResult[dsDto.DataSourceResponse], error) {
	panic("GetDataSources not expected on the auth-gate path")
}

func (m *mockDataSourceService) CreateDataSource(_ context.Context, _ *reqCtx.RequestContext, _ *dsDto.DataSourceCreateDTO) (*dsDto.DataSourceResponse, error) {
	panic("CreateDataSource not expected on the auth-gate path")
}

func (m *mockDataSourceService) UpdateDataSourceById(_ context.Context, _ *string, _ *dsDto.DataSourceUpdateDTO) (*dsDto.DataSourceResponse, error) {
	panic("UpdateDataSourceById not expected on the auth-gate path")
}

func (m *mockDataSourceService) DeleteDataSourceById(_ context.Context, _ *string) (map[string]bool, error) {
	panic("DeleteDataSourceById not expected on the auth-gate path")
}

// mockEventService captures PublishAuthFailure calls and panics on the other
// port methods. Concurrency-safe because PublishAuthFailure happens inside
// the synchronous middleware body — but kept locked for safety.
type mockEventService struct {
	mu                   sync.Mutex
	publishAuthFailureN  int
	lastErrorMsg         string
}

func (m *mockEventService) PublishAuthFailure(_ *dsDto.DataSourceResponse, _ map[string]any, _ string, errorMsg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishAuthFailureN++
	m.lastErrorMsg = errorMsg
}

func (m *mockEventService) ProcessEvent(_ context.Context, _ map[string]any, _ *dsDto.DataSourceResponse) (map[string]bool, error) {
	panic("ProcessEvent not expected on the auth-gate path")
}

func (m *mockEventService) ProcessHeartbeat(_ context.Context, _ *dsDto.DataSourceResponse, _ string) error {
	panic("ProcessHeartbeat not expected on the auth-gate path")
}

// newGateMetrics returns a minimal HttpGatewayMetrics wired to a fresh
// prometheus registry — only the 3 auth metrics touched by the gate are
// populated so the assertions can inspect the "disabled" label.
func newGateMetrics() *bootstrap.HttpGatewayMetrics {
	return &bootstrap.HttpGatewayMetrics{
		EventAuthTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "test_event_auth_total"},
			[]string{"type", "result"},
		),
		EventAuthDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "test_event_auth_duration_seconds"},
			[]string{"type"},
		),
		EventAuthFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "test_event_auth_failures_total"},
			[]string{"type"},
		),
	}
}

// newGateApp builds a Fiber app whose POST /test pipeline pre-populates
// c.Locals("queryDTO", &EvenIdentificationDto{Ds: ds}) so the production
// middleware can resolve the DataSource as if the validation middleware ran.
// Returns the app + the dsDto pointer for cleanup-free test calls.
func newGateApp(dsService dsPort.DataSourceServicePort, eventService eventsPort.EventServicePort, m *bootstrap.HttpGatewayMetrics) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var sc *customErrors.ServerCustomError
			if errors.As(err, &sc) {
				return c.Status(sc.Code).JSON(fiber.Map{"errors": sc.Errors})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Post("/test",
		func(c *fiber.Ctx) error {
			ds := c.Query("ds")
			c.Locals("queryDTO", &dtos.EvenIdentificationDto{Ds: &ds})
			return c.Next()
		},
		CustomAuthMiddleware(dsService, eventService, m),
		func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		},
	)
	return app
}

// counterValue reads a prometheus.CounterVec sample for a given label set.
// Returns 0 when the label combination has not been observed yet.
func counterValue(t *testing.T, vec *prometheus.CounterVec, labels ...string) float64 {
	t.Helper()
	c, err := vec.GetMetricWithLabelValues(labels...)
	if err != nil {
		t.Fatalf("GetMetricWithLabelValues(%v): %v", labels, err)
	}
	pb := &dto.Metric{}
	if err := c.Write(pb); err != nil {
		t.Fatalf("write metric: %v", err)
	}
	return pb.GetCounter().GetValue()
}

// boolPtr returns a pointer to b — Fiber DataSourceResponse.Enabled is *bool.
func boolPtr(b bool) *bool { return &b }

// orgPtr is a helper to allocate an ObjectID pointer used by DataSourceResponse.
func orgPtr(t *testing.T, hex string) *model.ObjectId {
	t.Helper()
	o, err := model.ToObjectID(hex)
	if err != nil {
		t.Fatalf("ToObjectID(%s): %v", hex, err)
	}
	return &o
}

// TestCustomAuthMiddleware_DisabledDataSource_Returns403 asserts that a
// disabled DataSource is rejected with 403 BEFORE any auth-strategy logic
// runs. PublishAuthFailure must be called for the audit trail and the
// "disabled" label must be incremented on all 3 existing auth metrics.
func TestCustomAuthMiddleware_DisabledDataSource_Returns403(t *testing.T) {
	dsService := &mockDataSourceService{
		response: &dsContract.DataSourceResponse{
			ID:      orgPtr(t, "68f5bbce1aef22967c3ebb30"),
			OrgId:   orgPtr(t, "68f5bbce1aef22967c3ebb31"),
			Enabled: boolPtr(false),
			Auth: &dsContract.DataSourceAuth{Type: "apiKey"},
		},
	}
	eventService := &mockEventService{}
	m := newGateMetrics()

	app := newGateApp(dsService, eventService, m)
	req := httptest.NewRequest("POST", "/test?ds=68f5bbce1aef22967c3ebb30", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("status: want 403, got %d (body=%s)", resp.StatusCode, body)
	}
	eventService.mu.Lock()
	defer eventService.mu.Unlock()
	if eventService.publishAuthFailureN != 1 {
		t.Errorf("PublishAuthFailure: want 1 call, got %d", eventService.publishAuthFailureN)
	}
	if eventService.lastErrorMsg != "DataSource is disabled" {
		t.Errorf("errorMsg: want %q, got %q", "DataSource is disabled", eventService.lastErrorMsg)
	}
	if got := counterValue(t, m.EventAuthTotal, "disabled", "failure"); got != 1 {
		t.Errorf(`EventAuthTotal{type="disabled",result="failure"}: want 1, got %v`, got)
	}
	if got := counterValue(t, m.EventAuthFailures, "disabled"); got != 1 {
		t.Errorf(`EventAuthFailures{type="disabled"}: want 1, got %v`, got)
	}
}

// TestCustomAuthMiddleware_NilEnabled_Returns403 asserts the defensive
// branch: when DataSourceResponse.Enabled is nil (legacy or partial DTO),
// the gate treats it as disabled. Same audit + metrics behavior.
func TestCustomAuthMiddleware_NilEnabled_Returns403(t *testing.T) {
	dsService := &mockDataSourceService{
		response: &dsContract.DataSourceResponse{
			ID:      orgPtr(t, "68f5bbce1aef22967c3ebb30"),
			OrgId:   orgPtr(t, "68f5bbce1aef22967c3ebb31"),
			Enabled: nil,
			Auth: &dsContract.DataSourceAuth{Type: "apiKey"},
		},
	}
	eventService := &mockEventService{}
	m := newGateMetrics()

	app := newGateApp(dsService, eventService, m)
	req := httptest.NewRequest("POST", "/test?ds=68f5bbce1aef22967c3ebb30", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("status: want 403 (nil Enabled treated as disabled), got %d", resp.StatusCode)
	}
}
