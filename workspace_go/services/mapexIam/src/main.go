package main

import (
	"fmt"
	"time"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	appModule "mapexIam/src/modules/app"
	"mapexIam/src/bootstrap"
)

// Main function to start the MapexOS service
func main() {

	/**
	* Initialize a DIG container (the import return a singleton)
	 */
	container.InitContainer()
	c := container.GetContainer()

	/**
	* Initialize configuration and logger
	 */
	bootstrap.InitConfig()
	bootstrap.InitLogger()

	/**
	* Provide the auth configuration to use in the middlewares
	 */
	c.Provide(config.GetAuthConfig)

	/**
	* Initialize metrics registry and service-specific metrics
	 */
	bootstrap.InitMetrics(c)

	/**
	* Initialize all infrastructure providers
	 */
	bootstrap.InitMongo(c)
	bootstrap.InitRedis(c)
	bootstrap.InitNATS(c)
	bootstrap.InitMiddlewares(c)

	/**
	* Create Fiber instance with global middlewares
	 */
	fiberInstance := bootstrap.InitFiber(c)

	/**
	* Initialize health check endpoint (before business modules, no auth required)
	 */
	bootstrap.InitHealth(c, fiberInstance)

	/*
	 * Create shutdown manager and register infrastructure hooks
	 */
	sm := shutdown.New()
	bootstrap.InitShutdown(c, sm, fiberInstance)

	/**
	* Initialize all business modules (repositories, services, consumers, routes)
	 */
	appModule.InitModule(fiberInstance)

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
