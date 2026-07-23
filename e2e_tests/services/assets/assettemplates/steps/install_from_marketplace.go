package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// installBody is the marketplace install/uninstall request body. shareWithChildren
// is false for the single-org e2e — the suite runs as one org (§13), so there are no
// child orgs to share the install with.
type installBody struct {
	ShareWithChildren bool `json:"shareWithChildren"`
}

// marketplaceInstallResponse decodes the installed-link id the install returns.
type marketplaceInstallResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// InstallFromMarketplace POSTs /api/v1/asset_templates/{vendor}/{slug}/install and
// publishes the caller-org installed-link id on BagKeyTemplateID — the same key a
// saga-created template writes — so a downstream CreateAsset binds an asset to the
// installed marketplace template (the uninstall guard then sees the reference).
//
// The bundle is served by the e2e marketplace mock, which advertises the sha256 the
// assets service hard-verifies before persisting anything. vendor/slug are
// compile-time wiring, not runtime bag values (§3).
//
// Reads (bag): none.
//
// Writes (bag):
//   - BagKeyTemplateID  string  Mongo ObjectID hex of the installed link
//
// Compensate: DELETE the install. 404-tolerant (an explicit UninstallFromMarketplace
// earlier in the journey may already have removed it) and a no-op when Do never ran.
func InstallFromMarketplace(vendor, slug string) saga.Step {
	path := "/api/v1/asset_templates/" + vendor + "/" + slug + "/install"
	return saga.Step{
		Name: fmt.Sprintf("assets/assettemplates.InstallFromMarketplace[%s/%s]", vendor, slug),
		Do: func(c *saga.Context) error {
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, path, installBody{ShareWithChildren: false})
			if err != nil {
				return fmt.Errorf("install %s/%s: %w", vendor, slug, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("install %s/%s: unexpected status %d body=%s", vendor, slug, resp.StatusCode, string(body))
			}
			var out marketplaceInstallResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode install %s/%s response: %w", vendor, slug, err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("install %s/%s: empty id in response", vendor, slug)
			}
			c.Set(BagKeyTemplateID, out.Data.ID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			// Do never ran (no installed id on the bag): nothing to undo.
			if _, ok := c.Get(BagKeyTemplateID); !ok {
				return nil
			}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodDelete, path, installBody{ShareWithChildren: false})
			if err != nil {
				return fmt.Errorf("uninstall %s/%s: %w", vendor, slug, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("uninstall %s/%s: unexpected status %d", vendor, slug, resp.StatusCode)
			}
			return nil
		},
	}
}
