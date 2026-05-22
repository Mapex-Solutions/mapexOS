package ports

import "mapexVault/src/modules/credentials/domain/entities"

// Port-level type aliases — expose domain entities through the port boundary.
type Credential = entities.Credential
type CredentialType = entities.CredentialType
type CredentialStatus = entities.CredentialStatus
type ProviderConfig = entities.ProviderConfig
type Connection = entities.Connection
type ConnectionStatus = entities.ConnectionStatus

// Re-export constants
const (
	CredentialManual      = entities.CredentialManual
	CredentialOAuth2      = entities.CredentialOAuth2
	CredentialUserAndPass = entities.CredentialUserAndPass

	CredentialStatusActive  = entities.CredentialStatusActive
	CredentialStatusExpired = entities.CredentialStatusExpired
	CredentialStatusRevoked = entities.CredentialStatusRevoked
	CredentialStatusError   = entities.CredentialStatusError

	ConnectionStatusActive  = entities.ConnectionStatusActive
	ConnectionStatusRevoked = entities.ConnectionStatusRevoked
	ConnectionStatusExpired = entities.ConnectionStatusExpired
)
