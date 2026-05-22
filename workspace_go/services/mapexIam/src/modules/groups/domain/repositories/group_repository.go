package repositories

import (
	"context"
	"mapexIam/src/modules/groups/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type GroupRepository interface {
	Create(ctx context.Context, u *entities.Group) (*entities.Group, error)
	FindById(ctx context.Context, groupId *string) (*entities.Group, error)
	FindByIdAndUpdate(ctx context.Context, groupId *string, payload map[string]any) (*entities.Group, error)
	DeleteById(ctx context.Context, groupId *string) error
	FindWithFilters(
		ctx context.Context,
		filters model.Map,
		pagination *model.PaginationOpts,
		projection model.Map,
	) (*model.PaginatedResult[entities.Group], error)
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
