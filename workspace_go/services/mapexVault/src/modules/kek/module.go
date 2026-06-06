package kek

import (
	"log"

	kekPorts "mapexVault/src/modules/kek/application/ports"
	service "mapexVault/src/modules/kek/application/services"
	kekEnvelope "mapexVault/src/modules/kek/infrastructure/envelope"
	mongoKEK "mapexVault/src/modules/kek/infrastructure/persistence/mongo"
	"mapexVault/src/modules/kek/interfaces/http/handlers"
	"mapexVault/src/modules/kek/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	apikeymw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"github.com/gofiber/fiber/v2"
)

// InitRepositories registers the kek repository + envelope adapter.
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(mongoKEK.NewEncryptionKeyRepository)
	c.Provide(kekEnvelope.NewEnvelopeAdapter)
	logger.Info("[MODULE:Kek] Repositories registered")
}

// InitServices registers the KekService.
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Kek] Services registered")
}

// InitInterfaces mounts the internal HTTP routes (API-key gated). KEK documents
// are populated by the mongodb-init container; no in-process bootstrap runs here.
func InitInterfaces() {
	c := container.GetContainer()
	if err := c.Invoke(func(params struct {
		container.In
		App     *fiber.App
		Service kekPorts.KekServicePort
	}) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")
		internalApiKey, _ := config.GetStringValue("internal_api_key")
		group := params.App.Group(
			"/internal/kek",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(internalApiKey),
		)
		h := handlers.NewKekInternalHandler(params.Service)
		routes.RegisterInternalRoutes(group, h)
	}); err != nil {
		log.Fatalf("failed to wire kek interfaces: %v", err)
	}
	logger.Info("[MODULE:Kek] Interfaces registered (/internal/kek)")
}
