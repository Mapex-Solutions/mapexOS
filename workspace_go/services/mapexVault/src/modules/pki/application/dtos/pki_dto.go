package dtos

import contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexVault/pki"

// Alias-only. Add fresh DTOs in
// packages/contracts/services/mapexVault/pki/ and re-alias here.

type (
	IntermediateCABundleResponse = contracts.IntermediateCABundleResponse
	SignServerRequest            = contracts.SignServerRequest
	SignServerResponse           = contracts.SignServerResponse
	CAChainResponse              = contracts.CAChainResponse
)
