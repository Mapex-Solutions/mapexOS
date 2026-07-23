package payloads

import (
	"encoding/json"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets_templates"
)

// The marketplace e2e journeys each install from the mock as the SAME seed-admin org
// (the suite is single-tenant by design), and a marketplace install is keyed by
// (marketplaceGuid, org), NOT by runID — so it is the one piece of state runID cannot
// isolate. To keep the journeys parallel-safe, each marketplace-installing journey
// serves and installs its OWN (vendor, slug, guid) template, so their installs, asset
// usage and installed-checks never collide:
//
//   - DTCo2*     -> phase1_install_check (install/uninstall + batch installed-check)
//   - Guard*     -> phase2_uninstall_guard (install + asset + usage-guarded uninstall)
//   - Pipeline*  -> marketplace_asset_pipeline (full telemetry pipeline)
//
// Each identity's three values MUST stay in lockstep with the bundle the mock serves:
// the mock advertises the guid in the X-Marketplace-Guid header, the assets service
// stores it on install, and the batch installed-check asserts against it.
const (
	// DTCo2Vendor is the marketplace vendor path segment shared by the DT CO2 test
	// templates (only the slug/guid differ per journey).
	DTCo2Vendor = "disruptive-technologies"

	// DTCo2Slug / DTCo2MarketplaceGuid identify phase1's install-check template.
	DTCo2Slug            = "co2-sensor"
	DTCo2MarketplaceGuid = "e2e-dt-co2-sensor-0001"

	// GuardSlug / GuardMarketplaceGuid identify phase2's uninstall-guard template.
	GuardSlug            = "co2-sensor-guard"
	GuardMarketplaceGuid = "e2e-dt-co2-guard-0002"

	// PipelineSlug / PipelineMarketplaceGuid identify the full-pipeline template.
	PipelineSlug            = "co2-sensor-pipeline"
	PipelineMarketplaceGuid = "e2e-dt-co2-pipeline-0003"
)

// passthroughConversion leaves the event untouched. It is enough for journeys that
// only install/uninstall or create-then-delete an asset without sending telemetry.
const passthroughConversion = "function convert(event) { return event; }"

// standardizedConversion emits the platform's canonical StandardizedPayload
// (`{ eventType, eventId, data, created }`) from the incoming telemetry, so an
// installed copy actually drives the event pipeline through js-executor -> router ->
// triggers. The payload field names match the shared SagaTelemetryEvent the datasource
// harness posts.
const standardizedConversion = `const result = {
  eventType: 'co2',
  eventId: payload.runId,
  data: { value: payload.value, unit: payload.unit },
  created: payload.timestamp,
};`

// DTCo2Bundle is phase1's install-check bundle: a minimal-but-valid installable DT CO2
// template with a pass-through conversion (phase1 never sends telemetry).
func DTCo2Bundle() []byte { return co2Bundle("DT CO2 Sensor", "Sensor de CO2 DT", "deviceId", passthroughConversion) }

// GuardBundle is phase2's uninstall-guard bundle. Pass-through is enough: phase2
// creates an asset only to trip the usage guard, then deletes it, without telemetry.
func GuardBundle() []byte {
	return co2Bundle("DT CO2 Sensor (Guard)", "Sensor de CO2 DT (Guarda)", "deviceId", passthroughConversion)
}

// PipelineBundle is the full-pipeline bundle: AssetIDPath and the standardizing
// conversion match the SagaTelemetryEvent harness so an installed copy drives the
// event pipeline end to end.
func PipelineBundle() []byte {
	return co2Bundle("DT CO2 Sensor (Pipeline)", "Sensor de CO2 DT (Pipeline)", "assetUUID", standardizedConversion)
}

// co2Bundle builds a DT-classified CO2 MarketplaceBundle for the given display name,
// asset-identity path and conversion script — the only bits that differ across the
// three e2e templates. The mock computes the sha256 of these exact bytes at serve
// time, so the checksum the assets service hard-verifies always matches.
func co2Bundle(nameEN, namePT, assetIDPath, scriptConversion string) []byte {
	bundle := contracts.MarketplaceBundle{
		Name:             map[string]string{"en-US": nameEN, "pt-BR": namePT},
		Description:      map[string]string{"en-US": "Disruptive Technologies CO2 sensor.", "pt-BR": "Sensor de CO2 da Disruptive Technologies."},
		CategorySlug:     "environmental-sensors",
		CategoryName:     "Environmental Sensors",
		ManufacturerSlug: "disruptive-technologies",
		ManufacturerName: "Disruptive Technologies",
		ModelSlug:        "dt-co2",
		ModelName:        "CO2 Sensor",
		Version:          "1.0.0",
		AssetIDPath:      assetIDPath,
		ScriptConversion: scriptConversion,
		AvailableFields:  []string{"co2", "temperature", "humidity"},
		DynamicFields:    []contracts.MarketplaceDynamicField{},
		NextFieldId:      1,
	}
	raw, err := json.Marshal(bundle)
	if err != nil {
		// The bundle is a static literal; marshaling it cannot fail in practice.
		panic("marshal CO2 marketplace bundle: " + err.Error())
	}
	return raw
}
