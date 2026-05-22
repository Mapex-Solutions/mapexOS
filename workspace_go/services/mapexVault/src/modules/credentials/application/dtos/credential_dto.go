package dtos

import (
	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexVault/credentials"
)

/**
 * DTO Aliases
 *
 * DTOs are defined ONCE in packages/contracts. This module re-exports them
 * via type aliases — zero re-definition, zero boundary violation (no
 * domain/entities import here).
 */

type (
	// Credential DTOs
	CreateCredentialDTO = contracts.CreateCredentialDTO
	UpdateCredentialDTO = contracts.UpdateCredentialDTO
	CredentialResponse  = contracts.CredentialResponse
	CredentialQueryDTO  = contracts.CredentialQueryDTO

	// OAuth2 DTOs
	OAuthCallbackDTO = contracts.OAuthCallbackDTO

	// Connection DTOs
	CreateConnectionDTO = contracts.CreateConnectionDTO
	UpsertConnectionDTO = contracts.UpsertConnectionDTO
	ConnectionResponse  = contracts.ConnectionResponse
	ConnectionQueryDTO  = contracts.ConnectionQueryDTO
)
