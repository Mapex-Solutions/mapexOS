// Package lorawan_gateway_presence exercises a LoRaWAN gateway asset's connectivity
// state end to end against the live stack, over BOTH radio transports (Semtech UDP
// and Basics Station key mode).
//
// A LoRaWAN gateway asset requires neither an asset template nor a route group
// (only end-device assets do), so this journey provisions the gateway alone.
//
// Outcome on PASS (per transport block):
//   - A gateway asset is provisioned and served in the LNS-facing L3 shape.
//   - The simulated gateway connects to mapexLNS; its keepalive/stats make the LNS
//     publish a "connect" presence advisory and the asset transitions to online.
//   - ForceOfflineByAdmin (internal ops endpoint) flips it back to offline.
//   - The gateway asset (labelled per transport) is torn down in compensation.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package lorawan_gateway_presence

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
)

// Items runs a UDP block and a Basics Station block, each provisioning a labelled
// gateway (no template / route group), connecting it to the LNS to go online, then
// forcing it offline. No shared route group or template — the gateway needs none.
func Items() []saga.Item {
	return []saga.Item{
		// ── UDP transport (eui-mode gateway) ──
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanGatewayFor("gwUdp"), "gwUdp"),
		assetAsserts.AssertAuthProjectionGatewayByLabel("gwUdp", assetPayloads.SagaGatewayFrequencyPlan),
		assetSteps.ConnectLorawanGatewayUDP("gwUdp"),
		assetAsserts.AssertHealthStatusByLabel("gwUdp", "online"),
		assetSteps.ForceOfflineByAdminByLabel("gwUdp", "saga-lorawan-gw-offline"),
		assetAsserts.AssertHealthStatusByLabel("gwUdp", "offline"),

		// ── Basics Station transport (key-mode gateway, bearer token) ──
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanGatewayKeyFor("gwBs"), "gwBs"),
		assetAsserts.AssertAuthProjectionGatewayByLabel("gwBs", assetPayloads.SagaGatewayFrequencyPlan),
		assetSteps.ConnectLorawanGatewayBasicStation("gwBs"),
		assetAsserts.AssertHealthStatusByLabel("gwBs", "online"),
		assetSteps.ForceOfflineByAdminByLabel("gwBs", "saga-lorawan-gw-offline"),
		assetAsserts.AssertHealthStatusByLabel("gwBs", "offline"),
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
