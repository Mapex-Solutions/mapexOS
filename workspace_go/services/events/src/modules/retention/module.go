package retention

import (
	"context"

	"events/src/modules/retention/application/ports"
	"events/src/modules/retention/domain/repositories"
	service "events/src/modules/retention/application/services"
	chAdapter "events/src/modules/retention/infrastructure/persistence/clickhouse"
	collection "events/src/modules/retention/infrastructure/persistence/mongo"
	routes "events/src/modules/retention/interfaces/http/routes"
	orgCreatedConsumer "events/src/modules/retention/interfaces/message/consumers/org_created"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gofiber/fiber/v2"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the retention repositories in the DIG container.
func InitRepositories() {
	c := container.GetContainer()

	// Provide MongoDB repository
	c.Provide(collection.New)

	// Provide the App Redis client as CacheRepository for domain services
	c.Provide(func(params struct {
		container.In
		RC *redisModel.RedisClient `name:"app"`
	}) repositories.CacheRepository {
		return params.RC
	})

	// Provide the ClickHouseConnPort adapter, wrapping the raw driver.Conn so
	// the application layer never depends on the concrete driver type.
	c.Provide(func(conn driver.Conn) ports.ClickHouseConnPort {
		return chAdapter.NewConnAdapter(conn)
	})

	logger.Info("[MODULE:Retention] Repositories registered")
}

// InitServices registers the retention services in the DIG container and
// seeds platform-level policies (asset_status_history default 7 days).
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)

	if err := c.Invoke(func(svc ports.RetentionServicePort) {
		if seedErr := svc.SeedPlatformPolicies(context.Background()); seedErr != nil {
			logger.Error(seedErr, "[MODULE:Retention] Failed to seed platform policies")
		}
	}); err != nil {
		logger.Error(err, "[MODULE:Retention] Failed to invoke platform seed")
	}

	logger.Info("[MODULE:Retention] Services registered")
}

// InitInterfaces registers the retention HTTP routes and NATS consumers.
func InitInterfaces() {
	c := container.GetContainer()

	// Register HTTP routes
	if err := c.Invoke(func(app *fiber.App, service ports.RetentionServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Create route group with auth middleware
		routesV1 := app.Group(
			"/api/v1/retention",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		routes.RegisterRoutes(routesV1, service)
		logger.Info("[MODULE:Retention] HTTP routes registered")

	}); err != nil {
		logger.Panic("[MODULE:Retention] failed to invoke retention module interfaces: " + err.Error())
	}

	// Start NATS consumer for organization created events
	if err := c.Invoke(func(bus *natsModel.Bus, service ports.RetentionServicePort) {
		orgCreatedConsumer.NewConsumer(bus, service)
	}); err != nil {
		logger.Error(err, "[MODULE:Retention] Failed to start retention org-created consumer")
	}

	logger.Info("[MODULE:Retention] Interfaces registered")
}
