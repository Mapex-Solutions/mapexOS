package phase1_connectivity

import "testing"

// TestJourney runs the Email trigger connectivity smoke end-to-end
// against the live stack. See package godoc for the PASS / FAIL
// outcome contract.
func TestJourney(t *testing.T) {
	Run(t)
}
