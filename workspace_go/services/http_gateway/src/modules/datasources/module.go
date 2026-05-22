package dataSources

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"http_gateway/src/modules/datasources/application/ports"
	service "http_gateway/src/modules/datasources/application/services"
	"http_gateway/src/modules/datasources/domain/repositories"
	redisCache "http_gateway/src/modules/datasources/infrastructure/cache/redis"
	collection "http_gateway/src/modules/datasources/infrastructure/persistence/mongo"
	routes "http_gateway/src/modules/datasources/interfaces/http/routes"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
)

// InitRepositories registers the datasources repositories in the DIG container
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(collection.New)

	// Provide CacheRepository using AppCache (service-private Redis)
	c.Provide(func(appCache common.AppCache) repositories.CacheRepository {
		return appCache
	})

	// Provide the Redis-backed DataSource cache key builder (keeps the Redis
	// key prefix out of application/constants).
	c.Provide(redisCache.NewDataSourceCacheKeyBuilder)

	logger.Info("[MODULE:DataSources] Repositories registered")
}

// InitServices registers the datasources services in the DIG container
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:DataSources] Services registered")
}

// InitInterfaces registers the datasources routes (HTTP)
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, service ports.DataSourceServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// External routes
		routesV1 := app.Group(
			"/api/v1/data_sources",

			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		// Register the routes and handlers
		routes.RegisterRoutes(routesV1, service)

	}); err != nil {
		logger.Panic(fmt.Sprintf("[MODULE:DataSources] failed to invoke module: %v", err))
	}

	logger.Info("[MODULE:DataSources] Routes registered")
}
