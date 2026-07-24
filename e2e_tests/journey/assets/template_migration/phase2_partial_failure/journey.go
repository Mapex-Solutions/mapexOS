// Package phase2_partial_failure drives asset-template migration with one asset that
// cannot migrate, proving the error path: the plan finalizes "completed_with_errors"
// with the real assets still switched and only the bad asset failed.
//
// Outcome on PASS:
//   - The plan reaches "completed_with_errors" with migrated==2 and failed==1.
//   - The two real assets' executions are "migrated" and both are re-pointed to the
//     target; the fabricated asset's execution is "failed" with a non-empty error and
//     a single attempt.
//   - Compensation removes the plan, assets, route group and templates.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package phase2_partial_failure

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

// fakeAssetID is a syntactically-valid but non-existent asset ObjectId. The migration
// create does not validate asset existence, so it is accepted and its execution fails
// at run time — the deterministic lever for the partial-failure path. If create ever
// starts validating asset existence, replace this with deleting a real asset after
// creating a future-scheduled plan and before it runs.
const fakeAssetID = "0123456789abcdef01234567"

// Items is the ordered slice of saga Items the journey runs.
func Items() []saga.Item {
	return []saga.Item{
		// Source template (assets start here) + target template (destination).
		tmplSteps.CreateTemplate(),
		tmplSteps.CreateTemplateWithLabel(tmplPayloads.SagaTemperatureTemplate, "mig-to"),

		// Route group — required by the non-gateway asset-create contract (setup only).
		rgSteps.CreateRouteGroup(),

		// Two real assets on the source template.
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaMqttTemperatureSensorFor("mig-1"), "mig-1"),
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaMqttTemperatureSensorFor("mig-2"), "mig-2"),

		// Immediate plan for the two real assets plus one fabricated, non-existent id.
		tmplSteps.CreateMigrationPlanWithExtra(
			tmplSteps.BagKeyTemplateID,
			tmplSteps.TemplateIDKey("mig-to"),
			[]string{fakeAssetID},
			assetSteps.AssetIDKey("mig-1"),
			assetSteps.AssetIDKey("mig-2"),
		),

		// The plan finalizes with the bad asset failed and the rest migrated.
		tmplAsserts.AssertMigrationPlanStatusEventually("completed_with_errors"),

		// The two real executions succeeded; the fabricated one failed.
		tmplAsserts.AssertMigrationExecutionStatus(assetSteps.AssetIDKey("mig-1"), "migrated"),
		tmplAsserts.AssertMigrationExecutionStatus(assetSteps.AssetIDKey("mig-2"), "migrated"),
		tmplAsserts.AssertMigrationExecutionStatusByID(fakeAssetID, "failed"),

		// The real assets still switched to the target despite the partial failure.
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
