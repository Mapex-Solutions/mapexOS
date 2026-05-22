package services

import (
	ctx "context"
	"fmt"
	"time"

	"events/src/modules/events/application/dtos"
	"events/src/modules/events/domain/entities"

	eventsContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// cursorQuery is the small structural interface implemented by every events
// query DTO. It exposes the four cursor inputs needed to build TimeCursorOpts.
type cursorQuery interface {
	GetLimit() int
	GetDirection() string
	GetSortAsc() bool
	GetCursor() *time.Time
}

// buildOrgFilterCH wraps orgfilter.BuildOrgFilterClickHouse so callers in
// service.go can use a shorter named-step call.
func (s *EventService) buildOrgFilterCH(rc *reqCtx.RequestContext, query interface{ GetIncludeChildren() bool }) (map[string]interface{}, error) {
	return orgfilter.BuildOrgFilterClickHouse(orgfilter.BuildFilterParams{
		ReqContext: rc,
		Query:      query,
	})
}

// buildEventCursorOpts assembles ClickHouse cursor options from any events
// query DTO that implements cursorQuery.
func (s *EventService) buildEventCursorOpts(query cursorQuery) *chModel.TimeCursorOpts {
	opts := &chModel.TimeCursorOpts{
		Limit:     int64(query.GetLimit()),
		Direction: query.GetDirection(),
		SortAsc:   query.GetSortAsc(),
	}
	if cur := query.GetCursor(); cur != nil {
		opts.Cursor = *cur
	}
	return opts
}

// execTimeOpToStringPtr converts the typed exec-time operator (used by
// js-exec queries) into the string pointer expected by the repository.
func (s *EventService) execTimeOpToStringPtr(op *eventsContracts.ExecTimeOperator) *string {
	if op == nil {
		return nil
	}
	str := string(*op)
	return &str
}

// buildEventsRawCursorResult maps RawEvent entities to DTOs and wraps the
// result with the cursor metadata.
func (s *EventService) buildEventsRawCursorResult(result *chModel.TimeCursorResult[entities.RawEvent]) *dtos.EventsRawCursorResultDto {
	dtoItems := make([]dtos.EventsRawResponseDto, 0, len(result.Items))
	for _, entity := range result.Items {
		dto, err := mapper.EntityToDto[entities.RawEvent, dtos.EventsRawResponseDto](&entity)
		if err != nil {
			logger.Error(err, "[SERVICE:Event] Failed to map entity to DTO, skipping")
			continue
		}
		dtoItems = append(dtoItems, *dto)
	}
	response := &dtos.EventsRawCursorResultDto{
		Items:       dtoItems,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}
	attachCursors(&response.NextCursor, &response.PrevCursor, result)
	return response
}

// buildEventsJsExecCursorResult maps JsExecEvent entities to DTOs and wraps
// the result with cursor metadata.
func (s *EventService) buildEventsJsExecCursorResult(result *chModel.TimeCursorResult[entities.JsExecEvent]) *dtos.EventsJsExecCursorResultDto {
	dtoItems := make([]dtos.EventsJsExecResponseDto, 0, len(result.Items))
	for _, entity := range result.Items {
		dto, err := mapper.EntityToDto[entities.JsExecEvent, dtos.EventsJsExecResponseDto](&entity)
		if err != nil {
			logger.Error(err, "[SERVICE:Event] Failed to map JS exec entity to DTO, skipping")
			continue
		}
		dtoItems = append(dtoItems, *dto)
	}
	response := &dtos.EventsJsExecCursorResultDto{
		Items:       dtoItems,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}
	attachCursors(&response.NextCursor, &response.PrevCursor, result)
	return response
}

// buildEventsDLQCursorResult maps DLQEvent entities to DTOs and wraps the
// result with cursor metadata.
func (s *EventService) buildEventsDLQCursorResult(result *chModel.TimeCursorResult[entities.DLQEvent]) *dtos.EventsDLQCursorResultDto {
	dtoItems := make([]dtos.EventsDLQResponseDto, 0, len(result.Items))
	for _, entity := range result.Items {
		dto, err := mapper.EntityToDto[entities.DLQEvent, dtos.EventsDLQResponseDto](&entity)
		if err != nil {
			logger.Error(err, "[SERVICE:Event] Failed to map DLQ entity to DTO, skipping")
			continue
		}
		dtoItems = append(dtoItems, *dto)
	}
	response := &dtos.EventsDLQCursorResultDto{
		Items:       dtoItems,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}
	attachCursors(&response.NextCursor, &response.PrevCursor, result)
	return response
}

// buildEventsDLQCountsResult sums the per-service-type counts and wraps them
// in the response DTO.
func (s *EventService) buildEventsDLQCountsResult(counts []entities.DLQServiceCount) *dtos.EventsDLQCountsResultDto {
	var total uint64
	dtoCounts := make([]dtos.EventsDLQServiceCountDto, 0, len(counts))
	for _, item := range counts {
		dtoCounts = append(dtoCounts, dtos.EventsDLQServiceCountDto{
			ServiceType: item.ServiceType,
			Count:       item.Count,
		})
		total += item.Count
	}
	return &dtos.EventsDLQCountsResultDto{Counts: dtoCounts, Total: total}
}

// buildEventsRouterCursorResult maps RouterEvent entities to DTOs (custom
// projection — Success int -> bool) and wraps the result with cursor metadata.
func (s *EventService) buildEventsRouterCursorResult(result *chModel.TimeCursorResult[entities.RouterEvent]) *dtos.EventsRouterCursorResultDto {
	dtoItems := make([]dtos.EventsRouterResponseDto, 0, len(result.Items))
	for _, entity := range result.Items {
		dtoItems = append(dtoItems, dtos.EventsRouterResponseDto{
			Created:        entity.Created,
			ThreadId:       entity.ThreadId,
			OrgId:          entity.OrgId,
			PathKey:        entity.PathKey,
			AssetId:        entity.AssetId,
			RouterId:       entity.RouterId,
			Name:           entity.Name,
			Description:    entity.Description,
			TotalRouters:   entity.TotalRouters,
			MatchedCount:   entity.MatchedCount,
			PublishedCount: entity.PublishedCount,
			Event:          entity.Event,
			Success:        entity.Success == 1,
			Error:          entity.Error,
			RetentionDays:  entity.RetentionDays,
		})
	}
	response := &dtos.EventsRouterCursorResultDto{
		Items:       dtoItems,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}
	attachCursors(&response.NextCursor, &response.PrevCursor, result)
	return response
}

// buildEventsBusinessRuleCursorResult maps BusinessRuleEvent entities to
// DTOs (custom projection — Matched int -> bool) and wraps the result with
// cursor metadata.
func (s *EventService) buildEventsBusinessRuleCursorResult(result *chModel.TimeCursorResult[entities.BusinessRuleEvent]) *dtos.EventsBusinessRuleCursorResultDto {
	dtoItems := make([]dtos.EventsBusinessRuleResponseDto, 0, len(result.Items))
	for _, entity := range result.Items {
		dtoItems = append(dtoItems, dtos.EventsBusinessRuleResponseDto{
			Created:                 entity.Created,
			ThreadId:                entity.ThreadId,
			OrgId:                   entity.OrgId,
			PathKey:                 entity.PathKey,
			RuleId:                  entity.RuleId,
			BusinessRuleId:          entity.BusinessRuleId,
			BusinessRuleName:        entity.BusinessRuleName,
			BusinessRuleDescription: entity.BusinessRuleDescription,
			Matched:                 entity.Matched == 1,
			DurationMs:              entity.DurationMs,
			ConditionsEvaluated:     int(entity.ConditionsEvaluated),
			ConditionsMatched:       int(entity.ConditionsMatched),
			GroupsEvaluated:         int(entity.GroupsEvaluated),
			MaxDepthReached:         int(entity.MaxDepthReached),
			FinalState:              entity.FinalState,
			StateChanges:            entity.StateChanges,
			EvaluationTree:          entity.EvaluationTree,
			ConditionLogs:           entity.ConditionLogs,
			ActionsToDispatch:       entity.ActionsToDispatch,
			RetentionDays:           entity.RetentionDays,
		})
	}
	response := &dtos.EventsBusinessRuleCursorResultDto{
		Items:       dtoItems,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}
	attachCursors(&response.NextCursor, &response.PrevCursor, result)
	return response
}

// buildEventsTriggerCursorResult maps TriggerEvent entities to DTOs (custom
// projection — Success int -> bool) and wraps the result with cursor metadata.
func (s *EventService) buildEventsTriggerCursorResult(result *chModel.TimeCursorResult[entities.TriggerEvent]) *dtos.EventsTriggerCursorResultDto {
	dtoItems := make([]dtos.EventsTriggerResponseDto, 0, len(result.Items))
	for _, entity := range result.Items {
		dtoItems = append(dtoItems, dtos.EventsTriggerResponseDto{
			Created:       entity.Created,
			OrgId:         entity.OrgId,
			PathKey:       entity.PathKey,
			TriggerId:     entity.TriggerId,
			TriggerName:   entity.TriggerName,
			TriggerType:   entity.TriggerType,
			Category:      entity.Category,
			Source:        entity.Source,
			Success:       entity.Success == 1,
			DurationMs:    entity.DurationMs,
			Error:         entity.Error,
			RequestData:   entity.RequestData,
			ResponseData:  entity.ResponseData,
			RetentionDays: entity.RetentionDays,
		})
	}
	response := &dtos.EventsTriggerCursorResultDto{
		Items:       dtoItems,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}
	attachCursors(&response.NextCursor, &response.PrevCursor, result)
	return response
}

// buildEventsWorkflowCursorResult maps WorkflowEvent entities to DTOs and
// converts the JSON-string columns into json.RawMessage, then wraps with
// cursor metadata.
func (s *EventService) buildEventsWorkflowCursorResult(result *chModel.TimeCursorResult[entities.WorkflowEvent]) *dtos.EventsWorkflowCursorResultDto {
	dtoItems := make([]dtos.EventsWorkflowResponseDto, 0, len(result.Items))
	for _, entity := range result.Items {
		dto, err := mapper.EntityToDto[entities.WorkflowEvent, dtos.EventsWorkflowResponseDto](&entity)
		if err != nil {
			continue
		}
		s.attachWorkflowJSONFields(dto, &entity)
		dtoItems = append(dtoItems, *dto)
	}
	response := &dtos.EventsWorkflowCursorResultDto{
		Items:       dtoItems,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}
	attachCursors(&response.NextCursor, &response.PrevCursor, result)
	return response
}

// buildWorkflowEventResponse projects a single WorkflowEvent entity (from
// the by-execution-id endpoint) into its response DTO, including the
// JSON-string -> json.RawMessage conversions.
func (s *EventService) buildWorkflowEventResponse(entity *entities.WorkflowEvent) (*dtos.EventsWorkflowResponseDto, error) {
	dto, err := mapper.EntityToDto[entities.WorkflowEvent, dtos.EventsWorkflowResponseDto](entity)
	if err != nil {
		return nil, fmt.Errorf("failed to map workflow event: %w", err)
	}
	s.attachWorkflowJSONFields(dto, entity)
	return dto, nil
}

// attachWorkflowJSONFields converts string JSON columns (ClickHouse stores
// them as strings) into json.RawMessage so the DTO carries valid JSON.
func (s *EventService) attachWorkflowJSONFields(dto *dtos.EventsWorkflowResponseDto, entity *entities.WorkflowEvent) {
	dto.ExecutionPath = toRawJSON(entity.ExecutionPath)
	dto.NodeOutputs = toRawJSON(entity.NodeOutputs)
	dto.ErrorInfo = toRawJSON(entity.ErrorInfo)
	dto.EventPayload = toRawJSON(entity.EventPayload)
	dto.State = toRawJSON(entity.State)
	dto.ExternalInputs = toRawJSON(entity.ExternalInputs)
}

// buildEventsStoreCursorResult projects EventStore entities (list view —
// no EVA fields) and wraps the result with cursor metadata.
func (s *EventService) buildEventsStoreCursorResult(result *chModel.TimeCursorResult[entities.Event]) *dtos.EventsStoreCursorResultDto {
	dtoItems := make([]dtos.EventsStoreResponseDto, 0, len(result.Items))
	for _, entity := range result.Items {
		dtoItems = append(dtoItems, dtos.EventsStoreResponseDto{
			Created:             entity.Created,
			EventTrackerId:      entity.EventTrackerId,
			ThreadId:            entity.ThreadId,
			AssetId:             entity.AssetId,
			AssetName:           entity.AssetName,
			AssetDescription:    entity.AssetDescription,
			TemplateName:        entity.TemplateName,
			TemplateDescription: entity.TemplateDescription,
			OrgId:               entity.OrgId,
			PathKey:             entity.PathKey,
			Source:              entity.Source,
			Payload:             entity.Payload,
			Metadata:            entity.Metadata,
		})
	}
	response := &dtos.EventsStoreCursorResultDto{
		Items:       dtoItems,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}
	attachCursors(&response.NextCursor, &response.PrevCursor, result)
	return response
}

// resolveEventStoreEvaFields builds the human-readable advancedSearch map
// for a single event store record by looking up its template's DynamicFields
// and mapping fieldId -> field name across all EVA value types.
func (s *EventService) resolveEventStoreEvaFields(c ctx.Context, entity *entities.Event) map[string]interface{} {
	advancedSearch := make(map[string]interface{})
	switch entity.Source {
	case "asset":
		if entity.AssetTemplateId == "" {
			return advancedSearch
		}
		template, err := s.deps.TemplateCache.GetTemplate(c, entity.AssetTemplateOrgId, entity.AssetTemplateId)
		if err != nil {
			logger.Warn(fmt.Sprintf("[SERVICE:Event] Failed to get template %s/%s for detail: %v",
				entity.AssetTemplateOrgId, entity.AssetTemplateId, err))
			return advancedSearch
		}
		fieldMap := make(map[uint16]string, len(template.DynamicFields))
		for _, f := range template.DynamicFields {
			fieldMap[f.FieldId] = f.Field
		}
		for fieldId, value := range entity.EvaNumber {
			if name, ok := fieldMap[fieldId]; ok {
				advancedSearch[name] = value
			}
		}
		for fieldId, value := range entity.EvaString {
			if name, ok := fieldMap[fieldId]; ok {
				advancedSearch[name] = value
			}
		}
		for fieldId, value := range entity.EvaBool {
			if name, ok := fieldMap[fieldId]; ok {
				advancedSearch[name] = value == 1
			}
		}
		for fieldId, value := range entity.EvaDate {
			if name, ok := fieldMap[fieldId]; ok {
				advancedSearch[name] = value
			}
		}
	case "rule":
		logger.Info(fmt.Sprintf("[SERVICE:Event] Rule source EVA resolution not yet implemented for event %s", entity.EventTrackerId))
	}
	return advancedSearch
}

// buildEventsStoreDetailResponse assembles the detail DTO from the entity +
// resolved EVA map. advancedSearch is only attached when non-empty so the
// JSON payload omits it for templates without dynamic fields.
func (s *EventService) buildEventsStoreDetailResponse(entity *entities.Event, advancedSearch map[string]interface{}) *dtos.EventsStoreDetailResponseDto {
	response := &dtos.EventsStoreDetailResponseDto{
		Created:             entity.Created,
		EventTrackerId:      entity.EventTrackerId,
		Source:              entity.Source,
		AssetId:             entity.AssetId,
		AssetName:           entity.AssetName,
		AssetDescription:    entity.AssetDescription,
		TemplateName:        entity.TemplateName,
		TemplateDescription: entity.TemplateDescription,
		OrgId:               entity.OrgId,
		PathKey:             entity.PathKey,
		Payload:             entity.Payload,
		Metadata:            entity.Metadata,
	}
	if len(advancedSearch) > 0 {
		response.AdvancedSearch = advancedSearch
	}
	return response
}

// attachCursors writes the next/prev cursor pointers when the result has
// non-zero values. Centralized so each cursor-result builder stays compact.
func attachCursors[T any](next, prev **time.Time, result *chModel.TimeCursorResult[T]) {
	if !result.NextCursor.IsZero() {
		c := result.NextCursor
		*next = &c
	}
	if !result.PrevCursor.IsZero() {
		c := result.PrevCursor
		*prev = &c
	}
}
