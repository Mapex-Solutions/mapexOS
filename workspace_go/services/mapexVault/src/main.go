package main

import (
	"fmt"
	"time"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	appModule "mapexVault/src/modules/app"
	"mapexVault/src/bootstrap"
)

func main() {
	container.InitContainer()
	c := container.GetContainer()

	bootstrap.InitConfig()
	bootstrap.InitLogger()

	c.Provide(config.GetAuthConfig)

	bootstrap.InitMetrics(c)

	bootstrap.InitMongo(c)
	bootstrap.InitRedis(c)
	bootstrap.InitNATS(c)
	bootstrap.InitEncryption(c)
	bootstrap.InitMiddlewares(c)

	fiberInstance := bootstrap.InitFiber(c)

	bootstrap.InitHealth(c, fiberInstance)

	sm := shutdown.New()
	bootstrap.InitShutdown(c, sm, fiberInstance)

	appModule.InitModule(fiberInstance)

	httpPort, _ := config.GetIntValue("http_port")
	httpAddress, _ := config.GetStringValue("http_address")
	address := fmt.Sprintf("%s:%d", httpAddress, httpPort)

	go func() {
		if err := fiberInstance.Listen(address); err != nil {
			logger.Error(err, "[APP:MAIN] HTTP server stopped")
		}
	}()

	sm.WaitForSignal(15 * time.Second)
}
