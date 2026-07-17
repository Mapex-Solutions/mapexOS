package steps

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// recordedRequest is one (method, path) the recorder server observed.
type recordedRequest struct {
	method string
	path   string
}

// recorderServer is an httptest server that records every request and replies with a
// caller-chosen status, so a step's Compensate can be exercised without a live stack.
type recorderServer struct {
	*httptest.Server
	mu       sync.Mutex
	requests []recordedRequest
}

// newRecorderServer starts a recorder replying with status to every request.
func newRecorderServer(status int) *recorderServer {
	rs := &recorderServer{}
	rs.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rs.mu.Lock()
		rs.requests = append(rs.requests, recordedRequest{r.Method, r.URL.Path})
		rs.mu.Unlock()
		w.WriteHeader(status)
	}))
	return rs
}

// count returns how many requests the server has seen.
func (rs *recorderServer) count() int {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return len(rs.requests)
}

// TestCreatePlan_Compensate_Contract proves the CreatePlan Compensate — which
// cancels the plan via DELETE /api/v1/ota/plans/{id} — is safe under the states
// the rollback can hit: the plan was never created, it exists, or it is already
// gone (404); and that a real error (500) surfaces.
func TestCreatePlan_Compensate_Contract(t *testing.T) {
	t.Run("bag key absent → no-op, no request", func(t *testing.T) {
		srv := newRecorderServer(http.StatusOK)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))

		if err := CreatePlan("src").Compensate(c); err != nil {
			t.Fatalf("Compensate with absent bag key = %v, want nil", err)
		}
		if srv.count() != 0 {
			t.Errorf("issued %d requests, want 0 when the plan was never created", srv.count())
		}
	})

	t.Run("plan present → one DELETE, nil", func(t *testing.T) {
		srv := newRecorderServer(http.StatusOK)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(BagKeyPlanID, "plan-1")

		if err := CreatePlan("src").Compensate(c); err != nil {
			t.Fatalf("Compensate present = %v, want nil", err)
		}
		if srv.count() != 1 {
			t.Fatalf("issued %d requests, want exactly 1 DELETE", srv.count())
		}
		if got := srv.requests[0]; got.method != http.MethodDelete || got.path != "/api/v1/ota/plans/plan-1" {
			t.Errorf("request = %+v, want DELETE /api/v1/ota/plans/plan-1", got)
		}
	})

	t.Run("already gone (404) → tolerated, nil", func(t *testing.T) {
		srv := newRecorderServer(http.StatusNotFound)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(BagKeyPlanID, "plan-1")

		if err := CreatePlan("src").Compensate(c); err != nil {
			t.Errorf("Compensate on 404 = %v, want nil (idempotent teardown tolerates not-found)", err)
		}
	})

	t.Run("server error (500) → surfaces the error", func(t *testing.T) {
		srv := newRecorderServer(http.StatusInternalServerError)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(BagKeyPlanID, "plan-1")

		if err := CreatePlan("src").Compensate(c); err == nil {
			t.Error("Compensate on 500 = nil, want a non-nil error")
		}
	})
}
