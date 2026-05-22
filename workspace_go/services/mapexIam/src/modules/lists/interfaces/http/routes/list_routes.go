package routes

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/lists/application/dtos"
	"mapexIam/src/modules/lists/application/ports"
	"mapexIam/src/modules/lists/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
)

func RegisterRoutes(group fiber.Router, service ports.ListServicePort) {

	/**
	* CRUD Routes
	 */

	// Get lists with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	listQueryDto := validation.NewValidation(nil, &dtos.ListQueryDTO{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(listQueryDto),  // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.ListList), // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),              // 3. Inject context (cache)
		handlers.GetLists(service),                     // 4. Handler
	)

	// Create a new list
	listCreateDto := validation.NewValidation(&dtos.ListCreateDTO{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(listCreateDto),   // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.ListCreate), // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                // 3. Inject context (cache)
		handlers.CreateList(service),                     // 4. Handler
	)

	// Get list by ID
	getListById := validation.NewValidation(nil, nil, &dtos.ListIdDTO{})
	group.Get("/:listId",
		validation.ValidationMiddleware(getListById),
		permissionMw.RequirePermission(perms.ListRead),
		handlers.GetListById(service),
	)

	// Update list by ID
	updateListById := validation.NewValidation(&dtos.ListUpdateDTO{}, nil, &dtos.ListIdDTO{})
	group.Patch("/:listId",
		validation.ValidationMiddleware(updateListById),
		permissionMw.RequirePermission(perms.ListUpdate),
		handlers.UpdateListById(service),
	)

	// Delete list by ID
	deleteListById := validation.NewValidation(nil, nil, &dtos.ListIdDTO{})
	group.Delete("/:listId",
		validation.ValidationMiddleware(deleteListById),
		permissionMw.RequirePermission(perms.ListDelete),
		handlers.DeleteListById(service),
	)
}
