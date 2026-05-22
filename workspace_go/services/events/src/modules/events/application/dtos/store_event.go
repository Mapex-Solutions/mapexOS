package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
)

/**
 * Event Store DTOs - Type aliases for contract DTOs
 */

type (
	// EventStoreDto is the payload received from Router service via NATS
	EventStoreDto = v1.EventStoreDTO

	// EventsStoreQueryDto is the query DTO for listing processed events with cursor pagination
	EventsStoreQueryDto = v1.EventsStoreQuery

	// EventsStoreResponseDto is the response DTO for a single processed event
	EventsStoreResponseDto = v1.EventsStoreResponse

	// EventsStoreCursorResultDto is the cursor-paginated result for processed events
	EventsStoreCursorResultDto = v1.EventsStoreCursorResult

	// EventsStoreDetailResponseDto is the detail response for a single event with resolved EVA fields
	EventsStoreDetailResponseDto = v1.EventsStoreDetailResponse
)
