package services

import (
	"context"
	"fmt"

	"events/src/modules/events/application/di"
	"events/src/modules/events/application/dtos"
	"events/src/modules/events/application/ports"
	"events/src/modules/events/domain/entities"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time check to ensure EventService implements EventServicePort.
var _ ports.EventServicePort = (*EventService)(nil)

// New creates and returns a new instance of EventService.
//
// This constructor follows Hexagonal Architecture by:
//   - Accepting dependencies through a DI struct (single parameter pattern)
//   - Returning the service port interface (not concrete type)
//   - Enabling loose coupling and testability
func New(deps di.EventServiceDependenciesInjection) ports.EventServicePort {
	return &EventService{deps: deps}
}

// ProcessEvent stores a single processed event in ClickHouse via the legacy
// single-message NATS path. Steps: parse the body, extract the required
// fields, resolve retention, build the entity, persist.
func (s *EventService) ProcessEvent(data []byte, index int, headers map[string][]string) error {
	logger.Info(fmt.Sprintf("[SERVICE:Event] Processing event %d", index))
	eventData, fields, err := s.parseLegacyEvent(data)
	if err != nil {
		return err
	}
	payloadJSON, metadataJSON, err := s.marshalLegacyEventPayloads(eventData, headers)
	if err != nil {
		return err
	}
	c := context.Background()
	event := s.buildLegacyEvent(fields, payloadJSON, metadataJSON, s.resolveRetentionOrDefault(c, fields.orgId, "events"))
	return s.persistLegacyEvent(c, event, fields)
}

// ProcessRawEventBatch persists a NATS batch of raw events. Three-phase
// pipeline: parallel parse/validate -> bulk insert -> Ack/Nack/Reject.
func (s *EventService) ProcessRawEventBatch(messages []*natsModel.Message) error {
	return orchestrateBatch(s, "raw", "eventsRaw", messages,
		s.processRawEventMessage,
		func(ctx context.Context, entities []*entities.RawEvent) error {
			return s.deps.EventRepo.SaveRawEventBatch(ctx, entities)
		},
	)
}

// GetEventsRaw retrieves raw events with cursor-based pagination scoped to
// the caller's org context.
func (s *EventService) GetEventsRaw(c context.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsRawQueryDto) (*dtos.EventsRawCursorResultDto, error) {
	orgFilter, err := s.buildOrgFilterCH(requestContext, query)
	if err != nil {
		return nil, err
	}
	cursorOpts := s.buildEventCursorOpts(query)
	result, err := s.deps.EventRepo.QueryEventsRawCursor(c, orgFilter, query.EventTrackerId, query.ThreadId, query.Source, query.Success, query.StartTime, query.EndTime, cursorOpts)
	if err != nil {
		return nil, err
	}
	return s.buildEventsRawCursorResult(result), nil
}

// ProcessJsExecEventBatch persists a NATS batch of JS-executor events.
// Same three-phase pipeline as ProcessRawEventBatch with JS-exec parser.
func (s *EventService) ProcessJsExecEventBatch(messages []*natsModel.Message) error {
	return orchestrateBatch(s, "jsexec", "eventsJsExecutor", messages,
		s.processJsExecEventMessage,
		func(ctx context.Context, entities []*entities.JsExecEvent) error {
			return s.deps.EventRepo.SaveJsExecEventBatch(ctx, entities)
		},
	)
}

// GetEventsJsExec retrieves JS Executor events with cursor-based pagination.
func (s *EventService) GetEventsJsExec(c context.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsJsExecQueryDto) (*dtos.EventsJsExecCursorResultDto, error) {
	orgFilter, err := s.buildOrgFilterCH(requestContext, query)
	if err != nil {
		return nil, err
	}
	cursorOpts := s.buildEventCursorOpts(query)
	execTimeOpStr := s.execTimeOpToStringPtr(query.ExecTimeOp)
	result, err := s.deps.EventRepo.QueryEventsJsExecCursor(c, orgFilter, query.EventTrackerId, query.ThreadId, query.Success, query.StartTime, query.EndTime, execTimeOpStr, query.ExecTimeValue, query.ExecTimeValueEnd, cursorOpts)
	if err != nil {
		return nil, err
	}
	return s.buildEventsJsExecCursorResult(result), nil
}

// ProcessDLQEventBatch persists a NATS batch of dead-lettered events. DLQ
// flow Acks parse failures (no DLQ-of-DLQ) and never Nacks.
func (s *EventService) ProcessDLQEventBatch(messages []*natsModel.Message) error {
	return orchestrateDLQBatch(s, "dlq", "eventsDLQ", messages,
		s.processDLQEventMessage,
		func(ctx context.Context, entities []*entities.DLQEvent) error {
			return s.deps.EventRepo.SaveDLQEventBatch(ctx, entities)
		},
	)
}

// GetEventsDLQ retrieves DLQ events with cursor-based pagination.
func (s *EventService) GetEventsDLQ(c context.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsDLQQueryDto) (*dtos.EventsDLQCursorResultDto, error) {
	orgFilter, err := s.buildOrgFilterCH(requestContext, query)
	if err != nil {
		return nil, err
	}
	cursorOpts := s.buildEventCursorOpts(query)
	result, err := s.deps.EventRepo.QueryEventsDLQCursor(c, orgFilter, query.EventTrackerId, query.ServiceName, query.ServiceType, query.EventType, query.LastError, query.StartTime, query.EndTime, cursorOpts)
	if err != nil {
		return nil, err
	}
	return s.buildEventsDLQCursorResult(result), nil
}

// GetEventsDLQCounts returns DLQ entry counts grouped by service type.
func (s *EventService) GetEventsDLQCounts(c context.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsDLQCountsQueryDto) (*dtos.EventsDLQCountsResultDto, error) {
	orgFilter, err := s.buildOrgFilterCH(requestContext, query)
	if err != nil {
		return nil, err
	}
	counts, err := s.deps.EventRepo.CountByServiceType(c, orgFilter, query.StartTime, query.EndTime)
	if err != nil {
		return nil, err
	}
	return s.buildEventsDLQCountsResult(counts), nil
}

// ProcessRouterEventBatch persists a NATS batch of router-execution events.
func (s *EventService) ProcessRouterEventBatch(messages []*natsModel.Message) error {
	return orchestrateBatch(s, "router", "eventsRouter", messages,
		s.processRouterEventMessage,
		func(ctx context.Context, entities []*entities.RouterEvent) error {
			return s.deps.EventRepo.SaveRouterEventBatch(ctx, entities)
		},
	)
}

// GetEventsRouter retrieves router events with cursor-based pagination.
func (s *EventService) GetEventsRouter(c context.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsRouterQueryDto) (*dtos.EventsRouterCursorResultDto, error) {
	orgFilter, err := s.buildOrgFilterCH(requestContext, query)
	if err != nil {
		return nil, err
	}
	cursorOpts := s.buildEventCursorOpts(query)
	result, err := s.deps.EventRepo.QueryEventsRouterCursor(c, orgFilter, query.EventTrackerId, query.ThreadId, query.AssetId, query.RouterId, query.Success, query.PublishedCount, query.StartTime, query.EndTime, cursorOpts)
	if err != nil {
		return nil, err
	}
	return s.buildEventsRouterCursorResult(result), nil
}

// ProcessBusinessRuleEventBatch persists a NATS batch of business-rule
// execution events.
func (s *EventService) ProcessBusinessRuleEventBatch(messages []*natsModel.Message) error {
	return orchestrateBatch(s, "businessrule", "eventsBusinessRule", messages,
		s.processBusinessRuleEventMessage,
		func(ctx context.Context, entities []*entities.BusinessRuleEvent) error {
			return s.deps.EventRepo.SaveBusinessRuleEventBatch(ctx, entities)
		},
	)
}

// GetEventsBusinessRule retrieves business rule events with cursor-based pagination.
func (s *EventService) GetEventsBusinessRule(c context.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsBusinessRuleQueryDto) (*dtos.EventsBusinessRuleCursorResultDto, error) {
	orgFilter, err := s.buildOrgFilterCH(requestContext, query)
	if err != nil {
		return nil, err
	}
	cursorOpts := s.buildEventCursorOpts(query)
	result, err := s.deps.EventRepo.QueryEventsBusinessRuleCursor(c, orgFilter, query.EventTrackerId, query.ThreadId, query.RuleId, query.BusinessRuleId, query.Matched, query.StartTime, query.EndTime, cursorOpts)
	if err != nil {
		return nil, err
	}
	return s.buildEventsBusinessRuleCursorResult(result), nil
}

// ProcessTriggerEventBatch persists a NATS batch of trigger-execution events.
func (s *EventService) ProcessTriggerEventBatch(messages []*natsModel.Message) error {
	return orchestrateBatch(s, "trigger", "eventsTrigger", messages,
		s.processTriggerEventMessage,
		func(ctx context.Context, entities []*entities.TriggerEvent) error {
			return s.deps.EventRepo.SaveTriggerEventBatch(ctx, entities)
		},
	)
}

// GetEventsTrigger retrieves trigger events with cursor-based pagination.
func (s *EventService) GetEventsTrigger(c context.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsTriggerQueryDto) (*dtos.EventsTriggerCursorResultDto, error) {
	orgFilter, err := s.buildOrgFilterCH(requestContext, query)
	if err != nil {
		return nil, err
	}
	cursorOpts := s.buildEventCursorOpts(query)
	result, err := s.deps.EventRepo.QueryEventsTriggerCursor(c, orgFilter, query.EventTrackerId, query.TriggerId, query.TriggerType, query.Category, query.Source, query.Success, query.StartTime, query.EndTime, cursorOpts)
	if err != nil {
		return nil, err
	}
	return s.buildEventsTriggerCursorResult(result), nil
}

// ProcessWorkflowEventBatch persists a NATS batch of workflow-execution events.
func (s *EventService) ProcessWorkflowEventBatch(messages []*natsModel.Message) error {
	return orchestrateBatch(s, "workflow", "eventsWorkflow", messages,
		s.processWorkflowEventMessage,
		func(ctx context.Context, entities []*entities.WorkflowEvent) error {
			return s.deps.EventRepo.SaveWorkflowEventBatch(ctx, entities)
		},
	)
}

// GetEventsWorkflow retrieves workflow events with cursor-based pagination.
func (s *EventService) GetEventsWorkflow(c context.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsWorkflowQueryDto) (*dtos.EventsWorkflowCursorResultDto, error) {
	orgFilter, err := s.buildOrgFilterCH(requestContext, query)
	if err != nil {
		return nil, err
	}
	cursorOpts := s.buildEventCursorOpts(query)
	result, err := s.deps.EventRepo.QueryEventsWorkflowCursor(c, orgFilter, query.EventTrackerId, query.WorkflowUUID, query.InstanceId, query.DefinitionId, query.Status, query.Success, query.StartTime, query.EndTime, cursorOpts)
	if err != nil {
		return nil, err
	}
	return s.buildEventsWorkflowCursorResult(result), nil
}

// GetWorkflowEventByExecutionId fetches a single workflow event by its
// execution_id from ClickHouse, scoped to the caller's org.
func (s *EventService) GetWorkflowEventByExecutionId(c context.Context, requestContext *reqCtx.RequestContext, executionId string) (*dtos.EventsWorkflowResponseDto, error) {
	entity, err := s.deps.EventRepo.FindWorkflowEventByExecutionId(c, *requestContext.OrgContext, executionId)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("workflow event not found: %s", executionId)
	}
	return s.buildWorkflowEventResponse(entity)
}

// ProcessEventStoreBatch persists a NATS batch of processed events with EVA
// field mapping (template lookup happens inside the per-message parser).
func (s *EventService) ProcessEventStoreBatch(messages []*natsModel.Message) error {
	return orchestrateBatch(s, "store", "events", messages,
		s.processEventStoreMessage,
		func(ctx context.Context, entities []*entities.Event) error {
			return s.deps.EventRepo.SaveEventStoreBatch(ctx, entities)
		},
	)
}

// GetEventsStore retrieves processed events with cursor-based pagination.
// List view: no EVA fields, just core data + payload.
func (s *EventService) GetEventsStore(c context.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsStoreQueryDto) (*dtos.EventsStoreCursorResultDto, error) {
	orgFilter, err := s.buildOrgFilterCH(requestContext, query)
	if err != nil {
		return nil, err
	}
	cursorOpts := s.buildEventCursorOpts(query)
	result, err := s.deps.EventRepo.QueryEventsStoreCursor(c, orgFilter, query.EventTrackerId, query.ThreadId, query.AssetId, query.AssetTemplateId, query.EventType, query.Source, query.StartTime, query.EndTime, query.EvaFilters, cursorOpts)
	if err != nil {
		return nil, err
	}
	return s.buildEventsStoreCursorResult(result), nil
}

// GetEventStoreDetail returns a single event with EVA fields resolved
// against the originating template (asset or rule).
func (s *EventService) GetEventStoreDetail(c context.Context, eventTrackerId string) (*dtos.EventsStoreDetailResponseDto, error) {
	entity, err := s.deps.EventRepo.GetEventStoreByTrackerId(c, eventTrackerId)
	if err != nil {
		return nil, err
	}
	advancedSearch := s.resolveEventStoreEvaFields(c, entity)
	return s.buildEventsStoreDetailResponse(entity, advancedSearch), nil
}

// HandleTemplateInvalidate processes one FANOUT message from
// mapexos.fanout.template.invalidate: parse the payload, build the
// {orgId}/{templateId} cache key, and clear the local TieredCache (L0+L1).
// The kit's SubscribeFanout callback Acks regardless of outcome.
func (s *EventService) HandleTemplateInvalidate(msg *natsModel.Message) {
	event, ok := s.parseTemplateInvalidatePayload(msg)
	if !ok {
		return
	}
	s.applyTemplateInvalidate(event)
}
