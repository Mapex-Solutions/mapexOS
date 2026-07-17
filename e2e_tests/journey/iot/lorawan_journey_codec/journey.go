// Package lorawan_journey_codec exercises the LoRaWAN device CODEC end to end
// against the live stack: a real device's uplink bytes are DECODED by the
// template's first script into the device's semantic fields — not left as raw
// bytes. It is the decode counterpart of the raw-ingest LoRaWAN journeys
// (lorawan_journey_bsp / lorawan_journey_otaa), which assert the undecoded frame;
// this one asserts the codec's structured output.
//
// The template's codec is the REAL Dragino LHT65N-S-DS decoder (decodeUplink) from
// the mapexMarketplace catalog, embedded in the SagaLorawanDecodedTemplate payload.
// The sensor fires a deterministic LHT65N fPort-2 real-time frame whose decoded
// values are fixed (temperature 21.00 C, humidity 60.0 %, battery 3.0 V "Good"),
// so the assert can pin exact fields.
//
// Flow: route group + decoded template → connect a UDP gateway (shared block) →
// provision + activate an ABP sensor → fire the LHT65N uplink → assert the
// js-executor produced a StandardizedPayload whose data carries the decoded
// fields → assert the sensor is online.
//
// Outcome on PASS:
//   - The uplink is decoded by the device codec: a successful js-executor event
//     for the sensor carries data.TempC_SHT=21, data.Hum_SHT=60, data.BatV=3,
//     data.Node_type="LHT65N" (observed through the public /events/jsexec feed).
//   - The sensor transitions to online.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log. A decode mismatch at
//     AssertLorawanEventDecoded points at the template codec (ScriptConversion) or
//     the js-executor decode path, not at LoRaWAN connectivity.
package lorawan_journey_codec

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	bootstrap "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"
	gateway "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/lorawan_gateway"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	assetAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/asserts"
	assetPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	templatePayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	eventAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/events/events/asserts"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// uplinkHex is a deterministic Dragino LHT65N fPort-2 real-time frame: battery
// 3.0 V (status Good), built-in SHT temperature 21.00 C, humidity 60.0 %, no
// external sensor. Verified against the device's own decodeUplink.
const uplinkHex = "CBB8083402580000000000"

// wantDecoded is the exact output the LHT65N codec must produce for uplinkHex.
var wantDecoded = map[string]any{
	"Node_type":  "LHT65N",
	"TempC_SHT":  21.0,
	"Hum_SHT":    60.0,
	"BatV":       3.0,
	"Bat_status": "Good",
}

// Items provisions the route group + decoded (real-codec) template, connects a UDP
// gateway (shared block), then provisions an ABP sensor, fires the LHT65N uplink,
// and asserts the codec decoded it into the device's semantic fields.
func Items() []saga.Item {
	items := []saga.Item{
		rgSteps.CreateRouteGroup(),
		templateSteps.CreateTemplateWith(templatePayloads.SagaLorawanDecodedTemplate),
	}
	items = append(items, gateway.ConnectUDPItems("gwUdp")...)
	items = append(items,
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanSensorABPFor("sn"), "sn"),
		assetSteps.ActivateLorawanSensorABP("gwUdp", "sn"),
		assetSteps.FireLorawanUplink("sn", uplinkHex),
		eventAsserts.AssertLorawanEventDecodedByLabel("sn", wantDecoded),
		assetAsserts.AssertHealthStatusByLabel("sn", "online"),
	)
	return items
}

// Run fronts the shared IAM bootstrap (common/journey/iam_bootstrap) and runs this
// journey's items under one rollback chain. The suite runner (journey/suite)
// provisions the stack once via infra.EnsureAll; Run never touches the environment.
func Run(t *testing.T) {
	t.Helper()
	runID := random.NewRunID()
	clients := bootstrap.NewClients()
	items := append(bootstrap.BootstrapItems(), Items()...)
	saga.Run(t, context.Background(), runID, clients, items...)
}
