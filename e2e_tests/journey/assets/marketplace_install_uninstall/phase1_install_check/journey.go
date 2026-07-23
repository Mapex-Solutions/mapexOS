// Package phase1_install_check exercises the asset-template marketplace install +
// batch installed-check + clean uninstall against the live stack, for the Disruptive
// Technologies CO2 sensor template served by the e2e marketplace mock.
//
// Outcome on PASS:
//   - Before install, the batch installed-check reports the DT CO2 guid as NOT installed.
//   - Installing it from the marketplace succeeds (the assets service hard-verifies the
//     bundle's sha256 against the mock's advertised checksum before persisting).
//   - The batch installed-check now reports the guid as installed — the exact data the
//     listing uses to show Uninstall instead of Install on a card.
//   - A clean uninstall (no asset references the template) succeeds, and the check
//     reports the guid not-installed again.
//   - Compensation leaves the stack clean.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package phase1_install_check

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	bootstrap "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"

	tmplAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/asserts"
	tmplPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
	tmplSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
)

// Items is the ordered slice of saga Items the journey runs.
//
//  1. AssertMarketplaceInstalled(false) -> the guid is not installed yet
//  2. InstallFromMarketplace           -> DT CO2 installed (sha256 hard-verified)
//  3. AssertMarketplaceInstalled(true)  -> the batch check now reports it installed
//  4. UninstallFromMarketplace          -> clean uninstall succeeds (no references)
//  5. AssertMarketplaceInstalled(false) -> the batch check reports it gone again
func Items() []saga.Item {
	return []saga.Item{
		tmplAsserts.AssertMarketplaceInstalled(tmplPayloads.DTCo2MarketplaceGuid, false),
		tmplSteps.InstallFromMarketplace(tmplPayloads.DTCo2Vendor, tmplPayloads.DTCo2Slug),
		tmplAsserts.AssertMarketplaceInstalled(tmplPayloads.DTCo2MarketplaceGuid, true),
		tmplSteps.UninstallFromMarketplace(tmplPayloads.DTCo2Vendor, tmplPayloads.DTCo2Slug),
		tmplAsserts.AssertMarketplaceInstalled(tmplPayloads.DTCo2MarketplaceGuid, false),
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
