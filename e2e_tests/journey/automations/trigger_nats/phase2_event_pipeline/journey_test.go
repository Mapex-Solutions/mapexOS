package phase2_event_pipeline

import "testing"

// TestJourney runs the NATS trigger event-pipeline smoke against the
// live stack.
func TestJourney(t *testing.T) {
	Run(t)
}
