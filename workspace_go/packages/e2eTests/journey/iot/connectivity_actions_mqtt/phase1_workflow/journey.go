// Package phase1_workflow exercises the healthmonitor → route group
// → workflow execution wiring end-to-end for MQTT-protocol assets.
// The phase provisions an online + an offline route group of
// kind=workflow pointing at the same workflow instance, then drives
// the asset through CONNECT (→ online → workflow execution recorded)
// and DISCONNECT (→ offline → workflow execution recorded).
//
// Outcome on PASS:
//   - Workflow definition + instance created and observable via the
//     workflow service.
//   - Online and offline route groups created, each carrying one
//     router of kind=workflow targeting the instance.
//   - Asset created with HealthMonitor.OnlineRouteGroupIds + Offline
//     RouteGroupIds bound to the two RGs.
//   - On CONNECT: healthStatus flips to "online" and a workflow
//     execution surfaces on GET /api/v1/events/workflow filtered by
//     instanceId.
//   - On DISCONNECT: healthStatus flips to "offline" and a second
//     workflow execution surfaces (search scoped by the disconnect
//     timestamp so the assert never observes the previous one).
//
// Outcome on FAIL:
//   - Saga step / assert name surfaces in the log so the operator
//     can chase the offending stage in the right service.
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

		templateSteps.CreateTemplate(),
		assetSteps.CreateConnectivityAsset(mqttConnectivityAsset()),

		// Warm-up: first-ever CONNECT moves the asset unknown→online but
		// the healthmonitor treats this as initial activation, not as a
		// transition, and does NOT publish to the route group. We
		// connect, assert sanity, then immediately disconnect to set up
		// the asset for the real transition cycle below.
		assetSteps.ConnectMqttPassword(),
		assetAsserts.AssertHealthStatusEventually("online"),

		// Real transition 1: online→offline fires the offline RG.
		assetSteps.DisconnectMqtt(),
		assetAsserts.AssertHealthStatusEventually("offline"),
		eventAsserts.AssertWorkflowEventReceivedAfter(assetSteps.BagKeyMqttDisconnectedAt),

		// Real transition 2: offline→online fires the online RG.
		assetSteps.ConnectMqttPassword(),
		assetAsserts.AssertHealthStatusEventually("online"),
		eventAsserts.AssertWorkflowEventReceivedAfter(assetSteps.BagKeyMqttConnectedAt),

		assetSteps.DeleteAsset(),
	}
}

// onlineWorkflowRG builds the online-transition route group with one
// router of kind=workflow targeting the instance the saga just
// created. Kept as a closure so the BagKeyInstanceID lookup happens at
// step execution time (after CreateInstance has published it) rather
// than at journey-assembly time.
func onlineWorkflowRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		instanceID := c.MustGetString(wfInstSteps.BagKeyInstanceID)
		return rgPayloads.SagaWorkflowRouteGroup(c.RunID, "online", instanceID)
	}
}

// offlineWorkflowRG mirrors onlineWorkflowRG for the offline transition.
// Uses the same instance — the journey distinguishes online from offline
// executions via the startTime filter on /api/v1/events/workflow.
func offlineWorkflowRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		instanceID := c.MustGetString(wfInstSteps.BagKeyInstanceID)
		return rgPayloads.SagaWorkflowRouteGroup(c.RunID, "offline", instanceID)
	}
}

// mqttConnectivityAsset closes over the bag so the AssetCreate body
// carries the two RG ids created earlier in the same phase.
func mqttConnectivityAsset() assetSteps.ConnectivityPayloadFn {
	return func(c *saga.Context) *assetPayloads.AssetCreateBuilder {
		templateID := c.MustGetString(templateSteps.BagKeyTemplateID)
		onlineRG := c.MustGetString(rgSteps.BagKeyOnlineRouteGroupID)
		offlineRG := c.MustGetString(rgSteps.BagKeyOfflineRouteGroupID)
		return assetPayloads.SagaMqttConnectivitySensor(c.RunID, templateID, onlineRG, offlineRG)
	}
}

// Run executes phase 0 (IAM bootstrap) + this phase as a single saga
// so the rollback chain unwinds in reverse.
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
