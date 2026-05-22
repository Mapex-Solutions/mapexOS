package bootstrap

import (
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
)

// InitMiddlewares initializes permission and coverage middlewares.
// Uses SharedCache for org/user validation in auth middleware.
func InitMiddlewares(c *dig.Container) {
	c.Invoke(func(sharedCache common.SharedCache) {
		// Build internal API base URL
		baseURL, _ := config.GetStringValue("mapexos_url")

		// Get internal API key
		apiKey, _ := config.GetStringValue("internal_api_key")

		// Initialize permission middleware with cache and internal API client
		permissionMw.InitPermissionMiddleware(sharedCache, baseURL, apiKey)
		logger.Info("[INFRA:Middlewares] Permission middleware initialized")

		// Initialize coverage middleware with shared cache
		coverageMw.InitCoverageMiddleware(sharedCache)

		// Initialize cache build client for lazy loading
		coverageMw.InitCacheBuildClient(baseURL, apiKey)

		logger.Info("[INFRA:Middlewares] Coverage middleware initialized")
	})
}
