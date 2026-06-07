// Package gateway_provisioning exercises provisioning a LoRaWAN GATEWAY as an
// asset end-to-end against the live stack. A gateway is radio infrastructure
// (kind=gateway, authMode=eui), never linked to a device, and carries no device
// keys.
//
// Outcome on PASS:
//   - A gateway asset is created (protocol=lorawan, kind=gateway, frequency plan).
//   - The internal asset-auth endpoint (the LNS Gateway Server's L3 source)
//     serves it as type=lorawan + lorawan.kind=gateway + the gateway block, in
//     the exact shape the LNS reads on connect.
//   - The asset is torn down (compensation), leaving the stack clean.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package gateway_provisioning

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
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// Items is the ordered slice of saga Items the journey runs.
//
//	1. CreateRouteGroup               -> route group the asset requires
//	2. CreateTemplate                 -> asset template the asset requires
//	3. CreateAssetWith (gateway)      -> lorawan gateway asset persisted (no keys, KEK skipped)
//	4. AssertAuthProjectionGateway    -> L3 source serves type=lorawan + kind=gateway + plan
//	5. DeleteAsset                    -> teardown
func Items() []saga.Item {
	return []saga.Item{
		// Route group bound to the asset (assets require at least one).
		rgSteps.CreateRouteGroup(),

		// Asset template the asset references.
		templateSteps.CreateTemplate(),

		// Gateway asset: protocol=lorawan, kind=gateway, authMode=eui, frequency plan.
		assetSteps.CreateAssetWith(assetPayloads.SagaLorawanGateway),

		// The LNS Gateway Server's L3 source publishes it as a gateway in the right shape.
		assetAsserts.AssertAuthProjectionGateway(assetPayloads.SagaGatewayFrequencyPlan),

		// Tear the asset down explicitly.
		assetSteps.DeleteAsset(),
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
