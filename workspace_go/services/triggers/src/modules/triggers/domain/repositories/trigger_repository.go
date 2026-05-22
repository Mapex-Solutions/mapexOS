package repositories

import (
	"context"

	"triggers/src/modules/triggers/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// TriggerRepository defines the interface for trigger data access
type TriggerRepository interface {
	Create(ctx context.Context, trigger *entities.Trigger) (*entities.Trigger, error)
	FindById(ctx context.Context, triggerId *string) (*entities.Trigger, error)
	FindByIdAndUpdate(ctx context.Context, triggerId *string, payload map[string]any) (*entities.Trigger, error)
	DeleteById(ctx context.Context, triggerId *string) error
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Trigger], error)
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
