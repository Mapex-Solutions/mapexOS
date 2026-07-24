package steps

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// recorderServer is a minimal test double that records each request and replies with
// a fixed status, so the Compensate contract can be asserted without the live stack.
type recorderServer struct {
	*httptest.Server
	status int
	mu     sync.Mutex
	method string
	path   string
	calls  int
}

func newRecorder(status int) *recorderServer {
	r := &recorderServer{status: status}
	r.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.mu.Lock()
		r.calls++
		r.method, r.path = req.Method, req.URL.Path
		r.mu.Unlock()
		w.WriteHeader(r.status)
	}))
	return r
}

func (r *recorderServer) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls
}

// TestCreateMigrationPlan_Compensate proves the Compensate is bag-driven and tolerant
// of both a missing plan (404) and a terminal, non-editable plan (409), so a fully
// green run never fails during rollback.
func TestCreateMigrationPlan_Compensate(t *testing.T) {
	step := CreateMigrationPlan("from", "to", "a1")

	t.Run("bag key absent → no-op, no request", func(t *testing.T) {
		srv := newRecorder(http.StatusOK)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))

		if err := step.Compensate(c); err != nil {
			t.Fatalf("Compensate with absent bag key = %v, want nil", err)
		}
		if srv.count() != 0 {
			t.Errorf("issued %d requests, want 0 when no plan was created", srv.count())
		}
	})

	t.Run("plan present → one DELETE, nil", func(t *testing.T) {
		srv := newRecorder(http.StatusOK)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(BagKeyMigrationPlanID, "plan-1")

		if err := step.Compensate(c); err != nil {
			t.Fatalf("Compensate present = %v, want nil", err)
		}
		if srv.count() != 1 {
			t.Fatalf("issued %d requests, want exactly 1 DELETE", srv.count())
		}
		if srv.method != http.MethodDelete || srv.path != "/api/v1/asset_templates/migrations/plan-1" {
			t.Errorf("request = %s %s, want DELETE /api/v1/asset_templates/migrations/plan-1", srv.method, srv.path)
		}
	})

	t.Run("already gone (404) → tolerated, nil", func(t *testing.T) {
		srv := newRecorder(http.StatusNotFound)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(BagKeyMigrationPlanID, "plan-1")

		if err := step.Compensate(c); err != nil {
			t.Errorf("Compensate on 404 = %v, want nil", err)
		}
	})

	t.Run("terminal plan (409) → tolerated, nil", func(t *testing.T) {
		srv := newRecorder(http.StatusConflict)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(BagKeyMigrationPlanID, "plan-1")

		if err := step.Compensate(c); err != nil {
			t.Errorf("Compensate on 409 (terminal plan) = %v, want nil", err)
		}
	})

	t.Run("other non-2xx (500) → error", func(t *testing.T) {
		srv := newRecorder(http.StatusInternalServerError)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(BagKeyMigrationPlanID, "plan-1")

		if err := step.Compensate(c); err == nil {
			t.Error("Compensate on 500 = nil, want an error")
		}
	})
}
