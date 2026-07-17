// Package ota_mqtt drives the OTA firmware-update rollout for an MQTT device end
// to end against the live stack — the PUSH transport. The broker delivers the
// ota_update command to a connected, presence-gated device on
// commands/{assetUUID}/ota_update; the device downloads + verifies the firmware
// from a presigned URL and reports its progress on events/{assetUUID}/ota_status.
// It shares the Layer-1 OTA building blocks with ota_http but never imports it.
//
// Ordering contract: the reconciler pushes ONCE, to an ONLINE device only, and
// advances the execution to INITIATED after dispatch. So the sim must connect,
// subscribe to the command topic, and be presence-online BEFORE CreatePlan, or the
// pushed command is missed.
//
// Outcome on PASS:
//   - Firmware is uploaded (presigned PUT) and marked READY; an OTA plan is created
//     migrating the asset from the source template to the firmware's target template.
//   - The reconciler pushes ota_update to the online device; the sim downloads +
//     verifies (size + base64 sha256) and drives the execution downloading→…→updated.
//   - The plan reaches COMPLETED, the execution UPDATED@100, and the asset's
//     template is switched source→target — all observed through the public API.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package ota_mqtt

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	bootstrap "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	assetAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/asserts"
	assetPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	templatePayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	otaAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/asserts"
	otaSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/steps"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// Items is the ordered slice of saga Items the ota_mqtt journey runs. The
// connect → subscribe → online sequence MUST precede CreatePlan (the push is
// once-only, to an online device).
func Items() []saga.Item {
	return []saga.Item{
		// Source template (the asset starts here) + target template (the migration destination).
		templateSteps.CreateTemplateWithLabel(templatePayloads.SagaOtaSourceTemplate, "ota-source"),
		templateSteps.CreateTemplateWithLabel(templatePayloads.SagaOtaTargetTemplate, "ota-target"),

		// Route group — the asset contract requires RouteGroupIds (min 1).
		rgSteps.CreateRouteGroup(),

		// MQTT connectivity asset (password + HealthMonitor) on the source template.
		assetSteps.CreateConnectivityAsset(otaMqttAsset()),

		// Presence + subscribe contract — all BEFORE CreatePlan:
		// connect the sim, arm the command subscription, and wait until it is online.
		assetSteps.ConnectMqttPassword(),
		otaSteps.SubscribeOtaCommand(),
		assetAsserts.AssertHealthStatusEventually("online"),

		// Firmware lifecycle: init (declare target template + size + base64 sha256) →
		// presigned PUT → complete (READY).
		otaSteps.InitFirmware(templateSteps.TemplateIDKey("ota-target")),
		otaSteps.UploadFirmwareBinary(),
		otaSteps.CompleteFirmware(),

		// OTA plan: creating it triggers the push to the online device on the next scan.
		otaSteps.CreatePlan(templateSteps.TemplateIDKey("ota-source"), assetSteps.BagKeyAssetID),

		// The sim receives the pushed command, downloads + verifies, and publishes status.
		otaSteps.RunMqttOtaDevice(),

		// Public-API outcome: plan COMPLETED, execution UPDATED@100, asset template switched.
		otaAsserts.AssertPlanStatusEventually("COMPLETED"),
		otaAsserts.AssertExecutionStateEventually(assetSteps.BagKeyAssetID, "UPDATED", 100),
		otaAsserts.AssertAssetTemplateSwitched(assetSteps.BagKeyAssetID, templateSteps.TemplateIDKey("ota-target")),
	}
}

// otaMqttAsset resolves the MQTT OTA device payload at execution time: it reads the
// source template id and the route group id from the bag and builds the asset on
// them. Returned as a ConnectivityPayloadFn so it carries no runtime data itself.
func otaMqttAsset() assetSteps.ConnectivityPayloadFn {
	return func(c *saga.Context) *assetPayloads.AssetCreateBuilder {
		templateID := c.MustGetString(templateSteps.TemplateIDKey("ota-source"))
		routeGroupID := c.MustGetString(rgSteps.BagKeyRouteGroupID)
		return assetPayloads.SagaOtaMqttDevice(c.RunID, templateID, routeGroupID)
	}
}

// Run fronts the shared IAM bootstrap (common/journey/iam_bootstrap) and runs this
// journey's items under one rollback chain. The suite runner (journey/suite)
// provisions the stack once via infra.EnsureAll; Run never touches the environment.
func Run(t *testing.T) {
	t.Helper()
	runID := random.NewRunID()
	clients := bootstrap.NewClients()
	items := append(bootstrap.BootstrapItems(), Items()...)
	saga.Run(t, context.Background(), runID, clients, items...)
}
