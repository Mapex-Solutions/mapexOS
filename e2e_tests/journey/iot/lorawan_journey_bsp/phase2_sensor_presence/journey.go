// Package phase2_sensor_presence exercises an ABP LoRaWAN sensor asset's
// connectivity state end to end against the live stack, over BOTH radio transports
// (Semtech UDP and Basics Station key mode). An ABP sensor (fixed session, no OTAA
// join) comes online from real data (an uplink is a presence signal) and is forced
// offline via the internal ops endpoint. The gateway is stood up by the shared
// common/journey/lorawan_gateway building block.
//
// Outcome on PASS (per transport block):
//   - A route group + LoRaWAN codec template + a gateway (to carry the uplink) +
//     an ABP sensor asset are provisioned.
//   - The gateway connects and the ABP sensor activates (no join); a fired uplink
//     drives the sensor asset online.
//   - ForceOfflineByAdmin flips it back to offline.
//   - Every asset (labelled per transport) is torn down in compensation.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package phase2_sensor_presence

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	bootstrap "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"
	gateway "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/lorawan_gateway"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	assetAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/asserts"
	assetPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	templatePayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// uplinkHex is the raw application payload the sensor fires to come online.
const uplinkHex = "0BB809F6025D0000000000"

// Items provisions the shared route group + LoRaWAN template, then for each
// transport connects a gateway (shared block) and runs the ABP sensor online-by-
// data then forces it offline.
func Items() []saga.Item {
	items := []saga.Item{
		rgSteps.CreateRouteGroup(),
		templateSteps.CreateTemplateWith(templatePayloads.SagaLorawanSensorTemplate),
	}
	// UDP transport.
	items = append(items, gateway.ConnectUDPItems("gwUdp")...)
	items = append(items,
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanSensorABPFor("snUdp"), "snUdp"),
		assetSteps.ActivateLorawanSensorABP("gwUdp", "snUdp"),
		assetSteps.FireLorawanUplink("snUdp", uplinkHex),
		assetAsserts.AssertHealthStatusByLabel("snUdp", "online"),
		assetSteps.ForceOfflineByAdminByLabel("snUdp", "saga-lorawan-abp-offline"),
		assetAsserts.AssertHealthStatusByLabel("snUdp", "offline"),
	)
	// Basics Station transport.
	items = append(items, gateway.ConnectBasicStationItems("gwBs")...)
	items = append(items,
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanSensorABPFor("snBs"), "snBs"),
		assetSteps.ActivateLorawanSensorABP("gwBs", "snBs"),
		assetSteps.FireLorawanUplink("snBs", uplinkHex),
		assetAsserts.AssertHealthStatusByLabel("snBs", "online"),
		assetSteps.ForceOfflineByAdminByLabel("snBs", "saga-lorawan-abp-offline"),
		assetAsserts.AssertHealthStatusByLabel("snBs", "offline"),
	)
	return items
}

// Run fronts the shared IAM bootstrap (common/journey/iam_bootstrap) and runs
// this phase's items under one rollback chain.
func Run(t *testing.T) {
	t.Helper()
	if err := utils.SetupE2EEnvironment(); err != nil {
		t.Fatalf("setup e2e environment: %v", err)
	}
	runID := random.NewRunID()
	clients := bootstrap.NewClients()
	items := append(bootstrap.BootstrapItems(), Items()...)
	saga.Run(t, context.Background(), runID, clients, items...)
}
