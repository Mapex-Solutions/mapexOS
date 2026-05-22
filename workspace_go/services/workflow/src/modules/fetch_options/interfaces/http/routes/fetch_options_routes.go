package routes

import (
	"workflow/src/modules/fetch_options/application/ports"
	"workflow/src/modules/fetch_options/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
)

// RegisterRoutes registers fetch-options HTTP routes.
func RegisterRoutes(router fiber.Router, service ports.FetchOptionsServicePort) {
	router.Post("/",
		permissionMw.RequirePermission(perms.CredentialRead),
		coverageMw.InjectRequestContext(),
		handlers.FetchOptions(service),
	)
}
