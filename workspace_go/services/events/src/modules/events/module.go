package events

import (
	ports "events/src/modules/events/application/ports"
	service "events/src/modules/events/application/services"
	repositories "events/src/modules/events/domain/repositories"
	templateCacheAdapter "events/src/modules/events/infrastructure/cache/tieredcache"
	clickhouseRepo "events/src/modules/events/infrastructure/persistence/clickhouse"
	consumers "events/src/modules/events/interfaces/message/consumers"
	routes "events/src/modules/events/interfaces/http/routes"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gofiber/fiber/v2"
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the events repositories in the DIG container
func InitRepositories() {
	c := container.GetContainer()

	// Register EventRepository (ClickHouse implementation)
	c.Provide(func(conn driver.Conn) repositories.EventRepository {
		return clickhouseRepo.NewEventRepository(conn)
	})

	logger.Info("[MODULE:Events] Repositories registered")
}

// InitServices registers the events services in the DIG container
func InitServices() {
	c := container.GetContainer()

	// Register TemplateCache as TemplateCachePort (wraps DIG-injected TieredCache)
	c.Provide(func(params struct {
		container.In
		Cache common.TieredCache `name:"templates"`
	}) ports.TemplateCachePort {
		return templateCacheAdapter.New(params.Cache)
	})

	// Register EventService
	c.Provide(service.New)

	logger.Info("[MODULE:Events] Services registered")
}

// InitInterfaces registers the events interfaces (HTTP routes and NATS consumers)
func InitInterfaces() {
	c := container.GetContainer()

	// Register HTTP routes
	if err := c.Invoke(func(app *fiber.App, service ports.EventServicePort) {

		// Set default timeout for this router
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// Create a new group for the events routes
		routesV1 := app.Group(
			"/api/v1/events",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)

		routes.RegisterRoutes(routesV1, service)
		logger.Info("[MODULE:Events] HTTP routes registered")

	}); err != nil {
		logger.Panic("[MODULE:Events] failed to invoke events module interfaces: " + err.Error())
	}

	// Start events.save consumer (processed events from router service)
	// Stores events in main events table (retention: 7-365 days)
	if err := c.Invoke(consumers.NewEventsSaveConsumer); err != nil {
		logger.Error(err, "[MODULE:Events] Failed to start events.save consumer")
	}

	// Start events.raw consumer (raw events from HTTP/MQTT gateways)
	// Stores events in events_raw table for debugging (retention: 1-3 days)
	if err := c.Invoke(consumers.NewEventsRawConsumer); err != nil {
		logger.Error(err, "[MODULE:Events] Failed to start events.raw consumer")
	}

	// Start events.logs.jsexecutor consumer (JS executor debug events)
	// Stores events in events_jsexecutor table for debugging (retention: 1-3 days)
	if err := c.Invoke(consumers.NewEventsJsExecConsumer); err != nil {
		logger.Error(err, "[MODULE:Events] Failed to start events.logs.jsexecutor consumer")
	}

	// Start DLQ consumer (Dead Letter Queue events)
	// Stores failed events in dlq table for investigation
	if err := c.Invoke(consumers.NewEventsDLQConsumer); err != nil {
		logger.Error(err, "[MODULE:Events] Failed to start DLQ consumer")
	}

	// Start events.router consumer (router execution history events)
	// Stores events in events_router table for UI visualization (retention: 1-30 days)
	if err := c.Invoke(consumers.NewEventsRouterConsumer); err != nil {
		logger.Error(err, "[MODULE:Events] Failed to start events.router consumer")
	}

	// Start events.businessrule consumer (business rule execution history events)
	// Stores events in events_businessrule table for UI visualization (retention: 1-30 days)
	if err := c.Invoke(consumers.NewEventsBusinessRuleConsumer); err != nil {
		logger.Error(err, "[MODULE:Events] Failed to start events.businessrule consumer")
	}

	// Start events.trigger consumer (trigger execution history events)
	// Stores events in events_trigger table for UI visualization (retention: 1-30 days)
	if err := c.Invoke(consumers.NewEventsTriggerConsumer); err != nil {
		logger.Error(err, "[MODULE:Events] Failed to start events.trigger consumer")
	}

	// Start events.workflow consumer (workflow execution history events)
	// Stores events in events_workflow table for UI visualization (retention: 1-365 days)
	if err := c.Invoke(consumers.NewEventsWorkflowConsumer); err != nil {
		logger.Error(err, "[MODULE:Events] Failed to start events.workflow consumer")
	}

	// Start FANOUT template invalidate consumer — clears TieredCache L0+L1 when
	// the assets service publishes mapexos.fanout.template.invalidate.
	if err := c.Invoke(consumers.NewTemplateInvalidateConsumer); err != nil {
		logger.Error(err, "[MODULE:Events] Failed to start template invalidate consumer")
	}

	logger.Info("[MODULE:Events] NATS consumers registered")
}
