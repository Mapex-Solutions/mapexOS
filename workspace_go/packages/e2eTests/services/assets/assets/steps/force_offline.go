package steps

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// ForceOfflineByAdmin POSTs the assets internal force-offline endpoint
// so the saga can exercise the offline-action route group without
// waiting the scheduler+threshold window the scanner enforces in
// production. The endpoint is gated by the assets internal API key —
// the same key shared by every L3 fallback consumer.
//
// Reads (bag):
//   - BagKeyAssetUUID  string  set by CreateAsset
//
// Writes (bag):
//   - BagKeyForceOfflineSentAt  time.Time  captured immediately before
//     the POST returns, used by AssertOfflineActionExecuted to scope the
//     events search window.
//
// Compensate: no-op. The endpoint is idempotent (the service skips when
// the asset is already alerted), and the route-group / asset cleanups
// happen in their own Compensate paths.
func ForceOfflineByAdmin(reason string) saga.Step {
	if reason == "" {
		reason = "saga-force-offline"
	}
	return saga.Step{
		Name: "assets/assets.ForceOfflineByAdmin",
		Do: func(c *saga.Context) error {
			uuid := c.MustGetString(BagKeyAssetUUID)
			payload := map[string]string{"reason": reason}
			headers := map[string]string{"X-API-Key": constants.InternalApiKey}

			sentAt := time.Now().UTC()
			resp, err := c.Clients.Assets.RawWithHeaders(
				c.Stdctx, http.MethodPost,
				"/internal/health-monitor/"+uuid+"/force-offline", payload, headers,
			)
			if err != nil {
				return fmt.Errorf("force offline: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("force offline: unexpected status %d", resp.StatusCode)
			}
			c.Set(BagKeyForceOfflineSentAt, sentAt)
			return nil
		},
	}
}
