package credentials

import (
	"context"
	"log"
	"time"

	"mapexVault/src/modules/credentials/application/constants"
	"mapexVault/src/modules/credentials/application/ports"
	service "mapexVault/src/modules/credentials/application/services"
	collection "mapexVault/src/modules/credentials/infrastructure/persistence/mongo"
	"mapexVault/src/modules/credentials/interfaces/http/routes"
	refreshConsumer "mapexVault/src/modules/credentials/interfaces/message/consumers/refresh"

	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
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

	// Reconcile leader: exactly one pod runs the reseed sweep on the
	// vault_reconcile_interval ticker plus an immediate reseed on election; the KV
	// lease elects it and fails over automatically (no self-rescheduling message,
	// so no dedup window can silently kill the loop).
	c.Provide(func(params struct {
		container.In
		Store   natsModel.KeyValueStore `name:"vault-leader"`
		Service ports.CredentialServicePort
	}) *natsModel.LeaderElection {
		interval, _ := config.GetIntValue("vault_reconcile_interval")
		if interval <= 0 {
			interval = constants.VaultReconcileDefaultIntervalSeconds
		}
		le, err := natsModel.NewLeaderElection(params.Store, natsModel.LeaderElectionConfig{
			Key:       constants.VaultReconcileLeaderKey,
			Interval:  time.Duration(interval) * time.Second,
			OnTick:    func(ctx context.Context) { params.Service.RunReconcile(ctx) },
			OnElected: func() { params.Service.RunReconcile(context.Background()) },
		})
		if err != nil {
			logger.Panic("[MODULE:Credentials] reconcile leader build failed: " + err.Error())
		}
		return le
	})

	logger.Info("[MODULE:Credentials] Services registered")
}

// InitInterfaces registers HTTP routes for credentials.
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(params struct {
		container.In
		App     *web.App
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
		Leader  *natsModel.LeaderElection
	}) {
		// Refresh consumer — pulls from vault.schedule.fired (VAULT-SCHEDULE stream).
		// Per-credential timers fire here; HandleRefreshMessage refreshes the token
		// and immediately re-arms the next timer before ack.
		if refreshConsumer.NewConsumer(params.Bus, params.Service) == nil {
			log.Fatalf("failed to start refresh consumer")
		}

		// Reseed watchdog — the elected leader runs RunReconcile on a ticker; the
		// per-credential timers that were missing get reseeded. Replaces the old
		// self-rescheduling reconcile message.
		params.Leader.Start(context.Background())
	}); err != nil {
		log.Fatalf("failed to invoke credentials consumer: %v", err)
	}

	logger.Info("[MODULE:Credentials] Interfaces registered (HTTP routes + refresh consumer + reconcile leader)")
}
