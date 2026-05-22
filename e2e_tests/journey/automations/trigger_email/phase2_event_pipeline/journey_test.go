package phase2_event_pipeline

import "testing"

// TestJourney runs the Email trigger event-pipeline smoke against
// the live stack.
func TestJourney(t *testing.T) {
	Run(t)
}
