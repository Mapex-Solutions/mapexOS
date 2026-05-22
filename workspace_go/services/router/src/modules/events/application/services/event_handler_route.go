package services

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

	"router/src/modules/events/application/constants"
	eventTypes "router/src/modules/events/application/types"
	domainServices "router/src/modules/events/domain/services"

	routegroupPorts "router/src/modules/routegroups/application/ports"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	routerEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/router/events"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/google/uuid"
)

// processBatchPhase1 runs the parallel worker pool that unmarshals, validates,
// routes, and buffers PublishCore calls for each message. Collects per-message
// results for deferred ACK/Nack/Reject in Phase 3.
func (s *EventService) processBatchPhase1(messages []*natsModel.Message) []messageResult {
	workers := runtime.NumCPU() * 2
	if workers > len(messages) {
		workers = len(messages)
	}

	results := make([]messageResult, len(messages))
	work := make(chan int, len(messages))
	for i := range messages {
		work <- i
	}
	close(work)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range work {
				results[idx] = s.processRouteMessageResult(idx, messages[idx])
			}
		}()
	}
	wg.Wait()
	return results
}

// processBatchPhase2Flush flushes the buffered fire-and-forget PublishCore calls
// emitted during Phase 1. A single FlushConnection replaces per-message flushes.
func (s *EventService) processBatchPhase2Flush() {
	if err := s.deps.NatsBus.FlushConnection(); err != nil {
		logger.Error(err, "[SERVICE:Event] Batch flush failed")
	}
}

// processBatchPhase3Ack applies ACK/Nack/Reject to each message based on the
// Phase 1 result and records per-message metrics.
func (s *EventService) processBatchPhase3Ack(results []messageResult) {
	for _, r := range results {
		switch r.action {
		case "ack":
			r.msg.Ack()
			s.deps.Metrics.MessagesTotal.WithLabelValues("ack").Inc()
		case "nack":
			r.msg.Nack(r.nackErr)
			s.deps.Metrics.MessagesTotal.WithLabelValues("nack").Inc()
		case "reject":
			r.msg.Reject(r.rejectReason)
			s.deps.Metrics.MessagesTotal.WithLabelValues("reject").Inc()
		}
		s.deps.Metrics.EventsProcessed.WithLabelValues(r.status).Inc()
		s.deps.Metrics.EventProcessingDuration.Observe(r.duration)
	}
}

// parseLegacyEventPayload decodes the legacy V1 route execution payload into
// the 5 fields required by execute. Returns an error for invalid JSON or any
// missing required field (orgId, assetId, event). eventTrackerId is optional.
// eventSource defaults to EventSourceAssetEvent when absent.
func (s *EventService) parseLegacyEventPayload(data []byte) (orgId, assetId string, event map[string]interface{}, eventTrackerId, eventSource string, err error) {
	var eventData map[string]interface{}
	if err = json.Unmarshal(data, &eventData); err != nil {
		err = fmt.Errorf("failed to parse event JSON: %w", err)
		return
	}

	var ok bool
	orgId, ok = eventData["orgId"].(string)
	if !ok || orgId == "" {
		err = fmt.Errorf("missing required field 'orgId'")
		return
	}

	assetId, ok = eventData["assetId"].(string)
	if !ok || assetId == "" {
		err = fmt.Errorf("missing required field 'assetId'")
		return
	}

	event, ok = eventData["event"].(map[string]interface{})
	if !ok {
		err = fmt.Errorf("missing required field 'event'")
		return
	}

	eventTrackerId, _ = eventData["eventTrackerId"].(string)

	eventSource, _ = eventData["eventSource"].(string)
	if eventSource == "" {
		eventSource = constants.EventSourceAssetEvent
	}
	return
}

// processRouteMessageResult parses JSON payload, validates required fields (orgId, assetUUID, event),
// then triggers routing execution. Returns a messageResult for deferred ACK/Nack/Reject in Phase 3.
// Uses Reject for validation errors (DLQ), Nack for processing failures (retry), Ack on success.
func (s *EventService) processRouteMessageResult(idx int, msg *natsModel.Message) messageResult {
	start := time.Now()

	var eventData map[string]interface{}
	if err := json.Unmarshal(msg.Data, &eventData); err != nil {
		return messageResult{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("invalid JSON: %s", err.Error()),
			status:       "error", duration: time.Since(start).Seconds(),
		}
	}

	orgId, ok := eventData["orgId"].(string)
	if !ok || orgId == "" {
		return messageResult{
			msg: msg, action: "reject", rejectReason: "missing required field 'orgId'",
			status: "error", duration: time.Since(start).Seconds(),
		}
	}
	msg.OrgId = orgId

	if pathKey, ok := eventData["pathKey"].(string); ok {
		msg.PathKey = pathKey
	}
	if trackerId, ok := eventData["eventTrackerId"].(string); ok {
		msg.EventTrackerId = trackerId
	}

	assetUUID, ok := eventData["assetUUID"].(string)
	if !ok || assetUUID == "" {
		return messageResult{
			msg: msg, action: "reject", rejectReason: "missing required field 'assetUUID'",
			status: "error", duration: time.Since(start).Seconds(),
		}
	}

	event, ok := eventData["event"].(map[string]interface{})
	if !ok {
		return messageResult{
			msg: msg, action: "reject", rejectReason: "missing required field 'event'",
			status: "error", duration: time.Since(start).Seconds(),
		}
	}

	// Extract eventTrackerId for end-to-end event tracking (optional - empty string if not present)
	eventTrackerId, _ := eventData["eventTrackerId"].(string)

	// eventSource discriminates which route groups on the asset to use:
	// - "assetEvent" (default): asset.RouteGroupIds (regular IoT event routing)
	// - "healthStatus": asset.HealthMonitor.OfflineRouteGroupIds or OnlineRouteGroupIds
	//   depending on the nested event.eventType ("offline" | "online")
	eventSource, _ := eventData["eventSource"].(string)
	if eventSource == "" {
		eventSource = constants.EventSourceAssetEvent
	}

	logger.Debug(fmt.Sprintf("[SERVICE:Event] Route event received: orgId=%s, assetUUID=%s, eventTrackerId=%s, eventSource=%s",
		orgId, assetUUID, eventTrackerId, eventSource))

	if err := s.execute(ctx.Background(), orgId, assetUUID, event, eventTrackerId, eventSource); err != nil {
		return messageResult{
			msg: msg, action: "nack", nackErr: err,
			status: "error", duration: time.Since(start).Seconds(),
		}
	}

	return messageResult{
		msg: msg, action: "ack",
		status: "success", duration: time.Since(start).Seconds(),
	}
}

// execute processes routing logic for an asset event.
//
// eventSource discriminates which route groups on the asset are used:
//   - constants.EventSourceAssetEvent (default): asset.RouteGroupIds
//   - constants.EventSourceHealthStatus: asset.HealthMonitor.{Offline,Online}RouteGroupIds
//     selected by the nested StandardizedPayload event.eventType
//     ("offline" | "online").
//
// Source of truth is always the asset in cache — payload never carries route IDs.
func (s *EventService) execute(context ctx.Context, orgId string, assetUUID string, event map[string]interface{}, eventTrackerId string, eventSource string) error {
	asset, err := s.fetchAssetFromCache(context, orgId, assetUUID)
	if err != nil {
		return err
	}

	groupIds := resolveRouteGroupIds(asset, eventSource, event)

	for _, groupId := range groupIds {
		if err := s.processRouteGroup(context, groupId, assetUUID, event, asset, eventTrackerId, eventSource); err != nil {
			logger.Warn(fmt.Sprintf("[SERVICE:Event] RouteGroup %s failed: %v (eventSource=%s)", groupId, err, eventSource))
		}
	}

	return nil
}

// fetchAssetFromCache retrieves asset data from TieredCache (L0 → L1 → L2 → Fallback).
func (s *EventService) fetchAssetFromCache(context ctx.Context, orgId string, assetUUID string) (*assetsContract.AssetReadModel, error) {
	cacheKey := orgId + "/" + assetUUID

	cacheStart := time.Now()
	data, tier, err := s.deps.AssetCache.Get(context, cacheKey)
	cacheDuration := time.Since(cacheStart).Seconds()

	// Resolve tier label for metrics
	tierLabel := tierLabels[tier]
	if tierLabel == "" {
		tierLabel = fmt.Sprintf("unknown_%d", tier)
	}

	if err != nil {
		s.deps.Metrics.AssetCacheTotal.WithLabelValues(tierLabel).Inc()
		s.deps.Metrics.AssetCacheDuration.WithLabelValues(tierLabel).Observe(cacheDuration)
		return nil, fmt.Errorf("asset not found in cache: %s - %w", cacheKey, err)
	}

	s.deps.Metrics.AssetCacheTotal.WithLabelValues(tierLabel).Inc()
	s.deps.Metrics.AssetCacheDuration.WithLabelValues(tierLabel).Observe(cacheDuration)

	var asset assetsContract.AssetReadModel
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %w", err)
	}

	return &asset, nil
}

// processRouteGroup fetches RouteGroup by ID, iterates routers, and publishes history.
// eventSource is threaded through so processRouter can filter disallowed kinds
// when eventSource == EventSourceHealthStatus.
func (s *EventService) processRouteGroup(
	context ctx.Context,
	groupId string,
	assetUUID string,
	event map[string]interface{},
	asset *assetsContract.AssetReadModel,
	eventTrackerId string,
	eventSource string,
) error {
	routeGroup, err := s.deps.RouteGroupService.GetRouteGroupEntityById(context, &groupId)
	if err != nil {
		return fmt.Errorf("RouteGroup not found: %s - %w", groupId, err)
	}
	if routeGroup == nil {
		return nil
	}

	if len(routeGroup.Routers) == 0 {
		return nil
	}

	routerResults := make([]eventTypes.RouterResultRecord, 0, len(routeGroup.Routers))

	for routerIdx := range routeGroup.Routers {
		result := s.processRouter(context, routerIdx, &routeGroup.Routers[routerIdx], assetUUID, event, asset, eventTrackerId, eventSource)
		routerResults = append(routerResults, result)
	}

	s.publishRouteHistory(context, routeGroup, assetUUID, asset, routerResults, eventTrackerId)

	return nil
}

// processRouter evaluates match conditions and publishes enriched event if matched.
// When eventSource == EventSourceHealthStatus, router kinds outside
// HealthStatusAllowedRouterKinds ({trigger, workflow}) are silently skipped
// at debug level — this guards the healthStatus track against misconfigured
// route groups that still contain save_event / notification / lake_house routers.
func (s *EventService) processRouter(
	context ctx.Context,
	routerIdx int,
	router *routegroupPorts.Router,
	assetUUID string,
	event map[string]interface{},
	asset *assetsContract.AssetReadModel,
	eventTrackerId string,
	eventSource string,
) eventTypes.RouterResultRecord {
	result := eventTypes.RouterResultRecord{
		Kind:       router.Kind,
		Matched:    false,
		Published:  false,
		Conditions: []eventTypes.ConditionResultRecord{},
	}

	// Kind filter for healthStatus eventSource — skip disallowed kinds silently.
	if eventSource == constants.EventSourceHealthStatus && !isAllowedKindForHealthStatus(router.Kind) {
		logger.Debug(fmt.Sprintf("[SERVICE:Event] skipping router kind=%s for eventSource=healthStatus (not allowed)", router.Kind))
		return result
	}

	evalResult, _ := s.matchEvaluator.Evaluate(event, toMatchConfig(router.Match))

	result.Matched = evalResult.ShouldProcess
	result.Conditions = s.convertConditionResults(evalResult.Conditions)

	// Track match evaluation metrics
	if router.Match == nil {
		s.deps.Metrics.MatchEvaluationsTotal.WithLabelValues("no_config").Inc()
	} else if result.Matched {
		s.deps.Metrics.MatchEvaluationsTotal.WithLabelValues("matched").Inc()
	} else {
		s.deps.Metrics.MatchEvaluationsTotal.WithLabelValues("unmatched").Inc()
	}
	s.deps.Metrics.MatchRulesEvaluatedTotal.Add(float64(len(evalResult.Conditions)))

	if !result.Matched {
		return result
	}

	// Get subject based on router kind
	subject, err := s.getSubjectForRouter(router)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to get subject for kind: %s", router.Kind))
		return result
	}

	// Build enriched event — same generic format for ALL router kinds.
	// Kind-specific fields are added by addRouterData().
	publishData := s.buildEnrichedEvent(assetUUID, event, asset, router, eventTrackerId)

	// MsgId = {eventTrackerId}-{routerIdx} for JetStream dedup.
	// routerIdx differentiates when one event fans out to N routers in the same RouteGroup.
	var msgId string
	if eventTrackerId != "" {
		msgId = fmt.Sprintf("%s-%d", eventTrackerId, routerIdx)
	}

	logger.Debug(fmt.Sprintf("[SERVICE:Event] Publishing: kind=%s, subject=%s, assetUUID=%s, eventTrackerId=%s",
		router.Kind, subject, assetUUID, eventTrackerId))

	publishStart := time.Now()
	if err := s.deps.NatsBus.PublishCore(natsModel.PublishCoreConfig{
		Subject: subject,
		Data:    publishData,
		MsgId:   msgId,
	}); err != nil {
		s.deps.Metrics.PublishDuration.WithLabelValues(router.Kind).Observe(time.Since(publishStart).Seconds())
		s.deps.Metrics.EventsPublished.WithLabelValues(router.Kind, "error").Inc()
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Publish failed: %s", subject))
		return result
	}

	s.deps.Metrics.PublishDuration.WithLabelValues(router.Kind).Observe(time.Since(publishStart).Seconds())
	s.deps.Metrics.EventsPublished.WithLabelValues(router.Kind, "success").Inc()

	result.Published = true

	logger.Debug(fmt.Sprintf("[SERVICE:Event] Dispatched: kind=%s, subject=%s, assetUUID=%s, eventTrackerId=%s",
		router.Kind, subject, assetUUID, eventTrackerId))

	return result
}

// getSubjectForRouter returns the NATS subject for a router based on its kind.
// For other kinds, it uses the static mapping from constants.
func (s *EventService) getSubjectForRouter(router *routegroupPorts.Router) (string, error) {
	switch router.Kind {

	case "trigger":
		if router.Trigger == nil || router.Trigger.TriggerId == "" {
			return "", fmt.Errorf("trigger router missing triggerId")
		}
		return constants.GetTriggerSubject(router.Trigger.TriggerId), nil

	case "workflow":
		return routerEvents.SubjectWorkflowExecution, nil

	default:
		return constants.GetSubjectByKind(router.Kind)
	}
}

// convertConditionResults converts domain ConditionResult to application ConditionResultRecord.
func (s *EventService) convertConditionResults(domainResults []domainServices.ConditionResult) []eventTypes.ConditionResultRecord {
	appResults := make([]eventTypes.ConditionResultRecord, len(domainResults))
	for i, dr := range domainResults {
		appResults[i] = eventTypes.ConditionResultRecord{
			Field:    dr.Field,
			Operator: dr.Operator,
			Expected: dr.Expected,
			Actual:   dr.Actual,
			Passed:   dr.Passed,
		}
	}
	return appResults
}

// publishRouteHistory sends routing execution results to events service for UI.
func (s *EventService) publishRouteHistory(
	context ctx.Context,
	routeGroup *routegroupPorts.RouteGroup,
	assetUUID string,
	asset *assetsContract.AssetReadModel,
	routerResults []eventTypes.RouterResultRecord,
	eventTrackerId string,
) {
	routerId := routeGroup.ID.Hex()
	name := routeGroup.Name
	description := routeGroup.Description

	historyEvent := eventTypes.RouterHistoryEvent{
		Created:        time.Now().UTC(),
		EventTrackerId: eventTrackerId,
		ThreadId:       assetUUID,
		OrgId:          asset.OrgId,
		PathKey:        asset.PathKey,
		AssetId:        asset.ID,
		RouterId:       routerId,
		Name:           name,
		Description:    description,
		Routers:        routerResults,
	}

	// MsgId = {eventTrackerId}-history for JetStream dedup
	var historyMsgId string
	if eventTrackerId != "" {
		historyMsgId = fmt.Sprintf("%s-history", eventTrackerId)
	}

	if err := s.deps.NatsBus.PublishCore(natsModel.PublishCoreConfig{
		Subject: routerEvents.SubjectRouterHistory,
		Data:    historyEvent,
		MsgId:   historyMsgId,
	}); err != nil {
		logger.Error(err, "[SERVICE:Event] Failed to publish route history")
	}
}

// buildEnrichedEvent creates a generic enriched event with asset context and tracking metadata.
// This is the ONLY event builder — used for ALL router kinds.
// Kind-specific fields are added by addRouterData.
//
// Base fields included in every enriched event:
//   - assetId, assetUUID, assetTemplateId, assetTemplateOrgId: asset context
//   - event: the original event payload from the sensor/gateway
//   - orgId, pathKey: multi-tenant context
//   - eventTrackerId: end-to-end tracking UUID (from HTTP Gateway)
//   - executionId: unique UUID per router dispatch (for per-execution tracking)
//   - created: timestamp of this dispatch
func (s *EventService) buildEnrichedEvent(
	assetUUID string,
	event map[string]interface{},
	asset *assetsContract.AssetReadModel,
	router *routegroupPorts.Router,
	eventTrackerId string,
) map[string]interface{} {
	enrichedEvent := map[string]interface{}{
		"assetId":            asset.ID,
		"assetUUID":          asset.UUID,
		"assetTemplateId":    asset.AssetTemplateID,
		"assetTemplateOrgId": asset.AssetTemplateOrgID,
		"assetName":          asset.Name,
		"assetDescription":   asset.Description,
		"event":              event,
		"orgId":              asset.OrgId,
		"pathKey":            asset.PathKey,
		"eventTrackerId":     eventTrackerId,
		"executionId":        uuid.New().String(),
		"created":            time.Now().UTC(),
	}

	// Enrich with template name from TemplateCache (non-blocking — empty string on failure)
	templateName := ""
	templateDescription := ""
	if asset.AssetTemplateID != "" {
		if tmpl, err := s.deps.TemplateCache.GetTemplate(ctx.Background(), asset.AssetTemplateOrgID, asset.AssetTemplateID); err == nil {
			templateName = tmpl.Name
			templateDescription = tmpl.Description
		}
	}
	enrichedEvent["templateName"] = templateName
	enrichedEvent["templateDescription"] = templateDescription

	s.addRouterData(enrichedEvent, router, asset)
	return enrichedEvent
}

// addRouterData enriches the event with kind-specific fields.
// Each kind adds only the extra data its consumer needs beyond the generic base.
// The generic base (built by buildEnrichedEvent) already contains:
// assetId, assetUUID, assetTemplateId, assetTemplateOrgId, event, orgId, pathKey,
// eventTrackerId, executionId, created.
func (s *EventService) addRouterData(enrichedEvent map[string]interface{}, router *routegroupPorts.Router, asset *assetsContract.AssetReadModel) {
	switch router.Kind {
	case "lake_house":
		if router.LakeHouse != nil {
			enrichedEvent["lakeHouseId"] = router.LakeHouse.LakeHouseId
			if router.LakeHouse.Metadata != nil {
				enrichedEvent["metadata"] = router.LakeHouse.Metadata
			}
		}
	case "notification":
		if router.Notification != nil {
			enrichedEvent["notificationId"] = router.Notification.NotificationId
			if router.Notification.Metadata != nil {
				enrichedEvent["metadata"] = router.Notification.Metadata
			}
		}
	case "trigger":
		if router.Trigger != nil {
			enrichedEvent["triggerId"] = router.Trigger.TriggerId
			enrichedEvent["source"] = "router"
			enrichedEvent["payload"] = enrichedEvent["event"]
			if router.Trigger.Metadata != nil {
				enrichedEvent["metadata"] = router.Trigger.Metadata
			}
		}
	case "save_event":
		enrichedEvent["source"] = "asset"
		enrichedEvent["threadId"] = asset.UUID
		if router.SaveEvent != nil && router.SaveEvent.Metadata != nil {
			enrichedEvent["metadata"] = router.SaveEvent.Metadata
		}
	case "workflow":
		logger.Debug(fmt.Sprintf("[SERVICE:Event] addRouterData workflow: Workflow=%+v, isNil=%t", router.Workflow, router.Workflow == nil))
		if router.Workflow != nil {
			enrichedEvent["mode"] = router.Workflow.Mode
			enrichedEvent["data"] = router.Workflow.Data
			if router.Workflow.Metadata != nil {
				enrichedEvent["metadata"] = router.Workflow.Metadata
			}
		}
	default:
		// Unknown router kind — fail loud. Without this guard the caller still
		// publishes whatever subject getSubjectForRouter returned (or an error
		// from it), but with NO kind-specific enrichment, so downstream
		// consumers get an event missing fields they expect — which can look
		// like a duplicate / wrong-target action. Log + return un-enriched so
		// nothing silently slips through.
		logger.Error(fmt.Errorf("unknown router kind: %s", router.Kind), fmt.Sprintf("[SERVICE:Event] addRouterData: unknown router kind=%s for assetUUID=%s — event will be published without kind-specific enrichment", router.Kind, asset.UUID))
	}
}
