package repositories

import (
	"assets/src/modules/assettemplates/domain/entities"
	"context"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type AssetTemplateRepository interface {
	Create(ctx context.Context, u *entities.Assettemplate) (*entities.Assettemplate, error)
	FindById(ctx context.Context, dataSourceId *string) (*entities.Assettemplate, error)

	// FindByMarketplaceGuidAndOrg returns the caller-org's link record for a
	// marketplace-installed template, or (nil, nil) when that org has not
	// installed it. The (marketplaceGuid, orgId) pair is the per-org install key.
	FindByMarketplaceGuidAndOrg(ctx context.Context, marketplaceGuid string, orgId model.ObjectId) (*entities.Assettemplate, error)

	// FindInstalledGuids returns the subset of the given marketplace guids the
	// caller org has installed, resolved in a single query (marketplaceGuid $in +
	// orgId). Empty input yields an empty result without querying.
	FindInstalledGuids(ctx context.Context, guids []string, orgId model.ObjectId) ([]string, error)
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
