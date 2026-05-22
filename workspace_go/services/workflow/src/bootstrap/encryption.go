package bootstrap

import (
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	envelope "github.com/Mapex-Solutions/mapexGoKit/utils/envelope"
)

// InitEncryption registers the EnvelopeService (AES-256-GCM envelope encryption)
// in the DIG container. Used by the credentials module for encrypting/decrypting
// sensitive credential data (bot tokens, API keys, OAuth secrets).
func InitEncryption(c *dig.Container) {
	masterKeyHex, _ := config.GetStringValue("credential_master_key")

	c.Provide(func() *envelope.EnvelopeService {
		svc, err := envelope.New(masterKeyHex)
		if err != nil {
			logger.Panic("Failed to initialize envelope encryption: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] Envelope encryption initialized")
		return svc
	})
}
