package bootstrap

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"

	chManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/manager"
	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/health"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/health/adapters"
)

// InitHealth registers the /health endpoint with all infrastructure checkers.
func InitHealth(c *dig.Container, app *fiber.App) {
	c.Invoke(func(params struct {
		container.In
		MongoMgr       *mongoManager.MongoManager
		RedisApp       *redisModel.RedisClient       `name:"app"`
		RedisShared    *redisModel.RedisClient       `name:"shared"`
		NATSClient     *natsModel.Client
		ClickHouseMgr  *chManager.ClickHouseManager
		MinIOTemplates *minioModel.MinIOClient        `name:"templates"`
	}) {
		serviceName, _ := config.GetStringValue("service_name")
		serviceVersion, _ := config.GetStringValue("service_version")

		health.RegisterRoutes(app, health.Config{
			ServiceName: serviceName,
			Version:     serviceVersion,
			CacheTTL:    10 * time.Second,
			Timeout:     5 * time.Second,
		},
			health.CheckerConfig{Checker: adapters.NewMongoAdapter(params.MongoMgr), Critical: true},
			health.CheckerConfig{Checker: adapters.NewRedisAdapter(params.RedisApp, "app"), Critical: true},
			health.CheckerConfig{Checker: adapters.NewRedisAdapter(params.RedisShared, "shared"), Critical: true},
			health.CheckerConfig{Checker: adapters.NewNATSAdapter(params.NATSClient, "core"), Critical: true},
			health.CheckerConfig{Checker: adapters.NewClickHouseAdapter(params.ClickHouseMgr), Critical: false},
			health.CheckerConfig{Checker: adapters.NewMinIOAdapter(params.MinIOTemplates, "templates"), Critical: false},
		)
	})
}
