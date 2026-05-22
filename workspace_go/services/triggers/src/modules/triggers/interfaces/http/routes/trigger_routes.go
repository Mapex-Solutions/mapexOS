package routes

import (
	"github.com/gofiber/fiber/v2"

	"triggers/src/modules/triggers/application/dtos"
	"triggers/src/modules/triggers/application/ports"
	"triggers/src/modules/triggers/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	permissions "github.com/Mapex-Solutions/MapexOS/permissions/triggers"
)

// RegisterRoutes registers trigger HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/v1/triggers
//
// HTTP Verbs follow REST conventions:
//   GET    /                - List triggers (paginated, filtered)
//   POST   /                - Create trigger
//   GET    /:id             - Get trigger by ID
//   PATCH  /:id             - Update trigger
//   DELETE /:id             - Delete trigger
//
// Middleware chain:
//   1. ValidationMiddleware - Validates and parses request body/query/params
//   2. RequirePermission - Checks user has required permission (uses cache)
//   3. InjectRequestContext - Injects RequestContext with org filtering (uses coverage cache)
//   4. Handler - Executes business logic via service port
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: Trigger service port interface implementation
func RegisterRoutes(group fiber.Router, service ports.TriggerServicePort) {

	/**
	* List Routes
	 */

	// Get triggers with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	// Includes hierarchical support via PathKey data from coverage cache
	triggerQueryDto := validation.NewValidation(nil, &dtos.TriggerQueryDto{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(triggerQueryDto),        // Validate DTO first (fail fast)
		permissionMw.RequirePermission(permissions.TriggerList), // Check permission (cache)
		coverageMw.InjectRequestContext(),                       // Inject context (cache)
		handlers.GetTriggers(service),                           // Handler
	)

	// Count triggers (counter endpoint with cache)
	group.Get("/counter",
		permissionMw.RequirePermission(permissions.TriggerList),
		coverageMw.InjectRequestContext(),
		handlers.GetTriggerCount(service),
	)

	/**
	* CRUD Routes
	 */

	// Create a new trigger
	triggerCreateDto := validation.NewValidation(&dtos.CreateTriggerDto{}, nil, nil)
	group.Post("/",
		validation.ValidationMiddleware(triggerCreateDto),         // Validate DTO
		permissionMw.RequirePermission(permissions.TriggerCreate), // Check permission
		coverageMw.InjectRequestContext(),                         // Inject context
		handlers.CreateTrigger(service),                           // Handler
	)

	// Get trigger by ID
	group.Get("/:id",
		permissionMw.RequirePermission(permissions.TriggerRead),
		handlers.GetTriggerById(service),
	)

	// Update trigger by ID
	triggerUpdateDto := validation.NewValidation(&dtos.UpdateTriggerDto{}, nil, nil)
	group.Patch("/:id",
		validation.ValidationMiddleware(triggerUpdateDto),          // Validate DTO first (fail fast)
		permissionMw.RequirePermission(permissions.TriggerUpdate), // Check permission (cache)
		coverageMw.InjectRequestContext(),                         // Need context for updatedBy
		handlers.UpdateTriggerById(service),
	)

	// Delete trigger by ID
	group.Delete("/:id",
		permissionMw.RequirePermission(permissions.TriggerDelete),
		handlers.DeleteTriggerById(service),
	)
}
