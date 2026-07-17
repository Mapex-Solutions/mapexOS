package services

import (
	"context"
	"fmt"
	"time"

	"mapexVault/src/modules/credentials/application/constants"
	"mapexVault/src/modules/credentials/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * Schedule Publish
 */

// publishRefreshSchedule publishes a NATS scheduled message to refresh a credential
// at tokenExpiresAt - 15 minutes. Called after every successful token acquisition and
// by the reconcile reseed.
func (s *CredentialService) publishRefreshSchedule(credentialId string, credentialType entities.CredentialType, tokenExpiresAt *time.Time) {
	if tokenExpiresAt == nil {
		return
	}

	refreshBuffer := time.Duration(constants.RefreshBufferMinutes) * time.Minute
	scheduleAt := tokenExpiresAt.Add(-refreshBuffer)
	if scheduleAt.Before(time.Now()) {
		// The refresh time is already past (e.g. the access token expired while the
		// service was down). Refresh near-immediately instead of dropping it — the
		// refresh token is still valid, so this recovers without any user re-login.
		scheduleAt = time.Now().Add(30 * time.Second)
	}

	subject := fmt.Sprintf("%s.%s", constants.VaultScheduleSubjectPrefix, credentialId)

	if err := s.deps.ScheduleManager.PurgeStreamSubject(constants.VaultScheduleStreamName, subject); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Credential] Failed to purge existing schedule for %s: %v", credentialId, err))
	}

	if err := s.deps.ScheduleManager.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       subject,
		TargetSubject: constants.VaultScheduleFiredSubject,
		ScheduleAt:    scheduleAt,
		Data: map[string]interface{}{
			"credentialId":   credentialId,
			"credentialType": string(credentialType),
		},
	}); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Credential] Failed to publish refresh schedule for %s: %v", credentialId, err))
	}

	logger.Info(fmt.Sprintf("[SERVICE:Credential] Scheduled refresh for %s at %s", credentialId, scheduleAt.UTC().Format(time.RFC3339)))
}

/**
 * Error Handling
 */

// markCredentialError sets credential status to error and publishes error event.
func (s *CredentialService) markCredentialError(cred *entities.Credential, refreshErr error) {
	id := cred.ID.Hex()
	_, _ = s.deps.CredentialRepo.FindByIdAndUpdate(context.Background(), &id, model.Map{
		"refreshError": refreshErr.Error(),
		"status":       string(entities.CredentialStatusError),
		"updated":      time.Now(),
	})
	s.publishVaultEvent(id, "error")
}
