// Package phase1_connectivity exercises the Email trigger smoke path
// via the connectivity pipeline: the platform's healthmonitor flips
// the saga asset between online and offline, the router matches the
// transition against the route groups, and the triggers service
// delivers a real message to the in-process SMTP sink the journey
// stood up. The saga validates both halves of the smoke:
//
//   - **Smoke (hit count).** The SMTP sink observes at least N
//     deliveries after N health transitions — confirms the trigger
//     fired and reached the SMTP layer end-to-end.
//   - **Content-key.** The most recent captured message carries the
//     expected from / to envelope and a subject fragment that
//     interpolates the run id — confirms the operator-supplied config
//     reached the executor unchanged (catches placeholder-rendering
//     and credential-mapping regressions).
//
// Outcome on PASS:
//   - Phase 0 outcomes hold (IAM bootstrap succeeds).
//   - The actor can create an Email trigger, two trigger-kind route
//     groups, and an HTTP connectivity asset bound to them.
//   - The first heartbeat warms the asset to "online" silently (no
//     trigger by design — first observation is unknown→online).
//   - Force-offline drives the asset to "offline", the offline route
//     group matches, and the triggers service delivers one message
//     to the SMTP sink.
//   - A new heartbeat drives the asset back to "online", the online
//     route group matches, and a second message is delivered.
//
// Outcome on FAIL:
//   - Logs identify the offending step or assertion by qualified name
//     (triggers/triggers.StartSmtpSink, triggers/triggers.CreateEmailTrigger,
//     triggers/triggers.AssertSmtpReceivedEventually, ...) so the failure
//     points at the breaking integration.
//   - "smtp hits: want >=1, got 0" typically means either the trigger
//     never matched the route group (router config), the executor
//     can't reach the sink (host/port resolution, see SmtpSinkHost),
//     or the SMTP server rejected the message (check SmtpServer logs).
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

// Items is the ordered slice of saga Items the phase runs. Each line
// carries a single comment above it explaining what that item proves
// or sets up, so readers don't have to chase the step package to know
// why this item is in the chain.
func Items() []saga.Item {
	return []saga.Item{
		// Boot an in-process SMTP server on SmtpSinkBindAddr; captures every successful delivery.
		triggerSteps.StartSmtpSink(),

		// Create an Email trigger pointing at the SMTP sink (overrides smtpHost/smtpPort to the saga sink).
		triggerSteps.CreateEmailTrigger(),

		// Route group that matches asset.health "online" transitions and points at the email trigger.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, onlineEmailTriggerRG()),

		// Route group that matches asset.health "offline" transitions and points at the email trigger.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOfflineRouteGroupID, offlineEmailTriggerRG()),

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

		// Real transition 1: admin forces the asset offline; the offline RG matches → email trigger fires.
		assetSteps.ForceOfflineByAdmin("saga-email-phase1-warmup"),

		// Confirm the healthmonitor saw the transition.
		assetAsserts.AssertHealthStatusEventually("offline"),

		// Smoke oracle: the SMTP sink has captured at least 1 delivery from the offline transition.
		triggerAsserts.AssertSmtpReceivedEventually(1),

		// Content-key oracle: envelope from/to and subject carry the values the trigger config emitted.
		triggerAsserts.AssertSmtpLastMessageContent("saga-no-reply@mapex.test", "saga-recipient@mapex.test", "Saga email smoke"),

		// Real transition 2: a fresh heartbeat brings the asset back online; the online RG matches → email trigger fires again.
		assetSteps.SendHttpHeartbeat(),

		// Confirm the asset is online again.
		assetAsserts.AssertHealthStatusEventually("online"),

		// Smoke oracle: the SMTP sink has now captured at least 2 deliveries (one per transition).
		triggerAsserts.AssertSmtpReceivedEventually(2),

		// Tear down the asset explicitly so the Compensate chain can verify cascade cleanup.
		assetSteps.DeleteAsset(),
	}
}

// onlineEmailTriggerRG builds a kind=trigger route group keyed to the
// "online" health transition, pointing at the trigger id Phase 1
// published on the bag.
func onlineEmailTriggerRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
		return rgPayloads.SagaTriggerRouteGroup(c.RunID, "online", triggerID)
	}
}

// offlineEmailTriggerRG builds the offline counterpart of the online
// route group above.
func offlineEmailTriggerRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
		return rgPayloads.SagaTriggerRouteGroup(c.RunID, "offline", triggerID)
	}
}

// httpConnectivityAsset binds the saga's HTTP connectivity asset to
// the online + offline trigger-kind route groups created above.
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
