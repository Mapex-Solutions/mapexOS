//go:build saga

package phase3_sensor_presence

import "testing"

// TestJourney provisions a LoRaWAN gateway + sensor, drives an uplink over UDP
// and Basics Station to bring the sensor online from data, then forces it
// offline — asserting the sensor presence lifecycle. Requires the live stack
// (assets + mapexIam + router + mapexLNS) per the journey README.
func TestJourney(t *testing.T) {
	Run(t)
}
