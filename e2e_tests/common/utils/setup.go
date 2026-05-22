package utils

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
)

var (
	setupOnce sync.Once
	setupErr  error
)

// SetupE2EEnvironment performs a lightweight readiness check against the
// running mapexIam service. It is idempotent: every test package that
// calls it observes the same outcome without mutating state.
//
// History: an earlier implementation invoked scripts/e2e-setup.sh which
// dropped collections and re-seeded data with hard-coded ids. Those ids
// drifted from the mongodb-init seed shipped with the docker-compose
// stacks (services_required and standalone/infra), and re-running the
// script broke with duplicate-key errors. The canonical mongodb-init
// seed already provides everything saga journeys need:
//
//   - The seed root organization (constants.MapexosOrgID).
//   - The seed admin user (constants.RootUserEmail) with
//     constants.RootUserPassword and a recursive-scope membership in
//     the root org.
//
// SetupE2EEnvironment therefore becomes a thin liveness check: the stack
// must be up. Whether the admin user can authenticate is verified by
// the saga's own SeedAdminLogin step — repeating the login here would
// duplicate that responsibility.
func SetupE2EEnvironment() error {
	setupOnce.Do(func() {
		setupErr = checkStackReady()
	})
	return setupErr
}

// ForceSetupE2EEnvironment re-runs the readiness check ignoring the
// sync.Once gate. Use it only when a test deliberately needs to reverify
// the stack mid-run; the regular SetupE2EEnvironment entry point is
// almost always what callers want.
func ForceSetupE2EEnvironment() error {
	return checkStackReady()
}

// checkStackReady probes mapexIam /health. A green response confirms the
// service is up and its dependencies are healthy. Authentication is
// verified later by the saga's SeedAdminLogin step, so checkStackReady
// stays read-only and side-effect-free against the live stack.
func checkStackReady() error {
	client := httpclient.New(httpclient.Config{
		BaseURL: constants.MapexosURL,
		Timeout: 5 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	resp, err := client.Raw(ctx, http.MethodGet, "/health", nil)
	if err != nil {
		return fmt.Errorf("e2e readiness: mapexIam health check failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("e2e readiness: unexpected status %d from /health", resp.StatusCode)
	}
	return nil
}
