package repositories

import (
	"context"

	"router/src/modules/routegroups/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type RouteGroupRepository interface {
	Create(ctx context.Context, u *entities.RouteGroup) (*entities.RouteGroup, error)
	FindById(ctx context.Context, listId *string) (*entities.RouteGroup, error)
	FindByIdAndUpdate(ctx context.Context, listId *string, payload map[string]any) (*entities.RouteGroup, error)
	DeleteById(ctx context.Context, listId *string) error
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.RouteGroup], error)
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
