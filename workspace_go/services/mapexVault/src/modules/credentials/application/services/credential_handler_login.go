package services

import (
	"context"
	"fmt"

	"mapexVault/src/modules/credentials/domain/entities"
)

// loadCredentialForLogin fetches the credential and validates it has a
// LoginConfig configured; the orchestration in HandleLogin assumes both
// preconditions have been checked when this returns nil error.
func (s *CredentialService) loadCredentialForLogin(ctx context.Context, credentialId string) (*entities.Credential, error) {
	cred, err := s.deps.CredentialRepo.FindById(ctx, &credentialId)
	if err != nil || cred == nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Credential %s not found for login", credentialId)
	}
	if cred.ProviderConfig == nil || cred.ProviderConfig.LoginConfig == nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Credential %s has no login config", credentialId)
	}
	return cred, nil
}

// loadAndDecryptForLogin combines the load + decrypt steps so HandleLogin
// can express the full flow in 5-15 lines.
func (s *CredentialService) loadAndDecryptForLogin(ctx context.Context, credentialId string) (*entities.Credential, map[string]interface{}, error) {
	cred, err := s.loadCredentialForLogin(ctx, credentialId)
	if err != nil {
		return nil, nil, err
	}
	data, err := decryptData(s.deps.Encryption, cred)
	if err != nil {
		return nil, nil, fmt.Errorf("[SERVICE:Credential] Failed to decrypt for login: %w", err)
	}
	return cred, data, nil
}
