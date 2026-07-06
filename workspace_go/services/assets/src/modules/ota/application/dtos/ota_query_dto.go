package dtos

import otaContract "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"

// Query and param DTOs are canonical in packages/contracts; re-aliased here so
// the module depends on the shared contract, never a local redefinition.
type (
	OTAPlanQueryDTO      = otaContract.OTAPlanQueryDTO
	OTAExecutionQueryDTO = otaContract.OTAExecutionQueryDTO
	OTAPlanUpdateRequest = otaContract.OTAPlanUpdateRequest
	OTAPlanIdDto         = otaContract.OTAPlanIdDto
	OTAFirmwareIdDto     = otaContract.OTAFirmwareIdDto
)
