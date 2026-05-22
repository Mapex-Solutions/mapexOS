package pki

import (
	"log"

	pkiPorts "mapexVault/src/modules/pki/application/ports"
	service "mapexVault/src/modules/pki/application/services"
	pkiCrypto "mapexVault/src/modules/pki/infrastructure/crypto"
	pkiEnvelope "mapexVault/src/modules/pki/infrastructure/envelope"
	mongoCA "mapexVault/src/modules/pki/infrastructure/persistence/mongo"
	"mapexVault/src/modules/pki/interfaces/http/handlers"
	"mapexVault/src/modules/pki/interfaces/http/routes"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	apikeymw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"github.com/gofiber/fiber/v2"
)

// InitRepositories registers the pki repository + crypto + envelope adapters.
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(mongoCA.NewCARepository)
	c.Provide(pkiCrypto.NewX509Signer)
	c.Provide(pkiEnvelope.NewEnvelopeAdapter)
	logger.Info("[MODULE:Pki] Repositories registered")
}

// InitServices registers the PkiService.
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Pki] Services registered")
}

// InitInterfaces mounts the internal HTTP routes (API-key gated). The
// CA collection is populated by the mongodb-init container via the
// seed-encryptor output produced by scripts/prebuild/pki/generate-pki.sh —
// no in-process bootstrap runs here.
func InitInterfaces() {
	c := container.GetContainer()
	if err := c.Invoke(func(params struct {
		container.In
		App     *fiber.App
		Service pkiPorts.PkiServicePort
	}) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")
		internalApiKey, _ := config.GetStringValue("internal_api_key")
		group := params.App.Group(
			"/internal/pki",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(internalApiKey),
		)
		h := handlers.NewPkiInternalHandler(params.Service)
		routes.RegisterInternalRoutes(group, h)
	}); err != nil {
		log.Fatalf("failed to wire pki interfaces: %v", err)
	}
	logger.Info("[MODULE:Pki] Interfaces registered (/internal/pki)")
}
