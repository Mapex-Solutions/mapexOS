package repositories

import "workflow/src/modules/runtime/domain/entities"

// ExecutionStateRepository provides hot state persistence for workflow executions.
// Implementations handle key formatting, serialization, and storage operations.
type ExecutionStateRepository interface {
	// Create persists the initial execution state (fails if already exists).
	Create(execution *entities.WorkflowExecution) error

	// Get retrieves the execution state from hot storage.
	Get(executionID string) (*entities.WorkflowExecution, error)

	// Save persists a checkpoint of the current execution state (overwrites).
	Save(execution *entities.WorkflowExecution) error

	// GetWithRevision retrieves the execution state plus KV revision for CAS operations.
	GetWithRevision(executionID string) (*entities.WorkflowExecution, uint64, error)

	// SaveWithRevision checkpoints using CAS — fails if revision changed since load.
	SaveWithRevision(execution *entities.WorkflowExecution, revision uint64) error
}
