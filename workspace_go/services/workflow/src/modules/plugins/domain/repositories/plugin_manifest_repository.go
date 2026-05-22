package repositories

import (
	"context"

	"workflow/src/modules/plugins/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// PluginManifestRepository defines the persistence contract for PluginManifest entities.
type PluginManifestRepository interface {
	Create(ctx context.Context, entity *entities.PluginManifest) (*entities.PluginManifest, error)
	FindById(ctx context.Context, id *string) (*entities.PluginManifest, error)
	FindByPluginId(ctx context.Context, pluginId string) (*entities.PluginManifest, error)
	FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.PluginManifest, error)
	DeleteById(ctx context.Context, id *string) error
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.PluginManifest], error)
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
