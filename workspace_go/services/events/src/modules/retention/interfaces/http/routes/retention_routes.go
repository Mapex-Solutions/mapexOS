package routes

import (
	"events/src/modules/retention/application/dtos"
	"events/src/modules/retention/application/ports"
	"events/src/modules/retention/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/events"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers retention policy HTTP routes. Base path: /api/v1/retention.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper: each NewValidation declares the input contract once (used to both
// validate and document), the module tag is declared once on Wrap, and each route's
// summary, description, and response type are attached via the fluent builder.
func RegisterRoutes(group web.Router, service ports.RetentionServicePort) {

	r := swagger.Wrap(group).Tag("Retention")

	// List retention policies (paginated, filtered).
	retentionQueryDto := validation.NewValidation(nil, &dtos.RetentionPolicyQueryDTO{}, nil)
	r.Get("/", retentionQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.RetentionList),
		coverageMw.InjectRequestContext(),
		handlers.GetRetentionPolicies(service),
	).
		Summary("List retention policies").
		Description("Returns a paginated, filterable list of retention policies scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.RetentionPolicyResponse]{})

	// Upsert a retention policy (create or update by org + type).
	upsertDto := validation.NewValidation(&dtos.RetentionPolicyUpsertDTO{}, nil, nil)
	r.Put("/", upsertDto, swagger.Expose,
		permissionMw.RequirePermission(perms.RetentionUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpsertRetentionPolicy(service),
	).
		Summary("Upsert retention policy").
		Description("Creates or updates the retention policy for the caller's organization and the given event type (one policy per org+type).").
		Returns(&dtos.RetentionPolicyResponse{})

	// Get retention policy by ID.
	getByIdDto := validation.NewValidation(nil, nil, &dtos.RetentionPolicyParamsDTO{})
	r.Get("/:retentionPolicyId", getByIdDto, swagger.Expose,
		permissionMw.RequirePermission(perms.RetentionRead),
		coverageMw.InjectRequestContext(),
		handlers.GetRetentionPolicyById(service),
	).
		Summary("Get retention policy by ID").
		Description("Retrieves a single retention policy by its MongoDB ObjectId.").
		Returns(&dtos.RetentionPolicyResponse{})

	// Delete retention policy by ID.
	deleteByIdDto := validation.NewValidation(nil, nil, &dtos.RetentionPolicyParamsDTO{})
	r.Delete("/:retentionPolicyId", deleteByIdDto, swagger.Expose,
		permissionMw.RequirePermission(perms.RetentionUpdate),
		coverageMw.InjectRequestContext(),
		handlers.DeleteRetentionPolicyById(service),
	).
		Summary("Delete retention policy").
		Description("Deletes a retention policy by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
