package dtos

import (
	triggersContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/triggers/triggers"
)

// Alias contracts DTOs (single source of truth)
// These aliases make the contracts available to the service with cleaner naming
type (
	CreateTriggerDto    = triggersContracts.TriggerCreate
	UpdateTriggerDto    = triggersContracts.TriggerUpdate
	TriggerResponse     = triggersContracts.TriggerResponse
	TriggerQueryDto     = triggersContracts.TriggerQuery
	TriggerListResponse = triggersContracts.TriggerListResponse
	TriggerExecuteEvent = triggersContracts.TriggerExecuteEvent
)
