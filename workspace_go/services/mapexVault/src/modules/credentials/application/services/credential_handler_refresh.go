package services

import (
	"context"
	"encoding/json"
	"fmt"

	"mapexVault/src/modules/credentials/domain/entities"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// parseRefreshPayload decodes the scheduled refresh message into the
// per-credential payload. Returns ok=false on invalid JSON (msg already
// Acked) so the caller can short-circuit.
func (s *CredentialService) parseRefreshPayload(msg *natsModel.Message) (string, bool) {
	var payload struct {
		CredentialId   string `json:"credentialId"`
		CredentialType string `json:"credentialType"`
	}
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Credential] Failed to unmarshal refresh message: %v", err))
		msg.Ack()
		return "", false
	}
	return payload.CredentialId, true
}

// loadCredentialForRefresh fetches the credential and filters out anything
// not active. Returns ok=false (msg already Acked) when the schedule
// should be silently dropped.
func (s *CredentialService) loadCredentialForRefresh(credentialId string, msg *natsModel.Message) (*entities.Credential, bool) {
	cred, err := s.deps.CredentialRepo.FindById(context.Background(), &credentialId)
	if err != nil || cred == nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Credential] Credential %s not found for scheduled refresh, skipping", credentialId))
		msg.Ack()
		return nil, false
	}
	if cred.Status != entities.CredentialStatusActive {
		logger.Info(fmt.Sprintf("[SERVICE:Credential] Credential %s is not active (status=%s), skipping refresh", credentialId, cred.Status))
		msg.Ack()
		return nil, false
	}
	return cred, true
}

// selectTokenRefreshConfig picks the RefreshConfig (preferred) or falls
// back to the LoginConfig when refresh is not configured. Returns
// ok=false (msg already Acked) when neither is present.
func (s *CredentialService) selectTokenRefreshConfig(cred *entities.Credential, credentialId string, msg *natsModel.Message) (*entities.TokenRequestConfig, bool) {
	if cred.ProviderConfig != nil && cred.ProviderConfig.RefreshConfig != nil {
		return cred.ProviderConfig.RefreshConfig, true
	}
	if cred.ProviderConfig != nil && cred.ProviderConfig.LoginConfig != nil {
		return cred.ProviderConfig.LoginConfig, true
	}
	logger.Warn(fmt.Sprintf("[SERVICE:Credential] Credential %s has no refresh or login config, skipping", credentialId))
	msg.Ack()
	return nil, false
}

// runTokenRefresh decrypts -> calls the provider -> persists the new
// tokens. Errors at any step mark the credential and Ack so the schedule
// stream does not retry a misconfigured credential indefinitely.
func (s *CredentialService) runTokenRefresh(cred *entities.Credential, cfg *entities.TokenRequestConfig, credentialId string, msg *natsModel.Message) {
	data, err := decryptData(s.deps.Encryption, cred)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Credential] Failed to decrypt credential %s for refresh", credentialId))
		s.markCredentialError(cred, err)
		msg.Ack()
		return
	}
	resp, err := s.executeTokenRequest(cred, cfg, data)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Credential] Scheduled refresh failed for %s", credentialId))
		s.markCredentialError(cred, err)
		msg.Ack()
		return
	}
	if err := s.updateCredentialTokens(cred, data, resp); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Credential] Failed to update tokens for %s after refresh", credentialId))
		s.markCredentialError(cred, err)
		msg.Ack()
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:Credential] Scheduled refresh completed for %s", credentialId))
	msg.Ack()
}
