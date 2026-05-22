package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/executions"
)

/**
 * Execution DTOs - Type aliases for contract DTOs
 */

type (
	// ExecutionQueryDTO is the query DTO for listing executions with pagination
	ExecutionQueryDTO = v1.ExecutionQuery

	// ExecutionResponseDTO is the response DTO for a single execution
	ExecutionResponseDTO = v1.ExecutionResponse

	// ExecutionIdDTO is the params DTO for /:executionId
	ExecutionIdDTO = v1.ExecutionId
)
