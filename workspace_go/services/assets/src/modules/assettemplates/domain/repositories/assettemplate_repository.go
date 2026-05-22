package repositories

import (
	"assets/src/modules/assettemplates/domain/entities"
	"context"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type AssetTemplateRepository interface {
	Create(ctx context.Context, u *entities.Assettemplate) (*entities.Assettemplate, error)
	FindById(ctx context.Context, dataSourceId *string) (*entities.Assettemplate, error)
	FindByIdAndUpdate(ctx context.Context, dataSourceId *string, payload map[string]any) (*entities.Assettemplate, error)
	DeleteById(ctx context.Context, dataSourceId *string) error
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.Assettemplate], error)
	UpdateMany(ctx context.Context, filter model.Map, update model.Map) (int64, error)

	// CountDocuments counts documents matching the provided filters.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellation signals
	//   - filters: A map of filters to apply to the count query
	//
	// Returns:
	//   - int64: The number of matching documents
	//   - error: If the count operation fails
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
