package credentials

import (
	"log"

	"mapexVault/src/modules/credentials/application/ports"
	service "mapexVault/src/modules/credentials/application/services"
	collection "mapexVault/src/modules/credentials/infrastructure/persistence/mongo"
	"mapexVault/src/modules/credentials/interfaces/http/routes"
	reconcileConsumer "mapexVault/src/modules/credentials/interfaces/message/consumers/reconcile"
	refreshConsumer "mapexVault/src/modules/credentials/interfaces/message/consumers/refresh"

	"github.com/gofiber/fiber/v2"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	common "github.com/Mapex-Solutions/mapexGoKit/microservices/common"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	apikeymw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers credential and connection repositories.
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(collection.NewCredentialRepository)
	c.Provide(collection.NewConnectionRepository)
	logger.Info("[MODULE:Credentials] Repositories registered")
}

// InitServices registers the credential service.
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Credentials] Services registered")
}

// InitInterfaces registers HTTP routes for credentials.
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(params struct {
		container.In
		App     *fiber.App
		Service ports.CredentialServicePort
	}) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")

		// External routes (JWT auth)
		externalRoutes := params.App.Group(
			"/api/v1/credentials",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		routes.RegisterRoutes(externalRoutes, params.Service)

		// Internal routes (API Key auth)
		internalApiKey, _ := config.GetStringValue("internal_api_key")
		internalRoutes := params.App.Group(
			"/internal/credentials",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(internalApiKey),
		)
		routes.RegisterInternalRoutes(internalRoutes, params.Service)

	}); err != nil {
		log.Fatalf("failed to invoke credentials module: %v", err)
	}

	// NATS consumers + lifecycle hooks
	if err := c.Invoke(func(params struct {
		container.In
		Bus     *natsModel.Bus `name:"core"`
		Service ports.CredentialServicePort
	}) {
		// Refresh consumer — pulls from vault.schedule.fired (VAULT-SCHEDULE stream).
		// Per-credential timers fire here; HandleRefreshMessage refreshes the token
		// and immediately re-arms the next timer before ack.
		if refreshConsumer.NewConsumer(params.Bus, params.Service) == nil {
			log.Fatalf("failed to start refresh consumer")
		}

		// Reconciler consumer — pulls from vault.reconcile.fired (VAULT-RECONCILER stream).
		// Safety-net loop that reseeds per-credential timers missing from VAULT-SCHEDULE.
		reconcileConsumer.NewConsumer(params.Bus, params.Service)

		// Run lifecycle hooks (OnMount): bootstrap seed + first reconcile timer.
		common.RunLifecycleHooks(params.Service, "Credentials")
	}); err != nil {
		log.Fatalf("failed to invoke credentials consumer: %v", err)
	}

	logger.Info("[MODULE:Credentials] Interfaces registered (HTTP routes + refresh + reconcile consumers)")
}
