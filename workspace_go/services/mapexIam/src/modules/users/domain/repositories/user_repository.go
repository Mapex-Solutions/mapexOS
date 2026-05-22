package repositories

import (
	"context"
	"mapexIam/src/modules/users/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *entities.User) (*entities.User, error)
	FindById(ctx context.Context, userId *string) (*entities.User, error)
	FindByIdAndUpdate(ctx context.Context, userId *string, payload map[string]any) (*entities.User, error)
	DeleteById(ctx context.Context, userId *string) error
	FindByEmail(ctx context.Context, email *string) (*entities.User, error)
	FindWithFilters(
		ctx context.Context,
		filters model.Map,
		pagination *model.PaginationOpts,
		projection model.Map,
	) (*model.PaginatedResult[entities.User], error)
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
