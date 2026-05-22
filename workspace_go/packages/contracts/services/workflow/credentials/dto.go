package credentials

import (
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

/**
 * DTOs (API Layer) — Credential management for encrypted secrets.
 *
 * Credentials store encrypted values (bot tokens, API keys, OAuth secrets)
 * that plugins need to connect to external services.
 * The plugin manifest defines the SCHEMA (CredentialDef), these DTOs
 * handle the CRUD for actual stored values.
 */

// CredentialId represents the params DTO for /:id
type CredentialId struct {
	Id string `params:"id" validate:"required"`
}

// CredentialCreate represents the payload to create an encrypted credential.
type CredentialCreate struct {
	Name           string                 `json:"name"           validate:"required,min=1,max=255"`
	PluginId       string                 `json:"pluginId"       validate:"required"`
	CredentialType string                 `json:"credentialType" validate:"required"`
	Data           map[string]interface{} `json:"data"           validate:"required"`
}

// CredentialUpdate represents the payload to partially update a credential.
// If Data is provided, the credential is re-encrypted with a new DEK.
type CredentialUpdate struct {
	Name *string                `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// CredentialQuery represents the query parameters for listing credentials.
// Embeds BaseQueryDTO for standard pagination and sorting.
type CredentialQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	PluginId       *string `query:"pluginId" validate:"omitempty"`
	CredentialType *string `query:"credentialType" validate:"omitempty"`
}

// CredentialResponse represents the API response for a credential.
// NEVER includes encrypted fields (EncryptedDEK, DEKNonce, EncryptedData, DataNonce).
type CredentialResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	PluginId       string    `json:"pluginId"`
	CredentialType string    `json:"credentialType"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
}

// CredentialTestResponse represents the result of testing a credential.
type CredentialTestResponse struct {
	Success bool `json:"success"`
}

/**
 * SCHEMA DTOs — Used by the UI to render credential forms dynamically.
 * These mirror the CredentialDef from the plugin manifest.
 */

// CredentialFieldSchema describes a single input field in the credential form.
type CredentialFieldSchema struct {
	Name        string           `json:"name"`
	DisplayName string           `json:"displayName"`
	Type        string           `json:"type"` // "string", "number", "boolean", "options"
	Required    *bool            `json:"required,omitempty"`
	IsSecret    *bool            `json:"isSecret,omitempty"`
	Hint        string           `json:"hint,omitempty"`
	Default     interface{}      `json:"default,omitempty"`
	Options     []PropertyOption `json:"options,omitempty"`
}

// PropertyOption is a selectable option for "options" type fields.
type PropertyOption struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

// CredentialTestSchema describes the test endpoint configuration.
type CredentialTestSchema struct {
	Method string                 `json:"method"`
	Path   string                 `json:"path"`
	Body   map[string]interface{} `json:"body,omitempty"`
}

// CredentialSchema is the full schema definition for a credential type.
// Returned by GET /api/v1/credentials/schema/:pluginId
type CredentialSchema struct {
	ID     string                  `json:"id"`
	Name   string                  `json:"name"`
	Fields []CredentialFieldSchema `json:"fields"`
	Test   *CredentialTestSchema   `json:"test,omitempty"`
}
