package ports

import (
	ctx "context"
	"http_gateway/src/modules/datasources/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// DataSourceServicePort defines the contract for data source business operations.
//
// This port interface enables Hexagonal Architecture by decoupling the business
// logic from its implementation. Handlers and routes depend on this interface
// rather than the concrete service implementation.
//
// Methods:
//   - GetDataSources: Lists data sources with filters, pagination, and projection
//   - CreateDataSource: Creates a new data source
//   - GetDataSourceById: Retrieves a data source by its ID
//   - UpdateDataSourceById: Updates an existing data source
//   - DeleteDataSourceById: Deletes a data source
type DataSourceServicePort interface {
	// GetDataSources retrieves a paginated and filtered list of data sources.
	// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
	//
	// The method receives RequestContext which provides:
	//   - ScopedOrgIds: All accessible organization IDs for the user
	//   - OrgContext: Optional org filter from X-Org-Context header
	//   - OrgContextData: Detailed org data with PathKey for hierarchical queries
	//   - CoverageOrgs: Full coverage data with organizational hierarchy
	//
	// Query supports hierarchical filtering via includeChildren parameter:
	//   - OrgContext + includeChildren=true: Returns org and all descendants (PathKey range)
	//   - OrgContext + includeChildren=false: Returns specific org only
	//   - No OrgContext: Returns all accessible orgs
	//
	// Module-specific filters from query:
	//   - Name: Partial match on data source name
	//   - Enabled: Filter by enabled status
	//   - Mode: Filter by mode (pull/push)
	//   - Protocol: Filter by protocol (http/mqtt)
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Context from coverage middleware with org filtering data
	//   - query: Query parameters with filters, pagination, projection, hierarchy
	//
	// Returns:
	//   - PaginatedResult: Contains data sources and pagination metadata
	//   - error: If query fails
	GetDataSources(ctx ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.DataSourceQueryDTO) (*model.PaginatedResult[dtos.DataSourceResponse], error)

	// CreateDataSource creates a new data source from the provided DTO.
	// Uses RequestContext to populate orgId and pathKey for multi-tenant support.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains OrgContext (selected orgId) and OrgContextData (pathKey)
	//   - dto: Data source creation data
	//
	// Returns:
	//   - DataSourceResponse: The created data source
	//   - error: If creation fails
	CreateDataSource(ctx ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.DataSourceCreateDTO) (*dtos.DataSourceResponse, error)

	// GetDataSourceById retrieves a data source by its unique identifier.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - dataSourceId: The unique identifier of the data source
	//
	// Returns:
	//   - DataSourceResponse: The data source if found
	//   - error: If not found or retrieval fails
	GetDataSourceById(ctx ctx.Context, dataSourceId *string) (*dtos.DataSourceResponse, error)

	// UpdateDataSourceById updates an existing data source.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - dataSourceId: The unique identifier of the data source
	//   - dto: Update data (only non-nil fields will be updated)
	//
	// Returns:
	//   - DataSourceResponse: The updated data source
	//   - error: If update fails or data source not found
	UpdateDataSourceById(ctx ctx.Context, dataSourceId *string, dto *dtos.DataSourceUpdateDTO) (*dtos.DataSourceResponse, error)

	// DeleteDataSourceById removes a data source by its unique identifier.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - dataSourceId: The unique identifier of the data source to delete
	//
	// Returns:
	//   - map[string]bool: Success indicator ({"success": true})
	//   - error: If deletion fails or data source not found
	DeleteDataSourceById(ctx ctx.Context, dataSourceId *string) (map[string]bool, error)
}
