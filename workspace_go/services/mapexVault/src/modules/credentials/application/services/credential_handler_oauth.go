package services

import (
	"fmt"
	"time"

	"mapexVault/src/modules/credentials/application/constants"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// publishVaultEvent publishes a credential event to the MAPEX-VAULT stream.
func (s *CredentialService) publishVaultEvent(credentialId string, action string) {
	event := map[string]interface{}{
		"credentialId": credentialId,
		"action":       action,
		"timestamp":    time.Now().Unix(),
	}

	if err := s.deps.Publisher.Publish(natsModel.PublishConfig{
		Subject: constants.VaultEventsSubject,
		Data:    event,
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Credential] Failed to publish vault event %s for %s", action, credentialId))
	}
}
