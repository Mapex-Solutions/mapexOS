package constants

import "time"

const (
	// LorawanKEKContext is the mapexVault KEK category for LoRaWAN device keys.
	// It must match the context seeded by the mongodb-init container
	// (kek-bootstrap.sh).
	LorawanKEKContext = "lorawan_device_keys"

	// Bootstrap timing for the OnMount KEK fetch, mirroring the mqttcerts CA
	// bootstrap: one short sync attempt, then a retry goroutine with exponential
	// backoff and no max attempts.
	LorawanKEKBootstrapInitialTimeout = 5 * time.Second
	LorawanKEKBootstrapBackoffMin     = 1 * time.Second
	LorawanKEKBootstrapBackoffMax     = 30 * time.Second
)
