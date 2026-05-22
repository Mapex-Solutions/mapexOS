// Package phase1_connectivity exercises the HTTP trigger smoke path
// via the connectivity pipeline: the healthmonitor flips the saga
// asset between online and offline, the router matches each
// transition against trigger-kind route groups, and the triggers
// service POSTs to an in-process HTTP sink the journey stood up.
//
// Outcome on PASS:
//   - Phase 0 outcomes hold (IAM bootstrap succeeds).
//   - The HTTP sink observes one POST per real health transition.
//
// Outcome on FAIL:
//   - "sink hits: want >=1, got 0" → trigger never matched the route
//     group (router config) or executor can't reach the sink
//     (host/port resolution).
package phase1_connectivity

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
	dsSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/steps"
	rgPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/payloads"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
	triggerAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/asserts"
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// Items is the ordered slice of saga Items the phase runs.
func Items() []saga.Item {
	return []saga.Item{
		// Boot an in-process HTTP server on TriggerSinkBindAddr; counts POSTs into BagKeyTriggerSinkHits.
		triggerSteps.StartTestSink(),

		// Create an HTTP-kind trigger pointing at the sink URL (overrides endpoint to TriggerSinkURL).
		triggerSteps.CreateTrigger(),

		// Route group that matches asset.health "online" transitions and points at the HTTP trigger.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, onlineTriggerRG()),

		// Route group that matches asset.health "offline" transitions and points at the HTTP trigger.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOfflineRouteGroupID, offlineTriggerRG()),

		// HTTP data source the heartbeat step will authenticate against.
		dsSteps.CreateDataSource(),

		// Asset template with the saga's StandardizedPayload transform script.
		templateSteps.CreateTemplate(),

		// HTTP connectivity asset wired to the two trigger-kind route groups.
		assetSteps.CreateConnectivityAsset(httpConnectivityAsset()),

		// Warm-up heartbeat: first observation is unknown→online and is silent by design (no trigger).
		assetSteps.SendHttpHeartbeat(),

		// Confirm the warm-up settled to "online" before forcing transitions.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Real transition 1: admin forces the asset offline → offline RG matches → HTTP trigger fires.
		assetSteps.ForceOfflineByAdmin("saga-http-phase1-warmup"),

		// Confirm the healthmonitor saw the transition.
		assetAsserts.AssertHealthStatusEventually("offline"),

		// Smoke oracle: the HTTP sink has captured at least 1 POST from the offline transition.
		triggerAsserts.AssertSinkHitEventually(1),

		// Real transition 2: a fresh heartbeat brings the asset back online → online RG fires the trigger again.
		assetSteps.SendHttpHeartbeat(),

		// Confirm the asset is online again.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Smoke oracle: the HTTP sink has now captured at least 2 POSTs.
		triggerAsserts.AssertSinkHitEventually(2),

		// Tear down the asset explicitly so the Compensate chain can verify cascade cleanup.
		assetSteps.DeleteAsset(),
	}
}

// onlineTriggerRG builds a kind=trigger route group keyed to the
// "online" health transition, pointing at the trigger id Phase 1
// published on the bag.
func onlineTriggerRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
		return rgPayloads.SagaTriggerRouteGroup(c.RunID, "online", triggerID)
	}
}

// offlineTriggerRG builds the offline counterpart of the online
// route group.
func offlineTriggerRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
		return rgPayloads.SagaTriggerRouteGroup(c.RunID, "offline", triggerID)
	}
}

// httpConnectivityAsset binds the saga's HTTP connectivity asset to
// the online + offline trigger-kind route groups.
func httpConnectivityAsset() assetSteps.ConnectivityPayloadFn {
	return func(c *saga.Context) *assetPayloads.AssetCreateBuilder {
		templateID := c.MustGetString(templateSteps.BagKeyTemplateID)
		onlineRG := c.MustGetString(rgSteps.BagKeyOnlineRouteGroupID)
		offlineRG := c.MustGetString(rgSteps.BagKeyOfflineRouteGroupID)
		return assetPayloads.SagaHttpConnectivitySensor(c.RunID, templateID, onlineRG, offlineRG)
	}
}

// Run wires Phase 0 (IAM bootstrap) in front of this phase's items
// and executes the resulting saga under one rollback chain.
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
