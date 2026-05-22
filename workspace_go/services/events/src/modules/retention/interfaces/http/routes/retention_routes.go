package routes

import (
	"github.com/gofiber/fiber/v2"

	"events/src/modules/retention/application/dtos"
	"events/src/modules/retention/application/ports"
	"events/src/modules/retention/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/events"
)

// RegisterRoutes registers retention policy HTTP routes.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation.
//
// Base path: /api/v1/retention
//
// HTTP Verbs follow REST conventions:
//
//	GET    /                       - List retention policies (paginated, filtered)
//	GET    /:retentionPolicyId     - Get retention policy by ID
//	PUT    /                       - Upsert retention policy (by org + type)
//	DELETE /:retentionPolicyId     - Delete retention policy
//
// Parameters:
//   - group: Fiber router group to register routes on
//   - service: Retention policy service port interface implementation
func RegisterRoutes(group fiber.Router, service ports.RetentionServicePort) {

	// List Routes — get retention policies with filters, pagination, and projection
	retentionQueryDto := validation.NewValidation(nil, &dtos.RetentionPolicyQueryDTO{}, nil)
	group.Get("/",
		validation.ValidationMiddleware(retentionQueryDto),     // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.RetentionList),    // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                      // 3. Inject context (cache)
		handlers.GetRetentionPolicies(service),                 // 4. Handler
	)

	// CRUD Routes — upsert retention policy (create or update by org + type)
	upsertDto := validation.NewValidation(&dtos.RetentionPolicyUpsertDTO{}, nil, nil)
	group.Put("/",
		validation.ValidationMiddleware(upsertDto),             // 1. Validate DTO
		permissionMw.RequirePermission(perms.RetentionUpdate),  // 2. Check permission
		coverageMw.InjectRequestContext(),                      // 3. Inject context
		handlers.UpsertRetentionPolicy(service),                // 4. Handler
	)

	// Get retention policy by ID
	getByIdDto := validation.NewValidation(nil, nil, &dtos.RetentionPolicyParamsDTO{})
	group.Get("/:retentionPolicyId",
		validation.ValidationMiddleware(getByIdDto),            // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.RetentionRead),    // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                      // 3. Inject context (cache)
		handlers.GetRetentionPolicyById(service),               // 4. Handler
	)

	// Delete retention policy by ID
	deleteByIdDto := validation.NewValidation(nil, nil, &dtos.RetentionPolicyParamsDTO{})
	group.Delete("/:retentionPolicyId",
		validation.ValidationMiddleware(deleteByIdDto),         // 1. Validate DTO first (fail fast)
		permissionMw.RequirePermission(perms.RetentionUpdate),  // 2. Check permission (cache)
		coverageMw.InjectRequestContext(),                      // 3. Inject context (cache)
		handlers.DeleteRetentionPolicyById(service),            // 4. Handler
	)
}
