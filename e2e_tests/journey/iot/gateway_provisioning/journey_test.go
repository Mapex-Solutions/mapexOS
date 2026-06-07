package gateway_provisioning

import "testing"

// TestJourney provisions a LoRaWAN gateway asset and verifies it is published as
// a gateway on the internal asset-auth endpoint. Requires the live stack (assets
// + mapexIam) per the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
