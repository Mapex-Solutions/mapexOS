package repositories

import (
	"context"
	"http_gateway/src/modules/datasources/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type DataSourceRepository interface {
	Create(ctx context.Context, u *entities.DataSource) (*entities.DataSource, error)
	FindById(ctx context.Context, dataSourceId *string) (*entities.DataSource, error)
	FindByIdAndUpdate(ctx context.Context, dataSourceId *string, payload map[string]any) (*entities.DataSource, error)
	DeleteById(ctx context.Context, dataSourceId *string) error
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.DataSource], error)
}
