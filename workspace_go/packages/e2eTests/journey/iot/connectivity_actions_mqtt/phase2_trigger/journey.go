// Package phase2_trigger exercises the healthmonitor → route group
// → trigger execution wiring end-to-end for MQTT-protocol assets.
// Mirrors phase1_workflow but swaps kind=workflow for kind=trigger,
// so the saga validates the second router kind permitted on the
// HealthMonitor surface.
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
	rgPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/payloads"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
	triggerAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/asserts"
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// Items is the ordered slice of saga Items the phase runs.
func Items() []saga.Item {
	return []saga.Item{
		triggerSteps.StartTestSink(),
		triggerSteps.CreateTrigger(),

		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, onlineTriggerRG()),
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOfflineRouteGroupID, offlineTriggerRG()),

		templateSteps.CreateTemplate(),
		assetSteps.CreateConnectivityAsset(mqttConnectivityAsset()),

		// Warm-up: first-ever CONNECT is unknown→online and is silent
		// by design — no alert publish, no trigger fire. The journey
		// disconnects right after so the next CONNECT is a real
		// offline→online transition.
		assetSteps.ConnectMqttPassword(),
		assetAsserts.AssertHealthStatusEventually("online"),

		// Real transition 1: online→offline fires the offline RG → trigger.
		assetSteps.DisconnectMqtt(),
		assetAsserts.AssertHealthStatusEventually("offline"),
		triggerAsserts.AssertSinkHitEventually(1),

		// Real transition 2: offline→online fires the online RG → trigger.
		assetSteps.ConnectMqttPassword(),
		assetAsserts.AssertHealthStatusEventually("online"),
		triggerAsserts.AssertSinkHitEventually(2),

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

func mqttConnectivityAsset() assetSteps.ConnectivityPayloadFn {
	return func(c *saga.Context) *assetPayloads.AssetCreateBuilder {
		templateID := c.MustGetString(templateSteps.BagKeyTemplateID)
		onlineRG := c.MustGetString(rgSteps.BagKeyOnlineRouteGroupID)
		offlineRG := c.MustGetString(rgSteps.BagKeyOfflineRouteGroupID)
		return assetPayloads.SagaMqttConnectivitySensor(c.RunID, templateID, onlineRG, offlineRG)
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
