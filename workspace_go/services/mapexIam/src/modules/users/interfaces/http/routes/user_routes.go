package routes

import (
	"mapexIam/src/modules/users/application/dtos"
	"mapexIam/src/modules/users/application/ports"
	"mapexIam/src/modules/users/interfaces/http/handlers"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers user HTTP routes. Base path: /api/v1/users.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper.
func RegisterRoutes(group web.Router, service ports.UserServicePort) {

	r := swagger.Wrap(group).Tag("Users")

	/** Self ("me") routes */

	noContract := validation.NewValidation(nil, nil, nil)
	r.Get("/me", noContract, swagger.Expose, handlers.Myself(service)).
		Summary("Get my profile").
		Description("Returns the authenticated caller's own user profile.").
		Returns(&dtos.UserResponse{})

	userUpdateDto := validation.NewValidation(&dtos.UserUpdateDTO{}, nil, nil)
	r.Patch("/me", userUpdateDto, swagger.Expose, handlers.UpdateMyself(service)).
		Summary("Update my profile").
		Description("Partially updates the authenticated caller's own user profile.").
		Returns(&dtos.UserResponse{})

	r.Patch("/me/tour", noContract, swagger.Expose, handlers.DisableMyTour(service)).
		Summary("Disable my onboarding tour").
		Description("Marks the authenticated caller's onboarding tour as completed so it is not shown again.").
		Returns(&dtos.UserResponse{})

	/** List + Counter */

	getUsersQuery := validation.NewValidation(nil, &dtos.UserQueryDto{}, nil)
	r.Get("/", getUsersQuery, swagger.Expose,
		permissionMw.RequirePermission(perms.UserList),
		coverageMw.InjectRequestContext(),
		handlers.GetUsers(service),
	).
		Summary("List users").
		Description("Returns a paginated, filterable list of users scoped to the caller's organization.").
		Returns(&model.PaginatedResult[dtos.UserResponse]{})

	r.Get("/counter", noContract, swagger.Expose,
		permissionMw.RequirePermission(perms.UserList),
		coverageMw.InjectRequestContext(),
		handlers.GetUserCount(service),
	).
		Summary("Count users").
		Description("Returns the total number of users for the caller's organization.").
		Returns(&contractsCommon.CounterResponse{})

	/** CRUD */

	userCreateDto := validation.NewValidation(&dtos.UserCreateDTO{}, nil, nil)
	r.Post("/", userCreateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.UserCreate),
		handlers.CreateUser(service),
	).
		Summary("Create user").
		Description("Creates a new user.").
		Returns(&dtos.UserResponse{})

	getUserById := validation.NewValidation(nil, nil, &dtos.UserIdDTO{})
	r.Get("/:userId", getUserById, swagger.Expose,
		permissionMw.RequirePermission(perms.UserRead),
		handlers.GetUserById(service),
	).
		Summary("Get user by ID").
		Description("Retrieves a single user by its MongoDB ObjectId.").
		Returns(&dtos.UserResponse{})

	updateUserById := validation.NewValidation(&dtos.UserUpdateDTO{}, nil, &dtos.UserIdDTO{})
	r.Patch("/:userId", updateUserById, swagger.Expose,
		permissionMw.RequirePermission(perms.UserUpdate),
		handlers.UpdateUserById(service),
	).
		Summary("Update user").
		Description("Partially updates an existing user. All body fields are optional.").
		Returns(&dtos.UserResponse{})

	deleteUserById := validation.NewValidation(nil, nil, &dtos.UserIdDTO{})
	r.Delete("/:userId", deleteUserById, swagger.Expose,
		permissionMw.RequirePermission(perms.UserDelete),
		handlers.DeleteUserById(service),
	).
		Summary("Delete user").
		Description("Deletes a user by its MongoDB ObjectId.").
		Returns(map[string]bool{})
}
