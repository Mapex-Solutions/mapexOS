package repositories

import (
	"context"

	"events/src/modules/retention/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// RetentionRepository defines the contract for retention policy persistence operations.
type RetentionRepository interface {
	// Create inserts a new retention policy into the repository.
	Create(ctx context.Context, policy *entities.RetentionPolicy) (*entities.RetentionPolicy, error)

	// FindById retrieves a retention policy by its unique identifier.
	FindById(ctx context.Context, policyId *string) (*entities.RetentionPolicy, error)

	// FindByOrgIdAndType retrieves a retention policy by organization ID and type.
	// Used by cache-aside pattern for per-type lookups.
	FindByOrgIdAndType(ctx context.Context, orgId *string, retentionType string) (*entities.RetentionPolicy, error)

	// Upsert creates or updates a retention policy by organization ID and type.
	// Uses FindOneAndUpdate with upsert=true on filter {orgId, type}.
	Upsert(ctx context.Context, orgId *model.ObjectId, retentionType string, policy *entities.RetentionPolicy) (*entities.RetentionPolicy, error)

	// DeleteById removes a retention policy by its ID.
	DeleteById(ctx context.Context, policyId *string) error

	// FindWithFilters retrieves a paginated list of retention policies based on filters.
	FindWithFilters(ctx context.Context, filters model.Map, pagination *model.PaginationOpts, projection model.Map) (*model.PaginatedResult[entities.RetentionPolicy], error)
}
