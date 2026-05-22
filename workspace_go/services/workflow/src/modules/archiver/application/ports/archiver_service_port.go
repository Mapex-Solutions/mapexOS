package ports

import (
	"context"

	"workflow/src/modules/archiver/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/*
 * ARCHIVER SERVICE PORT
 * Defines the contract for the Archiver service consumed by the NATS consumer and HTTP routes.
 */

// ArchiverServicePort defines the Archiver's batch processing and query interface.
type ArchiverServicePort interface {
	// ProcessStateBatch processes a batch of WORKFLOW-STATE lifecycle events.
	// Classifies events by type and persists to MongoDB:
	//   - "created"   → BulkInsertLightweight (listing stub)
	//   - "completed", "failed", "cancelled" → KV Get + BulkUpsertFull + KV Delete
	ProcessStateBatch(messages []*natsModel.Message)

	// GetExecutions retrieves workflow executions from MongoDB (hot storage) with pagination.
	// Uses RequestContext for context-aware organization filtering.
	GetExecutions(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.ExecutionQueryDTO) (*model.PaginatedResult[dtos.ExecutionResponseDTO], error)

	// GetExecutionById retrieves a single workflow execution by ID.
	GetExecutionById(ctx context.Context, executionId string) (*dtos.ExecutionResponseDTO, error)
}
