package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
)

/**
 * Events JS Executor DTOs - Type aliases for contract DTOs
 */

type (
	// EventsJsExecQueryDto is the query DTO for listing JS executor events with cursor pagination
	EventsJsExecQueryDto = v1.EventsJsExecQuery

	// EventsJsExecResponseDto is the response DTO for a single JS executor event
	EventsJsExecResponseDto = v1.EventsJsExecResponse

	// EventsJsExecCursorResultDto is the cursor-paginated result for JS executor events
	EventsJsExecCursorResultDto = v1.EventsJsExecCursorResult

	// JsExecEventDto is the payload received from JS Executor service via NATS
	JsExecEventDto = v1.JsExecEventDTO

	// JsExecEventFlatDto is the flat DTO for entity mapping via mapper.DtoToEntity[]
	JsExecEventFlatDto = v1.JsExecEventFlatDTO
)
