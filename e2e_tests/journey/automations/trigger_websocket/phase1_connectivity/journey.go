// Package phase1_connectivity exercises the WebSocket trigger
// smoke path via the connectivity pipeline. Uses a minimal
// in-process WS sink (gorilla/websocket) plus the events_trigger
// oracle to validate the publish path.
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
	eventAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/events/events/asserts"
	dsSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/steps"
	rgPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/payloads"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// Items is the ordered slice of saga Items the phase runs.
func Items() []saga.Item {
	return []saga.Item{
		// Boot an in-process WS server on WsSinkBindAddr; /ws is the upgrade endpoint.
		triggerSteps.StartWebsocketSink(),

		// Create a WebSocket trigger pointing at the sink's /ws endpoint.
		triggerSteps.CreateWebsocketTrigger(),

		// Route group for online health transitions → WS trigger.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, onlineTriggerRG()),

		// Route group for offline health transitions → WS trigger.
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

		// Force offline → offline RG fires WS trigger → connects + writes to sink.
		assetSteps.ForceOfflineByAdmin("saga-websocket-phase1-warmup"),

		// Confirm offline observed.
		assetAsserts.AssertHealthStatusEventually("offline"),

		// Smoke: events_trigger records at least 1 successful execution.
		eventAsserts.AssertTriggerExecutedSuccessfullyEventually(1),

		// Content-key: resolved config carries the saga-scoped runID in the message.
		eventAsserts.AssertLastTriggerRequestDataContains("\"saga\":\"websocket\""),

		// Heartbeat brings asset back online; online RG fires WS trigger again.
		assetSteps.SendHttpHeartbeat(),

		// Confirm online observed again.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Smoke: at least 2 successful executions recorded across both transitions.
		eventAsserts.AssertTriggerExecutedSuccessfullyEventually(2),

		// Explicit asset teardown.
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
