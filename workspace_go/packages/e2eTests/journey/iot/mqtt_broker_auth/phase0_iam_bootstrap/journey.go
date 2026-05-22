// Package phase0 holds the IAM bootstrap stage of the mqtt_broker_auth
// journey.
//
// The phase authenticates as the seed admin user (admin@mapex.local) and
// preconfigures every per-service HTTP client in the saga ClientSet with
// the resulting bearer plus the seed root org as X-Org-Context. After the
// phase finishes, every downstream phase can drive any service endpoint
// as the seeded admin without re-issuing logins or maintaining its own
// context.
//
// Phase 0 leaves the bag populated with:
//
//	iam.userJWT             bearer attached to every HTTP client
//	iam.organizationID      MapexosOrgID — the seed root used as saga org
//
// Phase 0 does NOT create a new tenant org, role, or user. The seed admin
// already carries the wildcard "mapex.*" permission anchored at the seed
// root (MapexosOrgID); reusing that anchor keeps the broker-auth saga
// focused on the MQTT pipeline and avoids the permission-cache rebuild
// path that needs work on the mapexIam side before brand-new child orgs
// can be driven by the bootstrap actor end to end. Tenant-isolation
// tests belong in their own journey (planned: iot/multi_tenant_isolation/).
//
// The login itself is the building-block step
// services/mapexIam/auth/steps.SeedAdminLogin — every journey reuses it
// so the auth flow lives in one place.
//
// Outcome (what passing this phase proves):
//   - The seed admin user (admin@mapex.local) can sign in via /auth/login.
//   - The login response still carries the access_token / refresh_token
//     envelope every downstream client expects.
//   - The JWT carries access to MapexosOrgID through the coverage
//     endpoint; the bootstrap actor is authorised against the seed root.
//   - Every per-service HTTP client in the ClientSet gets the bearer and
//     X-Org-Context propagated; subsequent phases need not touch headers.
//
// Failure typically points at one of:
//   - mongodb-init seed missing or stale (no admin user, no seed root).
//   - login response envelope changed (access_token field renamed).
//   - coverage build job failed on login; the user exists but
//     AssertJwtHasOrgContext sees zero accessible orgs.
package phase0

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	authAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/auth/asserts"
	authSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/auth/steps"
)

// BootstrapItems returns the ordered Items that authenticate as the seed
// admin and publish the (jwt, orgID) pair every downstream phase reads
// from the bag. Composing these into a parent saga.Run lets downstream
// phases append their own steps under the same rollback chain.
func BootstrapItems() []saga.Item {
	return []saga.Item{
		authSteps.SeedAdminLogin(),
		authAsserts.AssertJwtValid(),
		authAsserts.AssertJwtHasOrgContext(),
	}
}

// NewClients constructs the per-service ClientSet wired against the URLs
// the e2eTests common/constants module exposes. The clients are returned
// without authentication — the saga's first step (authSteps.SeedAdminLogin)
// publishes the bearer and propagates it across every client in the set.
func NewClients() saga.ClientSet {
	return saga.NewClientSet(saga.ClientURLs{
		MapexIam: constants.MapexosURL,
		Assets:   constants.AssetsURL,
		Router:   constants.RouterURL,
		Gateway:  constants.GatewayURL,
		Events:   constants.EventsURL,
		Triggers: constants.TriggersURL,
		Workflow: constants.WorkflowURL,
	})
}

// Run executes Phase 0 as a stand-alone saga. PhaseN+1 packages do not
// invoke Run; they compose BootstrapItems into their own saga.Run call
// instead.
func Run(t *testing.T) {
	t.Helper()
	if err := utils.SetupE2EEnvironment(); err != nil {
		t.Fatalf("setup e2e environment: %v", err)
	}
	runID := random.NewRunID()
	saga.Run(t, context.Background(), runID, NewClients(), BootstrapItems()...)
}
