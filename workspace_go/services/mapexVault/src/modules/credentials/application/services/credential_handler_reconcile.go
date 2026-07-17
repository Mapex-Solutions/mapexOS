package services

import (
	"fmt"

	"mapexVault/src/modules/credentials/application/constants"
	"mapexVault/src/modules/credentials/domain/entities"

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
