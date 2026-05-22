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
	// Initialize Redis - App Cache (service-private DB)
	appRedisCfg := config.GetRedisConfig()
	c.Provide(func() *redisModel.RedisClient {
		rm, err := redisModel.New(appRedisCfg)
		if err != nil {
			logger.Panic(err.Error())
		}
		return rm
	}, container.Name("app"))

	// Initialize Redis - Shared Cache (cross-service DB for authorization)
	sharedRedisCfg := config.GetSharedRedisConfig()
	c.Provide(func() *redisModel.RedisClient {
		rm, err := redisModel.New(sharedRedisCfg)
		if err != nil {
			logger.Panic(err.Error())
		}
		return rm
	}, container.Name("shared"))

	// Provide AppCache interface (service-private cache)
	c.Provide(func(params struct {
		container.In
		RC *redisModel.RedisClient `name:"app"`
	}) common.AppCache {
		return params.RC
	})

	// Provide SharedCache interface (cross-service authorization cache)
	c.Provide(func(params struct {
		container.In
		RC *redisModel.RedisClient `name:"shared"`
	}) common.SharedCache {
		return params.RC
	})
}
