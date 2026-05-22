package routes

import (
	"github.com/gofiber/fiber/v2"

	"http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/datasources/application/ports"
	"http_gateway/src/modules/datasources/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/http_gateway"
)

// RegisterRoutes registers data source HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/v1/data_sources
//
// HTTP Verbs follow REST conventions:
//
//	GET    /                     - List data sources (paginated, filtered)
//	POST   /                     - Create data source
//	GET    /:dataSourceId        - Get data source by ID
//	PATCH  /:dataSourceId        - Update data source
//	DELETE /:dataSourceId        - Delete data source
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: Data source service port interface implementation
func RegisterRoutes(group fiber.Router, service ports.DataSourceServicePort) {

	/**
	* List Routes
	 */

	// Get data sources with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	// Includes hierarchical support via PathKey data from coverage cache
	dataSourceQueryDto := validation.NewValidation(nil, &dtos.DataSourceQueryDTO{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(dataSourceQueryDto),  // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.DatasourceList), // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                    // 3. Inject context (cache)
		handlers.GetDataSources(service),                     // 4. Handler
	)

	/**
	* CRUD Routes
	 */

	// Create a new dataSource
	dataSourceCreateDto := validation.NewValidation(&dtos.DataSourceCreateDTO{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(dataSourceCreateDto),   // 1. Validate DTO
		permissionMw.RequirePermission(perms.DatasourceCreate), // 2. Check permission
		coverageMw.InjectRequestContext(),                      // 3. Inject context
		handlers.CreateDataSource(service),                     // 4. Handler
	)

	// Get dataSource by ID
	getDataSourceById := validation.NewValidation(nil, nil, &dtos.DataSourceIdDto{})
	group.Get("/:dataSourceId",
		permissionMw.RequirePermission(perms.DatasourceRead),
		validation.ValidationMiddleware(getDataSourceById),
		handlers.GetDataSourceById(service),
	)

	// Update dataSource by ID
	updateDataSourceById := validation.NewValidation(&dtos.DataSourceUpdateDTO{}, nil, &dtos.DataSourceIdDto{})
	group.Patch("/:dataSourceId",
		permissionMw.RequirePermission(perms.DatasourceUpdate),
		validation.ValidationMiddleware(updateDataSourceById),
		handlers.UpdateDataSourceById(service),
	)

	// Delete dataSource by ID
	deleteDataSourceById := validation.NewValidation(nil, nil, &dtos.DataSourceIdDto{})
	group.Delete("/:dataSourceId",
		permissionMw.RequirePermission(perms.DatasourceDelete),
		validation.ValidationMiddleware(deleteDataSourceById),
		handlers.DeleteDataSourceById(service),
	)
}
