package fetch_options

import (
	"log"

	fetchPorts "workflow/src/modules/fetch_options/application/ports"
	"workflow/src/modules/fetch_options/application/services"
	"workflow/src/modules/fetch_options/interfaces/http/routes"

	"github.com/gofiber/fiber/v2"

	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitServices registers the FetchOptionsService in DIG.
func InitServices() {
	c := container.GetContainer()
	c.Provide(services.New)
	logger.Info("[MODULE:FetchOptions] Services registered")
}

// InitInterfaces registers HTTP routes.
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(params struct {
		container.In
		App     *fiber.App
		Service fetchPorts.FetchOptionsServicePort
	}) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		routeGroup := params.App.Group(
			"/api/v1/load_options",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		routes.RegisterRoutes(routeGroup, params.Service)

	}); err != nil {
		log.Fatalf("failed to invoke fetch_options module: %v", err)
	}

	logger.Info("[MODULE:FetchOptions] Interfaces registered")
}
