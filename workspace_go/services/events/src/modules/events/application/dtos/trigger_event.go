package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
)

/**
 * Events Trigger DTOs - Type aliases for contract DTOs
 */

type (
	// EventsTriggerQueryDto is the query DTO for listing trigger events with cursor pagination
	EventsTriggerQueryDto = v1.EventsTriggerQuery

	// EventsTriggerResponseDto is the response DTO for a single trigger event
	EventsTriggerResponseDto = v1.EventsTriggerResponse

	// EventsTriggerCursorResultDto is the cursor-paginated result for trigger events
	EventsTriggerCursorResultDto = v1.EventsTriggerCursorResult

	// TriggerEventIncomingDto is the incoming DTO from NATS (triggers service)
	TriggerEventIncomingDto = v1.TriggerEventDTO
)
