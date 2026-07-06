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
			sentAt := time.Now().UTC()
			if err := postForceOffline(c, c.MustGetString(BagKeyAssetUUID), reason); err != nil {
				return err
			}
			c.Set(BagKeyForceOfflineSentAt, sentAt)
			return nil
		},
	}
}

// ForceOfflineByAdminByLabel is the label-scoped variant of ForceOfflineByAdmin: it
// forces the asset at AssetUUIDKey(label) offline, so a journey that provisions
// several assets can flip a specific one. It does not record a sent-at timestamp
// (presence-only journeys assert the health status directly, not an events window).
//
// Reads (bag):
//   - AssetUUIDKey(label)  string  set by CreateAssetWithLabel
func ForceOfflineByAdminByLabel(label, reason string) saga.Step {
	if reason == "" {
		reason = "saga-force-offline"
	}
	return saga.Step{
		Name: "assets/assets.ForceOfflineByAdmin[" + label + "]",
		Do: func(c *saga.Context) error {
			return postForceOffline(c, c.MustGetString(AssetUUIDKey(label)), reason)
		},
	}
}

// postForceOffline POSTs the internal force-offline endpoint for an asset uuid.
// Shared by the keyed and non-keyed variants; the endpoint is idempotent (a no-op
// when the asset is already alerted) and gated by the assets internal API key.
func postForceOffline(c *saga.Context, uuid, reason string) error {
	payload := map[string]string{"reason": reason}
	headers := map[string]string{"X-API-Key": constants.InternalApiKey}
	resp, err := c.Clients.Assets.RawWithHeaders(
		c.Stdctx, http.MethodPost,
		"/internal/health_monitor/"+uuid+"/force_offline", payload, headers,
	)
	if err != nil {
		return fmt.Errorf("force offline: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("force offline: unexpected status %d", resp.StatusCode)
	}
	return nil
}
