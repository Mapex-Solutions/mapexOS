//go:build saga

package gateway_key_provisioning

import "testing"

// TestJourney provisions a key-mode LoRaWAN gateway asset and verifies the
// supplied token is hashed and reflected on the internal asset-auth endpoint
// (the LNS Gateway Server's L3 source). Requires the live stack (assets +
// mapexIam) per the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
