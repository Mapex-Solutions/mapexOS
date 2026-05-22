//go:build saga

package phase3_cascade

import "testing"

// TestPhase3_Cascade_Saga validates the broker's L1→L2→L3 cascade plus
// the fanout invalidate lazy-pull path.
//
// Outcome on PASS: all forced-miss + log-grep assertions complete
// inside their timeouts.
// Outcome on FAIL: the saga rollback log + the broker's tail logs are
// useful diagnostic outputs.
func TestPhase3_Cascade_Saga(t *testing.T) {
	Run(t)
}
