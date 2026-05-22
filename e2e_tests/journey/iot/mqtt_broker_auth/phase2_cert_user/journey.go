// Package phase2_cert_user exercises the MQTT cert (mTLS) auth path
// end-to-end against the live broker. The flow mirrors the password
// phase but presents an issued device cert on the wire instead of a
// plaintext password.
//
// Outcome on PASS:
//   - Asset is created with protocol=mqtt + authType=cert.
//   - POST /api/v1/mqtt_certs returns a signed cert + private key
//     + CA chain; assets MS persists asset.currentCert + emits the
//     fanout invalidation; the PEM bundle lands on the saga bag.
//   - mTLS CONNECT (cert presented, username=assetUUID) succeeds
//     against the broker's 8883 listener.
//   - Asset's healthStatus flips to "online" within the polling
//     window (presence advisory consumed, status persisted).
//   - PUBLISH to events/{assetUUID}/temperature is accepted by the
//     broker ACL.
//   - The events service exposes the row via /api/v1/events/raw
//     within the polling window.
//   - DELETE /api/v1/assets/{id} succeeds, FANOUT invalidation fires.
//   - A fresh mTLS CONNECT with the SAME cert is rejected by the
//     broker (asset gone → currentCert serial gone → cert-mode auth
//     fails).
//
// Outcome on FAIL:
//   - Saga step / assert name in the log identifies the offending
//     stage (e.g. assets/assets.IssueCert,
//     assets/assets.ConnectMqttCert) so the operator can chase it.
package phase2_cert_user

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	phase0 "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/mqtt_broker_auth/phase0_iam_bootstrap"

	assetAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/asserts"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	eventAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/events/events/asserts"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// Items is the ordered slice of saga Items driving the cert path.
// CreateAssetWith picks the cert payload variant so the asset's
// authType is "cert" + certTTL is set; IssueCert then handles the
// PEM bundle the mTLS connect step consumes.
//
// Step / Assert ordering:
//
//  1. CreateRouteGroup                → route group for the asset
//  2. CreateTemplate                  → asset template
//  3. CreateAssetWith (cert variant)  → asset persisted (authType=cert)
//  4. IssueCert                       → PEM bundle on bag + currentCert in Mongo
//  5. ConnectMqttCert                 → mTLS handshake against :8883
//  6. AssertHealthStatusEventually    → presence flowed end-to-end
//  7. PublishTelemetry                → publish accepted by ACL
//  8. AssertRawEventReceivedAfter     → events service has the row
//  9. DisconnectMqtt                  → clean MQTT teardown
//  10. DeleteAsset                    → fanout invalidation fires
//  11. AssertConnectDeniedCert        → cache eviction proved on mTLS path
func Items() []saga.Item {
	return []saga.Item{
		// Route group bound to the asset (telemetry path).
		rgSteps.CreateRouteGroup(),

		// Asset template (temperature schema).
		templateSteps.CreateTemplate(),

		// Asset persisted with protocol=mqtt + authType=cert + certTTL using the cert payload variant.
		assetSteps.CreateAssetWith(payloads.SagaMqttCertTemperatureSensor),

		// POST /api/v1/mqtt_certs returns the PEM bundle; assets persists currentCert + emits FANOUT invalidation.
		assetSteps.IssueCert(),

		// mTLS CONNECT against :8883 with the freshly issued cert (username=assetUUID).
		assetSteps.ConnectMqttCert(),

		// Presence advisory consumed end-to-end; read model surfaces healthStatus=online.
		assetAsserts.AssertHealthStatusEventually("online"),

		// PUBLISH on events/{assetUUID}/temperature; broker ACL must accept the bare-assetUUID topic.
		assetSteps.PublishTelemetry(),

		// Events service surfaces the row via /api/v1/events/raw scoped after the publish timestamp.
		eventAsserts.AssertRawEventReceivedAfter(assetSteps.BagKeyTelemetrySentAt),

		// Clean MQTT teardown.
		assetSteps.DisconnectMqtt(),

		// DELETE the asset; FANOUT invalidation must reach the broker plugin L1 cache.
		assetSteps.DeleteAsset(),

		// Fresh mTLS CONNECT with the same cert must be rejected — proves the cache eviction fired on the cert path.
		assetAsserts.AssertConnectDeniedCert(),
	}
}

// Run executes phase 0 (IAM bootstrap) + this phase as a single saga.
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
