package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	dsDto "http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/events/application/di"
	"http_gateway/src/modules/events/application/ports"
)

// Compile-time check to ensure EventService implements EventServicePort interface.
// This will cause a compilation error if the interface is not fully implemented.
var _ ports.EventServicePort = (*EventService)(nil)

// New creates and returns a new instance of EventService.
//
// This constructor follows Hexagonal Architecture by:
//   - Accepting dependencies through a DI struct (single parameter pattern)
//   - Returning the service port interface (not concrete type)
//   - Enabling loose coupling and testability
//
// Parameters:
//   - deps: Aggregated dependencies (NATS bus) injected by dig
//
// Returns:
//   - EventServicePort: The service port interface implementation
func New(deps di.EventServiceDependenciesInjection) ports.EventServicePort {
	return &EventService{
		deps: deps,
	}
}

// ProcessEvent processes an event by encoding it and publishing it to a NATS subject.
//
// This method receives event data from external sources (webhooks, APIs) and forwards
// it to the processing pipeline via NATS messaging for asynchronous processing.
//
// The event is published to processor.js.execute where JS-Executor will:
//   - Execute configured JavaScript transforms
//   - Route to appropriate route groups
//   - Forward to downstream services
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines
//   - event: A map representing the event data to be processed
//   - dataSource: A pointer to a DataSourceResponse containing information about the data source
//
// Returns:
//   - A map indicating the success status of the operation
//   - An error if the operation fails (NATS publish failure, encoding error, etc.)
func (s *EventService) ProcessEvent(ctx context.Context, event map[string]any, dataSource *dsDto.DataSourceResponse) (map[string]bool, error) {
	start := time.Now()

	// Generate eventTrackerId for HTTP events
	// This allows end-to-end tracking across the entire pipeline
	eventTrackerId := uuid.New().String()

	// Publish to JS-Executor (pipeline)
	if err := s.publishToJSExecutor(ctx, dataSource, event, eventTrackerId); err != nil {
		// Metrics: NATS publish failed — count processed+published as error
		s.deps.Metrics.EventsProcessed.WithLabelValues("error").Inc()
		s.deps.Metrics.EventsPublished.WithLabelValues("processor.js.execute", "error").Inc()
		return nil, err
	}

	// Metrics: event fully processed — count success, record handler-only latency (excludes auth)
	s.deps.Metrics.EventsProcessed.WithLabelValues("success").Inc()
	s.deps.Metrics.EventsPublished.WithLabelValues("processor.js.execute", "success").Inc()
	s.deps.Metrics.EventProcessingDuration.Observe(time.Since(start).Seconds())

	return map[string]bool{"success": true}, nil
}

// PublishAuthFailure publishes a raw event with success=false when authentication
// fails (or when CustomAuthMiddleware rejects a disabled DataSource — see
// TKT-2026-0036). This is called by the auth middleware on every rejection
// path for the security audit trail. Thin orchestration: build payload →
// fire-and-forget publish (privates in event_handler_authfailure.go).
func (s *EventService) PublishAuthFailure(dataSource *dsDto.DataSourceResponse, event map[string]any, eventTrackerId string, errorMsg string) {
	payload := s.buildAuthFailurePayload(dataSource, event, eventTrackerId, errorMsg)
	s.publishAuthFailureFireAndForget(payload)
}

// ProcessHeartbeat handles POST /api/v1/heartbeat?ds={ds}.
//
// Used by explicit-mode HTTP-protocol assets (HealthMonitorConfig.heartbeatMode='explicit'
// + protocol='http'): the device POSTs `{ assetUUID }` and the platform
// publishes a fire-and-forget heartbeat to mapexos.asset.heartbeat.{orgId}.
// orgId and pathKey come from the resolved DataSource (server-side, c.Locals);
// assetUUID comes from the request body — a compromised body cannot spoof a
// different tenant because the authentication is per-DataSource (TKT-2026-0036).
//
// Steps: validate dataSource + assetUUID -> compose payload -> fire-and-forget
// core publish -> emit success metric.
func (s *EventService) ProcessHeartbeat(ctx context.Context, dataSource *dsDto.DataSourceResponse, assetUUID string) error {
	start := time.Now()
	if err := s.validateHeartbeat(start, dataSource, assetUUID); err != nil {
		return err
	}
	plan := s.buildHeartbeatPayload(dataSource, assetUUID)
	if err := s.publishHeartbeatCore(ctx, start, plan); err != nil {
		return err
	}
	s.recordHeartbeatResult(start, "success")
	return nil
}

