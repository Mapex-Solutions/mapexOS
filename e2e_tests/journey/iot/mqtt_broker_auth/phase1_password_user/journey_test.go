//go:build saga

package phase1_password_user

import "testing"

// TestPhase1_PasswordUser_Saga runs the full password auth phase against
// the live services_required stack.
//
// Outcome on PASS:
//   - All steps complete (IAM, asset, connect, presence-connect-asserted,
//     publish, data-asserted, disconnect, presence-disconnect-asserted).
//
// Outcome on FAIL:
//   - The failing step's name appears in the saga rollback log so the
//     offending integration point is unambiguous.
func TestPhase1_PasswordUser_Saga(t *testing.T) {
	Run(t)
}
