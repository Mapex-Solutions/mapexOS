package routes

import (
	"http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/datasources/application/ports"
	"http_gateway/src/modules/datasources/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/http_gateway"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers data source HTTP routes. Base path: /api/v1/data_sources.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper: each NewValidation declares the input contract once (used to both
// validate and document), the module tag is declared once on Wrap, and each route's
// summary, description, and response type are attached via the fluent builder.
func RegisterRoutes(group web.Router, service ports.DataSourceServicePort) {

	r := swagger.Wrap(group).Tag("Data Sources")

	// List data sources with filters, pagination, and projection. coverage
	// middleware injects context-aware org filtering (hierarchical via PathKey).
	dataSourceQueryDto := validation.NewValidation(nil, &dtos.DataSourceQueryDTO{}, nil)
	r.Get("/", dataSourceQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.DatasourceList),
		coverageMw.InjectRequestContext(),
		handlers.GetDataSources(service),
	).
		Summary("List data sources").
		Description("Returns a paginated, filterable list of data sources scoped to the caller's organization. Supports hierarchical queries via includeChildren.").
		Returns(&model.PaginatedResult[dtos.DataSourceResponse]{})

	// Create a new data source.
	dataSourceCreateDto := validation.NewValidation(&dtos.DataSourceCreateDTO{}, nil, nil)
	r.Post("/", dataSourceCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.DatasourceCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateDataSource(service),
	).
		Summary("Create data source").
		Description("Creates a new data source (ingestion endpoint with its auth method and asset binding). Organization scoping is applied automatically from the request context.").
		Returns(&dtos.DataSourceResponse{})

	// Get data source by ID.
	getDataSourceById := validation.NewValidation(nil, nil, &dtos.DataSourceIdDto{})
	r.Get("/:dataSourceId", getDataSourceById, swagger.Expose,
		permissionMw.RequirePermission(perms.DatasourceRead),
		handlers.GetDataSourceById(service),
	).
		Summary("Get data source by ID").
		Description("Retrieves a single data source by its MongoDB ObjectId.").
		Returns(&dtos.DataSourceResponse{})

	// Update data source by ID.
	updateDataSourceById := validation.NewValidation(&dtos.DataSourceUpdateDTO{}, nil, &dtos.DataSourceIdDto{})
	r.Patch("/:dataSourceId", updateDataSourceById, swagger.Expose,
		permissionMw.RequirePermission(perms.DatasourceUpdate),
		handlers.UpdateDataSourceById(service),
	).
		Summary("Update data source").
		Description("Partially updates an existing data source. All body fields are optional; only provided fields are changed.").
		Returns(&dtos.DataSourceResponse{})

	// Delete data source by ID.
	deleteDataSourceById := validation.NewValidation(nil, nil, &dtos.DataSourceIdDto{})
	r.Delete("/:dataSourceId", deleteDataSourceById, swagger.Expose,
		permissionMw.RequirePermission(perms.DatasourceDelete),
		handlers.DeleteDataSourceById(service),
	).
		Summary("Delete data source").
		Description("Deletes a data source by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
