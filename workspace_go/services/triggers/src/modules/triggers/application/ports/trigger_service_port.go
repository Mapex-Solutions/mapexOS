package ports

import (
	ctx "context"

	"triggers/src/modules/triggers/application/dtos"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// TriggerServicePort defines the contract for trigger business operations.
//
// This port interface enables Hexagonal Architecture by decoupling the business
// logic from its implementation. Handlers and routes depend on this interface
// rather than the concrete service implementation.
//
// Methods:
//   - CreateTrigger: Creates a new trigger
//   - GetTriggerById: Retrieves a trigger by its ID
//   - UpdateTriggerById: Updates an existing trigger
//   - DeleteTriggerById: Deletes a trigger
//   - GetTriggers: Retrieves a paginated list of triggers with filters
type TriggerServicePort interface {
	// CreateTrigger creates a new trigger from the provided DTO.
	// Uses RequestContext to populate orgId and pathKey for multi-tenant support.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains OrgContext (selected orgId) and OrgContextData (pathKey)
	//   - dto: Trigger creation data
	//
	// Returns:
	//   - TriggerResponse: The created trigger
	//   - error: If creation fails
	CreateTrigger(ctx ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateTriggerDto) (*dtos.TriggerResponse, error)

	// GetTriggerById retrieves a trigger by its unique identifier.
	//
	// This method uses a cache-aside pattern: checks cache first, then DB if cache miss.
	// Optionally accepts a *CacheMetrics pointer to report hit/miss info (for Prometheus).
	// Existing callers without metrics continue to work unchanged.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - triggerId: The unique identifier of the trigger
	//   - metrics: Optional CacheMetrics pointer — populated with hit/miss when provided
	//
	// Returns:
	//   - TriggerResponse: The trigger if found
	//   - error: If not found or retrieval fails
	GetTriggerById(ctx ctx.Context, triggerId *string, metrics ...*common.CacheMetrics) (*dtos.TriggerResponse, error)

	// UpdateTriggerById updates an existing trigger.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains UserContext for updatedBy field
	//   - triggerId: The unique identifier of the trigger
	//   - dto: Update data (only non-nil fields will be updated)
	//
	// Returns:
	//   - TriggerResponse: The updated trigger
	//   - error: If update fails or trigger not found
	UpdateTriggerById(ctx ctx.Context, requestContext *reqCtx.RequestContext, triggerId *string, dto *dtos.UpdateTriggerDto) (*dtos.TriggerResponse, error)

	// GetTriggers retrieves a paginated and filtered list of triggers.
	// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
	//
	// This method supports multi-tenant isolation by filtering results based on the
	// RequestContext parameter provided by the InjectRequestContext middleware.
	// Only triggers belonging to organizations accessible by the user will be returned.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Request context with org access data (from coverage middleware)
	//   - query: Filters, pagination, and projection options
	//
	// Returns:
	//   - PaginatedResult: Matching triggers and pagination metadata
	//   - error: If query fails
	GetTriggers(ctx ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.TriggerQueryDto) (*model.PaginatedResult[dtos.TriggerResponse], error)

	// DeleteTriggerById removes a trigger by its unique identifier.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - triggerId: The unique identifier of the trigger to delete
	//
	// Returns:
	//   - map[string]bool: Success indicator ({"success": true})
	//   - error: If deletion fails or trigger not found
	DeleteTriggerById(ctx ctx.Context, triggerId *string) (map[string]bool, error)

	// CountTriggers returns the total count of triggers for the given org context.
	// Implements cache-aside: check Redis first, fallback to MongoDB CountDocuments.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains org access data from coverage middleware
	//
	// Returns:
	//   - int64: Total count of matching triggers
	//   - error: If query fails
	CountTriggers(ctx ctx.Context, requestContext *reqCtx.RequestContext) (int64, error)
}
