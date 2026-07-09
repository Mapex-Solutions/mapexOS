//go:build saga

package phase2_sensor_uplink

import "testing"

// TestJourney provisions a LoRaWAN gateway + sensor, drives a real uplink through
// mapexLNS over UDP and Basics Station, and verifies MapexOS ingests it (raw
// event bytes + metadata) and marks the sensor online. Requires the live stack
// (assets + mapexIam + router + events + mapexLNS) per the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
