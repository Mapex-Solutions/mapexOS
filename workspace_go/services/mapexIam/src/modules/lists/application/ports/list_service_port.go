package ports

import (
	"context"

	"mapexIam/src/modules/lists/application/dtos"
	"mapexIam/src/modules/lists/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// ListServicePort defines the output port (Hexagonal Architecture) for list-related operations.
// This port allows other application services and interface layers to depend on
// list operations without coupling to the concrete ListService implementation.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: This interface (defines the contract)
//   - Adapter: ListService (implements the contract)
//
// Used by:
//   - HTTP Handlers (for HTTP interface layer)
type ListServicePort interface {
	// CreateList creates a new list entity.
	// Uses RequestContext to populate orgId and pathKey for multi-tenant support.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - requestContext: Contains OrgContext (selected orgId) and OrgContextData (pathKey)
	//   - dto: List creation data
	//
	// Returns:
	//   - *dtos.ListResponse: Created list data
	//   - error: Error if creation fails
	CreateList(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.ListCreateDTO) (*dtos.ListResponse, error)

	// GetListById retrieves a list entity by its unique identifier.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - listId: Unique identifier of the list
	//
	// Returns:
	//   - *dtos.ListResponse: List data DTO
	//   - error: Error if list not found or database error
	GetListById(ctx context.Context, listId *string) (*dtos.ListResponse, error)

	// UpdateListById updates an existing list entity.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - listId: Unique identifier of the list
	//   - dto: Fields to update
	//
	// Returns:
	//   - *dtos.ListResponse: Updated list data
	//   - error: Error if update fails
	UpdateListById(ctx context.Context, listId *string, dto *dtos.ListUpdateDTO) (*dtos.ListResponse, error)

	// DeleteListById removes a list entity by its unique ID.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - listId: Unique identifier of the list
	//
	// Returns:
	//   - map[string]bool: Success status
	//   - error: Error if deletion fails
	DeleteListById(ctx context.Context, listId *string) (map[string]bool, error)

	// GetListByEmail retrieves a list entity by its email address.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - email: Email address of the list
	//
	// Returns:
	//   - *entities.List: List entity
	//   - error: Error if not found or database error
	GetListByEmail(ctx context.Context, email *string) (*entities.List, error)

	// GetLists retrieves a paginated and filtered list of lists.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Filters, pagination, and projection options
	//
	// Returns:
	//   - *model.PaginatedResult[dtos.ListResponse]: Paginated lists with DTOs
	//   - error: Error if query fails
	GetLists(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.ListQueryDTO) (*model.PaginatedResult[dtos.ListResponse], error)
}
