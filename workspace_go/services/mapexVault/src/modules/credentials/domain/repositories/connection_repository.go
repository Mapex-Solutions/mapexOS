package repositories

import (
	"context"

	"mapexVault/src/modules/credentials/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// ConnectionRepository defines the persistence interface for connection entities.
type ConnectionRepository interface {
	// Create persists a new connection.
	Create(ctx context.Context, connection *entities.Connection) (*entities.Connection, error)

	// FindById retrieves a connection by its MongoDB _id.
	FindById(ctx context.Context, id *string) (*entities.Connection, error)

	// FindByIdAndUpdate updates a connection by its MongoDB _id.
	FindByIdAndUpdate(ctx context.Context, id *string, update map[string]any) (*entities.Connection, error)

	// DeleteById removes a connection by its MongoDB _id.
	DeleteById(ctx context.Context, id *string) error

	// FindWithFilters queries connections with org filter and pagination.
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, sort model.Map) (*model.PaginatedResult[entities.Connection], error)

	// UpsertByAccount creates or updates a connection by provider + accountId + orgId.
	// Used when a user reconnects the same external account — overwrites instead of duplicating.
	UpsertByAccount(ctx context.Context, provider string, accountId string, orgId *model.ObjectId, connection *entities.Connection) (*entities.Connection, error)

	// CountDocuments returns the count of documents matching the filter.
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
