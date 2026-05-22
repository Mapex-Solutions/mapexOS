package repositories

import (
	"assets/src/modules/assets/domain/entities"
	"context"
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type AssetRepository interface {
	Create(ctx context.Context, u *entities.Asset) (*entities.Asset, error)
	FindById(ctx context.Context, assetId *string) (*entities.Asset, error)
	FindByAssetUUID(ctx context.Context, assetUUID *string) (*entities.Asset, error)
	FindByMqttUsername(ctx context.Context, username string) (*entities.Asset, error)
	FindByIdAndUpdate(ctx context.Context, assetId *string, payload map[string]any) (*entities.Asset, error)
	DeleteById(ctx context.Context, assetId *string) error
	FindWithFilters(
		ctx context.Context,
		filters model.Map,
		pagination *model.PaginationOpts,
		projection model.Map,
	) (*model.PaginatedResult[entities.Asset], error)

	// FindWithFiltersAndTemplate retrieves assets with template data joined via $lookup aggregation.
	// This method performs a single MongoDB query that joins assets with their templates,
	// avoiding the N+1 query problem and improving performance significantly.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellation signals
	//   - assetFilters: MongoDB filters to apply to assets (orgId, status, name, etc)
	//   - templateFilters: Filters to apply to the joined template data (categoryId, manufacturerId, modelId)
	//   - pagination: Pagination options (page, perPage)
	//   - sort: Sort options (field and direction)
	//
	// Returns:
	//   - PaginatedResult with AssetWithTemplate entities (assets + template classification data)
	//   - Error if query fails
	FindWithFiltersAndTemplate(
		ctx context.Context,
		assetFilters model.Map,
		templateFilters model.Map,
		pagination *model.PaginationOpts,
		sort model.Map,
	) (*model.PaginatedResult[entities.AssetWithTemplate], error)

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

	// UpdateHealthStatusWithChangedAt atomically updates healthStatus and healthStatusChangedAt
	// for an asset identified by assetUUID. Both fields are written in a single $set so a
	// caller that receives nil error is guaranteed the pair flipped together.
	UpdateHealthStatusWithChangedAt(ctx context.Context, assetUUID *string, status string, changedAt time.Time) error
}
