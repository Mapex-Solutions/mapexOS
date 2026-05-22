package bootstrap

import (
	"fmt"

	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
)

// InitMiddlewares initializes permission and coverage middlewares.
func InitMiddlewares(c *dig.Container) {
	// Initialize Permission Middleware (Singleton)
	c.Invoke(func(sharedCache common.SharedCache) {
		// Build internal API base URL
		httpAddress, _ := config.GetStringValue("http_address")
		httpPort, _ := config.GetIntValue("http_port")
		baseURL := fmt.Sprintf("http://%s:%d", httpAddress, httpPort)

		// Get internal API key
		apiKey, _ := config.GetStringValue("internal_api_key")

		// Initialize middleware with cache and internal API client
		permissionMw.InitPermissionMiddleware(sharedCache, baseURL, apiKey)
		logger.Info("[APP:BOOTSTRAP] Permission middleware initialized")
	})

	// Initialize Coverage Middleware (Singleton)
	c.Invoke(func(sharedCache common.SharedCache) {
		// Build internal API base URL
		httpAddress, _ := config.GetStringValue("http_address")
		httpPort, _ := config.GetIntValue("http_port")
		baseURL := fmt.Sprintf("http://%s:%d", httpAddress, httpPort)

		// Get internal API key
		apiKey, _ := config.GetStringValue("internal_api_key")

		// Initialize middleware with shared cache
		coverageMw.InitCoverageMiddleware(sharedCache)

		// Initialize cache build client for lazy loading
		coverageMw.InitCacheBuildClient(baseURL, apiKey)

		logger.Info("[APP:BOOTSTRAP] Coverage middleware initialized")
	})
}
