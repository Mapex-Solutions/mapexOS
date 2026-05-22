package steps

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	dsPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/payloads"
	dsSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/steps"
)

// SendHttpHeartbeat POSTs /api/v1/heartbeat against the http_gateway
// so the saga's HTTP-protocol asset transitions to online through the
// same channel a real device would use. Body shape matches the
// healthmonitor explicit-mode contract: `{ "assetUUID": "..." }`.
//
// Reads (bag):
//   - BagKeyAssetUUID                       string  set by CreateAsset
//   - dsSteps.BagKeyDataSourceID            string  set by CreateDataSource
//   - dsSteps.BagKeyDataSourceApiKey        string  set by CreateDataSource
//
// Writes (bag):
//   - BagKeyHeartbeatSentAt  time.Time  captured immediately before
//     the POST returns, used by AssertOnlineActionExecuted to scope the
//     events search window.
//
// Compensate: no-op. Heartbeats are best-effort state pings; nothing
// to undo on the http_gateway side.
func SendHttpHeartbeat() saga.Step {
	return saga.Step{
		Name: "assets/assets.SendHttpHeartbeat",
		Do: func(c *saga.Context) error {
			uuid := c.MustGetString(BagKeyAssetUUID)
			dsID := c.MustGetString(dsSteps.BagKeyDataSourceID)
			apiKey := c.MustGetString(dsSteps.BagKeyDataSourceApiKey)

			body := map[string]string{"assetUUID": uuid}
			headers := map[string]string{dsPayloads.SagaApiKeyHeaderName: apiKey}

			sentAt := time.Now().UTC()
			resp, err := c.Clients.Gateway.RawWithHeaders(
				c.Stdctx, http.MethodPost,
				"/api/v1/heartbeat?ds="+dsID, body, headers,
			)
			if err != nil {
				return fmt.Errorf("send heartbeat: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("send heartbeat: unexpected status %d", resp.StatusCode)
			}
			c.Set(BagKeyHeartbeatSentAt, sentAt)
			return nil
		},
	}
}
