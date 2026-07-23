// Package asserts holds read-only oracles for the assets/assettemplates module,
// verifying marketplace install state through the public HTTP API only (§6).
package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// installedCheckRequest is the batch installed-check body.
type installedCheckRequest struct {
	MarketplaceGuids []string `json:"marketplaceGuids"`
}

// installedCheckResponse decodes the subset of the batch response the assert needs.
type installedCheckResponse struct {
	Data struct {
		Installed []string `json:"installed"`
	} `json:"data"`
}

// AssertMarketplaceInstalled polls POST /api/v1/asset_templates/marketplace/installed
// with the single guid and asserts the batch result contains it (want=true) or omits
// it (want=false). This is the exact call the marketplace listing makes to toggle a
// card between Install and Uninstall, so the assert proves the toggle's data source.
//
// Install and uninstall are synchronous, so a single round usually suffices; the
// short poll budget only absorbs a slow or first-run stack.
func AssertMarketplaceInstalled(guid string, want bool) saga.Assert {
	return marketplaceInstalledAssert(guid, want, constants.ScaleTimeout(10*time.Second), 500*time.Millisecond)
}

// marketplaceInstalledAssert is the shared polling body: it queries the batch
// installed-check until the guid's presence equals want or the timeout elapses,
// surfacing both the last-seen presence and the last fetch error on timeout.
func marketplaceInstalledAssert(guid string, want bool, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("assets/assettemplates.AssertMarketplaceInstalled[%s=%t]", guid, want),
		Check: func(c *saga.Context) error {
			deadline := time.Now().Add(timeout)
			var lastSeen bool
			var lastErr error
			for {
				got, err := installedContainsGuid(c, guid)
				if err != nil {
					lastErr = err
				} else {
					lastSeen = got
					if got == want {
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("installed-check for guid %q: want present=%t, last present=%t (lastErr=%v)", guid, want, lastSeen, lastErr)
				}
				time.Sleep(tick)
			}
		},
	}
}

// installedContainsGuid runs one batch installed-check and reports whether the guid
// is in the installed set.
func installedContainsGuid(c *saga.Context, guid string) (bool, error) {
	resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/asset_templates/marketplace/installed", installedCheckRequest{MarketplaceGuids: []string{guid}})
	if err != nil {
		return false, fmt.Errorf("installed-check request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("installed-check: unexpected status %d", resp.StatusCode)
	}
	var out installedCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return false, fmt.Errorf("decode installed-check: %w", err)
	}
	return slices.Contains(out.Data.Installed, guid), nil
}
