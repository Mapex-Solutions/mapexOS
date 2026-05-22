//go:build saga

package phase2_cert_user

import "testing"

// TestPhase2_CertUser_Saga validates the cert auth path end-to-end
// including revocation enforcement.
//
// Outcome on PASS: all phase-2 items complete; the revoked-cert
// reconnect fails as expected.
// Outcome on FAIL: the saga rollback log identifies the offending step.
func TestPhase2_CertUser_Saga(t *testing.T) {
	Run(t)
}
