package services

import (
	"time"

	mqttDi "assets/src/modules/mqttcerts/application/di"
)

// MqttCertsService implements MqttCertsServicePort + common.Mountable.
type MqttCertsService struct {
	deps mqttDi.MqttCertsServiceDI
}

type signedCertBundle struct {
	certPEM     []byte
	keyPEM      []byte
	serialHex   string
	fingerprint string
	subjectCN   string
	issuedAt    time.Time
	expiresAt   time.Time
}
