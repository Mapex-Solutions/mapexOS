package repositories

import (
	"context"
	"mapexIam/src/modules/lists/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type ListRepository interface {
	Create(ctx context.Context, u *entities.List) (*entities.List, error)
	FindById(ctx context.Context, listId *string) (*entities.List, error)
	FindByIds(ctx context.Context, listIds []string) ([]*entities.List, error)
	FindByIdAndUpdate(ctx context.Context, listId *string, payload map[string]any) (*entities.List, error)
	DeleteById(ctx context.Context, listId *string) error
	FindByEmail(ctx context.Context, email *string) (*entities.List, error)
	FindWithFilters(
		ctx context.Context,
		filters model.Map,
		pagination *model.PaginationOpts,
		projection model.Map,
	) (*model.PaginatedResult[entities.List], error)
}
