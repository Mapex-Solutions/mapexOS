package payloads

import (
	"time"
)

// SagaTelemetryEvent returns the body the saga POSTs to
// /api/v1/events?ds=<dsID> to drive the telemetry pipeline. Shape
// matches what SagaTemperatureTemplate's ScriptConversion expects:
//
//	payload.runId      → eventId
//	payload.value      → data.value
//	payload.unit       → data.unit
//	payload.timestamp  → created
//
// assetUUID is required by the asset-resolution middleware on the
// gateway side; the template's AssetIDPath ("assetUUID") tells the
// js-executor which event field carries the device id.
func SagaTelemetryEvent(runID, assetUUID string, value float64) map[string]any {
	return map[string]any{
		"assetUUID": assetUUID,
		"runId":     runID,
		"value":     value,
		"unit":      "C",
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	}
}
