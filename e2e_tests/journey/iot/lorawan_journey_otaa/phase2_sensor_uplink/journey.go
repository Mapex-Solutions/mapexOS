// Package phase2_sensor_uplink exercises the full OTAA LoRaWAN ingestion path end
// to end against the live stack: a simulated gateway + OTAA sensor send a real
// uplink, through mapexLNS, into MapexOS — over BOTH radio transports (Semtech UDP
// and Basics Station key mode).
//
// The gateway is stood up by the shared common/journey/lorawan_gateway building
// block; this phase adds the OTAA sensor: create → OTAA-join → fire → assert.
// Each transport uses its own labelled gateway + sensor so bag keys and Compensate
// teardown never collide.
//
// Outcome on PASS (per transport block):
//   - A LoRaWAN gateway asset and an OTAA sensor asset are provisioned; the
//     gateway's L3 auth projection is served in the shape the LNS reads.
//   - The gateway connects to mapexLNS and the sensor OTAA-joins.
//   - A fired uplink is ingested: a raw event appears for the sensor's threadId
//     carrying the undecoded application bytes (== the fired hex) plus
//     fPort/fCnt/rxInfo, and the sensor asset transitions to online.
//   - Every asset (labelled per transport) is torn down in compensation.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log. A payload mismatch
//     at AssertLorawanUplinkIngested points at the LNS session-key (AppSKey)
//     FRMPayload decrypt, not at this phase's wiring.
package phase2_sensor_uplink

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

// uplinkHex is the raw application payload the sensor fires — a representative
// Milesight/Dragino-style frame. The LNS decrypts the LoRaWAN layer, so the raw
// event MapexOS stores carries exactly these bytes (codec decode is downstream).
const uplinkHex = "0BB809F6025D0000000000"

// Items provisions the shared route group + LoRaWAN codec template, then for each
// transport connects a gateway (shared block) and runs the OTAA sensor: create,
// OTAA-join, fire an uplink, assert ingestion + online presence.
func Items() []saga.Item {
	items := []saga.Item{
		rgSteps.CreateRouteGroup(),
		templateSteps.CreateTemplateWith(templatePayloads.SagaLorawanSensorTemplate),
	}
	// UDP transport.
	items = append(items, gateway.ConnectUDPItems("gwUdp")...)
	items = append(items,
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanSensorOTAAFor("snUdp"), "snUdp"),
		assetSteps.JoinLorawanSensor("gwUdp", "snUdp"),
		assetSteps.FireLorawanUplink("snUdp", uplinkHex),
		eventAsserts.AssertLorawanUplinkIngestedByLabel("snUdp", uplinkHex),
		assetAsserts.AssertHealthStatusByLabel("snUdp", "online"),
	)
	// Basics Station transport.
	items = append(items, gateway.ConnectBasicStationItems("gwBs")...)
	items = append(items,
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanSensorOTAAFor("snBs"), "snBs"),
		assetSteps.JoinLorawanSensor("gwBs", "snBs"),
		assetSteps.FireLorawanUplink("snBs", uplinkHex),
		eventAsserts.AssertLorawanUplinkIngestedByLabel("snBs", uplinkHex),
		assetAsserts.AssertHealthStatusByLabel("snBs", "online"),
	)
	return items
}

// Run fronts the shared IAM bootstrap (common/journey/iam_bootstrap) and runs
// this phase's items under one rollback chain. The suite runner (journey/suite)
// provisions the stack once via infra.EnsureAll; Run never touches the environment.
func Run(t *testing.T) {
	t.Helper()
	runID := random.NewRunID()
	clients := bootstrap.NewClients()
	items := append(bootstrap.BootstrapItems(), Items()...)
	saga.Run(t, context.Background(), runID, clients, items...)
}
