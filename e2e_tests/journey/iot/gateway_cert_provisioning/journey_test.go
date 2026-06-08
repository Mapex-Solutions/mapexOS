package gateway_cert_provisioning

import "testing"

// TestJourney provisions a cert-mode LoRaWAN gateway asset, issues its mTLS cert,
// and verifies the issued serial is reflected on the internal asset-auth endpoint
// (the LNS Gateway Server's L3 source). Requires the live stack (assets + mapexIam
// + mapexVault PKI) per the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
