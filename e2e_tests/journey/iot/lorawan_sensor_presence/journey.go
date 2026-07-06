// Package lorawan_sensor_presence exercises a LoRaWAN sensor asset's connectivity
// state end to end against the live stack, over BOTH radio transports (Semtech UDP
// and Basics Station key mode). A sensor comes online from real data (an uplink is
// a presence signal) and is forced offline via the internal ops endpoint.
//
// Outcome on PASS (per transport block):
//   - A route group + LoRaWAN codec template + a gateway (to carry the uplink) + a
//     sensor asset are provisioned.
//   - The simulated gateway connects and the sensor OTAA-joins; a fired uplink
//     drives the sensor asset online.
//   - ForceOfflineByAdmin flips it back to offline.
//   - Every asset (labelled per transport) is torn down in compensation.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package lorawan_sensor_presence

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
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// uplinkHex is the raw application payload the sensor fires to come online.
const uplinkHex = "0BB809F6025D0000000000"

// Items runs a shared route group + LoRaWAN template, then a UDP block and a Basics
// Station block, each provisioning a labelled gateway + sensor, firing an uplink to
// go online, then forcing the sensor offline.
func Items() []saga.Item {
	return []saga.Item{
		rgSteps.CreateRouteGroup(),
		templateSteps.CreateTemplateWith(templatePayloads.SagaLorawanSensorTemplate),

		// ── UDP transport ──
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanGatewayFor("gwUdp"), "gwUdp"),
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanSensorOTAAFor("snUdp"), "snUdp"),
		assetSteps.ConnectLorawanGatewayUDP("gwUdp"),
		assetSteps.JoinLorawanSensor("gwUdp", "snUdp"),
		assetSteps.FireLorawanUplink("snUdp", uplinkHex),
		assetAsserts.AssertHealthStatusByLabel("snUdp", "online"),
		assetSteps.ForceOfflineByAdminByLabel("snUdp", "saga-lorawan-sensor-offline"),
		assetAsserts.AssertHealthStatusByLabel("snUdp", "offline"),

		// ── Basics Station transport (key-mode gateway) ──
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanGatewayKeyFor("gwBs"), "gwBs"),
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanSensorOTAAFor("snBs"), "snBs"),
		assetSteps.ConnectLorawanGatewayBasicStation("gwBs"),
		assetSteps.JoinLorawanSensor("gwBs", "snBs"),
		assetSteps.FireLorawanUplink("snBs", uplinkHex),
		assetAsserts.AssertHealthStatusByLabel("snBs", "online"),
		assetSteps.ForceOfflineByAdminByLabel("snBs", "saga-lorawan-sensor-offline"),
		assetAsserts.AssertHealthStatusByLabel("snBs", "offline"),
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
