package entities

import (
	"time"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexVault/credentials"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/**
 * Credential Types
 *
 * Domain re-exports of cross-service contract types. The authoritative shapes
 * live in packages/contracts/services/mapexVault/credentials.
 */

type CredentialType = contracts.CredentialType

const (
	CredentialManual      = contracts.CredentialManual
	CredentialOAuth2      = contracts.CredentialOAuth2
	CredentialUserAndPass = contracts.CredentialUserAndPass
)

/**
 * Credential Status
 */

type CredentialStatus = contracts.CredentialStatus

const (
	CredentialStatusActive  = contracts.CredentialStatusActive
	CredentialStatusExpired = contracts.CredentialStatusExpired
	CredentialStatusRevoked = contracts.CredentialStatusRevoked
	CredentialStatusError   = contracts.CredentialStatusError
)

/**
 * Credential Entity
 * Stores encrypted credential data with token lifecycle management.
 * Sensitive data is stored using envelope encryption (Master Key → DEK → Data).
 * Encrypted fields are tagged json:"-" to prevent accidental API exposure.
 */

type Credential struct {
	ID              model.ObjectId   `bson:"_id,omitempty" json:"id"`
	Name            string           `bson:"name" json:"name"`
	Type            CredentialType   `bson:"type" json:"type"`
	PluginId        string           `bson:"pluginId" json:"pluginId"`
	CredentialDefId string           `bson:"credentialDefId" json:"credentialDefId"`
	OrgId           *model.ObjectId  `bson:"orgId,omitempty" json:"orgId,omitempty"`
	PathKey         string           `bson:"pathKey" json:"pathKey"`
	IsTemplate      bool             `bson:"isTemplate" json:"isTemplate"`
	Status          CredentialStatus `bson:"status" json:"status"`

	// Envelope encryption (NEVER exposed via API)
	EncryptedDEK  []byte `bson:"encryptedDEK" json:"-"`
	DEKNonce      []byte `bson:"dekNonce" json:"-"`
	EncryptedData []byte `bson:"encryptedData" json:"-"`
	DataNonce     []byte `bson:"dataNonce" json:"-"`

	// Token lifecycle (computed by vault, NOT from user input)
	TokenExpiresAt  *time.Time `bson:"tokenExpiresAt,omitempty" json:"tokenExpiresAt,omitempty"`
	LastRefreshedAt *time.Time `bson:"lastRefreshedAt,omitempty" json:"lastRefreshedAt,omitempty"`
	RefreshError    string     `bson:"refreshError,omitempty" json:"refreshError,omitempty"`

	// Provider config (NOT encrypted — needed by refresh consumer without decrypt)
	ProviderConfig *ProviderConfig `bson:"providerConfig,omitempty" json:"providerConfig,omitempty"`

	Created time.Time `bson:"created" json:"created"`
	Updated time.Time `bson:"updated" json:"updated"`
}

/**
 * TokenRequestConfig
 * Alias of contracts.TokenRequestConfig — kept here so repositories and
 * services continue to use the domain name while the wire shape lives in
 * packages/contracts.
 */

type TokenRequestConfig = contracts.TokenRequestConfig

/**
 * ProviderConfig
 * Alias of contracts.ProviderConfig — same rationale as TokenRequestConfig.
 */

type ProviderConfig = contracts.ProviderConfig
