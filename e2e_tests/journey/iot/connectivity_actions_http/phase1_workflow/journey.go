// Package phase1_workflow exercises the healthmonitor → route group
// → workflow execution wiring end-to-end for HTTP-protocol assets.
// The phase provisions an HTTP DataSource (push-mode + apiKey auth)
// so the saga can issue an explicit heartbeat to drive online, then
// hits the assets internal force-offline endpoint to drive offline
// within the run budget (the scheduled scan path would force the
// saga to wait the configured scan interval, default 600s).
//
// Outcome on PASS:
//   - DataSource provisioned (id + apiKey on bag).
//   - Workflow definition + instance + two RGs (online / offline).
//   - Asset created in HTTP protocol with HealthMonitor explicit mode
//     and the two RGs bound.
//   - POST /api/v1/heartbeat → healthStatus=online → workflow execution
//     surfaces on /api/v1/events/workflow.
//   - POST /internal/health-monitor/:uuid/force-offline → healthStatus=
//     offline → second workflow execution surfaces (scoped by the
//     force-offline timestamp).
package phase1_workflow

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
	eventAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/events/events/asserts"
	dsSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/steps"
	rgPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/payloads"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
	wfDefSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/workflow/definitions/steps"
	wfInstSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/workflow/instances/steps"
)

// Items is the ordered slice of saga Items the phase runs. Each line
// carries a single comment above it explaining what that item proves
// or sets up, so readers don't have to chase the step package to know
// why this item is in the chain.
func Items() []saga.Item {
	return []saga.Item{
		// Workflow definition that the online and offline route groups will target.
		wfDefSteps.CreateDefinition(),

		// Instance of the workflow definition; route groups bind to the instance id.
		wfInstSteps.CreateInstance(),

		// Route group kind=workflow that matches asset.health "online" transitions.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, onlineWorkflowRG()),

		// Route group kind=workflow that matches asset.health "offline" transitions.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOfflineRouteGroupID, offlineWorkflowRG()),

		// HTTP data source (push-mode + apiKey auth) the heartbeat step will authenticate against.
		dsSteps.CreateDataSource(),

		// Asset template with the saga's StandardizedPayload transform script.
		templateSteps.CreateTemplate(),

		// HTTP connectivity asset wired to the two workflow-kind route groups.
		assetSteps.CreateConnectivityAsset(httpConnectivityAsset()),

		// Warm-up heartbeat: first-ever observation is unknown→online and is silent by design (no workflow fires).
		assetSteps.SendHttpHeartbeat(),

		// Confirm the warm-up settled to "online" before forcing transitions.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Real transition 1: admin force-offline drives online→offline; offline RG matches → workflow fires.
		assetSteps.ForceOfflineByAdmin("saga-http-phase1-workflow-warmup"),

		// Confirm the healthmonitor saw the transition.
		assetAsserts.AssertHealthStatusEventually("offline"),

		// Smoke oracle: events service surfaces a workflow execution scoped after the force-offline timestamp.
		eventAsserts.AssertWorkflowEventReceivedAfter(assetSteps.BagKeyForceOfflineSentAt),

		// Real transition 2: a fresh heartbeat drives offline→online; online RG matches → workflow fires again.
		assetSteps.SendHttpHeartbeat(),

		// Confirm the asset is online again.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Smoke oracle: events service surfaces a second workflow execution scoped after the heartbeat timestamp.
		eventAsserts.AssertWorkflowEventReceivedAfter(assetSteps.BagKeyHeartbeatSentAt),

		// Tear down the asset explicitly so the Compensate chain can verify cascade cleanup.
		assetSteps.DeleteAsset(),
	}
}

func onlineWorkflowRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		instanceID := c.MustGetString(wfInstSteps.BagKeyInstanceID)
		return rgPayloads.SagaWorkflowRouteGroup(c.RunID, "online", instanceID)
	}
}

func offlineWorkflowRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		instanceID := c.MustGetString(wfInstSteps.BagKeyInstanceID)
		return rgPayloads.SagaWorkflowRouteGroup(c.RunID, "offline", instanceID)
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
