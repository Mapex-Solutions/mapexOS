package services

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"time"

	"events/src/modules/events/application/dtos"
	"events/src/modules/events/domain/entities"
	domainservices "events/src/modules/events/domain/services"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/validator"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// legacyEventFields holds the required+optional projection of the raw legacy
// event body used by ProcessEvent. assetId is required; the others are
// optional and default to empty strings when absent.
type legacyEventFields struct {
	assetId   string
	orgId     string
	pathKey   string
	eventType string
	source    string
}

// parseLegacyEvent unmarshals the raw NATS body and extracts the required
// and optional fields in one step. Returns wrapped errors so the public
// caller surfaces stable parse / missing-field messages.
func (s *EventService) parseLegacyEvent(data []byte) (map[string]interface{}, legacyEventFields, error) {
	var eventData map[string]interface{}
	if err := json.Unmarshal(data, &eventData); err != nil {
		logger.Error(err, "[SERVICE:Event] Failed to parse event JSON")
		return nil, legacyEventFields{}, fmt.Errorf("failed to parse event JSON: %w", err)
	}
	fields := legacyEventFields{}
	fields.assetId, _ = eventData["assetId"].(string)
	fields.orgId, _ = eventData["orgId"].(string)
	fields.pathKey, _ = eventData["pathKey"].(string)
	fields.eventType, _ = eventData["eventType"].(string)
	fields.source, _ = eventData["source"].(string)
	if fields.assetId == "" {
		logger.Error(nil, "[SERVICE:Event] Event missing required field 'assetId'")
		return nil, fields, fmt.Errorf("event missing required field 'assetId'")
	}
	return eventData, fields, nil
}

// marshalLegacyEventPayloads serializes the inner event body and the NATS
// headers as JSON strings stored verbatim in ClickHouse.
func (s *EventService) marshalLegacyEventPayloads(eventData map[string]interface{}, headers map[string][]string) (string, string, error) {
	eventPayload, ok := eventData["event"].(map[string]interface{})
	if !ok {
		eventPayload = eventData
	}
	payloadJSON, err := json.Marshal(eventPayload)
	if err != nil {
		logger.Error(err, "[SERVICE:Event] Failed to marshal event payload")
		return "", "", fmt.Errorf("failed to marshal event payload: %w", err)
	}
	metadataJSON, err := json.Marshal(headers)
	if err != nil {
		logger.Error(err, "[SERVICE:Event] Failed to marshal metadata")
		return "", "", fmt.Errorf("failed to marshal metadata: %w", err)
	}
	return string(payloadJSON), string(metadataJSON), nil
}

// resolveRetentionOrDefault fetches the retention-days for the (org, table)
// pair, falling back to the safe default of 1 day on lookup failure (logged).
func (s *EventService) resolveRetentionOrDefault(c ctx.Context, orgId, tableName string) uint16 {
	retentionDays, err := s.getRetentionDays(c, orgId, tableName)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Critical error fetching retention for org %s", orgId))
		return 1
	}
	return retentionDays
}

// buildLegacyEvent assembles the ClickHouse Event entity from the parsed
// fields and the marshalled JSON strings.
func (s *EventService) buildLegacyEvent(fields legacyEventFields, payloadJSON, metadataJSON string, retentionDays uint16) *entities.Event {
	return &entities.Event{
		Created:       time.Now().UTC(),
		AssetId:       fields.assetId,
		OrgId:         fields.orgId,
		PathKey:       fields.pathKey,
		EventType:     fields.eventType,
		Source:        fields.source,
		Payload:       payloadJSON,
		Metadata:      metadataJSON,
		RetentionDays: retentionDays,
	}
}

// persistLegacyEvent writes the entity to ClickHouse and emits the success
// log. Failures are logged with assetId for triage and bubbled to the caller.
func (s *EventService) persistLegacyEvent(c ctx.Context, event *entities.Event, fields legacyEventFields) error {
	if err := s.deps.EventRepo.Save(c, event); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to save event for assetId: %s", fields.assetId))
		return err
	}
	logger.Info(fmt.Sprintf("[SERVICE:Event] Successfully saved event for assetId: %s, type: %s", fields.assetId, fields.eventType))
	return nil
}

/* PER-MESSAGE PARSERS — used by the batch orchestrator */

// processRawEventMessage parses, validates, and maps a single raw event message.
func (s *EventService) processRawEventMessage(idx int, msg *natsModel.Message) messageResult[entities.RawEvent] {
	setTenantContext(msg)

	var dto dtos.RawEventDto
	if err := validator.UnmarshalAndValidate(msg.Data, &dto); err != nil {
		return messageResult[entities.RawEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("validation failed: %s", err.Error()),
		}
	}

	msg.OrgId = dto.OrgId
	msg.PathKey = dto.PathKey
	msg.EventTrackerId = dto.EventTrackerId

	event, err := mapper.DtoToEntity[dtos.RawEventDto, entities.RawEvent](&dto)
	if err != nil {
		return messageResult[entities.RawEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("mapping failed: %s", err.Error()),
		}
	}

	if event.Created.IsZero() {
		event.Created = time.Now().UTC()
	}
	event.Metadata = s.headersToMetadata(msg.Headers)

	retentionDays, _ := s.getRetentionDays(ctx.Background(), dto.OrgId, "eventsRaw")
	if retentionDays == 0 {
		retentionDays = 1
	}
	event.RetentionDays = retentionDays

	return messageResult[entities.RawEvent]{
		msg: msg, action: "pending", entity: event,
	}
}

// processJsExecEventMessage parses, validates, and maps a single JS executor event message.
func (s *EventService) processJsExecEventMessage(idx int, msg *natsModel.Message) messageResult[entities.JsExecEvent] {
	setTenantContext(msg)

	var incomingDto dtos.JsExecEventDto
	if err := validator.UnmarshalAndValidate(msg.Data, &incomingDto); err != nil {
		return messageResult[entities.JsExecEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("validation failed: %s", err.Error()),
		}
	}

	msg.OrgId = incomingDto.OrgId
	msg.PathKey = incomingDto.PathKey
	msg.EventTrackerId = incomingDto.EventTrackerId

	flatDto := incomingDto.ToFlatDTO()
	event, err := mapper.DtoToEntity[dtos.JsExecEventFlatDto, entities.JsExecEvent](flatDto)
	if err != nil {
		return messageResult[entities.JsExecEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("mapping failed: %s", err.Error()),
		}
	}

	if event.Created.IsZero() {
		event.Created = time.Now().UTC()
	}
	if flatDto.Event != nil {
		payloadJSON, err := json.Marshal(flatDto.Event)
		if err != nil {
			return messageResult[entities.JsExecEvent]{
				msg: msg, action: "reject",
				rejectReason: fmt.Sprintf("marshal failed: %s", err.Error()),
			}
		}
		event.Event = string(payloadJSON)
	}

	retentionDays, _ := s.getRetentionDays(ctx.Background(), flatDto.OrgId, "eventsJsExecutor")
	if retentionDays == 0 {
		retentionDays = 1
	}
	event.RetentionDays = retentionDays

	return messageResult[entities.JsExecEvent]{
		msg: msg, action: "pending", entity: event,
	}
}

// processDLQEventMessage parses and maps a single DLQ event message.
func (s *EventService) processDLQEventMessage(idx int, msg *natsModel.Message) messageResult[entities.DLQEvent] {
	var incomingDto dtos.DLQEventIncomingDto
	if err := json.Unmarshal(msg.Data, &incomingDto); err != nil {
		return messageResult[entities.DLQEvent]{
			msg: msg, action: "ack_skip",
		}
	}

	headersJSON := ""
	if incomingDto.OriginalHeaders != nil {
		if headerBytes, err := json.Marshal(incomingDto.OriginalHeaders); err == nil {
			headersJSON = string(headerBytes)
		}
	}

	originalDataStr := ""
	if incomingDto.OriginalData != nil {
		originalDataStr = string(incomingDto.OriginalData)
	}

	event := &entities.DLQEvent{
		Created:         incomingDto.SentToDLQAt,
		EventTrackerId:  incomingDto.EventTrackerId,
		ID:              incomingDto.ID,
		OrgId:           incomingDto.OrgId,
		PathKey:         incomingDto.PathKey,
		ServiceName:     incomingDto.ServiceName,
		ServiceType:     incomingDto.ServiceType,
		EventType:       incomingDto.EventType,
		OriginalSubject: incomingDto.OriginalSubject,
		OriginalStream:  incomingDto.OriginalStream,
		OriginalData:    originalDataStr,
		OriginalHeaders: headersJSON,
		LastError:       incomingDto.LastError,
		ErrorCount:      uint32(incomingDto.ErrorCount),
		FirstDelivery:   incomingDto.FirstDelivery,
		LastDelivery:    incomingDto.LastDelivery,
		TotalDeliveries: uint32(incomingDto.TotalDeliveries),
		ConsumerName:    incomingDto.ConsumerName,
		RetentionDays:   30,
	}

	return messageResult[entities.DLQEvent]{
		msg: msg, action: "pending", entity: event,
	}
}

// processRouterEventMessage parses, validates, and maps a single router event message.
func (s *EventService) processRouterEventMessage(idx int, msg *natsModel.Message) messageResult[entities.RouterEvent] {
	setTenantContext(msg)

	var incomingDto dtos.RouterEventIncomingDto
	if err := json.Unmarshal(msg.Data, &incomingDto); err != nil {
		return messageResult[entities.RouterEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("parse failed: %s", err.Error()),
		}
	}

	msg.OrgId = incomingDto.OrgId
	msg.PathKey = incomingDto.PathKey
	msg.EventTrackerId = incomingDto.EventTrackerId

	totalRouters := uint8(len(incomingDto.Routers))
	matchedCount := uint8(0)
	publishedCount := uint8(0)
	for _, router := range incomingDto.Routers {
		if router.Matched {
			matchedCount++
		}
		if router.Published {
			publishedCount++
		}
	}

	routersJSON, err := json.Marshal(incomingDto.Routers)
	if err != nil {
		return messageResult[entities.RouterEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("marshal failed: %s", err.Error()),
		}
	}

	success := uint8(0)
	if publishedCount > 0 {
		success = 1
	}

	retentionDays, _ := s.getRetentionDays(ctx.Background(), incomingDto.OrgId, "eventsRouter")
	if retentionDays == 0 {
		retentionDays = 7
	}

	event := &entities.RouterEvent{
		Created:        incomingDto.Created,
		EventTrackerId: incomingDto.EventTrackerId,
		ThreadId:       incomingDto.ThreadId,
		OrgId:          incomingDto.OrgId,
		PathKey:        incomingDto.PathKey,
		AssetId:        incomingDto.AssetId,
		RouterId:       incomingDto.RouterId,
		Name:           incomingDto.Name,
		Description:    incomingDto.Description,
		TotalRouters:   totalRouters,
		MatchedCount:   matchedCount,
		PublishedCount: publishedCount,
		Event:          string(routersJSON),
		Success:        success,
		Error:          "",
		RetentionDays:  retentionDays,
	}

	return messageResult[entities.RouterEvent]{
		msg: msg, action: "pending", entity: event,
	}
}

// processBusinessRuleEventMessage parses, validates, and maps a single business rule event message.
func (s *EventService) processBusinessRuleEventMessage(idx int, msg *natsModel.Message) messageResult[entities.BusinessRuleEvent] {
	setTenantContext(msg)

	var incomingDto dtos.BusinessRuleEventIncomingDto
	if err := json.Unmarshal(msg.Data, &incomingDto); err != nil {
		return messageResult[entities.BusinessRuleEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("parse failed: %s", err.Error()),
		}
	}

	msg.OrgId = incomingDto.OrgId
	msg.PathKey = incomingDto.PathKey
	msg.EventTrackerId = incomingDto.EventTrackerId

	finalStateJSON := ""
	if incomingDto.FinalState != nil {
		if jsonBytes, err := json.Marshal(incomingDto.FinalState); err == nil {
			finalStateJSON = string(jsonBytes)
		}
	}
	stateChangesJSON := ""
	if incomingDto.StateChanges != nil {
		if jsonBytes, err := json.Marshal(incomingDto.StateChanges); err == nil {
			stateChangesJSON = string(jsonBytes)
		}
	}
	evaluationTreeJSON := ""
	if incomingDto.EvaluationTree != nil {
		if jsonBytes, err := json.Marshal(incomingDto.EvaluationTree); err == nil {
			evaluationTreeJSON = string(jsonBytes)
		}
	}
	conditionLogsJSON := ""
	if incomingDto.ConditionLogs != nil {
		if jsonBytes, err := json.Marshal(incomingDto.ConditionLogs); err == nil {
			conditionLogsJSON = string(jsonBytes)
		}
	}
	actionsJSON := ""
	if incomingDto.ActionsToDispatch != nil {
		if jsonBytes, err := json.Marshal(incomingDto.ActionsToDispatch); err == nil {
			actionsJSON = string(jsonBytes)
		}
	}

	matched := uint8(0)
	if incomingDto.Matched {
		matched = 1
	}

	retentionDays, _ := s.getRetentionDays(ctx.Background(), incomingDto.OrgId, "eventsBusinessRule")
	if retentionDays == 0 {
		retentionDays = 7
	}

	event := &entities.BusinessRuleEvent{
		Created:                 incomingDto.Created,
		EventTrackerId:          incomingDto.EventTrackerId,
		ThreadId:                incomingDto.ThreadId,
		OrgId:                   incomingDto.OrgId,
		PathKey:                 incomingDto.PathKey,
		RuleId:                  incomingDto.RuleId,
		BusinessRuleId:          incomingDto.BusinessRuleId,
		BusinessRuleName:        incomingDto.BusinessRuleName,
		BusinessRuleDescription: incomingDto.BusinessRuleDescription,
		Matched:                 matched,
		DurationMs:              incomingDto.DurationMs,
		ConditionsEvaluated:     uint16(incomingDto.ConditionsEvaluated),
		ConditionsMatched:       uint16(incomingDto.ConditionsMatched),
		GroupsEvaluated:         uint16(incomingDto.GroupsEvaluated),
		MaxDepthReached:         uint16(incomingDto.MaxDepthReached),
		FinalState:              finalStateJSON,
		StateChanges:            stateChangesJSON,
		EvaluationTree:          evaluationTreeJSON,
		ConditionLogs:           conditionLogsJSON,
		ActionsToDispatch:       actionsJSON,
		RetentionDays:           retentionDays,
	}

	return messageResult[entities.BusinessRuleEvent]{
		msg: msg, action: "pending", entity: event,
	}
}

// processTriggerEventMessage parses, validates, and maps a single trigger event message.
func (s *EventService) processTriggerEventMessage(idx int, msg *natsModel.Message) messageResult[entities.TriggerEvent] {
	setTenantContext(msg)

	var incomingDto dtos.TriggerEventIncomingDto
	if err := json.Unmarshal(msg.Data, &incomingDto); err != nil {
		return messageResult[entities.TriggerEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("parse failed: %s", err.Error()),
		}
	}

	msg.OrgId = incomingDto.OrgId
	msg.PathKey = incomingDto.PathKey
	msg.EventTrackerId = incomingDto.EventTrackerId

	success := uint8(0)
	if incomingDto.Success {
		success = 1
	}

	retentionDays, _ := s.getRetentionDays(ctx.Background(), incomingDto.OrgId, "eventsTrigger")
	if retentionDays == 0 {
		retentionDays = 7
	}

	event := &entities.TriggerEvent{
		Created:        incomingDto.Created,
		EventTrackerId: incomingDto.EventTrackerId,
		OrgId:          incomingDto.OrgId,
		PathKey:        incomingDto.PathKey,
		TriggerId:      incomingDto.TriggerId,
		TriggerName:    incomingDto.TriggerName,
		TriggerType:    incomingDto.TriggerType,
		Category:       incomingDto.Category,
		Source:         incomingDto.Source,
		Success:        success,
		DurationMs:     incomingDto.DurationMs,
		Error:          incomingDto.Error,
		RequestData:    incomingDto.RequestData,
		ResponseData:   incomingDto.ResponseData,
		RetentionDays:  retentionDays,
	}

	return messageResult[entities.TriggerEvent]{
		msg: msg, action: "pending", entity: event,
	}
}

// processWorkflowEventMessage parses and maps a single workflow event message.
func (s *EventService) processWorkflowEventMessage(idx int, msg *natsModel.Message) messageResult[entities.WorkflowEvent] {
	setTenantContext(msg)

	var incomingDto dtos.WorkflowEventIncomingDto
	if err := json.Unmarshal(msg.Data, &incomingDto); err != nil {
		return messageResult[entities.WorkflowEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("parse failed: %s", err.Error()),
		}
	}

	msg.OrgId = incomingDto.OrgId
	msg.PathKey = incomingDto.PathKey
	msg.EventTrackerId = incomingDto.EventTrackerId

	var retentionDays uint16
	if incomingDto.RetentionDays > 0 {
		retentionDays = uint16(incomingDto.RetentionDays)
	} else if incomingDto.RetentionDays == 0 {
		retentionDays = 0
	} else {
		orgRetention, _ := s.getRetentionDays(ctx.Background(), incomingDto.OrgId, "eventsWorkflow")
		if orgRetention > 0 {
			retentionDays = orgRetention
		} else {
			retentionDays = 7
		}
	}

	event := &entities.WorkflowEvent{
		Created:           incomingDto.Created,
		Finished:          incomingDto.Finished,
		ExecutionId:       incomingDto.ExecutionId,
		EventTrackerId:    incomingDto.EventTrackerId,
		OrgId:             incomingDto.OrgId,
		PathKey:           incomingDto.PathKey,
		WorkflowUUID:      incomingDto.WorkflowUUID,
		InstanceId:        incomingDto.InstanceId,
		DefinitionId:      incomingDto.DefinitionId,
		WorkflowName:      incomingDto.WorkflowName,
		InstanceName:      incomingDto.InstanceName,
		DefinitionName:    incomingDto.DefinitionName,
		Status:            incomingDto.Status,
		Success:           incomingDto.Success,
		DurationMs:        incomingDto.DurationMs,
		ErrorMessage:      incomingDto.ErrorMessage,
		ExecutionPath:     incomingDto.ExecutionPath,
		NodeOutputs:       incomingDto.NodeOutputs,
		ErrorInfo:         incomingDto.ErrorInfo,
		EventPayload:      incomingDto.EventPayload,
		TriggerSource:     incomingDto.TriggerSource,
		ParentExecutionId: incomingDto.ParentExecutionId,
		Depth:             incomingDto.Depth,
		RetentionDays:     retentionDays,
		State:             incomingDto.State,
		ExternalInputs:    incomingDto.ExternalInputs,
	}

	return messageResult[entities.WorkflowEvent]{
		msg: msg, action: "pending", entity: event,
	}
}

// processEventStoreMessage parses, validates, maps, and resolves EVA fields
// for a single event store message.
func (s *EventService) processEventStoreMessage(idx int, msg *natsModel.Message) messageResult[entities.Event] {
	setTenantContext(msg)

	var incomingDto dtos.EventStoreDto
	if err := json.Unmarshal(msg.Data, &incomingDto); err != nil {
		return messageResult[entities.Event]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("parse failed: %s", err.Error()),
		}
	}

	if incomingDto.AssetId == "" {
		return messageResult[entities.Event]{
			msg: msg, action: "reject",
			rejectReason: "missing required field 'assetId'",
		}
	}
	if incomingDto.OrgId == "" {
		return messageResult[entities.Event]{
			msg: msg, action: "reject",
			rejectReason: "missing required field 'orgId'",
		}
	}

	msg.OrgId = incomingDto.OrgId
	msg.PathKey = incomingDto.PathKey
	msg.EventTrackerId = incomingDto.EventTrackerId

	event := entities.NewEvent()
	event.Created = incomingDto.Created
	if event.Created.IsZero() {
		event.Created = time.Now().UTC()
	}
	event.EventTrackerId = incomingDto.EventTrackerId
	event.ThreadId = incomingDto.ThreadId
	event.AssetId = incomingDto.AssetId
	event.AssetTemplateId = incomingDto.AssetTemplateId
	event.AssetTemplateOrgId = incomingDto.AssetTemplateOrgId
	event.AssetName = incomingDto.AssetName
	event.AssetDescription = incomingDto.AssetDescription
	event.TemplateName = incomingDto.TemplateName
	event.TemplateDescription = incomingDto.TemplateDescription
	event.OrgId = incomingDto.OrgId
	event.PathKey = incomingDto.PathKey
	event.Source = incomingDto.Source

	if incomingDto.Event != nil {
		payloadJSON, err := json.Marshal(incomingDto.Event)
		if err == nil {
			event.Payload = string(payloadJSON)
		}
	}
	if incomingDto.Metadata != nil {
		metadataJSON, err := json.Marshal(incomingDto.Metadata)
		if err == nil {
			event.Metadata = string(metadataJSON)
		}
	}

	c := ctx.Background()
	if incomingDto.AssetTemplateId != "" && incomingDto.Event != nil {
		template, err := s.deps.TemplateCache.GetTemplate(c, incomingDto.AssetTemplateOrgId, incomingDto.AssetTemplateId)
		if err != nil {
			logger.Warn(fmt.Sprintf("[SERVICE:Event] Failed to get template %s/%s for EVA mapping: %v", incomingDto.AssetTemplateOrgId, incomingDto.AssetTemplateId, err))
			s.deps.Metrics.TemplateCacheTotal.WithLabelValues("error").Inc()
		} else {
			s.deps.Metrics.TemplateCacheTotal.WithLabelValues("hit").Inc()
			domainservices.MapEvaFields(event, template.DynamicFields, incomingDto.Event)
			s.deps.Metrics.EvaFieldsMapped.Add(float64(len(template.DynamicFields)))
		}
	}

	retentionDays, _ := s.getRetentionDays(c, incomingDto.OrgId, "events")
	if retentionDays == 0 {
		retentionDays = 30
	}
	event.RetentionDays = retentionDays

	return messageResult[entities.Event]{
		msg: msg, action: "pending", entity: event,
	}
}
