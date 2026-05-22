package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/instances"
)

// Instance DTO type aliases re-export contract types for the instances module.
type (
	InstanceIdDTO      = v1.InstanceId
	InstanceCreateDTO  = v1.InstanceCreate
	InstanceUpdateDTO  = v1.InstanceUpdate
	InstanceQueryDTO   = v1.InstanceQuery
	InstanceResponse   = v1.InstanceResponse
	ExecuteRequestDTO  = v1.ExecuteRequest
	ExecuteResponseDTO = v1.ExecuteResponse
)
