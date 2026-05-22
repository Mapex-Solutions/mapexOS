package services

import (
	"encoding/json"
	"fmt"

	"mapexVault/src/modules/credentials/application/dtos"
	"mapexVault/src/modules/credentials/domain/entities"

	envelope "github.com/Mapex-Solutions/mapexGoKit/utils/envelope"
)

// encryptData encrypts the credential data using envelope encryption.
func encryptData(encryption *envelope.EnvelopeService, data map[string]interface{}) (*envelope.EncryptedEnvelope, error) {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to marshal credential data: %w", err)
	}

	env, err := encryption.Encrypt(plaintext)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to encrypt credential data: %w", err)
	}

	return env, nil
}

// decryptData decrypts the credential data from envelope encryption.
func decryptData(encryption *envelope.EnvelopeService, cred *entities.Credential) (map[string]interface{}, error) {
	env := &envelope.EncryptedEnvelope{
		EncryptedDEK:  cred.EncryptedDEK,
		DEKNonce:      cred.DEKNonce,
		EncryptedData: cred.EncryptedData,
		DataNonce:     cred.DataNonce,
	}

	plaintext, err := encryption.Decrypt(env)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to decrypt credential data: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to unmarshal decrypted data: %w", err)
	}

	return data, nil
}

// toCredentialResponse converts a credential entity to a response DTO.
func toCredentialResponse(cred *entities.Credential) *dtos.CredentialResponse {
	return &dtos.CredentialResponse{
		ID:              cred.ID,
		Name:            cred.Name,
		Type:            cred.Type,
		PluginId:        cred.PluginId,
		CredentialDefId: cred.CredentialDefId,
		OrgId:           cred.OrgId,
		PathKey:         cred.PathKey,
		IsTemplate:      cred.IsTemplate,
		Status:          cred.Status,
		TokenExpiresAt:  cred.TokenExpiresAt,
		LastRefreshedAt: cred.LastRefreshedAt,
		RefreshError:    cred.RefreshError,
		ProviderConfig:  cred.ProviderConfig,
		Created:         cred.Created,
		Updated:         cred.Updated,
	}
}

// toConnectionResponse converts a connection entity to a response DTO.
func toConnectionResponse(conn *entities.Connection) *dtos.ConnectionResponse {
	return &dtos.ConnectionResponse{
		ID:           conn.ID,
		Provider:     conn.Provider,
		AccountId:    conn.AccountId,
		AccountName:  conn.AccountName,
		Status:       string(conn.Status),
		CredentialId: conn.CredentialId,
		UserId:       conn.UserId,
		OrgId:        conn.OrgId,
		PathKey:      conn.PathKey,
		Scopes:       conn.Scopes,
		ConnectedAt:  conn.ConnectedAt,
		LastUsedAt:   conn.LastUsedAt,
		Created:      conn.Created,
		Updated:      conn.Updated,
	}
}
