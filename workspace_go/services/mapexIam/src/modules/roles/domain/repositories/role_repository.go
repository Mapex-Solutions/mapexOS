package repositories

import (
	"context"
	"mapexIam/src/modules/roles/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type RoleRepository interface {
	Create(ctx context.Context, u *entities.Role) (*entities.Role, error)
	FindById(ctx context.Context, roleId *string) (*entities.Role, error)
	FindByIdAndUpdate(ctx context.Context, roleId *string, payload map[string]any) (*entities.Role, error)
	DeleteById(ctx context.Context, roleId *string) error
	FindWithFilters(
		ctx context.Context,
		filters model.Map,
		pagination *model.PaginationOpts,
		projection model.Map,
	) (*model.PaginatedResult[entities.Role], error)
}
