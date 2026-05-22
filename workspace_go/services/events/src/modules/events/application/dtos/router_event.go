package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
)

/**
 * Events Router DTOs - Type aliases for contract DTOs
 */

type (
	// EventsRouterQueryDto is the query DTO for listing router events with cursor pagination
	EventsRouterQueryDto = v1.EventsRouterQuery

	// EventsRouterResponseDto is the response DTO for a single router event
	EventsRouterResponseDto = v1.EventsRouterResponse

	// EventsRouterCursorResultDto is the cursor-paginated result for router events
	EventsRouterCursorResultDto = v1.EventsRouterCursorResult

	// RouterEventIncomingDto is the incoming DTO from NATS (router service)
	RouterEventIncomingDto = v1.RouterEventDTO
)
