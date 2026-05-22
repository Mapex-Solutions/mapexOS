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
// at tokenExpiresAt - 15 minutes. Called after every successful token acquisition.
func (s *CredentialService) publishRefreshSchedule(credentialId string, credentialType entities.CredentialType, tokenExpiresAt *time.Time) {
	if tokenExpiresAt == nil {
		return
	}

	refreshBuffer := time.Duration(constants.RefreshBufferMinutes) * time.Minute
	scheduleAt := tokenExpiresAt.Add(-refreshBuffer)
	if scheduleAt.Before(time.Now()) {
		return
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

/**
 * Bootstrap Seed
 */

// bootstrapSeed queries existing credentials and publishes initial schedules on startup.
// For credentials with tokenExpiresAt in the future: schedule at tokenExpiresAt - 15min.
// For credentials already expired: schedule at now + 30s (near-immediate refresh).
//
// Called by OnMount lifecycle hook after the service is fully wired.
func (s *CredentialService) bootstrapSeed() {
	credentials, err := s.deps.CredentialRepo.FindActiveWithTokenExpiry(context.Background())
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Credential] Bootstrap seed failed to query credentials: %v", err))
		return
	}

	if len(credentials) == 0 {
		logger.Info("[SERVICE:Credential] Bootstrap seed: no credentials with token expiry found")
		return
	}

	refreshBuffer := time.Duration(constants.RefreshBufferMinutes) * time.Minute
	now := time.Now()
	count := 0

	for _, cred := range credentials {
		if cred.TokenExpiresAt == nil {
			continue
		}

		scheduleAt := cred.TokenExpiresAt.Add(-refreshBuffer)
		if scheduleAt.Before(now) {
			scheduleAt = now.Add(30 * time.Second)
		}

		subject := fmt.Sprintf("%s.%s", constants.VaultScheduleSubjectPrefix, cred.ID.Hex())

		if err := s.deps.ScheduleManager.PurgeStreamSubject(constants.VaultScheduleStreamName, subject); err != nil {
			logger.Warn(fmt.Sprintf("[SERVICE:Credential] Bootstrap seed: failed to purge schedule for %s: %v", cred.ID.Hex(), err))
		}

		if err := s.deps.ScheduleManager.PublishScheduled(natsModel.ScheduledPublishConfig{
			Subject:       subject,
			TargetSubject: constants.VaultScheduleFiredSubject,
			ScheduleAt:    scheduleAt,
			Data: map[string]interface{}{
				"credentialId":   cred.ID.Hex(),
				"credentialType": string(cred.Type),
			},
		}); err != nil {
			logger.Warn(fmt.Sprintf("[SERVICE:Credential] Bootstrap seed: failed to publish schedule for %s: %v", cred.ID.Hex(), err))
			continue
		}
		count++
	}

	logger.Info(fmt.Sprintf("[SERVICE:Credential] Bootstrap seed: published %d schedules for existing credentials", count))
}
