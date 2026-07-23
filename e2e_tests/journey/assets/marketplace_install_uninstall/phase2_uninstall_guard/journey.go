// Package phase2_uninstall_guard exercises the usage guard on marketplace uninstall:
// an installed DT CO2 template that an asset references cannot be removed, and the
// uninstall succeeds only once that asset is gone.
//
// Outcome on PASS:
//   - The DT CO2 template installs and an asset is created referencing it (the install
//     step publishes the installed-link id on the template bag key CreateAsset reads).
//   - Uninstalling the template is refused with HTTP 403 carrying TEMPLATE_IN_USE.
//   - After the referencing asset is deleted, uninstalling the template succeeds.
//   - Compensation leaves the stack clean.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package phase2_uninstall_guard

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	bootstrap "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"

	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	tmplAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/asserts"
	tmplPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
	tmplSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

// Items is the ordered slice of saga Items the journey runs.
//
//  1. CreateRouteGroup          -> route group the asset requires
//  2. InstallFromMarketplace    -> DT CO2 installed; installed-link id on the template bag key
//  3. CreateAsset               -> asset bound to the installed marketplace template
//  4. AssertUninstallBlocked    -> uninstall refused with 403 TEMPLATE_IN_USE (the guard)
//  5. DeleteAsset               -> remove the only reference to the template
//  6. UninstallFromMarketplace  -> uninstall now succeeds
func Items() []saga.Item {
	return []saga.Item{
		rgSteps.CreateRouteGroup(),
		tmplSteps.InstallFromMarketplace(tmplPayloads.DTCo2Vendor, tmplPayloads.GuardSlug),
		assetSteps.CreateAsset(),
		tmplAsserts.AssertUninstallBlocked(tmplPayloads.DTCo2Vendor, tmplPayloads.GuardSlug),
		assetSteps.DeleteAsset(),
		tmplSteps.UninstallFromMarketplace(tmplPayloads.DTCo2Vendor, tmplPayloads.GuardSlug),
	}
}

// Run executes phase 0 (IAM bootstrap) + this journey as a single saga.
func Run(t *testing.T) {
	t.Helper()
	runID := random.NewRunID()
	clients := bootstrap.NewClients()
	items := append(bootstrap.BootstrapItems(), Items()...)
	saga.Run(t, context.Background(), runID, clients, items...)
}
