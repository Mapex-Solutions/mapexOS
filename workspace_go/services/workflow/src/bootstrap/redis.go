package bootstrap

import (
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
)

// InitRedis registers Shared Redis client for authorization middleware.
// Workflow service does NOT use App Redis — it uses TieredCache + MinIO for caching.
// Only Shared Redis (DB 5) is needed for permission/coverage middleware.
func InitRedis(c *dig.Container) {
	// Initialize Redis - Shared Cache (cross-service DB for authorization)
	// Used by: Permission middleware for org/user validation
	sharedRedisCfg := config.GetSharedRedisConfig()
	c.Provide(func() *redisModel.RedisClient {
		rm, err := redisModel.New(sharedRedisCfg)
		if err != nil {
			logger.Panic(err.Error())
		}
		return rm
	}, container.Name("shared"))

	// Provide SharedCache interface (cross-service authorization cache)
	c.Provide(func(params struct {
		container.In
		RC *redisModel.RedisClient `name:"shared"`
	}) common.SharedCache {
		return params.RC
	})
}
