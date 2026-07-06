package repositories

import (
	"context"

	"assets/src/modules/assettemplates/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// MigrationPlanRepository persists template migration plans (collection
// template_migration_plans).
type MigrationPlanRepository interface {
	Create(ctx context.Context, p *entities.MigrationPlan) (*entities.MigrationPlan, error)
	FindById(ctx context.Context, id *string) (*entities.MigrationPlan, error)
	FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.MigrationPlan, error)
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.MigrationPlan], error)
	// IncrementCounter atomically adjusts a nested counter (e.g.
	// "counters.migrated") by delta.
	IncrementCounter(ctx context.Context, planID string, field string, delta int) error
}

// MigrationExecutionRepository persists per-asset migration executions
// (collection template_migration_executions, indexed by planId and status).
type MigrationExecutionRepository interface {
	BulkInsert(ctx context.Context, execs []*entities.MigrationExecution) (int, error)
	FindByPlan(ctx context.Context, planID *string, filters model.Map, pagination *model.PaginationOpts) (*model.PaginatedResult[entities.MigrationExecution], error)
	FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.MigrationExecution, error)
	CountByPlanAndStatus(ctx context.Context, planID string, status entities.ExecStatus) (int64, error)
	// UpdateByPlanAndAsset sets fields on the single execution identified by its
	// plan and asset, the natural key the run loop holds without a prior read.
	UpdateByPlanAndAsset(ctx context.Context, planID, assetID model.ObjectId, fields map[string]any) error
	// CancelPending transitions every still-pending execution of a plan to
	// cancelled so a cancelled plan's assets no longer read as queued.
	CancelPending(ctx context.Context, planID model.ObjectId) (int64, error)
}
