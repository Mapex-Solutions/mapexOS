package services

import (
	"context"

	dsDto "http_gateway/src/modules/datasources/application/dtos"
	eventsConstants "http_gateway/src/modules/events/application/constants"

	processorContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/events"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// publishToJSExecutor sends the event to JS-Executor for script processing.
//
// This is the main event processing pipeline entry point. The JS-Executor will:
//   - Execute configured JavaScript transforms
//   - Route to appropriate route groups
//   - Forward to downstream services
//
// Payload includes sourceType="http" to indicate the origin gateway.
// JS-Executor uses this to determine how to resolve the assetUUID.
//
// Note: JS-Executor fetches pathKey, name, description from Asset cache (source of truth).
// We only send orgId and assetBind (required for cache lookup and asset resolution).
// This reduces payload size and ensures metadata consistency.
//
// Parameters:
//   - ctx: Context for cancellation and deadlines
//   - dataSource: The data source configuration
//   - event: The raw event data to process
//   - eventTrackerId: UUID for tracking event across the pipeline
//
// Returns:
//   - error: Returns error if NATS publish fails (critical path)
func (s *EventService) publishToJSExecutor(ctx context.Context, dataSource *dsDto.DataSourceResponse, event map[string]any, eventTrackerId string) error {
	// Build minimal dataSource for JS-Executor.
	// Only send fields required for processing:
	// - orgId: For Asset cache lookup key {orgId}/{assetUUID}
	// - assetBind: To resolve which Asset the event belongs to
	// Other fields (name, description, pathKey) are fetched from Asset cache.
	//
	// Payload is typed via the cross-service contract
	// packages/contracts/services/http_gateway/events.ProcessorExecutePayload;
	// JSON wire shape is identical to the previous ad-hoc map[string]any.
	payload := processorContracts.ProcessorExecutePayload{
		SourceType: "http",
		DataSource: processorContracts.ProcessorExecuteDataSource{
			OrgId:     dataSource.OrgId,
			AssetBind: dataSource.AssetBind,
		},
		Event:          event,
		EventTrackerId: eventTrackerId,
	}

	if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{
		Ctx:     ctx,
		Subject: eventsConstants.ProcessorJsExecuteSubject,
		Data:    payload,
		Headers: nil,
	}); err != nil {
		return &customErrors.ServerCustomError{Code: status.INTERNAL_SERVER_ERROR, Errors: []string{err.Error()}}
	}

	return nil
}
