package routes

import (
	"mapexIam/src/modules/lists/application/dtos"
	"mapexIam/src/modules/lists/application/ports"
	"mapexIam/src/modules/lists/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers list HTTP routes. Base path: /api/v1/lists.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper.
func RegisterRoutes(group web.Router, service ports.ListServicePort) {

	r := swagger.Wrap(group).Tag("Lists")

	listQueryDto := validation.NewValidation(nil, &dtos.ListQueryDTO{}, nil)
	r.Get("/", listQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.ListList),
		coverageMw.InjectRequestContext(),
		handlers.GetLists(service),
	).
		Summary("List lists").
		Description("Returns a paginated, filterable list of reusable lists scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.ListResponse]{})

	listCreateDto := validation.NewValidation(&dtos.ListCreateDTO{}, nil, nil)
	r.Post("/", listCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.ListCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateList(service),
	).
		Summary("Create list").
		Description("Creates a new reusable list (key/value option set).").
		Returns(&dtos.ListResponse{})

	getListById := validation.NewValidation(nil, nil, &dtos.ListIdDTO{})
	r.Get("/:listId", getListById, swagger.Expose,
		permissionMw.RequirePermission(perms.ListRead),
		handlers.GetListById(service),
	).
		Summary("Get list by ID").
		Description("Retrieves a single list by its MongoDB ObjectId.").
		Returns(&dtos.ListResponse{})

	updateListById := validation.NewValidation(&dtos.ListUpdateDTO{}, nil, &dtos.ListIdDTO{})
	r.Patch("/:listId", updateListById, swagger.Expose,
		permissionMw.RequirePermission(perms.ListUpdate),
		handlers.UpdateListById(service),
	).
		Summary("Update list").
		Description("Partially updates an existing list. All body fields are optional.").
		Returns(&dtos.ListResponse{})

	deleteListById := validation.NewValidation(nil, nil, &dtos.ListIdDTO{})
	r.Delete("/:listId", deleteListById, swagger.Expose,
		permissionMw.RequirePermission(perms.ListDelete),
		handlers.DeleteListById(service),
	).
		Summary("Delete list").
		Description("Deletes a list by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
