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
func InitMiddlewares(c *dig.Container) {
	c.Invoke(func(sharedCache common.SharedCache) {
		baseURL, _ := config.GetStringValue("mapexos_url")
		apiKey, _ := config.GetStringValue("internal_api_key")

		permissionMw.InitPermissionMiddleware(sharedCache, baseURL, apiKey)
		logger.Info("[APP:BOOTSTRAP] Permission middleware initialized")

		coverageMw.InitCoverageMiddleware(sharedCache)
		coverageMw.InitCacheBuildClient(baseURL, apiKey)
		logger.Info("[APP:BOOTSTRAP] Coverage middleware initialized")
	})
}
