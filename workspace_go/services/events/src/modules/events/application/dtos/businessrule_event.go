package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
)

/**
 * Events BusinessRule DTOs - Type aliases for contract DTOs
 */

type (
	// EventsBusinessRuleQueryDto is the query DTO for listing business rule events with cursor pagination
	EventsBusinessRuleQueryDto = v1.EventsBusinessRuleQuery

	// EventsBusinessRuleResponseDto is the response DTO for a single business rule event
	EventsBusinessRuleResponseDto = v1.EventsBusinessRuleResponse

	// EventsBusinessRuleCursorResultDto is the cursor-paginated result for business rule events
	EventsBusinessRuleCursorResultDto = v1.EventsBusinessRuleCursorResult

	// BusinessRuleEventIncomingDto is the incoming DTO from NATS (workflow service)
	BusinessRuleEventIncomingDto = v1.BusinessRuleEventDTO
)
