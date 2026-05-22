package handlers

import (
	"github.com/gofiber/fiber/v2"

	"http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/datasources/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetDataSources returns a Fiber handler that retrieves a paginated and filtered list of data sources.
// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
//
// The handler extracts RequestContext from c.Locals("requestContext") which was set by InjectRequestContext middleware.
// This provides access to:
//   - ScopedOrgIds: All accessible organization IDs
//   - OrgContext: Optional org filter from X-Org-Context header
//   - OrgContextData: Detailed org data with PathKey
//   - CoverageOrgs: Full coverage data with hierarchical information
//
// Query supports hierarchical filtering via includeChildren parameter:
//   - OrgContext + includeChildren=true: Returns org and all descendants (PathKey range)
//   - OrgContext + includeChildren=false: Returns specific org only
//   - No OrgContext: Returns all accessible orgs
//
// It expects a validated DTO of type dtos.DataSourceQueryDTO in the "queryDTO" context key
// (populated by requestValidation middleware) containing optional filters such as:
//   - name, enabled, mode, protocol (filters)
//   - page, perPage (pagination)
//   - projection (field selection)
//   - includeChildren (hierarchical query flag)
//
// Returns:
//   - 200 OK with paginated data source list
//   - 400 Bad Request if query validation fails
//   - 500 Internal Server Error on service failure or requestContext not found
func GetDataSources(service ports.DataSourceServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.DataSourceQueryDTO](c, "queryDTO")
		retData, err := service.GetDataSources(ctx, requestContext, queryData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// CreateDataSource returns a Fiber handler that creates a new data source.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It uses RequestContext (injected by coverage middleware) which contains:
//   - OrgContext: The selected organization ID from X-Org-Context header
//   - OrgContextData: Organization data including PathKey for hierarchical filtering
//
// The handler passes the full RequestContext to the service layer, which extracts
// the needed fields (orgId, pathKey) for multi-tenant support.
//
// It expects a validated DTO of type dtos.DataSourceCreateDTO to be stored
// in the Fiber context under the key "bodyDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The DataSourceServicePort interface for data source business operations
//
// Returns:
//   - A Fiber handler function that processes the data source creation request
func CreateDataSource(service ports.DataSourceServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.DataSourceCreateDTO](c, "bodyDTO")

		// Pass requestContext to service (contains OrgContext and OrgContextData)
		retData, err := service.CreateDataSource(ctx, requestContext, bodyData)

		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// GetDataSourceById returns a Fiber handler that retrieves a data source by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects a validated DTO of type dtos.DataSourceIdDto to be stored
// in the Fiber context under the key "paramsDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The DataSourceServicePort interface for data source business operations
//
// Returns:
//   - A Fiber handler function that processes the data source retrieval request
func GetDataSourceById(service ports.DataSourceServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		dataSource, _ := requestValidation.GetDTO[*dtos.DataSourceIdDto](c, "paramsDTO")
		retData, err := service.GetDataSourceById(ctx, &dataSource.DataSourceId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// UpdateDataSourceById returns a Fiber handler that updates a data source by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects two validated DTOs:
//   - dtos.DataSourceIdDto stored in the Fiber context under the key "paramsDTO"
//   - dtos.DataSourceUpdateDTO stored in the Fiber context under the key "bodyDTO"
//
// (Both are usually populated by requestValidation middleware)
//
// Parameters:
//   - service: The DataSourceServicePort interface for data source business operations
//
// Returns:
//   - A Fiber handler function that processes the data source update request
func UpdateDataSourceById(service ports.DataSourceServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		dataSource, _ := requestValidation.GetDTO[*dtos.DataSourceIdDto](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.DataSourceUpdateDTO](c, "bodyDTO")
		retData, err := service.UpdateDataSourceById(ctx, &dataSource.DataSourceId, bodyData)

		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// DeleteDataSourceById returns a Fiber handler that deletes a data source by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects a validated DTO of type dtos.DataSourceIdDto to be stored
// in the Fiber context under the key "paramsDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The DataSourceServicePort interface for data source business operations
//
// Returns:
//   - A Fiber handler function that processes the data source deletion request
func DeleteDataSourceById(service ports.DataSourceServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		dataSource, _ := requestValidation.GetDTO[*dtos.DataSourceIdDto](c, "paramsDTO")
		retData, err := service.DeleteDataSourceById(ctx, &dataSource.DataSourceId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}
