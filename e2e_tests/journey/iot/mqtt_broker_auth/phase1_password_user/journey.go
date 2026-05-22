// Package phase1_password_user exercises the MQTT password auth path
// end-to-end against the live broker. The phase covers the seven-step
// flow the operator runs by hand when validating an asset: create the
// asset (password mode) → CONNECT → assert health=online (presence
// flowed) → publish telemetry → assert the events service surfaced
// the row → delete the asset → re-CONNECT must be denied (proves the
// FANOUT-driven L1 invalidation in the broker plugin).
//
// Outcome on PASS:
//   - Asset is created with protocol=mqtt + authType=password and the
//     plaintext password the saga supplied lands on the bag.
//   - MQTT CONNECT (username=assetUUID, password=plaintext) succeeds.
//   - Asset's healthStatus flips to "online" within the polling
//     window (presence advisory consumed, status persisted, read
//     model surfaces it via GET /api/v1/assets/:id).
//   - PUBLISH to events/{assetUUID}/temperature is accepted by the
//     broker ACL (bare-assetUUID topic shape).
//   - The events service exposes the row via /api/v1/events/raw
//     within the polling window.
//   - DELETE /api/v1/assets/{id} succeeds and propagates the FANOUT
//     invalidation to the broker plugin.
//   - A fresh CONNECT with the same credentials is rejected by the
//     broker.
//
// Outcome on FAIL:
//   - Saga step / assert name in the log identifies the offending
//     stage (e.g. assets/assets.ConnectMqttPassword,
//     events/events.AssertRawEventReceived) so the operator can chase
//     it in the right service log without re-running.
package phase1_password_user

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	phase0 "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/mqtt_broker_auth/phase0_iam_bootstrap"

	assetAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/asserts"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	eventAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/events/events/asserts"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// Items is the ordered slice of saga Items the journey runs. The
// MQTT phase requires a route group + template + asset in place;
// it provisions them inline rather than importing the
// mqtt_full_pipeline setup so this journey stays independent.
//
// Step / Assert ordering — keep aligned with the journey README:
//
//  1. CreateRouteGroup                → route group for the asset
//  2. CreateTemplate                  → asset template
//  3. CreateAsset                     → asset persisted (password)
//  4. ConnectMqttPassword             → broker accepts the device
//  5. AssertHealthStatusEventually    → presence flowed end-to-end
//  6. PublishTelemetry                → publish accepted by ACL
//  7. AssertRawEventReceivedAfter     → events service has the row
//  8. DisconnectMqtt                  → clean MQTT teardown
//  9. DeleteAsset                     → fanout invalidation fires
//  10. AssertConnectDeniedPassword    → cache eviction proved
func Items() []saga.Item {
	return []saga.Item{
		// Route group bound to the asset (telemetry path).
		rgSteps.CreateRouteGroup(),

		// Asset template (temperature schema).
		templateSteps.CreateTemplate(),

		// Asset persisted with protocol=mqtt + authType=password; plaintext password lands on the bag.
		assetSteps.CreateAsset(),

		// MQTT CONNECT against the password listener (username=assetUUID, password=plaintext).
		assetSteps.ConnectMqttPassword(),

		// Presence advisory consumed end-to-end; read model surfaces healthStatus=online.
		assetAsserts.AssertHealthStatusEventually("online"),

		// PUBLISH on events/{assetUUID}/temperature; broker ACL must accept the bare-assetUUID topic.
		assetSteps.PublishTelemetry(),

		// Events service surfaces the row via /api/v1/events/raw scoped after the publish timestamp.
		eventAsserts.AssertRawEventReceivedAfter(assetSteps.BagKeyTelemetrySentAt),

		// Clean MQTT teardown (no presence side-effect needed for the deny probe below).
		assetSteps.DisconnectMqtt(),

		// DELETE the asset; FANOUT invalidation must reach the broker plugin L1 cache.
		assetSteps.DeleteAsset(),

		// Fresh CONNECT with the same credentials must be rejected — proves the cache eviction fired.
		assetAsserts.AssertConnectDeniedPassword(),
	}
}

// Run executes phase 0 (IAM bootstrap) + this phase as a single saga
// so the rollback chain unwinds in reverse: deny-probe is a no-op
// rollback, DeleteAsset is a no-op rollback (already deleted),
// CreateAsset's compensate is short-circuited by BagKeyAssetDeleted,
// then template + route group cleanups run as usual.
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
