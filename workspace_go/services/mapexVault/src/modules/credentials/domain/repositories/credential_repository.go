package repositories

import (
	"context"

	"mapexVault/src/modules/credentials/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// CredentialRepository defines the persistence interface for credential entities.
type CredentialRepository interface {
	// Create persists a new credential.
	Create(ctx context.Context, credential *entities.Credential) (*entities.Credential, error)

	// FindById retrieves a credential by its MongoDB _id.
	FindById(ctx context.Context, id *string) (*entities.Credential, error)

	// FindByIdAndUpdate updates a credential by its MongoDB _id.
	FindByIdAndUpdate(ctx context.Context, id *string, update map[string]any) (*entities.Credential, error)

	// DeleteById removes a credential by its MongoDB _id.
	DeleteById(ctx context.Context, id *string) error

	// FindWithFilters queries credentials with org filter and pagination.
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, sort model.Map) (*model.PaginatedResult[entities.Credential], error)

	// FindActiveWithTokenExpiry returns all active oauth2/userAndPass credentials
	// with non-nil tokenExpiresAt. Used by bootstrap seed to publish initial schedules.
	FindActiveWithTokenExpiry(ctx context.Context) ([]entities.Credential, error)

	// CountDocuments returns the count of documents matching the filter.
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
