package asset_status

import (
	ports "events/src/modules/asset_status/application/ports"
	service "events/src/modules/asset_status/application/services"
	repositories "events/src/modules/asset_status/domain/repositories"
	clickhouseRepo "events/src/modules/asset_status/infrastructure/persistence/clickhouse"
	routes "events/src/modules/asset_status/interfaces/http/routes"
	consumers "events/src/modules/asset_status/interfaces/message/consumers"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gofiber/fiber/v2"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the ClickHouse-backed asset_status repository.
func InitRepositories() {
	c := container.GetContainer()

	c.Provide(func(conn driver.Conn) repositories.AssetStatusRepository {
		return clickhouseRepo.NewAssetStatusRepository(conn)
	})

	logger.Info("[MODULE:AssetStatus] Repositories registered")
}

// InitServices registers the application service.
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:AssetStatus] Services registered")
}

// InitInterfaces registers HTTP routes + the NATS save consumer.
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *fiber.App, svc ports.AssetStatusServicePort) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")
		routesV1 := app.Group(
			"/api/v1/events",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		routes.RegisterRoutes(routesV1, svc)
		logger.Info("[MODULE:AssetStatus] HTTP routes registered")
	}); err != nil {
		logger.Panic("[MODULE:AssetStatus] Failed to register interfaces: " + err.Error())
	}

	if err := c.Invoke(consumers.NewAssetStatusSaveConsumer); err != nil {
		logger.Error(err, "[MODULE:AssetStatus] Failed to start asset_status save consumer")
	}

	logger.Info("[MODULE:AssetStatus] NATS consumers registered")
}
