package routes

import (
	"github.com/gofiber/fiber/v2"

	"http_gateway/src/bootstrap"
	dsPort "http_gateway/src/modules/datasources/application/ports"
	"http_gateway/src/modules/events/application/dtos"
	"http_gateway/src/modules/events/application/ports"
	"http_gateway/src/modules/events/interfaces/http/handlers"
	"http_gateway/src/modules/events/interfaces/http/middlewares"

	ctxInjector "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/contextInjector"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
)

// RegisterRoutes registers all event-related HTTP routes under their
// respective base paths. Hexagonal: accepts service port interfaces, not
// concrete implementations.
//
// Routes registered:
//
//	POST /api/v1/events     - Webhook receiver (publishes to processor.js.execute).
//	                          Body is arbitrary device telemetry; only the query
//	                          (?ds={dataSourceId}) is validated upstream.
//	POST /api/v1/heartbeat  - Explicit-mode HTTP heartbeat (publishes to
//	                          mapexos.asset.heartbeat.{orgId}). Body is required:
//	                          { "assetUUID": "<v>" } — TKT-2026-0036 reformulation
//	                          dropped the legacy AssetBind.Type='fixedAssetId'
//	                          constraint, so any DataSource shape works as long
//	                          as the auth chain accepts it.
//
// Both routes share the same auth middleware chain (CustomAuthMiddleware) so
// DataSource resolution + per-DS auth (apiKey/jwt/oauth2/ip_whitelist) works
// identically. The middleware also rejects requests on disabled DataSources
// (403 — TKT-2026-0036).
//
// Parameters:
//   - app: Fiber app used to mount the per-path groups
//   - ctxTimeout: Timeout (seconds) configured on the request context
//   - service: Event service port interface
//   - dtService: Data source service port (used by CustomAuthMiddleware)
//   - m: Service-specific metrics for instrumentation
func RegisterRoutes(app *fiber.App, ctxTimeout int, service ports.EventServicePort, dtService dsPort.DataSourceServicePort, m *bootstrap.HttpGatewayMetrics) {
	// /events — query-only validation (body is arbitrary).
	eventIdentificationDto := validation.NewValidation(nil, &dtos.EvenIdentificationDto{}, nil)

	// /heartbeat — body { assetUUID } + query (?ds={dataSourceId}) in one validator.
	// mapexGoKit signature: NewValidation(bodyDTO, queryDTO, paramsDTO).
	heartbeatValidation := validation.NewValidation(
		&dtos.HeartbeatRequestDTO{},
		&dtos.EvenIdentificationDto{},
		nil,
	)

	eventsV1 := app.Group("/api/v1/events", ctxInjector.ContextInjector(ctxTimeout))
	eventsV1.Post(
		"/",
		validation.ValidationMiddleware(eventIdentificationDto),
		middlewares.CustomAuthMiddleware(dtService, service, m),
		handlers.ProcessEvent(service, m),
	)

	heartbeatV1 := app.Group("/api/v1/heartbeat", ctxInjector.ContextInjector(ctxTimeout))
	heartbeatV1.Post(
		"/",
		validation.ValidationMiddleware(heartbeatValidation),
		middlewares.CustomAuthMiddleware(dtService, service, m),
		handlers.ProcessHeartbeat(service, m),
	)
}
