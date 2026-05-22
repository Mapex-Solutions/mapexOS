package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
)

/**
 * Events Workflow DTOs - Type aliases for contract DTOs
 */

type (
	// EventsWorkflowExecutionIdParamDto is the params DTO for /:executionId
	EventsWorkflowExecutionIdParamDto = v1.EventsWorkflowExecutionIdParam

	// EventsWorkflowQueryDto is the query DTO for listing workflow events with cursor pagination
	EventsWorkflowQueryDto = v1.EventsWorkflowQuery

	// EventsWorkflowResponseDto is the response DTO for a single workflow event
	EventsWorkflowResponseDto = v1.EventsWorkflowResponse

	// EventsWorkflowCursorResultDto is the cursor-paginated result for workflow events
	EventsWorkflowCursorResultDto = v1.EventsWorkflowCursorResult

	// WorkflowEventIncomingDto is the incoming DTO from NATS (workflow archiver)
	WorkflowEventIncomingDto = v1.WorkflowEventDTO
)
