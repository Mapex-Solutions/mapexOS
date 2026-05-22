// Package phase2_trigger exercises the healthmonitor → route group
// → trigger execution wiring end-to-end for HTTP-protocol assets.
// Mirrors phase1_workflow but swaps kind=workflow for kind=trigger,
// so the saga validates the second router kind permitted on the
// HealthMonitor surface for HTTP assets.
package phase2_trigger

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

// Items is the ordered slice of saga Items the phase runs. Each line
// carries a single comment above it explaining what that item proves
// or sets up, so readers don't have to chase the step package to know
// why this item is in the chain.
func Items() []saga.Item {
	return []saga.Item{
		// Boot an in-process HTTP sink server; captures every trigger fire as a sink hit.
		triggerSteps.StartTestSink(),

		// Create a trigger pointing at the in-process sink so executions are observable on the sink.
		triggerSteps.CreateTrigger(),

		// Route group kind=trigger that matches asset.health "online" transitions.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, onlineTriggerRG()),

		// Route group kind=trigger that matches asset.health "offline" transitions.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOfflineRouteGroupID, offlineTriggerRG()),

		// HTTP data source (push-mode + apiKey auth) the heartbeat step will authenticate against.
		dsSteps.CreateDataSource(),

		// Asset template with the saga's StandardizedPayload transform script.
		templateSteps.CreateTemplate(),

		// HTTP connectivity asset wired to the two trigger-kind route groups.
		assetSteps.CreateConnectivityAsset(httpConnectivityAsset()),

		// Warm-up heartbeat: first-ever observation is unknown→online and is silent by design (no trigger fires).
		assetSteps.SendHttpHeartbeat(),

		// Confirm the warm-up settled to "online" before forcing transitions.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Real transition 1: admin force-offline drives online→offline; offline RG matches → trigger fires.
		assetSteps.ForceOfflineByAdmin("saga-http-phase2-trigger-warmup"),

		// Confirm the healthmonitor saw the transition.
		assetAsserts.AssertHealthStatusEventually("offline"),

		// Smoke oracle: the sink has captured at least 1 hit from the offline transition.
		triggerAsserts.AssertSinkHitEventually(1),

		// Real transition 2: a fresh heartbeat drives offline→online; online RG matches → trigger fires again.
		assetSteps.SendHttpHeartbeat(),

		// Confirm the asset is online again.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Smoke oracle: the sink has now captured at least 2 hits (one per transition).
		triggerAsserts.AssertSinkHitEventually(2),

		// Tear down the asset explicitly so the Compensate chain can verify cascade cleanup.
		assetSteps.DeleteAsset(),
	}
}

func onlineTriggerRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
		return rgPayloads.SagaTriggerRouteGroup(c.RunID, "online", triggerID)
	}
}

func offlineTriggerRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
		return rgPayloads.SagaTriggerRouteGroup(c.RunID, "offline", triggerID)
	}
}

func httpConnectivityAsset() assetSteps.ConnectivityPayloadFn {
	return func(c *saga.Context) *assetPayloads.AssetCreateBuilder {
		templateID := c.MustGetString(templateSteps.BagKeyTemplateID)
		onlineRG := c.MustGetString(rgSteps.BagKeyOnlineRouteGroupID)
		offlineRG := c.MustGetString(rgSteps.BagKeyOfflineRouteGroupID)
		return assetPayloads.SagaHttpConnectivitySensor(c.RunID, templateID, onlineRG, offlineRG)
	}
}

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
