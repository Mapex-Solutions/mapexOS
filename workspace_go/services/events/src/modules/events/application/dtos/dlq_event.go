package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
)

/**
 * Events DLQ DTOs - Type aliases for contract DTOs
 */

type (
	// EventsDLQQueryDto is the query DTO for listing DLQ events with cursor pagination
	EventsDLQQueryDto = v1.EventsDLQQuery

	// EventsDLQResponseDto is the response DTO for a single DLQ event
	EventsDLQResponseDto = v1.EventsDLQResponse

	// EventsDLQCursorResultDto is the cursor-paginated result for DLQ events
	EventsDLQCursorResultDto = v1.EventsDLQCursorResult

	// DLQEventIncomingDto is the incoming DTO from NATS (MAPEXOS-DLQ stream)
	DLQEventIncomingDto = v1.DLQEventIncomingDTO

	// EventsDLQCountsQueryDto is the query DTO for counting DLQ by service type
	EventsDLQCountsQueryDto = v1.EventsDLQCountsQuery

	// EventsDLQServiceCountDto is a single service type count
	EventsDLQServiceCountDto = v1.EventsDLQServiceCount

	// EventsDLQCountsResultDto is the counts response
	EventsDLQCountsResultDto = v1.EventsDLQCountsResult
)
