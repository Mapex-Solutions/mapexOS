package repositories

import (
	"context"

	"workflow/src/modules/instances/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// InstanceRepository defines the persistence contract for WorkflowInstance entities.
type InstanceRepository interface {
	// Create inserts a new WorkflowInstance entity.
	Create(ctx context.Context, instance *entities.WorkflowInstance) (*entities.WorkflowInstance, error)

	// FindById retrieves a WorkflowInstance by its MongoDB ID.
	FindById(ctx context.Context, id *string) (*entities.WorkflowInstance, error)

	// FindByIdAndUpdate updates a WorkflowInstance and returns the updated document.
	FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.WorkflowInstance, error)

	// DeleteById removes a WorkflowInstance by its MongoDB ID.
	DeleteById(ctx context.Context, id *string) error

	// FindWithFilters retrieves a paginated list of WorkflowInstance entities.
	FindWithFilters(
		ctx context.Context,
		filters model.Map,
		pagination *model.PaginationOpts,
		projection model.Map,
	) (*model.PaginatedResult[entities.WorkflowInstance], error)

	// CountDocuments counts instances matching the provided filters.
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
