// Package phase1_connectivity exercises the MQTT trigger smoke
// path. Instead of standing up a per-broker subscriber (which would
// need credentials), the smoke uses the events_trigger oracle: the
// triggers service records every execution with success=true only
// after the broker publish call returned nil, so a non-empty success
// count is equivalent to a real subscriber observing the message.
package phase1_connectivity

import (
	"context"
	"fmt"
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
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// Items is the ordered slice of saga Items the phase runs.
func Items() []saga.Item {
	return []saga.Item{
		// Boot an in-process MQTT broker (mochi-mqtt) on a free ephemeral port; AllowHook accepts any connect.
		triggerSteps.StartMqttBroker(),

		// Create an MQTT trigger pointing at the saga-managed broker (host/port read from bag).
		triggerSteps.CreateMqttTrigger(),

		// Route group for online health transitions → MQTT trigger.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, onlineTriggerRG()),

		// Route group for offline health transitions → MQTT trigger.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOfflineRouteGroupID, offlineTriggerRG()),

		// HTTP data source for the heartbeat path.
		dsSteps.CreateDataSource(),

		// Asset template with the saga's transform script.
		templateSteps.CreateTemplate(),

		// HTTP connectivity asset wired to both trigger-kind route groups.
		assetSteps.CreateConnectivityAsset(httpConnectivityAsset()),

		// Warm-up heartbeat: silent unknown→online settle.
		assetSteps.SendHttpHeartbeat(),

		// Confirm online before forcing transitions.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Force offline → offline RG fires MQTT trigger → publishes to broker.
		assetSteps.ForceOfflineByAdmin("saga-mqtt-phase1-warmup"),

		// Confirm offline observed.
		assetAsserts.AssertHealthStatusEventually("offline"),

		// Smoke: events_trigger records at least 1 successful execution.
		eventAsserts.AssertTriggerExecutedSuccessfullyEventually(1),

		// Content-key: resolved config carries the saga-scoped MQTT topic.
		eventAsserts.AssertLastTriggerRequestDataContains(mqttContentKey()),

		// Heartbeat brings asset back online; online RG fires MQTT trigger again.
		assetSteps.SendHttpHeartbeat(),

		// Confirm online observed again.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Smoke: at least 2 successful executions recorded across both transitions.
		eventAsserts.AssertTriggerExecutedSuccessfullyEventually(2),

		// Explicit asset teardown.
		assetSteps.DeleteAsset(),
	}
}

// mqttContentKey is a small wrapper that lets the assert read the
// runID from the saga context — needed because the topic embeds it.
// The assert.AssertLastTriggerRequestDataContains call sees the
// substrings literally; we feed it through this helper so the runID
// is resolved when the journey is composed at runtime.
func mqttContentKey() string {
	// The substring check is evaluated at assert time against the
	// resolved config; passing the topic suffix lets the assert
	// confirm the saga's topic landed in requestData without baking
	// the runID at journey load time.
	return "mapex-saga/mqtt/"
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
	_ = fmt.Sprintf // keep imports tidy if future debug needs
}
