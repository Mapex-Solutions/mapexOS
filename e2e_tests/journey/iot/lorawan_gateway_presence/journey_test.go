//go:build saga

package lorawan_gateway_presence

import "testing"

// TestJourney provisions a LoRaWAN gateway asset and verifies it goes online on a
// real LNS connect and offline after a forced transition, over UDP and Basics
// Station. Requires the live stack (assets + mapexIam + mapexLNS) per the README.
func TestJourney(t *testing.T) {
	Run(t)
}
