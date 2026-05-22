package services

import (
	"context"

	"workflow/src/modules/archiver/application/di"
	"workflow/src/modules/archiver/application/dtos"
	"workflow/src/modules/archiver/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

/*
 * ARCHIVER SERVICE
 * Consumes WORKFLOW-STATE lifecycle events and persists to MongoDB.
 * This is the ONLY module that writes to MongoDB for workflow instances.
 * Runtime NEVER touches MongoDB — only KV and NATS streams.
 *
 * Event types (4-write lifecycle):
 *   created   → BulkInsertLightweight (listing stub ~200B)
 *   waiting   → BulkUpdateWaiting (status, activeNodeIds ~150B)
 *   resumed   → BulkUpdateResumed (status=running ~100B)
 *   completed → KV Get FULL → BulkUpsertFull (~5-25KB) → KV Delete
 *   failed    → KV Get FULL → BulkUpsertFull (~5-25KB) → KV Delete
 *   cancelled → KV Get FULL → BulkUpsertFull (~5-25KB) → KV Delete
 */

// Compile-time check
var _ ports.ArchiverServicePort = (*ArchiverService)(nil)

// New creates and returns a new instance of ArchiverService.
func New(deps di.ArchiverServiceDependenciesInjection) ports.ArchiverServicePort {
	return &ArchiverService{deps: deps}
}

// ProcessStateBatch processes a batch of WORKFLOW-STATE events. Steps:
// honour backpressure (sleep when MongoDB P99 latency is high) -> classify
// each message into one of the 4 batch buckets (created/waiting/resumed/
// terminal) -> run the 4 bulk-write batches against MongoDB and ack/nack
// per batch result.
func (s *ArchiverService) ProcessStateBatch(messages []*natsModel.Message) {
	s.applyArchiverBackpressure()
	classified := s.classifyStateEvents(messages)
	s.runArchiverWriteBatches(context.Background(), classified)
}

// GetExecutions retrieves workflow executions from MongoDB with org filtering
// and pagination. Steps: build the org-scoped filter -> overlay the per-field
// filters from the query DTO -> paginate -> map entities to response DTOs.
func (s *ArchiverService) GetExecutions(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.ExecutionQueryDTO) (*model.PaginatedResult[dtos.ExecutionResponseDTO], error) {
	filters, err := s.buildExecutionListFilters(requestContext, query)
	if err != nil {
		return nil, err
	}
	pagination := s.buildExecutionListPagination(query)
	result, err := s.deps.ArchiveRepo.FindExecutions(ctx, filters, pagination)
	if err != nil {
		return nil, err
	}
	return s.mapExecutionListResult(result), nil
}

// GetExecutionById retrieves a single execution by MongoDB _id. For
// non-terminal executions (running/waiting), enriches with full state from
// NATS KV using the workflowUUID as the KV key.
func (s *ArchiverService) GetExecutionById(ctx context.Context, executionId string) (*dtos.ExecutionResponseDTO, error) {
	exec, err := s.fetchExecutionOr404(ctx, executionId)
	if err != nil {
		return nil, err
	}
	exec = s.enrichWithKVState(exec)
	dto := mapExecutionToDTO(exec)
	return &dto, nil
}
