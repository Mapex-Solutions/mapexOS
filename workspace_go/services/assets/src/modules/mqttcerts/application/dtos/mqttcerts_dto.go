package dtos

import contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/mqttcerts"

// Alias-only — canonical DTO shape lives in packages/contracts.

type (
	IssueCertRequest    = contracts.IssueCertRequest
	IssueCertResponse   = contracts.IssueCertResponse
	RevokedCertResponse = contracts.RevokedCertResponse
	ListRevokedQuery    = contracts.ListRevokedQuery
)
