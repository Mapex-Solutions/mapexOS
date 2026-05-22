package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/definitions"
)

// Definition DTO type aliases re-export contract types for the definitions module.
type (
	DefinitionCreateDTO = v1.DefinitionCreate
	DefinitionUpdateDTO = v1.DefinitionUpdate
	DefinitionQueryDTO  = v1.DefinitionQuery
	DefinitionIdDTO     = v1.DefinitionId
	DefinitionResponse  = v1.DefinitionResponse
)
