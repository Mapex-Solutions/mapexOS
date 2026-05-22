package healthmonitor

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"assets/src/modules/healthmonitor/application/di"
	"assets/src/modules/healthmonitor/application/ports"
	service "assets/src/modules/healthmonitor/application/services"
	redisRepo "assets/src/modules/healthmonitor/infrastructure/persistence/redis"
	natsMessaging "assets/src/modules/healthmonitor/infrastructure/messaging/nats"
	httpRoutes "assets/src/modules/healthmonitor/interfaces/http/routes"
	"assets/src/modules/healthmonitor/interfaces/message/consumers/heartbeat"
	"assets/src/modules/healthmonitor/interfaces/message/consumers/presence"
	"assets/src/modules/healthmonitor/interfaces/message/consumers/scan"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
	common "github.com/Mapex-Solutions/mapexGoKit/microservices/common"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	apikeymw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the health monitor repositories in the DIG container.
func InitRepositories() {
	c := container.GetContainer()

	c.Provide(func(rc *redisModel.RedisClient) ports.HealthRepository {
		return redisRepo.New(rc)
	})

	logger.Info("[MODULE:HealthMonitor] Repositories registered")
}

// InitServices registers the health monitor services in the DIG container.
func InitServices() {
	c := container.GetContainer()

	c.Provide(func(params struct {
		container.In
		Publisher natsModel.Publisher `name:"core"`
	}) ports.AlertPublisherPort {
		return natsMessaging.NewAlertPublisher(params.Publisher)
	})

	batchSize, _ := config.GetIntValue("health_monitor_batch_size")
	if batchSize <= 0 {
		batchSize = 500
	}

	c.Provide(func(deps di.HealthMonitorServiceDI) ports.HealthMonitorServicePort {
		return service.New(deps, batchSize)
	})

	c.Provide(func(svc ports.HealthMonitorServicePort) ports.HealthLifecyclePort {
		return svc.(ports.HealthLifecyclePort)
	})

	c.Provide(func(svc ports.HealthMonitorServicePort) ports.PresencePort {
		return svc.(ports.PresencePort)
	})

	c.Provide(func(svc ports.HealthMonitorServicePort) ports.HealthAdminPort {
		return svc.(ports.HealthAdminPort)
	})

	logger.Info("[MODULE:HealthMonitor] Services registered")
}

// InitInterfaces registers the health monitor consumers (NATS) and the
// internal HTTP routes (API key gated). The internal routes are used by
// e2e journeys to drive offline transitions without waiting the
// configured scan cycle; nothing on the production data path depends
// on them.
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(params struct {
		container.In
		Bus     *natsModel.Bus `name:"core"`
		Service ports.HealthMonitorServicePort
	}) {
		heartbeat.NewConsumer(params.Bus, params.Service)
		presence.NewConsumer(params.Bus, params.Service)
		scan.NewConsumer(params.Bus, params.Service)
		common.RunLifecycleHooks(params.Service, "HealthMonitor")
	}); err != nil {
		log.Fatalf("failed to invoke healthmonitor consumers: %v", err)
	}

	if err := c.Invoke(func(app *fiber.App, adminSvc ports.HealthAdminPort) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")
		internalApiKey, _ := config.GetStringValue("internal_api_key")
		internalRoutes := app.Group(
			"/internal/health-monitor",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(internalApiKey),
		)
		httpRoutes.RegisterInternalRoutes(internalRoutes, adminSvc)
	}); err != nil {
		log.Fatalf("failed to invoke healthmonitor internal routes: %v", err)
	}

	logger.Info("[MODULE:HealthMonitor] Interfaces registered (consumers + internal routes)")
}
