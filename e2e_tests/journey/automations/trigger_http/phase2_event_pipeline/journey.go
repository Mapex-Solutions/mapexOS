// Package phase2_event_pipeline exercises the HTTP trigger via the
// telemetry path: POST /api/v1/events → gateway → js-executor →
// router → trigger → HTTP sink. Validates the part of the stack
// phase1_connectivity does NOT touch — namely the js-executor's
// template script transform sitting between the gateway and router.
package phase2_event_pipeline

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	phase0 "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/mqtt_broker_auth/phase0_iam_bootstrap"

	assetPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	eventAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/events/events/asserts"
	dsSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/steps"
	rgPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/payloads"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
	triggerAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/asserts"
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// Items is the ordered slice of saga Items the phase runs.
func Items() []saga.Item {
	return []saga.Item{
		// Boot the in-process HTTP sink the trigger executor will POST to.
		triggerSteps.StartTestSink(),

		// Create an HTTP-kind trigger pointing at the sink.
		triggerSteps.CreateTrigger(),

		// Single trigger-kind route group; reused as online/offline/general so telemetry fires it.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, eventTriggerRG()),
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOfflineRouteGroupID, eventTriggerRG()),

		// HTTP data source — its API key authenticates the POST /events call.
		dsSteps.CreateDataSource(),

		// Asset template with the StandardizedPayload transform script (js-executor runs this).
		templateSteps.CreateTemplate(),

		// HTTP connectivity asset wired to the same trigger RG for both health flows.
		assetSteps.CreateConnectivityAsset(httpEventAsset()),

		// POST a telemetry event to the gateway with the saga's runID embedded in the body.
		dsSteps.PostRawEvent(),

		// Smoke: events_trigger records ≥ 1 successful execution from the telemetry path.
		eventAsserts.AssertTriggerExecutedSuccessfullyEventually(1),

		// Content-key: the HTTP sink received the POST the trigger emitted.
		triggerAsserts.AssertSinkHitEventually(1),

		// Explicit asset teardown so Compensate verifies cascade cleanup.
		assetSteps.DeleteAsset(),
	}
}

// eventTriggerRG builds a single kind=trigger route group reused for
// online, offline, and general telemetry routing. Match conditions
// are intentionally empty — every event the asset emits fires this
// route group and therefore the trigger.
func eventTriggerRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
		return rgPayloads.SagaTriggerRouteGroup(c.RunID, "event", triggerID)
	}
}

// httpEventAsset binds the connectivity asset to the same trigger RG
// for both health flows; this leaves the RouteGroupIds and the
// HealthMonitor's Online/Offline lists pointing at the trigger so the
// telemetry path fires it.
func httpEventAsset() assetSteps.ConnectivityPayloadFn {
	return func(c *saga.Context) *assetPayloads.AssetCreateBuilder {
		templateID := c.MustGetString(templateSteps.BagKeyTemplateID)
		onlineRG := c.MustGetString(rgSteps.BagKeyOnlineRouteGroupID)
		offlineRG := c.MustGetString(rgSteps.BagKeyOfflineRouteGroupID)
		return assetPayloads.SagaHttpConnectivitySensor(c.RunID, templateID, onlineRG, offlineRG)
	}
}

// Run wires Phase 0 (IAM bootstrap) in front of this phase's items.
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
