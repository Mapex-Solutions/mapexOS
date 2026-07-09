//go:build saga

package phase1_gateway_connectivity

import "testing"

// TestJourney provisions a LoRaWAN gateway (UDP + Basics Station), connects it
// to mapexLNS to go online, and forces it offline — asserting the presence
// lifecycle. Requires the live stack (assets + mapexIam + mapexLNS) per the
// journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
