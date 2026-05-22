//go:build saga

package phase0

import "testing"

// TestPhase0_IAMBootstrap_Saga executes Phase 0 of the mqtt_broker_auth
// journey against a running stack. It is gated by the saga build tag so
// the standard unit suite (go test ./...) ignores it; CI invokes it via
// go test -tags=saga.
//
// Required environment:
//   - mapexIam reachable at MAPEXOS_URL (default http://localhost:5000)
//   - root user seeded (root@mapex.global / mapex123)
//   - mongodb-init seed already applied (provides RootRoleID and the seed
//     organization referenced by SagaTestCustomerOrg as parent)
//
// Outcome on PASS:
//   - One scratch customer org created under the seed root and torn down.
//   - One IoT-admin role created inside that org with full IoT-pipeline
//     permissions (asset/template/route group CRUD + IAM read paths).
//   - One user provisioned via the onboarding orchestrator alongside a
//     fresh group bound to the role; user is enabled and can sign in.
//   - Login returns a structurally valid 3-segment JWT carrying access
//     to the new org through the coverage endpoint.
//   - Live stack ends with zero saga-tagged rows: org cascade-delete
//     removes role, group, membership, and user.
//
// Outcome on FAIL:
//   - Logs identify which step or assertion failed by qualified name
//     (e.g. mapexIam/onboarding.CreateUserWithMemberships) so the
//     failure points at the breaking integration without source diving.
func TestPhase0_IAMBootstrap_Saga(t *testing.T) {
	Run(t)
}
