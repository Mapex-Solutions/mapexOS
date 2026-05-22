package services

import (
	"fmt"
	"time"

	"mapexVault/src/modules/credentials/application/constants"
	"mapexVault/src/modules/credentials/domain/entities"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// checkAndReseed returns true when a missing timer was republished.
// For each credential it checks if a pending schedule exists on
// VAULT-SCHEDULE under subject vault.schedule.{credentialId}; if not, it
// republishes via the existing publishRefreshSchedule helper.
func (s *CredentialService) checkAndReseed(cred *entities.Credential) bool {
	if cred.TokenExpiresAt == nil {
		return false
	}
	credentialId := cred.ID.Hex()
	subject := fmt.Sprintf("%s.%s", constants.VaultScheduleSubjectPrefix, credentialId)
	pending, err := s.deps.ScheduleManager.HasPendingMessages(constants.VaultScheduleStreamName, subject)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Credential] Reconciler failed to check pending schedule for %s: %v", credentialId, err))
		return false
	}
	if pending {
		return false
	}
	s.publishRefreshSchedule(credentialId, cred.Type, cred.TokenExpiresAt)
	logger.Info(fmt.Sprintf("[SERVICE:Credential] Reconciler reseeded missing schedule for %s", credentialId))
	return true
}

// reseedMissingSchedules iterates active credentials and reseeds any whose
// refresh timer is missing. Returns (checked, reseeded) counts for the log.
func (s *CredentialService) reseedMissingSchedules(credentials []entities.Credential) (int, int) {
	checked := len(credentials)
	reseeded := 0
	for i := range credentials {
		if s.checkAndReseed(&credentials[i]) {
			reseeded++
		}
	}
	return checked, reseeded
}

// scheduleNextReconcile publishes the next reconcile timer if none is
// pending. Layer 1: HasPendingMessages avoids double-publishing across
// pods. Layer 2: stream Duplicates=10s + fixed MsgId catches the race.
func (s *CredentialService) scheduleNextReconcile() {
	pending, err := s.deps.ScheduleManager.HasPendingMessages(
		constants.VaultReconcilerStreamName,
		constants.VaultReconcileScheduleSubject,
	)
	if err != nil {
		logger.Error(err, "[SERVICE:Credential] Failed to check pending reconcile messages")
		return
	}
	if pending {
		logger.Debug("[SERVICE:Credential] Reconcile already scheduled, skipping")
		return
	}
	interval, _ := config.GetIntValue("vault_reconcile_interval")
	if interval <= 0 {
		interval = constants.VaultReconcileDefaultIntervalSeconds
	}
	scheduleAt := time.Now().Add(time.Duration(interval) * time.Second)
	if err := s.deps.ScheduleManager.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       constants.VaultReconcileScheduleSubject,
		TargetSubject: constants.VaultReconcileFiredSubject,
		ScheduleAt:    scheduleAt,
		Data:          map[string]string{"trigger": "scheduled"},
		MsgId:         constants.VaultReconcileMsgId,
	}); err != nil {
		logger.Error(err, "[SERVICE:Credential] Failed to schedule next reconcile")
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:Credential] Next reconcile scheduled at %s", scheduleAt.Format(time.RFC3339)))
}
