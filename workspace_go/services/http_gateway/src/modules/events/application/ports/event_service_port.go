package ports

import (
	ctx "context"
	dsDto "http_gateway/src/modules/datasources/application/dtos"
)

// EventServicePort defines the business operations for processing events.
//
// Following Hexagonal Architecture, this port interface:
//   - Defines the contract for event processing operations
//   - Abstracts the service implementation from its consumers
//   - Enables dependency inversion and testability
//   - Decouples the application layer from infrastructure concerns
//
// This service receives events from external sources (webhooks, APIs)
// and publishes them to NATS for asynchronous processing.
type EventServicePort interface {
	// ProcessEvent processes an event by encoding it and publishing it to NATS.
	//
	// This method receives event data from external sources and forwards it to
	// the processing pipeline via NATS messaging.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines
	//   - event: A map representing the event data to be processed
	//   - dataSource: The data source information associated with this event
	//
	// Returns:
	//   - A map indicating the success status of the operation
	//   - An error if the operation fails (NATS publish failure, encoding error, etc.)
	ProcessEvent(ctx ctx.Context, event map[string]any, dataSource *dsDto.DataSourceResponse) (map[string]bool, error)

	// PublishAuthFailure publishes a raw event with success=false when authentication fails.
	//
	// This method is called by the auth middleware when authentication validation fails.
	// It always publishes for security monitoring purposes.
	//
	// Parameters:
	//   - dataSource: The data source information (may have partial data if lookup failed)
	//   - event: The event payload (may be nil if body parsing failed)
	//   - eventTrackerId: UUID for tracking event across the pipeline
	//   - errorMsg: The authentication error message
	PublishAuthFailure(dataSource *dsDto.DataSourceResponse, event map[string]any, eventTrackerId string, errorMsg string)

	// ProcessHeartbeat publishes a fire-and-forget heartbeat to
	// mapexos.asset.heartbeat.{orgId} for the asset whose UUID is carried
	// in the request body. Used by POST /api/v1/heartbeat?ds={dataSourceId}
	// (HTTP path of the explicit-mode heartbeat — TKT-2026-0036 reformulation).
	//
	// The published payload mirrors the shape js-executor publishes today
	// ({orgId, assetUUID, pathKey, ts}) so the assets/healthmonitor consumer
	// is origin-agnostic. orgId and pathKey come from the resolved DataSource
	// (server-side, c.Locals); assetUUID comes from the request body — a
	// compromised body cannot spoof a different tenant.
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and deadlines
	//   - dataSource: Resolved DataSource (set by CustomAuthMiddleware) — provides orgId + pathKey
	//   - assetUUID: Body-supplied asset identifier (validated upstream by the request DTO)
	//
	// Returns:
	//   - error if the DataSource is missing required fields, the assetUUID is empty, or the publish fails
	ProcessHeartbeat(ctx ctx.Context, dataSource *dsDto.DataSourceResponse, assetUUID string) error
}
