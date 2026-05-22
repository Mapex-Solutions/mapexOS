package phase1_workflow

import "testing"

// TestJourney runs the phase end-to-end against the live stack.
// Requires the standalone compose to be up and the IAM seed to have
// completed; see the journey README for the full setup checklist.
func TestJourney(t *testing.T) {
	Run(t)
}
