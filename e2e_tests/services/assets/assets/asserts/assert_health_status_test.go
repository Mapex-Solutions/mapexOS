package asserts

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// TestAssertHealthStatus_TimeoutSurfacesLastError proves that when the endpoint keeps
// failing, the timeout message carries the real cause (the last fetch error) rather
// than a misleading empty last-seen value.
func TestAssertHealthStatus_TimeoutSurfacesLastError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
	c.Set(assetSteps.BagKeyAssetID, "asset-123")

	a := AssertHealthStatusEventuallyWithTimeout("online", 150*time.Millisecond, 30*time.Millisecond)
	err := a.Check(c)
	if err == nil {
		t.Fatal("expected a timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "last error") {
		t.Errorf("error %q does not surface the last fetch error", err.Error())
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error %q does not carry the 500 status behind the failure", err.Error())
	}
}
