// Package gateway_cert_provisioning exercises provisioning a cert-mode LoRaWAN
// GATEWAY as an asset end-to-end against the live stack, then issuing its mTLS
// client certificate and proving the issued serial is reflected on the LNS
// Gateway Server's L3 source. This guards the MapexOS server side of the gateway
// mTLS flow; the actual Basics Station handshake is exercised separately by the
// device simulator.
//
// Outcome on PASS:
//   - A cert-mode gateway asset is created (protocol=lorawan, kind=gateway,
//     authMode=cert, frequency plan, certTTL).
//   - POST /api/v1/gateway_certs returns a signed cert + key + CA chain
//     (CN=assetUUID); the assets MS persists currentCert and rewrites the
//     mapex-asset-auth projection synchronously.
//   - The internal asset-auth endpoint serves type=lorawan + kind=gateway +
//     gateway.authMode=cert + frequencyPlanId + currentCertSerial equal to the
//     issued serial, the exact shape the LNS pins on connect.
//   - The asset is torn down (compensation), leaving the stack clean.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log (e.g.
//     assets/assets.IssueGatewayCert, assets/assets.AssertAuthProjectionGatewayCert).
package gateway_cert_provisioning

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
//  1. CreateRouteGroup                  -> route group the asset requires
//  2. CreateTemplate                    -> asset template the asset requires
//  3. CreateAssetWith (gateway cert)    -> lorawan gateway asset persisted (authMode=cert)
//  4. IssueGatewayCert                  -> signed cert + serial on the bag, currentCert persisted
//  5. AssertAuthProjectionGatewayCert   -> L3 source reflects authMode=cert + the issued serial
//  6. DeleteAsset                       -> teardown
func Items() []saga.Item {
	return []saga.Item{
		// Route group bound to the asset (assets require at least one).
		rgSteps.CreateRouteGroup(),

		// Asset template the asset references.
		templateSteps.CreateTemplate(),

		// Gateway asset: protocol=lorawan, kind=gateway, authMode=cert, frequency plan.
		assetSteps.CreateAssetWith(assetPayloads.SagaLorawanGatewayCert),

		// Issue the gateway's mTLS cert; the serial lands on the bag.
		assetSteps.IssueGatewayCert(),

		// The L3 source reflects the cert: authMode=cert + currentCertSerial = issued serial.
		assetAsserts.AssertAuthProjectionGatewayCert(assetPayloads.SagaGatewayFrequencyPlan),

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
