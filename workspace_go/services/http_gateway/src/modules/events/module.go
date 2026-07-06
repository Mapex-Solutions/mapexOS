package events

import (
	"fmt"

	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"http_gateway/src/bootstrap"
	dsPort "http_gateway/src/modules/datasources/application/ports"
	port "http_gateway/src/modules/events/application/ports"
	service "http_gateway/src/modules/events/application/services"
	assetsClient "http_gateway/src/modules/events/infrastructure/httpclient/assets"
	routes "http_gateway/src/modules/events/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitServices registers the events services in the DIG container.
func InitServices() {
	c := container.GetContainer()
	c.Provide(assetsClient.NewOTAJobsPort)
	c.Provide(service.New)
	logger.Info("[MODULE:Events] Services registered")
}

// InitInterfaces registers the events HTTP routes (webhook + heartbeat).
//
// Following Hexagonal Architecture, this function:
//   - Accepts service port interfaces (not concrete implementations)
//   - Enables dependency inversion and testability
//   - Decouples route registration from service implementation details
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *web.App, service port.EventServicePort, dsService dsPort.DataSourceServicePort, m *bootstrap.HttpGatewayMetrics) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")
		routes.RegisterRoutes(app, ctxTimeout, service, dsService, m)
	}); err != nil {
		logger.Panic(fmt.Sprintf("[MODULE:Events] failed to invoke module: %v", err))
	}

	logger.Info("[MODULE:Events] Routes registered")
}
