package main

import (
	"fmt"
	"time"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	appModule "workflow/src/modules/app"
	"workflow/src/bootstrap"
	runtimePorts "workflow/src/modules/runtime/application/ports"
)

func main() {

	/*
	* Initialize a DIG container (the import return a singleton)
	 */
	container.InitContainer()
	c := container.GetContainer()

	/*
	* Initialize configuration and logger
	 */
	bootstrap.InitConfig()
	bootstrap.InitLogger()

	/*
	* Provide the auth configuration to use in the middlewares
	 */
	c.Provide(config.GetAuthConfig)

	/*
	* Initialize metrics registry and service-specific metrics
	 */
	bootstrap.InitMetrics(c)

	/*
	* Initialize all infrastructure providers
	 */
	bootstrap.InitMongo(c)
	bootstrap.InitRedis(c)
	bootstrap.InitNATS(c)
	bootstrap.InitTieredCache(c)
	bootstrap.InitEncryption(c)
	bootstrap.InitVaultClient(c)
	bootstrap.InitMiddlewares(c)

	/*
	* Create Fiber instance with global middlewares
	 */
	fiberInstance := bootstrap.InitFiber(c)

	/*
	* Initialize health check endpoint (before business modules, no auth required)
	 */
	bootstrap.InitHealth(c, fiberInstance)

	/*
	* Create shutdown manager and provide to DIG container
	 */
	sm := shutdown.New()
	c.Provide(func() *shutdown.ShutdownManager { return sm })

	/*
	* Create consumer registry and provide to DIG container
	 */
	registry := &bootstrap.ConsumerRegistry{}
	c.Provide(func() *bootstrap.ConsumerRegistry { return registry })

	/*
	* Initialize all business modules (repositories, services, consumers, routes)
	 */
	appModule.InitModule(fiberInstance)

	/*
	* Register shutdown hooks (after modules so RuntimeService is available)
	 */
	var drainer bootstrap.WalkerDrainer
	c.Invoke(func(svc runtimePorts.RuntimeServicePort) {
		if d, ok := svc.(bootstrap.WalkerDrainer); ok {
			drainer = d
		}
	})
	bootstrap.InitShutdown(c, sm, fiberInstance, registry, drainer)

	/*
	* Start the HTTP server (non-blocking)
	 */
	httpPort, _ := config.GetIntValue("http_port")
	httpAddress, _ := config.GetStringValue("http_address")
	address := fmt.Sprintf("%s:%d", httpAddress, httpPort)

	go func() {
		if err := fiberInstance.Listen(address); err != nil {
			logger.Error(err, "[APP:MAIN] HTTP server stopped")
		}
	}()

	/*
	* Block until SIGTERM/SIGINT, then graceful shutdown
	 */
	sm.WaitForSignal(15 * time.Second)
}
