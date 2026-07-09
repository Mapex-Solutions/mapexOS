//go:build saga

package phase2_sensor_presence

import "testing"

// TestJourney provisions a LoRaWAN gateway + ABP sensor, drives an uplink over UDP
// and Basics Station (no OTAA join — ABP fixed session) to bring the sensor online
// from data, then forces it offline — asserting the ABP sensor presence lifecycle.
// Requires the live stack (assets + mapexIam + router + mapexLNS) per the journey
// README.
func TestJourney(t *testing.T) {
	Run(t)
}
