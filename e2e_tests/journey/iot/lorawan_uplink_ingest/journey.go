// Package lorawan_uplink_ingest exercises the full LoRaWAN ingestion path end to
// end against the live stack: a simulated gateway + sensor send a real uplink,
// through mapexLNS, into MapexOS — over BOTH radio transports (Semtech UDP and
// Basics Station key mode).
//
// Outcome on PASS (per transport block):
//   - A LoRaWAN gateway asset and a sensor asset are provisioned.
//   - The gateway's L3 auth projection is served in the shape the LNS reads.
//   - The simulated gateway connects to mapexLNS and the sensor OTAA-joins.
//   - A fired uplink is ingested by MapexOS: a raw event appears for the sensor's
//     threadId carrying the undecoded application bytes (== the fired hex) plus
//     fPort/fCnt/rxInfo, and the sensor asset transitions to online.
//   - Every asset (labelled per transport) is torn down in compensation.
//
// Each transport uses its own labelled gateway + sensor (gwUdp/snUdp, gwBs/snBs)
// so their bag keys and Compensate teardown never collide.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package lorawan_uplink_ingest

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	phase0 "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/mqtt_broker_auth/phase0_iam_bootstrap"

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
// event MapexOS stores carries exactly these bytes (the codec decode is downstream).
const uplinkHex = "0BB809F6025D0000000000"

// Items is the ordered saga: a shared route group + LoRaWAN codec template, then a
// UDP block and a Basics Station block, each provisioning a labelled gateway +
// sensor, connecting, joining, firing, and asserting ingestion + online presence.
func Items() []saga.Item {
	return []saga.Item{
		rgSteps.CreateRouteGroup(),
		templateSteps.CreateTemplateWith(templatePayloads.SagaLorawanSensorTemplate),

		// ── UDP transport (eui-mode gateway) ──
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanGatewayFor("gwUdp"), "gwUdp"),
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanSensorOTAAFor("snUdp"), "snUdp"),
		assetAsserts.AssertAuthProjectionGatewayByLabel("gwUdp", assetPayloads.SagaGatewayFrequencyPlan),
		assetSteps.ConnectLorawanGatewayUDP("gwUdp"),
		assetSteps.JoinLorawanSensor("gwUdp", "snUdp"),
		assetSteps.FireLorawanUplink("snUdp", uplinkHex),
		eventAsserts.AssertLorawanUplinkIngestedByLabel("snUdp", uplinkHex),
		assetAsserts.AssertHealthStatusByLabel("snUdp", "online"),

		// ── Basics Station transport (key-mode gateway, bearer token) ──
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanGatewayKeyFor("gwBs"), "gwBs"),
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanSensorOTAAFor("snBs"), "snBs"),
		assetAsserts.AssertAuthProjectionGatewayByLabel("gwBs", assetPayloads.SagaGatewayFrequencyPlan),
		assetSteps.ConnectLorawanGatewayBasicStation("gwBs"),
		assetSteps.JoinLorawanSensor("gwBs", "snBs"),
		assetSteps.FireLorawanUplink("snBs", uplinkHex),
		eventAsserts.AssertLorawanUplinkIngestedByLabel("snBs", uplinkHex),
		assetAsserts.AssertHealthStatusByLabel("snBs", "online"),
	}
}

// Run executes phase 0 (IAM bootstrap) + this journey as a single saga.
func Run(t *testing.T) {
	t.Helper()
	if err := utils.SetupE2EEnvironment(); err != nil {
		t.Fatalf("setup e2e environment: %v", err)
	}
	runID := random.NewRunID()
	clients := phase0.NewClients()
	items := append(phase0.BootstrapItems(), Items()...)
	saga.Run(t, context.Background(), runID, clients, items...)
}
