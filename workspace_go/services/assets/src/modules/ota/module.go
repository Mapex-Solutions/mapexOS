package ota

import (
	"log"
	"time"

	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	assetsPorts "assets/src/modules/assets/application/ports"
	hmPorts "assets/src/modules/healthmonitor/application/ports"
	di "assets/src/modules/ota/application/di"
	"assets/src/modules/ota/application/ports"
	service "assets/src/modules/ota/application/services"
	adapters "assets/src/modules/ota/infrastructure/adapters"
	natsAdapter "assets/src/modules/ota/infrastructure/messaging/nats"
	collection "assets/src/modules/ota/infrastructure/persistence/mongo"
	redisAdapter "assets/src/modules/ota/infrastructure/persistence/redis"
	storage "assets/src/modules/ota/infrastructure/storage"
	httpRoutes "assets/src/modules/ota/interfaces/http/routes"
	statusConsumer "assets/src/modules/ota/interfaces/message/consumers/status"
	timersConsumer "assets/src/modules/ota/interfaces/message/consumers/timers"

	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	apikeymw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/apiKey"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitRepositories registers the OTA Mongo repositories.
func InitRepositories() {
	c := container.GetContainer()
	c.Provide(collection.NewFirmwareRepository)
	c.Provide(collection.NewPlanRepository)
	c.Provide(collection.NewExecutionRepository)
	logger.Info("[MODULE:OTA] Repositories registered")
}

// InitServices registers the OTA adapters + application services.
func InitServices() {
	c := container.GetContainer()

	// Firmware object store — an S3-compatible client bound to the firmware
	// bucket. Configured from the dedicated firmware_store_* keys (independent
	// of the TieredCache MinIO), so it can point at a different bucket or a
	// real S3 endpoint without touching the other caches.
	c.Provide(func() ports.FirmwareStorePort {
		bucket, _ := config.GetStringValue("firmware_store_bucket")
		if bucket == "" {
			bucket = "mapex-firmware"
		}
		// Endpoint/region/bucket are firmware-specific (browser-reachable); the
		// credential + auth mode are shared with the caches (svc-assets).
		endpoint, _ := config.GetStringValue("firmware_store_endpoint")
		region, _ := config.GetStringValue("firmware_store_region")
		useSSL, _ := config.GetConfigValue("firmware_store_use_ssl").(bool)
		accessKey, _ := config.GetStringValue("object_store_access_key")
		secretKey, _ := config.GetStringValue("object_store_secret_key")
		authIsNeeded := config.GetConfigValue("object_store_auth_is_needed").(bool)
		mc, err := minioModel.New(minioModel.Config{
			Endpoint:        endpoint,
			AuthIsNeeded:    authIsNeeded,
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
			UseSSL:          useSSL,
			Region:          region,
			BucketName:      bucket,
		})
		if err != nil {
			logger.Panic("[MODULE:OTA] firmware store client failed: " + err.Error())
		}
		return storage.NewMinIOFirmwareStore(mc)
	})

	// NATS scheduler + firmware abandon scheduler.
	c.Provide(func(params struct {
		container.In
		SM natsModel.ScheduleManager `name:"core"`
	}) ports.OTASchedulerPort {
		return natsAdapter.NewOTAScheduler(params.SM)
	})
	c.Provide(func(params struct {
		container.In
		SM natsModel.ScheduleManager `name:"core"`
	}) ports.FirmwareSchedulerPort {
		return natsAdapter.NewFirmwareScheduler(params.SM)
	})

	// Edge dispatch + status history publisher.
	c.Provide(func(params struct {
		container.In
		Pub natsModel.Publisher `name:"core"`
	}) ports.EdgeDispatchPort {
		return natsAdapter.NewEdgeDispatchAdapter(params.Pub)
	})
	c.Provide(func(params struct {
		container.In
		Pub natsModel.Publisher `name:"core"`
	}) ports.StatusHistoryPublisherPort {
		return natsAdapter.NewStatusHistoryPublisher(params.Pub)
	})

	// Redis live state.
	c.Provide(func(rc *redisModel.RedisClient) ports.LiveStatePort {
		return redisAdapter.NewLiveStateAdapter(rc, durationCfg("ota_live_state_ttl", 24*time.Hour))
	})

	// Cross-module read/write adapters (presence, asset read, template switch).
	c.Provide(func(repo hmPorts.HealthRepository) ports.PresenceReaderPort {
		return adapters.NewPresenceReaderAdapter(repo)
	})
	c.Provide(func(svc assetsPorts.AssetServicePort) ports.AssetReaderPort {
		return adapters.NewAssetReaderAdapter(svc)
	})
	c.Provide(func(svc assetsPorts.AssetServicePort) ports.AssetTemplateSwitcherPort {
		return adapters.NewTemplateSwitcherAdapter(svc)
	})

	// Firmware service.
	c.Provide(func(deps di.FirmwareServiceDI) ports.OTAFirmwareServicePort {
		return service.NewFirmwareService(service.FirmwareServiceDeps{
			FirmwareRepo: deps.FirmwareRepo,
			Store:        deps.Store,
			Scheduler:    deps.Scheduler,
			PresignTTL:   durationCfg("ota_presigned_url_ttl", 30*time.Minute),
			AbandonTTL:   durationCfg("ota_orphan_gc_ttl", 24*time.Hour),
		})
	})

	// Plan service.
	c.Provide(func(deps di.PlanServiceDI) ports.OTAPlanServicePort {
		return service.NewPlanService(service.PlanServiceDeps{
			PlanRepo:      deps.PlanRepo,
			ExecutionRepo: deps.ExecutionRepo,
			FirmwareRepo:  deps.FirmwareRepo,
			Scheduler:     deps.Scheduler,
			Store:         deps.Store,
			PresignTTL:    durationCfg("ota_presigned_url_ttl", 30*time.Minute),
		})
	})

	// Reconciler engine.
	c.Provide(func(deps di.ReconcilerDI) *service.Reconciler {
		return service.NewReconciler(service.ReconcilerDeps{
			PlanRepo:      deps.PlanRepo,
			ExecutionRepo: deps.ExecutionRepo,
			FirmwareRepo:  deps.FirmwareRepo,
			Store:         deps.Store,
			Presence:      deps.Presence,
			Edge:          deps.Edge,
			AssetReader:   deps.AssetReader,
			PresignTTL:    durationCfg("ota_presigned_url_ttl", 30*time.Minute),
			MaxAttempts:   intCfg("ota_max_attempts", 2),
			DefaultRate:   intCfg("ota_rate_per_minute", 60),
		})
	})

	// Reconciler timers (start/close + global scan). The concrete *Reconciler is
	// injected alongside the DI struct.
	c.Provide(func(deps di.ReconcilerTimersDI, rec *service.Reconciler) *service.ReconcilerTimers {
		return service.NewReconcilerTimers(service.ReconcilerTimersDeps{
			PlanRepo:      deps.PlanRepo,
			ExecutionRepo: deps.ExecutionRepo,
			FirmwareRepo:  deps.FirmwareRepo,
			Store:         deps.Store,
			Scheduler:     deps.Scheduler,
			Reconciler:    rec,
			ScanInterval:  durationCfg("ota_scan_interval", time.Minute),
		})
	})

	// Status handler.
	c.Provide(func(deps di.StatusHandlerDI) *service.StatusHandler {
		return service.NewStatusHandler(service.StatusHandlerDeps{
			ExecutionRepo:    deps.ExecutionRepo,
			PlanRepo:         deps.PlanRepo,
			LiveState:        deps.LiveState,
			History:          deps.History,
			TemplateSwitcher: deps.TemplateSwitcher,
		})
	})

	// Device poll port — the reconciler serves the HTTP device poll (fresh URL
	// minting + attempt accounting live with the dispatch semantics).
	c.Provide(func(rec *service.Reconciler) ports.DeviceJobPort {
		return rec
	})

	logger.Info("[MODULE:OTA] Services registered")
}

// InitInterfaces registers the OTA HTTP routes, starts the consumers, and kicks
// the single global pacing scan.
func InitInterfaces() {
	c := container.GetContainer()

	if err := c.Invoke(func(app *web.App, firmwareSvc ports.OTAFirmwareServicePort, planSvc ports.OTAPlanServicePort) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")
		routesV1 := app.Group(
			"/api/v1/ota",
			ctxInjector.ContextInjector(ctxTimeout),
			authmw.AuthMiddleware(config.GetAuthConfig()),
		)
		httpRoutes.RegisterRoutes(routesV1, firmwareSvc, planSvc)
	}); err != nil {
		log.Fatalf("failed to invoke OTA routes: %v", err)
	}

	// Internal routes (API key gated) — the HTTP gateway relays the device poll
	// here; device-facing auth stays on the gateway.
	if err := c.Invoke(func(app *web.App, jobs ports.DeviceJobPort) {
		ctxTimeout, _ := config.GetIntValue("ctx_timeout")
		internalApiKey, _ := config.GetStringValue("internal_api_key")
		internalRoutes := app.Group(
			"/internal/ota",
			ctxInjector.ContextInjector(ctxTimeout),
			apikeymw.ApiKeyAuthMiddleware(internalApiKey),
		)
		httpRoutes.RegisterInternalRoutes(internalRoutes, jobs)
	}); err != nil {
		log.Fatalf("failed to invoke OTA internal routes: %v", err)
	}

	if err := c.Invoke(func(params struct {
		container.In
		Bus         *natsModel.Bus `name:"core"`
		Timers      *service.ReconcilerTimers
		Status      *service.StatusHandler
		FirmwareSvc ports.OTAFirmwareServicePort
		Scheduler   ports.OTASchedulerPort
	}) {
		timersConsumer.NewConsumer(params.Bus, params.Timers, params.FirmwareSvc)
		statusConsumer.NewConsumer(params.Bus, params.Status)
		// Kick the single global pacing scan (self-reschedules thereafter).
		_ = params.Scheduler.ScheduleScan(time.Now().Add(durationCfg("ota_scan_interval", time.Minute)))
	}); err != nil {
		log.Fatalf("failed to invoke OTA consumers: %v", err)
	}

	logger.Info("[MODULE:OTA] Interfaces registered (routes + consumers)")
}

// durationCfg reads a config value in SECONDS, falling back to def.
func durationCfg(key string, def time.Duration) time.Duration {
	secs, _ := config.GetIntValue(key)
	if secs <= 0 {
		return def
	}
	return time.Duration(secs) * time.Second
}

// intCfg reads a positive int config value, falling back to def.
func intCfg(key string, def int) int {
	v, _ := config.GetIntValue(key)
	if v <= 0 {
		return def
	}
	return v
}
