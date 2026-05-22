package phase1_connectivity

import "testing"

// TestJourney runs the NATS trigger connectivity smoke against the
// live stack.
func TestJourney(t *testing.T) {
	Run(t)
}
