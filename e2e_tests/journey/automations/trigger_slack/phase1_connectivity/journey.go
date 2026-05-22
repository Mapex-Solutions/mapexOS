// Package phase1_connectivity exercises the Slack trigger smoke
// path via the connectivity pipeline. The Slack executor POSTs the
// webhook URL as a real HTTP request, so the saga reuses the
// in-process HTTP sink (StartTestSink) instead of standing up a
// separate Slack-shaped server.
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
		// Boot the shared HTTP sink (Slack webhooks are POST HTTP under the hood).
		triggerSteps.StartTestSink(),

		// Create a Slack-kind trigger whose webhookUrl points at the HTTP sink.
		triggerSteps.CreateSlackTrigger(),

		// Route group for online health transitions → Slack trigger.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, onlineTriggerRG()),

		// Route group for offline health transitions → Slack trigger.
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

		// Force offline → offline RG fires Slack trigger → POSTs to HTTP sink.
		assetSteps.ForceOfflineByAdmin("saga-slack-phase1-warmup"),

		// Confirm offline observed.
		assetAsserts.AssertHealthStatusEventually("offline"),

		// Smoke: sink received the Slack webhook POST.
		triggerAsserts.AssertSinkHitEventually(1),

		// Force online → online RG fires Slack trigger again.
		assetSteps.SendHttpHeartbeat(),

		// Confirm online observed again.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Smoke: sink received the second Slack POST.
		triggerAsserts.AssertSinkHitEventually(2),

		// Explicit asset teardown so Compensate verifies cascade cleanup.
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
