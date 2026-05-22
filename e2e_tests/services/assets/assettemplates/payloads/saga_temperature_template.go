// Package payloads holds canonical AssetTemplateCreate fixtures for the
// assets/assettemplates module.
//
// The asset template carries the script the js-executor runs on every
// telemetry event — without a template the asset cannot exist (asset
// validation requires assetTemplateId, mongoid-shaped). Saga journeys
// reuse SagaTemperatureTemplate as the simplest template that still lets
// MQTT telemetry flow through the script-processor → router → events
// pipeline.
package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets_templates"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// AssetTemplateCreateBuilder wraps the contract DTO with fluent overrides.
type AssetTemplateCreateBuilder struct {
	spec contracts.AssetTemplateCreate
}

// Build returns the contracts payload ready for
// POST /api/v1/asset_templates.
func (b *AssetTemplateCreateBuilder) Build() contracts.AssetTemplateCreate {
	return b.spec
}

// WithName overrides the template name.
func (b *AssetTemplateCreateBuilder) WithName(name string) *AssetTemplateCreateBuilder {
	b.spec.Name = name
	return b
}

// SagaTemperatureTemplate returns the canonical template for an MQTT
// temperature sensor saga. ScriptConversion produces the platform's
// canonical StandardizedPayload shape — `{ eventType, eventId, data,
// created }` — so the engine validation passes and downstream
// consumers (router, events) see a record they can persist.
//
// Mapping from the saga's device JSON (PublishTelemetry step):
//
//	{ runId, timestamp, unit, value }  →
//	{
//	  eventType: 'temperature',
//	  eventId:    payload.runId,
//	  data:       { value, unit },
//	  created:    payload.timestamp,
//	}
//
// Defaults:
//   - Enabled:      true
//   - AssetIDPath:  "assetUUID" (mqtt username = device id)
func SagaTemperatureTemplate(runID string) *AssetTemplateCreateBuilder {
	const scriptConversion = `const result = {
  eventType: 'temperature',
  eventId: payload.runId,
  data: { value: payload.value, unit: payload.unit },
  created: payload.timestamp,
};`
	return &AssetTemplateCreateBuilder{
		spec: contracts.AssetTemplateCreate{
			Name:             fmt.Sprintf("saga-temperature-%s", runID),
			Description:      zerovalue.Ptr("Saga temperature sensor template (StandardizedPayload transform)"),
			Enabled:          true,
			AssetIDPath:      "assetUUID",
			ScriptConversion: scriptConversion,
		},
	}
}
