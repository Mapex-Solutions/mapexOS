package phase1_connectivity

import "testing"

// TestJourney runs the RabbitMQ trigger connectivity smoke against
// the live stack. First run pulls the rabbitmq image (~30 s).
func TestJourney(t *testing.T) {
	Run(t)
}
