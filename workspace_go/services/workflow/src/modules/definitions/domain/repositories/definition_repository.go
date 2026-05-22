package repositories

import (
	"context"

	"workflow/src/modules/definitions/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// DefinitionRepository defines the persistence contract for WorkflowDefinition entities.
type DefinitionRepository interface {
	// Create inserts a new WorkflowDefinition entity.
	Create(ctx context.Context, def *entities.WorkflowDefinition) (*entities.WorkflowDefinition, error)

	// FindById retrieves a WorkflowDefinition by its MongoDB ID.
	FindById(ctx context.Context, id *string) (*entities.WorkflowDefinition, error)

	// FindByIdAndUpdate updates a WorkflowDefinition and returns the updated document.
	FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.WorkflowDefinition, error)

	// DeleteById removes a WorkflowDefinition by its MongoDB ID.
	DeleteById(ctx context.Context, id *string) error

	// FindWithFilters retrieves a paginated list of WorkflowDefinition entities.
	FindWithFilters(
		ctx context.Context,
		filters model.Map,
		pagination *model.PaginationOpts,
		projection model.Map,
	) (*model.PaginatedResult[entities.WorkflowDefinition], error)

	// CountDocuments counts documents matching the provided filters.
	CountDocuments(ctx context.Context, filters model.Map) (int64, error)
}
