package bootstrap

import (
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
)

// InitRedis registers App and Shared Redis clients and cache interfaces.
func InitRedis(c *dig.Container) {
	// Initialize Redis - App Cache (DB 0: MQTT Auth username index + auth cache)
	appRedisCfg := config.GetRedisConfig()
	c.Provide(func() *redisModel.RedisClient {
		rm, err := redisModel.New(appRedisCfg)
		if err != nil {
			logger.Panic(err.Error())
		}
		return rm
	})

	// Provide AppCache interface (Redis DB 0)
	// Used for: username → UUID index for MQTT Auth Callout
	c.Provide(func(rc *redisModel.RedisClient) common.AppCache {
		return rc
	})

	// Initialize Redis - Shared Cache (cross-service DB for authorization)
	// Used by: Permission middleware for org/user validation
	// NOT used by: Asset/AssetTemplate modules (they use TieredCache + MinIO)
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
