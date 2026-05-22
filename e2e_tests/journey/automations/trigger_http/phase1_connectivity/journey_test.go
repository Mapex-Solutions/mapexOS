package phase1_connectivity

import "testing"

// TestJourney runs the HTTP trigger connectivity smoke end-to-end
// against the live stack.
func TestJourney(t *testing.T) {
	Run(t)
}
