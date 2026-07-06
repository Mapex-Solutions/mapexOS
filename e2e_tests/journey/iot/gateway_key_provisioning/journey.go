// Package gateway_key_provisioning exercises provisioning a key-mode LoRaWAN
// GATEWAY as an asset end-to-end against the live stack, proving the supplied
// Basics Station token is hashed and reflected on the LNS Gateway Server's L3
// source. This guards the MapexOS server side of the gateway token flow; the
// actual Basics Station handshake is exercised separately by the device
// simulator.
//
// Outcome on PASS:
//   - A key-mode gateway asset is created (protocol=lorawan, kind=gateway,
//     authMode=key, frequency plan, a request-only token).
//   - The assets MS bcrypt-hashes the token and rewrites the mapex-asset-auth
//     projection; the plaintext is never persisted or returned.
//   - The internal asset-auth endpoint serves type=lorawan + kind=gateway +
//     gateway.authMode=key + frequencyPlanId + a non-empty apiKeyHash that
//     bcrypt-matches the registered token, the exact shape the LNS compares on
//     connect.
//   - The asset is torn down (compensation), leaving the stack clean.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log (e.g.
//     assets/assets.CreateAsset, assets/assets.AssertAuthProjectionGatewayKey).
package gateway_key_provisioning

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
//  1. CreateRouteGroup                 -> route group the asset requires
//  2. CreateTemplate                   -> asset template the asset requires
//  3. CreateAssetWith (gateway key)    -> lorawan gateway asset persisted (authMode=key)
//  4. AssertAuthProjectionGatewayKey   -> L3 source reflects authMode=key + apiKeyHash
//  5. DeleteAsset                      -> teardown
func Items() []saga.Item {
	return []saga.Item{
		// Route group bound to the asset (assets require at least one).
		rgSteps.CreateRouteGroup(),

		// Asset template the asset references.
		templateSteps.CreateTemplate(),

		// Gateway asset: protocol=lorawan, kind=gateway, authMode=key, frequency plan.
		assetSteps.CreateAssetWith(assetPayloads.SagaLorawanGatewayKey),

		// The L3 source reflects the key mode: authMode=key + a matching apiKeyHash.
		assetAsserts.AssertAuthProjectionGatewayKey(assetPayloads.SagaGatewayFrequencyPlan),

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
