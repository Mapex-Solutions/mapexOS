package repositories

import (
	"context"

	runtimePorts "workflow/src/modules/runtime/application/ports"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/*
 * ARCHIVE REPOSITORY
 * Provides BulkWrite operations for persisting workflow executions to MongoDB.
 * Writes to the "executions" collection with batch-oriented operations.
 */

// ArchiveRepository defines the persistence interface for the Archiver module.
// All writes are batch-oriented (BulkWrite) to minimize MongoDB round-trips.
type ArchiveRepository interface {
	// BulkInsertLightweight inserts minimal execution stubs for listing visibility.
	// Called on "created" events (~200 bytes per document).
	BulkInsertLightweight(ctx context.Context, stubs []LightweightExecution) error

	// BulkUpsertFull upserts complete execution documents for terminal workflows.
	// Called on "completed", "failed", "cancelled" events (~5-25KB per document).
	BulkUpsertFull(ctx context.Context, executions []*runtimePorts.WorkflowExecution) error

	// BulkUpdateWaiting updates status and activeNodeIds for waiting executions.
	// Called on "waiting" events (~150 bytes per update).
	BulkUpdateWaiting(ctx context.Context, updates []WaitingUpdate) error

	// BulkUpdateResumed sets status to running for resumed executions.
	// Called on "resumed" events (~100 bytes per update).
	BulkUpdateResumed(ctx context.Context, executionIDs []string) error

	// FindExecutions queries workflow executions with org filter, status filter, and pagination.
	// Used by the HTTP API for listing execution history from MongoDB (hot storage).
	FindExecutions(ctx context.Context, filters model.Map, pagination *model.PaginationOpts) (*model.PaginatedResult[runtimePorts.WorkflowExecution], error)

	// FindExecutionById retrieves a single execution by its MongoDB _id.
	FindExecutionById(ctx context.Context, executionId string) (*runtimePorts.WorkflowExecution, error)
}
