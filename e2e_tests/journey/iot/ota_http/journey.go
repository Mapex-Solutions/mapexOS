// Package ota_http drives the OTA firmware-update rollout for an HTTP device end
// to end against the live stack — the POLL transport. The device pulls its job
// from the gateway (GET /api/v1/ota/jobs), downloads + verifies the firmware from
// a presigned URL, and reports its progress to POST /api/v1/ota/status; there is
// no presence gate (the device shows up by polling, no push race). It shares the
// Layer-1 OTA building blocks with ota_mqtt but never imports it.
//
// Outcome on PASS:
//   - Firmware is uploaded (presigned PUT) and marked READY; an OTA plan is created
//     migrating the asset from the source template to the firmware's target template.
//   - The HTTP device polls its job, downloads + verifies (size + base64 sha256),
//     and drives the execution downloading→…→updated via POST /ota/status.
//   - The plan reaches COMPLETED, the execution UPDATED@100, and the asset's
//     template is switched source→target — all observed through the public API.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package ota_http

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	bootstrap "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	assetPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	templatePayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	otaAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/asserts"
	otaSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/steps"
	dsSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/steps"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// Items is the ordered slice of saga Items the ota_http journey runs. The HTTP
// device polls for its job, so there is no subscribe/presence ordering — the
// device simply runs after the plan exists.
func Items() []saga.Item {
	return []saga.Item{
		// Source template (the asset starts here) + target template (the migration destination).
		templateSteps.CreateTemplateWithLabel(templatePayloads.SagaOtaSourceTemplate, "ota-source"),
		templateSteps.CreateTemplateWithLabel(templatePayloads.SagaOtaTargetTemplate, "ota-target"),

		// Route group — the asset contract requires RouteGroupIds (min 1).
		rgSteps.CreateRouteGroup(),

		// HTTP connectivity asset on the source template, bound to the route group.
		assetSteps.CreateConnectivityAsset(otaHttpAsset()),

		// Data source — its id + api key authenticate the device's gateway calls (?ds=).
		dsSteps.CreateDataSource(),

		// Firmware lifecycle: init (declare target template + size + base64 sha256) →
		// presigned PUT → complete (READY).
		otaSteps.InitFirmware(templateSteps.TemplateIDKey("ota-target")),
		otaSteps.UploadFirmwareBinary(),
		otaSteps.CompleteFirmware(),

		// OTA plan: migrate the asset from the source template to the firmware's target.
		otaSteps.CreatePlan(templateSteps.TemplateIDKey("ota-source"), assetSteps.BagKeyAssetID),

		// The HTTP device polls its job, downloads + verifies, and posts the status progression.
		otaSteps.RunHttpOtaDevice(),

		// Public-API outcome: plan COMPLETED, execution UPDATED@100, asset template switched.
		otaAsserts.AssertPlanStatusEventually("COMPLETED"),
		otaAsserts.AssertExecutionStateEventually(assetSteps.BagKeyAssetID, "UPDATED", 100),
		otaAsserts.AssertAssetTemplateSwitched(assetSteps.BagKeyAssetID, templateSteps.TemplateIDKey("ota-target")),
	}
}

// otaHttpAsset resolves the HTTP OTA device payload at execution time: it reads
// the source template id and the route group id from the bag and builds the asset
// on them. Returned as a ConnectivityPayloadFn so it carries no runtime data itself.
func otaHttpAsset() assetSteps.ConnectivityPayloadFn {
	return func(c *saga.Context) *assetPayloads.AssetCreateBuilder {
		templateID := c.MustGetString(templateSteps.TemplateIDKey("ota-source"))
		routeGroupID := c.MustGetString(rgSteps.BagKeyRouteGroupID)
		return assetPayloads.SagaOtaHttpDevice(c.RunID, templateID, routeGroupID)
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
