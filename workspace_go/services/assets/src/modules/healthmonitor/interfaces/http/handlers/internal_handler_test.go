package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"assets/src/modules/healthmonitor/application/ports"
)

// mockHealthAdmin captures the arguments to ForceOfflineByAssetUUID so
// handler tests can assert path-param decoding + body propagation
// without spinning up the real service.
type mockHealthAdmin struct {
	calls            int
	lastAssetUUID    string
	lastReason       string
	returnErr        error
}

var _ ports.HealthAdminPort = (*mockHealthAdmin)(nil)

func (m *mockHealthAdmin) ForceOfflineByAssetUUID(_ context.Context, assetUUID, reason string) error {
	m.calls++
	m.lastAssetUUID = assetUUID
	m.lastReason = reason
	return m.returnErr
}

// newAppWithRoute builds a minimal Fiber app with the route under test
// registered at /:assetUUID/force-offline. The 'app.Test' driver
// short-circuits the real server, so each subtest stays isolated.
func newAppWithRoute(svc ports.HealthAdminPort) *fiber.App {
	app := fiber.New()
	app.Post("/:assetUUID/force-offline", ForceOffline(svc))
	return app
}

// TestForceOffline_HappyPath_NoBody — 204 + service called with empty
// reason (handler does NOT default; defaulting lives in the service).
func TestForceOffline_HappyPath_NoBody(t *testing.T) {
	svc := &mockHealthAdmin{}
	app := newAppWithRoute(svc)

	req := httptest.NewRequest(http.MethodPost, "/asset-abc/force-offline", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: unexpected err: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status: want 204, got %d", resp.StatusCode)
	}
	if svc.calls != 1 {
		t.Fatalf("service calls: want 1, got %d", svc.calls)
	}
	if svc.lastAssetUUID != "asset-abc" {
		t.Errorf("assetUUID: want %q, got %q", "asset-abc", svc.lastAssetUUID)
	}
	if svc.lastReason != "" {
		t.Errorf("reason: handler must not default — want empty string, got %q", svc.lastReason)
	}
}

// TestForceOffline_HappyPath_WithReason — body reason flows through to
// the service unchanged.
func TestForceOffline_HappyPath_WithReason(t *testing.T) {
	svc := &mockHealthAdmin{}
	app := newAppWithRoute(svc)

	body := strings.NewReader(`{"reason":"ci-saga-mqtt-trigger"}`)
	req := httptest.NewRequest(http.MethodPost, "/uuid-1/force-offline", body)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: unexpected err: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status: want 204, got %d", resp.StatusCode)
	}
	if svc.lastReason != "ci-saga-mqtt-trigger" {
		t.Errorf("reason: want %q, got %q", "ci-saga-mqtt-trigger", svc.lastReason)
	}
}

// TestForceOffline_InvalidJSON_BadRequest — broken JSON body returns
// 400 before the service is touched. Guards against the handler
// silently consuming bad input on the wire.
func TestForceOffline_InvalidJSON_BadRequest(t *testing.T) {
	svc := &mockHealthAdmin{}
	app := newAppWithRoute(svc)

	body := bytes.NewReader([]byte(`{not-json`))
	req := httptest.NewRequest(http.MethodPost, "/uuid-1/force-offline", body)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: unexpected err: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", resp.StatusCode)
	}
	if svc.calls != 0 {
		t.Errorf("service must not be called on bad body, got %d calls", svc.calls)
	}
}

// TestForceOffline_ServiceError_PropagatesToFiber — service-side error
// is returned by the handler so Fiber's default error handler responds
// with 500. Keeps the contract that admin endpoints surface failures
// instead of swallowing them.
func TestForceOffline_ServiceError_PropagatesToFiber(t *testing.T) {
	svc := &mockHealthAdmin{returnErr: errors.New("redis down")}
	app := newAppWithRoute(svc)

	req := httptest.NewRequest(http.MethodPost, "/uuid-1/force-offline", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: unexpected err: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status: want 500 (Fiber default for plain error), got %d", resp.StatusCode)
	}
}
