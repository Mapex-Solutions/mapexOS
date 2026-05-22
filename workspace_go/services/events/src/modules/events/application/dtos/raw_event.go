package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
)

/**
 * Events Raw DTOs - Type aliases for contract DTOs
 */

type (
	// EventsRawQueryDto is the query DTO for listing raw events with cursor pagination
	EventsRawQueryDto = v1.EventsRawQuery

	// EventsRawResponseDto is the response DTO for a single raw event
	EventsRawResponseDto = v1.EventsRawResponse

	// EventsRawCursorResultDto is the cursor-paginated result for raw events
	EventsRawCursorResultDto = v1.EventsRawCursorResult

	// RawEventDto is the incoming DTO from NATS (HTTP/MQTT gateways)
	RawEventDto = v1.RawEventDTO
)
