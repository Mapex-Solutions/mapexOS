package collection

import (
	"context"
	"fmt"
	"time"

	"workflow/src/modules/archiver/domain/repositories"
	"workflow/src/modules/archiver/infrastructure/persistence/mongo/constants"
	runtimePorts "workflow/src/modules/runtime/application/ports"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/*
 * ARCHIVE REPOSITORY IMPLEMENTATION
 * Writes to the "executions" collection.
 * All writes are BulkWrite for throughput (batch from NATS consumer).
 */

// Compile-time check
var _ repositories.ArchiveRepository = (*repository)(nil)

// New creates the archive repository for executions collection.
func New(m *manager.MongoManager) RepositoryOut {
	mdl := model.New[runtimePorts.WorkflowExecution](m.GetDatabase(), constants.CollectionName, model.Config{
		Indexes: constants.Indexes,
	})
	repo := &repository{model: mdl}
	return RepositoryOut{
		ArchiveRepo: repo,
	}
}

// BulkInsertLightweight inserts minimal stubs for "created" events.
func (r *repository) BulkInsertLightweight(ctx context.Context, stubs []repositories.LightweightExecution) error {
	if len(stubs) == 0 {
		return nil
	}

	models := make([]model.WriteModel, 0, len(stubs))
	for i := range stubs {
		models = append(models, model.NewInsertOneModel().SetDocument(stubs[i]))
	}

	opts := model.BulkWrite().SetOrdered(false)
	result, err := r.model.DIRECT().BulkWrite(ctx, models, opts)
	if err != nil {
		if !model.IsDuplicateKeyError(err) {
			return fmt.Errorf("[INFRA:Archive] BulkInsertLightweight failed: %w", err)
		}
		logger.Debug(fmt.Sprintf("[INFRA:Archive] BulkInsertLightweight: %d inserted, some duplicates ignored", result.InsertedCount))
		return nil
	}

	logger.Debug(fmt.Sprintf("[INFRA:Archive] BulkInsertLightweight: %d inserted", result.InsertedCount))
	return nil
}

// BulkUpsertFull replaces complete execution documents for terminal events.
func (r *repository) BulkUpsertFull(ctx context.Context, executions []*runtimePorts.WorkflowExecution) error {
	if len(executions) == 0 {
		return nil
	}

	models := make([]model.WriteModel, 0, len(executions))
	for _, exec := range executions {
		filter := model.Map{"workflowUUID": exec.WorkflowUUID}
		// Use $set instead of Replace to preserve the auto-generated _id
		models = append(models, model.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(model.Map{"$set": exec}).
			SetUpsert(true))
	}

	opts := model.BulkWrite().SetOrdered(false)
	result, err := r.model.DIRECT().BulkWrite(ctx, models, opts)
	if err != nil {
		return fmt.Errorf("[INFRA:Archive] BulkUpsertFull failed: %w", err)
	}

	logger.Debug(fmt.Sprintf("[INFRA:Archive] BulkUpsertFull: %d upserted, %d modified",
		result.UpsertedCount, result.ModifiedCount))
	return nil
}

// BulkUpdateWaiting updates status and activeNodeIds for waiting executions.
func (r *repository) BulkUpdateWaiting(ctx context.Context, updates []repositories.WaitingUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	models := make([]model.WriteModel, 0, len(updates))
	for _, u := range updates {
		update := model.Map{
			"$set": model.Map{
				"status":        u.Status,
				"activeNodeIds": u.ActiveNodeIDs,
				"updated":       u.Updated,
			},
		}

		models = append(models, model.NewUpdateOneModel().
			SetFilter(model.Map{"workflowUUID": u.WorkflowUUID}).
			SetUpdate(update))
	}

	if len(models) == 0 {
		return nil
	}

	opts := model.BulkWrite().SetOrdered(false)
	result, err := r.model.DIRECT().BulkWrite(ctx, models, opts)
	if err != nil {
		return fmt.Errorf("[INFRA:Archive] BulkUpdateWaiting failed: %w", err)
	}

	logger.Debug(fmt.Sprintf("[INFRA:Archive] BulkUpdateWaiting: %d modified", result.ModifiedCount))
	return nil
}

// BulkUpdateResumed sets status to running for resumed executions.
func (r *repository) BulkUpdateResumed(ctx context.Context, executionIDs []string) error {
	if len(executionIDs) == 0 {
		return nil
	}

	models := make([]model.WriteModel, 0, len(executionIDs))
	for _, id := range executionIDs {
		models = append(models, model.NewUpdateOneModel().
			SetFilter(model.Map{"workflowUUID": id}).
			SetUpdate(model.Map{
				"$set": model.Map{"status": "running", "updated": time.Now()},
			}))
	}

	if len(models) == 0 {
		return nil
	}

	opts := model.BulkWrite().SetOrdered(false)
	result, err := r.model.DIRECT().BulkWrite(ctx, models, opts)
	if err != nil {
		return fmt.Errorf("[INFRA:Archive] BulkUpdateResumed failed: %w", err)
	}

	logger.Debug(fmt.Sprintf("[INFRA:Archive] BulkUpdateResumed: %d modified", result.ModifiedCount))
	return nil
}

// FindExecutions queries executions with org filter, status, and pagination.
func (r *repository) FindExecutions(ctx context.Context, filters model.Map, pagination *model.PaginationOpts) (*model.PaginatedResult[runtimePorts.WorkflowExecution], error) {
	return r.model.FindByOffset(ctx, filters, pagination)
}

// FindExecutionById retrieves a single execution by its MongoDB _id.
func (r *repository) FindExecutionById(ctx context.Context, executionId string) (*runtimePorts.WorkflowExecution, error) {
	objId, err := model.ToObjectID(executionId)
	if err != nil {
		return nil, fmt.Errorf("[INFRA:Archive] invalid executionId: %w", err)
	}
	filter := model.Map{"_id": objId}
	return r.model.FindOne(ctx, &filter)
}

