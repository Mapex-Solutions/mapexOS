package routes

import (
	"http_gateway/src/bootstrap"
	dsPort "http_gateway/src/modules/datasources/application/ports"
	"http_gateway/src/modules/events/application/dtos"
	"http_gateway/src/modules/events/application/ports"
	"http_gateway/src/modules/events/interfaces/http/handlers"
	"http_gateway/src/modules/events/interfaces/http/middlewares"

	downlink "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/downlink"
	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers the inbound ingestion routes under their respective
// base paths. Hexagonal: accepts service port interfaces, not concrete
// implementations.
//
// Routes registered:
//
//	POST /api/v1/events     - Telemetry webhook receiver (publishes to
//	                          processor.js.execute). Body is arbitrary device
//	                          telemetry; only the query (?ds={dataSourceId}) is
//	                          validated.
//	POST /api/v1/heartbeat  - Explicit-mode HTTP heartbeat (publishes to
//	                          mapexos.asset.heartbeat.{orgId}). Body is required:
//	                          { "assetUUID": "<v>" }.
//
// Both routes share the same auth middleware chain (CustomAuthMiddleware) so
// DataSource resolution + per-DS auth (apiKey/jwt/oauth2/ip_whitelist) works
// identically. The middleware also rejects requests on disabled DataSources (403).
//
// Routes are registered through the swagger wrapper: each NewValidation declares
// the input contract once (used to both validate and document); the module tag is
// declared once on Wrap; summary, description, and response type are attached via
// the fluent builder.
//
// Parameters:
//   - app: Fiber app used to mount the per-path groups
//   - ctxTimeout: Timeout (seconds) configured on the request context
//   - service: Event service port interface
//   - dtService: Data source service port (used by CustomAuthMiddleware)
//   - m: Service-specific metrics for instrumentation
func RegisterRoutes(app *web.App, ctxTimeout int, service ports.EventServicePort, dtService dsPort.DataSourceServicePort, m *bootstrap.HttpGatewayMetrics) {

	// /events — query-only validation (body is arbitrary device telemetry).
	eventIdentificationDto := validation.NewValidation(nil, &dtos.EvenIdentificationDto{}, nil)
	eventsV1 := app.Group("/api/v1/events", ctxInjector.ContextInjector(ctxTimeout))
	swagger.Wrap(eventsV1).Tag("Ingestion").
		Post("/", eventIdentificationDto, swagger.Expose,
			middlewares.CustomAuthMiddleware(dtService, service, m),
			handlers.ProcessEvent(service, m),
		).
		Summary("Ingest telemetry event").
		Description("Webhook receiver for device telemetry. Authenticated per-DataSource via the ?ds={dataSourceId} query parameter; the body is arbitrary device telemetry. Accepted events are published to NATS for downstream processing.").
		Returns(map[string]bool{})

	// /heartbeat — body { assetUUID } + query (?ds={dataSourceId}) in one validator.
	// mapexGoKit signature: NewValidation(bodyDTO, queryDTO, paramsDTO).
	heartbeatValidation := validation.NewValidation(
		&dtos.HeartbeatRequestDTO{},
		&dtos.EvenIdentificationDto{},
		nil,
	)
	heartbeatV1 := app.Group("/api/v1/heartbeat", ctxInjector.ContextInjector(ctxTimeout))
	swagger.Wrap(heartbeatV1).Tag("Ingestion").
		Post("/", heartbeatValidation, swagger.Expose,
			middlewares.CustomAuthMiddleware(dtService, service, m),
			handlers.ProcessHeartbeat(service, m),
		).
		Summary("Send asset heartbeat").
		Description("Explicit-mode HTTP heartbeat for an asset. Authenticated per-DataSource via the ?ds={dataSourceId} query parameter; the body carries { assetUUID }. Publishes a fire-and-forget heartbeat to NATS to drive the asset's online state.").
		Returns(map[string]bool{})

	// /ota — device-facing OTA endpoints (poll + status report). Same
	// per-DataSource auth chain as /events and /heartbeat.
	otaV1 := app.Group("/api/v1/ota", ctxInjector.ContextInjector(ctxTimeout))
	ota := swagger.Wrap(otaV1).Tag("Ingestion")

	otaStatusValidation := validation.NewValidation(
		&dtos.OTAStatusRequestDTO{},
		&dtos.EvenIdentificationDto{},
		nil,
	)
	ota.Post("/status", otaStatusValidation, swagger.Expose,
		middlewares.CustomAuthMiddleware(dtService, service, m),
		handlers.ProcessOTAStatus(service, m),
	).
		Summary("Report OTA status").
		Description("Device OTA progress report. Authenticated per-DataSource via the ?ds={dataSourceId} query parameter; the body carries { assetUUID, executionId, status, progress, error?, message? }. The report is normalized into the OTA status advisory consumed by the Asset MS.").
		Returns(map[string]bool{})

	otaJobsValidation := validation.NewValidation(nil, &dtos.EvenIdentificationDto{}, nil)
	ota.Get("/jobs", otaJobsValidation, swagger.Expose,
		middlewares.CustomAuthMiddleware(dtService, service, m),
		handlers.GetOTAJob(service, m),
	).
		Summary("Poll pending OTA job").
		Description("Device poll for its pending OTA firmware-update command. Authenticated per-DataSource via the ?ds={dataSourceId} query parameter plus the &assetUUID={assetUUID} identity. Returns the command with a freshly minted presigned download URL, or 204 when the device has nothing actionable.").
		Returns(&downlink.OTAUpdateCommand{})
}
