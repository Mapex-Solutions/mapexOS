package steps

import (
	"net/http"
	"testing"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// TestCreateAsset_Compensate_Contract proves the CreateAsset Compensate — which now
// actually runs after the saga rollback fix — is safe under the three states the
// rollback can hit: the asset was never created, it exists, or it is already gone.
func TestCreateAsset_Compensate_Contract(t *testing.T) {
	const label = "x"

	t.Run("bag key absent → no-op, no request", func(t *testing.T) {
		srv := newRecorderServer(http.StatusOK)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))

		if err := CreateAssetWithLabel(nil, label).Compensate(c); err != nil {
			t.Fatalf("Compensate with absent bag key = %v, want nil", err)
		}
		if srv.count() != 0 {
			t.Errorf("issued %d requests, want 0 when the asset was never created", srv.count())
		}
	})

	t.Run("entity present → one DELETE, nil", func(t *testing.T) {
		srv := newRecorderServer(http.StatusOK)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(AssetIDKey(label), "asset-1")

		if err := CreateAssetWithLabel(nil, label).Compensate(c); err != nil {
			t.Fatalf("Compensate present = %v, want nil", err)
		}
		if srv.count() != 1 {
			t.Fatalf("issued %d requests, want exactly 1 DELETE", srv.count())
		}
		if got := srv.requests[0]; got.method != http.MethodDelete || got.path != "/api/v1/assets/asset-1" {
			t.Errorf("request = %+v, want DELETE /api/v1/assets/asset-1", got)
		}
	})

	t.Run("already gone (404) → tolerated, nil", func(t *testing.T) {
		srv := newRecorderServer(http.StatusNotFound)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(AssetIDKey(label), "asset-1")

		if err := CreateAssetWithLabel(nil, label).Compensate(c); err != nil {
			t.Errorf("Compensate on 404 = %v, want nil (idempotent teardown tolerates not-found)", err)
		}
	})

	t.Run("idempotent → running twice stays safe", func(t *testing.T) {
		srv := newRecorderServer(http.StatusOK)
		defer srv.Close()
		c := saga.NewTestContext(t, saga.NewClientSet(saga.ClientURLs{Assets: srv.URL}))
		c.Set(AssetIDKey(label), "asset-1")

		step := CreateAssetWithLabel(nil, label)
		if err := step.Compensate(c); err != nil {
			t.Fatalf("first Compensate = %v", err)
		}
		if err := step.Compensate(c); err != nil {
			t.Fatalf("second Compensate = %v", err)
		}
		if srv.count() != 2 {
			t.Errorf("issued %d requests over two runs, want 2", srv.count())
		}
	})
}
