package nats

import (
	"encoding/json"
	"fmt"

	"workflow/src/modules/runtime/domain/entities"
	"workflow/src/modules/runtime/domain/repositories"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time check
var _ repositories.ExecutionStateRepository = (*ExecutionStateRepository)(nil)

// NewExecutionStateRepository creates a new ExecutionStateRepository backed by the given NATS KV store.
func NewExecutionStateRepository(kvStore natsModel.KeyValueStore) repositories.ExecutionStateRepository {
	return &ExecutionStateRepository{kvStore: kvStore}
}

// Create persists a new workflow execution to the NATS KV store.
func (r *ExecutionStateRepository) Create(execution *entities.WorkflowExecution) error {
	key := formatKey(execution.WorkflowUUID)
	data, err := json.Marshal(execution)
	if err != nil {
		return fmt.Errorf("[INFRA:ExecutionStateRepo] failed to marshal execution: %w", err)
	}
	if _, err := r.kvStore.Create(key, data); err != nil {
		return fmt.Errorf("[INFRA:ExecutionStateRepo] failed to create KV entry: %w", err)
	}
	logger.Debug(fmt.Sprintf("[INFRA:ExecutionStateRepo] Created execution %s", execution.WorkflowUUID))
	return nil
}

// Get retrieves a workflow execution from the NATS KV store by its UUID.
func (r *ExecutionStateRepository) Get(executionID string) (*entities.WorkflowExecution, error) {
	key := formatKey(executionID)
	entry, err := r.kvStore.Get(key)
	if err != nil {
		return nil, fmt.Errorf("[INFRA:ExecutionStateRepo] failed to get execution %s: %w", executionID, err)
	}
	var execution entities.WorkflowExecution
	if err := json.Unmarshal(entry.Value, &execution); err != nil {
		return nil, fmt.Errorf("[INFRA:ExecutionStateRepo] failed to unmarshal execution %s: %w", executionID, err)
	}
	return &execution, nil
}

// Save checkpoints the workflow execution state to the NATS KV store.
func (r *ExecutionStateRepository) Save(execution *entities.WorkflowExecution) error {
	key := formatKey(execution.WorkflowUUID)
	data, err := json.Marshal(execution)
	if err != nil {
		return fmt.Errorf("[INFRA:ExecutionStateRepo] failed to marshal execution: %w", err)
	}
	if _, err := r.kvStore.Put(key, data); err != nil {
		return fmt.Errorf("[INFRA:ExecutionStateRepo] KV checkpoint failed: %w", err)
	}
	return nil
}

// GetWithRevision retrieves the execution state plus KV revision for CAS operations.
func (r *ExecutionStateRepository) GetWithRevision(executionID string) (*entities.WorkflowExecution, uint64, error) {
	key := formatKey(executionID)
	entry, err := r.kvStore.Get(key)
	if err != nil {
		return nil, 0, fmt.Errorf("[INFRA:ExecutionStateRepo] failed to get execution %s: %w", executionID, err)
	}
	var execution entities.WorkflowExecution
	if err := json.Unmarshal(entry.Value, &execution); err != nil {
		return nil, 0, fmt.Errorf("[INFRA:ExecutionStateRepo] failed to unmarshal execution %s: %w", executionID, err)
	}
	return &execution, entry.Revision, nil
}

// SaveWithRevision checkpoints using CAS — fails if revision changed since load.
func (r *ExecutionStateRepository) SaveWithRevision(execution *entities.WorkflowExecution, revision uint64) error {
	key := formatKey(execution.WorkflowUUID)
	data, err := json.Marshal(execution)
	if err != nil {
		return fmt.Errorf("[INFRA:ExecutionStateRepo] failed to marshal execution: %w", err)
	}
	if _, err := r.kvStore.Update(key, data, revision); err != nil {
		return fmt.Errorf("[INFRA:ExecutionStateRepo] CAS checkpoint failed (revision %d): %w", revision, err)
	}
	return nil
}

func formatKey(executionID string) string {
	return fmt.Sprintf("exec.%s", executionID)
}
