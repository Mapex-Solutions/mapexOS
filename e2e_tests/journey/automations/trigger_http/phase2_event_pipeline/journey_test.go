package phase2_event_pipeline

import "testing"

// TestJourney runs the HTTP trigger event-pipeline smoke against the
// live stack.
func TestJourney(t *testing.T) {
	Run(t)
}
