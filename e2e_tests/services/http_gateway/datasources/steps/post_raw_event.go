package steps

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	dsPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/payloads"
)

// BagKeyAssetUUID is the asset UUID bag key written by the assets
// CreateConnectivityAsset step. Re-declared here so the gateway step
// can read it without importing the assets package (which would
// create a cycle in some compositions).
const BagKeyAssetUUID = "assets.assetUUID"

// PostRawEvent POSTs a SagaTelemetryEvent body to
// /api/v1/events?ds=<dsID> using the data source API key. The request
// drives the platform's full ingestion pipeline:
//
//	gateway → NATS → js-executor (template script) → NATS → router → trigger
//
// Used by phase2_event_pipeline journeys to validate the telemetry
// path independent of the healthmonitor's force-online/force-offline
// shortcuts exercised in phase1.
//
// Reads (bag):
//   - BagKeyAssetUUID        string  set by CreateConnectivityAsset
//   - BagKeyDataSourceID     string  set by CreateDataSource
//   - BagKeyDataSourceApiKey string  set by CreateDataSource
//
// Writes (bag):
//   - BagKeyTelemetrySentAt  time.Time  captured before the POST returns;
//     callers can scope events_trigger queries to events after this.
//
// Compensate: no-op. Telemetry events are best-effort signals;
// nothing to undo on the gateway side.
func PostRawEvent() saga.Step {
	return saga.Step{
		Name: "http_gateway/datasources.PostRawEvent",
		Do: func(c *saga.Context) error {
			uuid := c.MustGetString(BagKeyAssetUUID)
			dsID := c.MustGetString(BagKeyDataSourceID)
			apiKey := c.MustGetString(BagKeyDataSourceApiKey)

			body := dsPayloads.SagaTelemetryEvent(c.RunID, uuid, 23.5)
			headers := map[string]string{dsPayloads.SagaApiKeyHeaderName: apiKey}

			sentAt := time.Now().UTC()
			resp, err := c.Clients.Gateway.RawWithHeaders(
				c.Stdctx, http.MethodPost,
				"/api/v1/events?ds="+dsID, body, headers,
			)
			if err != nil {
				return fmt.Errorf("post raw event: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("post raw event: unexpected status %d", resp.StatusCode)
			}
			c.Set(BagKeyTelemetrySentAt, sentAt)
			return nil
		},
	}
}
