// Package marketplace_asset_pipeline drives the full telemetry pipeline on top of
// a marketplace-installed asset template: install the DT CO2 template from the
// e2e marketplace mock, bind an asset to the installed link, POST a telemetry
// event, and watch every consuming microservice pick it up end-to-end
// (http_gateway -> events -> js-executor -> router -> triggers -> HTTP sink).
//
// This is the cross-service proof of the marketplace durability model: the
// install persists the heavy body ONCE as a shared content document and links
// the caller org to it; here the installed link's scripts must actually reach
// js-executor and drive the transform, or the trigger never fires.
//
// Outcome on PASS:
//   - InstallFromMarketplace persists the DT CO2 link (sha256 hard-verified) and
//     publishes the installed template id on the same bag key a saga-created
//     template uses, so the asset binds straight to the marketplace template.
//   - The telemetry POST is ingested by http_gateway and, hop by hop:
//       * events stores the raw event (AssertRawEventReceivedAfter);
//       * js-executor runs the installed template's ScriptConversion and events stores
//         the decoded StandardizedPayload with the expected values
//         (AssertDecodedEventReceivedAfter) — the asset template is truly exercised;
//       * router matches and triggers records a successful execution read back from
//         the events /events/trigger store (AssertTriggerExecutedSuccessfullyEventually);
//       * the in-process HTTP sink receives the trigger's POST (AssertSinkHitEventually).
//     Together these prove http_gateway, events, js-executor, router, and triggers all
//     consumed the event — the flow is proven 100% end-to-end.
//   - Compensation deletes the asset (cascade) and uninstalls the template, leaving
//     the stack clean.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log; a missing raw event
//     points at gateway/events, a missing trigger execution at js-executor/router,
//     and a missing sink hit at the triggers executor.
package marketplace_asset_pipeline

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	bootstrap "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"

	assetPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	tmplPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	eventAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/events/events/asserts"
	dsSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/steps"
	rgPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/payloads"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
	triggerAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/asserts"
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// Items is the ordered slice of saga Items the journey runs.
func Items() []saga.Item {
	return []saga.Item{
		// Boot the in-process HTTP sink the trigger executor will POST to.
		triggerSteps.StartTestSink(),

		// Create an HTTP-kind trigger pointing at the sink.
		triggerSteps.CreateTrigger(),

		// Single trigger-kind route group; reused as online/offline so telemetry fires it.
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOnlineRouteGroupID, eventTriggerRG()),
		rgSteps.CreateRouteGroupAt(rgSteps.BagKeyOfflineRouteGroupID, eventTriggerRG()),

		// HTTP data source — its API key authenticates the POST /events call.
		dsSteps.CreateDataSource(),

		// Install the DT CO2 template from the marketplace mock: persists the shared
		// content once, links this org to it, and publishes the installed template id.
		templateSteps.InstallFromMarketplace(tmplPayloads.DTCo2Vendor, tmplPayloads.PipelineSlug),

		// HTTP connectivity asset bound to the INSTALLED marketplace template.
		assetSteps.CreateConnectivityAsset(marketplaceEventAsset()),

		// POST a telemetry event to the gateway with the saga's runID embedded in the body.
		dsSteps.PostRawEvent(),

		// events MS ingested the raw telemetry event (queryable by the asset's threadId).
		eventAsserts.AssertRawEventReceivedAfter(dsSteps.BagKeyTelemetrySentAt),

		// events MS ran the INSTALLED marketplace template's ScriptConversion in
		// js-executor and stored the decoded StandardizedPayload; asserting the exact
		// decoded values proves the marketplace link's shared script reached js-executor
		// and the asset template was truly exercised — not just installed.
		eventAsserts.AssertDecodedEventReceivedAfter(dsSteps.BagKeyTelemetrySentAt, map[string]any{"value": 23.5, "unit": "C"}),

		// router matched and the triggers executor recorded a successful execution,
		// read back from the events service's /events/trigger store.
		eventAsserts.AssertTriggerExecutedSuccessfullyEventually(1),

		// Content-key: the HTTP sink received the POST the trigger emitted (100% e2e).
		triggerAsserts.AssertSinkHitEventually(1),

		// Explicit asset teardown so Compensate verifies cascade cleanup; the install
		// Compensate then uninstalls the template.
		assetSteps.DeleteAsset(),
	}
}

// eventTriggerRG builds a single kind=trigger route group reused for online and
// offline routing. Match conditions are intentionally empty — every event the asset
// emits fires this route group and therefore the trigger.
func eventTriggerRG() rgSteps.BuilderFn {
	return func(c *saga.Context) *rgPayloads.RouteGroupCreateBuilder {
		triggerID := c.MustGetString(triggerSteps.BagKeyTriggerID)
		return rgPayloads.SagaTriggerRouteGroup(c.RunID, "event", triggerID)
	}
}

// marketplaceEventAsset binds the connectivity asset to the marketplace-installed
// template (read from the bag) and to the same trigger RG for both health flows, so
// the telemetry path fires the trigger.
func marketplaceEventAsset() assetSteps.ConnectivityPayloadFn {
	return func(c *saga.Context) *assetPayloads.AssetCreateBuilder {
		templateID := c.MustGetString(templateSteps.BagKeyTemplateID)
		onlineRG := c.MustGetString(rgSteps.BagKeyOnlineRouteGroupID)
		offlineRG := c.MustGetString(rgSteps.BagKeyOfflineRouteGroupID)
		return assetPayloads.SagaHttpConnectivitySensor(c.RunID, templateID, onlineRG, offlineRG)
	}
}

// Run wires Phase 0 (IAM bootstrap) in front of this journey's items.
func Run(t *testing.T) {
	t.Helper()
	runID := random.NewRunID()
	clients := bootstrap.NewClients()
	items := append(bootstrap.BootstrapItems(), Items()...)
	saga.Run(t, context.Background(), runID, clients, items...)
}
