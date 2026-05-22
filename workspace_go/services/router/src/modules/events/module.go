package events

import (
	"router/src/modules/events/application/ports"
	service "router/src/modules/events/application/services"
	templateCache "router/src/modules/events/infrastructure/cache/tieredcache"
	consumers "router/src/modules/events/interfaces/message/consumers"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitServices registers the events services in the DIG container.
func InitServices() {
	c := container.GetContainer()

	// Provide TemplateCache (wraps TieredCache for template metadata enrichment)
	c.Provide(func(params struct {
		container.In
		Cache common.TieredCache `name:"templates"`
	}) ports.TemplateCachePort {
		return templateCache.New(params.Cache)
	})

	c.Provide(service.New)
	logger.Info("[MODULE:Events] Services registered")
}

// InitInterfaces registers the events consumers (NATS)
func InitInterfaces() {
	c := container.GetContainer()

	// Start route execution consumer (WorkQueue pattern - load balanced)
	if err := c.Invoke(consumers.NewRouteExecuteConsumer); err != nil {
		logger.Error(err, "[CONSUMER:RouteExecute] Failed to start route execute consumer")
	}

	// Start asset invalidation consumer (FANOUT pattern - broadcast to all instances)
	if err := c.Invoke(consumers.NewAssetInvalidateConsumer); err != nil {
		logger.Error(err, "[CONSUMER:AssetInvalidate] Failed to start asset invalidate consumer")
	}

	// Start template invalidation consumer (FANOUT pattern - broadcast to all instances)
	if err := c.Invoke(consumers.NewTemplateInvalidateConsumer); err != nil {
		logger.Error(err, "[CONSUMER:TemplateInvalidate] Failed to start template invalidate consumer")
	}

	logger.Info("[MODULE:Events] Consumers registered")
}
