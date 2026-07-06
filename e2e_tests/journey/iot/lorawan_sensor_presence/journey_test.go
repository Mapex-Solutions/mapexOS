//go:build saga

package lorawan_sensor_presence

import "testing"

// TestJourney provisions a LoRaWAN sensor asset, drives it online with a real
// uplink over UDP and Basics Station, then forces it offline. Requires the live
// stack (assets + mapexIam + router + events + mapexLNS) per the README.
func TestJourney(t *testing.T) {
	Run(t)
}
