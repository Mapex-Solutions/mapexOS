// Package phase1_immediate_complete drives the asset-template migration feature end
// to end: it creates a source and a target template, two assets on the source, then
// a migration plan (immediate) and watches the plan run and re-point BOTH assets onto
// the target — all through the public API.
//
// Outcome on PASS:
//   - The plan reaches "complete" with migrated==total and failed==0.
//   - Each per-asset execution is "migrated".
//   - Each asset's assetTemplateId now equals the target template id.
//   - Compensation removes the plan, assets, route group and templates.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log; a plan stuck short of
//     "complete" points at the timer/consumer run path, an unswitched asset at the
//     template-switch adapter.
package phase1_immediate_complete

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	bootstrap "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	assetPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	tmplAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/asserts"
	tmplPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
	tmplSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	otaAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/asserts"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// Items is the ordered slice of saga Items the journey runs.
func Items() []saga.Item {
	return []saga.Item{
		// Source template (assets start here) on the shared key CreateAssetWithLabel binds to.
		tmplSteps.CreateTemplate(),

		// Target template (the migration destination).
		tmplSteps.CreateTemplateWithLabel(tmplPayloads.SagaTemperatureTemplate, "mig-to"),

		// Route group — the non-gateway asset-create contract requires one (setup only).
		rgSteps.CreateRouteGroup(),

		// Two assets on the source template, so the plan proves a batch migration.
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaMqttTemperatureSensorFor("mig-1"), "mig-1"),
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaMqttTemperatureSensorFor("mig-2"), "mig-2"),

		// Immediate migration plan moving both assets from the source to the target.
		tmplSteps.CreateMigrationPlan(
			tmplSteps.BagKeyTemplateID,
			tmplSteps.TemplateIDKey("mig-to"),
			assetSteps.AssetIDKey("mig-1"),
			assetSteps.AssetIDKey("mig-2"),
		),

		// The plan ran to completion: every asset migrated, none failed.
		tmplAsserts.AssertMigrationPlanStatusEventually("complete"),

		// Each per-asset execution succeeded.
		tmplAsserts.AssertMigrationExecutionStatus(assetSteps.AssetIDKey("mig-1"), "migrated"),
		tmplAsserts.AssertMigrationExecutionStatus(assetSteps.AssetIDKey("mig-2"), "migrated"),

		// The observable outcome: each asset is now on the target template.
		otaAsserts.AssertAssetTemplateSwitched(assetSteps.AssetIDKey("mig-1"), tmplSteps.TemplateIDKey("mig-to")),
		otaAsserts.AssertAssetTemplateSwitched(assetSteps.AssetIDKey("mig-2"), tmplSteps.TemplateIDKey("mig-to")),
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
