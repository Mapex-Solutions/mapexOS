package steps

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// UninstallFromMarketplace DELETEs the caller-org install of {vendor}/{slug} and
// expects success — the explicit clean-uninstall path, exercised when no asset
// references the template. It proves the happy uninstall (HTTP 200), the counterpart
// to the guarded-403 path AssertUninstallBlocked verifies.
//
// vendor/slug are compile-time wiring (§3); the org comes from the auth headers.
//
// Reads (bag): none.
//
// Writes (bag): none.
//
// Compensate: documented no-op. Uninstall is itself a teardown — re-installing is not
// a meaningful inverse, and the InstallFromMarketplace step already owns removing the
// link (its Compensate is 404-tolerant, so a later rollback simply finds it gone).
func UninstallFromMarketplace(vendor, slug string) saga.Step {
	path := "/api/v1/asset_templates/" + vendor + "/" + slug + "/install"
	return saga.Step{
		Name: fmt.Sprintf("assets/assettemplates.UninstallFromMarketplace[%s/%s]", vendor, slug),
		Do: func(c *saga.Context) error {
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodDelete, path, installBody{ShareWithChildren: false})
			if err != nil {
				return fmt.Errorf("uninstall %s/%s: %w", vendor, slug, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("uninstall %s/%s: unexpected status %d body=%s", vendor, slug, resp.StatusCode, string(body))
			}
			return nil
		},
		Compensate: func(_ *saga.Context) error {
			return nil
		},
	}
}
