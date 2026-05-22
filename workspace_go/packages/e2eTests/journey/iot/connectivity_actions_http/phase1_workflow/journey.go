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

// Items is the ordered slice of saga Items the phase runs.
func Items() []saga.Item {
	return []saga.Item{
		wfDefSteps.CreateDefinition(),
		wfInstSteps.CreateInstance(),

		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, onlineWorkflowRG()),
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOfflineRouteGroupID, offlineWorkflowRG()),

		dsSteps.CreateDataSource(),
		templateSteps.CreateTemplate(),
		assetSteps.CreateConnectivityAsset(httpConnectivityAsset()),

		// Warm-up: first-ever heartbeat is unknown→online and is silent
		// by design. The journey forces an offline transition right
		// after so the next heartbeat is a real offline→online cycle.
		assetSteps.SendHttpHeartbeat(),
		assetAsserts.AssertHealthStatusEventually("online"),

		// Real transition 1: online→offline fires the offline RG.
		assetSteps.ForceOfflineByAdmin("saga-http-phase1-workflow-warmup"),
		assetAsserts.AssertHealthStatusEventually("offline"),
		eventAsserts.AssertWorkflowEventReceivedAfter(assetSteps.BagKeyForceOfflineSentAt),

		// Real transition 2: offline→online fires the online RG.
		assetSteps.SendHttpHeartbeat(),
		assetAsserts.AssertHealthStatusEventually("online"),
		eventAsserts.AssertWorkflowEventReceivedAfter(assetSteps.BagKeyHeartbeatSentAt),

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
