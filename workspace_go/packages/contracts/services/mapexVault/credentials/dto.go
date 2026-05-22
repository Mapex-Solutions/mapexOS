package credentials

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/**
 * Credential Types
 *
 * Cross-service contract types. Authoritative home for credential type discriminators,
 * status enums, and provider configuration shapes. Services alias these via
 * application/dtos and domain/entities packages to keep a single source of truth.
 */

// CredentialType is the discriminator for how the vault manages the credential.
type CredentialType string

const (
	CredentialManual      CredentialType = "manual"
	CredentialOAuth2      CredentialType = "oauth2"
	CredentialUserAndPass CredentialType = "userAndPass"
)

// CredentialStatus is the lifecycle status of a credential.
type CredentialStatus string

const (
	CredentialStatusActive  CredentialStatus = "active"
	CredentialStatusExpired CredentialStatus = "expired"
	CredentialStatusRevoked CredentialStatus = "revoked"
	CredentialStatusError   CredentialStatus = "error"
)

// ConnectionStatus is the lifecycle status of a provider connection.
type ConnectionStatus string

const (
	ConnectionStatusActive  ConnectionStatus = "active"
	ConnectionStatusRevoked ConnectionStatus = "revoked"
	ConnectionStatusExpired ConnectionStatus = "expired"
)

/**
 * TokenRequestConfig
 *
 * Fully describes one HTTP request for token acquisition or renewal. Stored
 * unencrypted in ProviderConfig — contains no secrets, only templates.
 * Template placeholders (e.g. {{credential.refreshToken}}) are resolved at
 * execution time using the decrypted credential data map.
 */
type TokenRequestConfig struct {
	// HTTP request
	Method      string                 `bson:"method,omitempty" json:"method,omitempty"`
	Url         string                 `bson:"url" json:"url"`
	ContentType string                 `bson:"contentType,omitempty" json:"contentType,omitempty"`
	Headers     map[string]string      `bson:"headers,omitempty" json:"headers,omitempty"`
	Body        map[string]interface{} `bson:"body,omitempty" json:"body,omitempty"`
	QueryParams map[string]string      `bson:"queryParams,omitempty" json:"queryParams,omitempty"`

	// Response extraction (dot-notation paths)
	AccessTokenPath  string `bson:"accessTokenPath,omitempty" json:"accessTokenPath,omitempty"`
	RefreshTokenPath string `bson:"refreshTokenPath,omitempty" json:"refreshTokenPath,omitempty"`
	ExpiresInPath    string `bson:"expiresInPath,omitempty" json:"expiresInPath,omitempty"`
}

/**
 * ProviderConfig
 *
 * Configuration for credential token lifecycle management. Defines how the
 * vault acquires and renews tokens via HTTP. Stored unencrypted because the
 * refresh consumer needs to read these fields to know WHERE and HOW to
 * refresh, without decrypting secrets.
 */
type ProviderConfig struct {
	LoginConfig   *TokenRequestConfig `bson:"loginConfig,omitempty" json:"loginConfig,omitempty"`
	RefreshConfig *TokenRequestConfig `bson:"refreshConfig,omitempty" json:"refreshConfig,omitempty"`
}

/**
 * Credential DTOs
 */

// CreateCredentialDTO is the API contract for creating a credential.
type CreateCredentialDTO struct {
	Name            string                 `json:"name" validate:"required"`
	Type            CredentialType         `json:"type" validate:"required,oneof=manual oauth2 userAndPass"`
	PluginId        string                 `json:"pluginId" validate:"required"`
	CredentialDefId string                 `json:"credentialDefId"`
	Data            map[string]interface{} `json:"data" validate:"required"`
	ProviderConfig  *ProviderConfig        `json:"providerConfig,omitempty"`
	IsTemplate      bool                   `json:"isTemplate"`
}

// UpdateCredentialDTO is the API contract for updating a credential.
type UpdateCredentialDTO struct {
	Name           *string                `json:"name,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
	ProviderConfig *ProviderConfig        `json:"providerConfig,omitempty"`
	IsTemplate     *bool                  `json:"isTemplate,omitempty"`
}

// CredentialResponse is the API response for a credential (no encrypted fields).
type CredentialResponse struct {
	ID              model.ObjectId   `json:"id"`
	Name            string           `json:"name"`
	Type            CredentialType   `json:"type"`
	PluginId        string           `json:"pluginId"`
	CredentialDefId string           `json:"credentialDefId"`
	OrgId           *model.ObjectId  `json:"orgId,omitempty"`
	PathKey         string           `json:"pathKey"`
	IsTemplate      bool             `json:"isTemplate"`
	Status          CredentialStatus `json:"status"`
	TokenExpiresAt  *time.Time       `json:"tokenExpiresAt,omitempty"`
	LastRefreshedAt *time.Time       `json:"lastRefreshedAt,omitempty"`
	RefreshError    string           `json:"refreshError,omitempty"`
	ProviderConfig  *ProviderConfig  `json:"providerConfig,omitempty"`
	Created         time.Time        `json:"created"`
	Updated         time.Time        `json:"updated"`
}

// CredentialQueryDTO is the query contract for listing credentials.
type CredentialQueryDTO struct {
	PluginId *string `query:"pluginId" validate:"omitempty"`
	Type     *string `query:"type" validate:"omitempty"`
	Status   *string `query:"status" validate:"omitempty"`
	Page     *int    `query:"page" validate:"omitempty,min=1"`
	PerPage  *int    `query:"perPage" validate:"omitempty,min=1,max=100"`
}

/**
 * OAuth2 DTOs
 */

// OAuthCallbackDTO is the input for OAuth2 callback processing.
type OAuthCallbackDTO struct {
	Provider        string `json:"provider" validate:"required"`
	Code            string `json:"code" validate:"required"`
	RedirectUri     string `json:"redirectUri" validate:"required"`
	PluginId        string `json:"pluginId" validate:"required"`
	CredentialDefId string `json:"credentialDefId"`
}

/**
 * Connection DTOs
 */

// CreateConnectionDTO is the input for creating a connection.
type CreateConnectionDTO struct {
	Provider     string   `json:"provider" validate:"required"`
	AccountId    string   `json:"accountId" validate:"required"`
	AccountName  string   `json:"accountName"`
	CredentialId string   `json:"credentialId" validate:"required"`
	Scopes       []string `json:"scopes"`
}

// UpsertConnectionDTO is the input for upserting a connection by account.
type UpsertConnectionDTO struct {
	Provider     string   `json:"provider" validate:"required"`
	AccountId    string   `json:"accountId" validate:"required"`
	AccountName  string   `json:"accountName"`
	CredentialId string   `json:"credentialId" validate:"required"`
	Scopes       []string `json:"scopes"`
}

// ConnectionResponse is the API response for a connection.
type ConnectionResponse struct {
	ID           model.ObjectId  `json:"id"`
	Provider     string          `json:"provider"`
	AccountId    string          `json:"accountId"`
	AccountName  string          `json:"accountName"`
	Status       string          `json:"status"`
	CredentialId model.ObjectId  `json:"credentialId"`
	UserId       model.ObjectId  `json:"userId"`
	OrgId        *model.ObjectId `json:"orgId,omitempty"`
	PathKey      string          `json:"pathKey"`
	Scopes       []string        `json:"scopes"`
	ConnectedAt  time.Time       `json:"connectedAt"`
	LastUsedAt   *time.Time      `json:"lastUsedAt,omitempty"`
	Created      time.Time       `json:"created"`
	Updated      time.Time       `json:"updated"`
}

// ConnectionQueryDTO is the query input for listing connections.
type ConnectionQueryDTO struct {
	Provider *string `query:"provider" validate:"omitempty"`
	Status   *string `query:"status" validate:"omitempty"`
	Page     *int    `query:"page" validate:"omitempty,min=1"`
	PerPage  *int    `query:"perPage" validate:"omitempty,min=1,max=100"`
}
