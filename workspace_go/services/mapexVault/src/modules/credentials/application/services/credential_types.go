package services

import (
	"time"

	"mapexVault/src/modules/credentials/application/di"
)

// CredentialService orchestrates credential management: CRUD, encryption,
// OAuth2 callbacks, userAndPass sessions, and NATS event publishing.
type CredentialService struct {
	deps di.CredentialServiceDependenciesInjection
}

// tokenResponse holds the extracted tokens + raw response from a provider's
// OAuth2/login endpoint after executing a TokenRequestConfig.
type tokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    *time.Time
	RawResponse  map[string]interface{}
}
