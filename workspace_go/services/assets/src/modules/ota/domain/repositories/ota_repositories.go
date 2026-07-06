package repositories

import (
	"context"

	"assets/src/modules/ota/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// FirmwareRepository persists firmware artifacts (collection ota_firmware).
type FirmwareRepository interface {
	Create(ctx context.Context, f *entities.Firmware) (*entities.Firmware, error)
	FindById(ctx context.Context, id *string) (*entities.Firmware, error)
	FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.Firmware, error)
	DeleteById(ctx context.Context, id *string) error
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Firmware], error)
}

// OTAPlanRepository persists OTA plans (collection ota_plans).
type OTAPlanRepository interface {
	Create(ctx context.Context, p *entities.OTAPlan) (*entities.OTAPlan, error)
	FindById(ctx context.Context, id *string) (*entities.OTAPlan, error)
	FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.OTAPlan, error)
	DeleteById(ctx context.Context, id *string) error
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.OTAPlan], error)
	// IncrementCounter atomically adjusts a nested counter (e.g.
	// "counters.succeeded") by delta.
	IncrementCounter(ctx context.Context, planID string, field string, delta int) error
}

// OTAExecutionRepository persists per-asset executions (collection
// ota_executions, indexed by planId and assetId).
type OTAExecutionRepository interface {
	BulkInsert(ctx context.Context, execs []*entities.OTAExecution) (int64, error)
	FindById(ctx context.Context, id *string) (*entities.OTAExecution, error)
	FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.OTAExecution, error)
	FindByPlan(ctx context.Context, planID *string, filters model.Map, pagination *model.PaginationOpts) (*model.PaginatedResult[entities.OTAExecution], error)
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.OTAExecution], error)
	CountByState(ctx context.Context, planID *string, state entities.ExecutionState) (int64, error)
	UpdateMany(ctx context.Context, filter model.Map, update model.Map) (int64, error)
}
