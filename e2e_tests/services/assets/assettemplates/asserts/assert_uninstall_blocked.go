package asserts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// templateInUseCode is the machine-readable error code the guard returns with the
// 403 (the contract's ErrCodeTemplateInUse). Asserted on the wire so a future rename
// that breaks the toggle's error handling is caught here.
const templateInUseCode = "TEMPLATE_IN_USE"

// uninstallBody is the DELETE request body the uninstall route validates.
type uninstallBody struct {
	ShareWithChildren bool `json:"shareWithChildren"`
}

// uninstallErrorResponse decodes the envelope's error list.
type uninstallErrorResponse struct {
	Errors []string `json:"errors"`
}

// AssertUninstallBlocked sends DELETE /api/v1/asset_templates/{vendor}/{slug}/install
// while an asset still references the installed template, and asserts the usage guard
// rejects it with HTTP 403 carrying the TEMPLATE_IN_USE code. Because the guard
// refuses the delete, no state changes — the check stays a read-only oracle.
//
// vendor/slug are compile-time wiring; the reference is provided by an upstream
// CreateAsset bound to the installed template (§3).
func AssertUninstallBlocked(vendor, slug string) saga.Assert {
	path := "/api/v1/asset_templates/" + vendor + "/" + slug + "/install"
	return saga.Assert{
		Name: fmt.Sprintf("assets/assettemplates.AssertUninstallBlocked[%s/%s]", vendor, slug),
		Check: func(c *saga.Context) error {
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodDelete, path, uninstallBody{ShareWithChildren: false})
			if err != nil {
				return fmt.Errorf("uninstall-blocked request %s/%s: %w", vendor, slug, err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusForbidden {
				return fmt.Errorf("uninstall of in-use template %s/%s: want status 403, got %d body=%s", vendor, slug, resp.StatusCode, string(body))
			}

			var out uninstallErrorResponse
			if err := json.Unmarshal(body, &out); err != nil {
				return fmt.Errorf("decode uninstall-blocked errors %s/%s: %w", vendor, slug, err)
			}
			if slices.Contains(out.Errors, templateInUseCode) {
				return nil
			}
			return fmt.Errorf("uninstall of in-use template %s/%s: 403 did not carry %q, errors=%v", vendor, slug, templateInUseCode, out.Errors)
		},
	}
}
