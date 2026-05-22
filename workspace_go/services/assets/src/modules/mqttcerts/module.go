package mqttcerts

import (
	"log"

	mqttPorts "assets/src/modules/mqttcerts/application/ports"
	service "assets/src/modules/mqttcerts/application/services"
	mqttCrypto "assets/src/modules/mqttcerts/infrastructure/crypto"
	vaultClient "assets/src/modules/mqttcerts/infrastructure/http/mapexvault_client"
	mongoRevoked "assets/src/modules/mqttcerts/infrastructure/persistence/mongo"
	mqttRam "assets/src/modules/mqttcerts/infrastructure/ram"
	"assets/src/modules/mqttcerts/interfaces/http/handlers"
	"assets/src/modules/mqttcerts/interfaces/http/routes"

	"github.com/gofiber/fiber/v2"

	common "github.com/Mapex-Solutions/mapexGoKit/microservices/common"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

func InitRepositories() {
	c := container.GetContainer()
	c.Provide(mongoRevoked.NewRevokedRepository)
	c.Provide(vaultClient.NewMapexVaultClient)
	c.Provide(mqttCrypto.NewX509Signer)
	c.Provide(mqttRam.NewInMemoryCAStore)
	logger.Info("[MODULE:MqttCerts] Repositories registered")
}

func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:MqttCerts] Services registered")
}

func InitInterfaces() {
	c := container.GetContainer()
	if err := c.Invoke(func(params struct {
		container.In
		App     *fiber.App
		Service mqttPorts.MqttCertsServicePort
	}) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")
		external := params.App.Group(
			"/api/v1/mqtt_certs",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		h := handlers.NewMqttCertsHandler(params.Service)
		routes.RegisterRoutes(external, h, params.Service)

		// Fire OnMount lifecycle hook (sync attempt + retry goroutine on fail).
		if m, ok := params.Service.(common.Mountable); ok {
			m.OnMount()
		}
	}); err != nil {
		log.Fatalf("failed to wire mqttcerts module: %v", err)
	}
	logger.Info("[MODULE:MqttCerts] Interfaces registered (/api/v1/mqtt_certs)")
}
