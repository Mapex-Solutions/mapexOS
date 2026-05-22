package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	triggerDtos "triggers/src/modules/triggers/application/dtos"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// runTriggerBatchPhase1 launches the parallel-execution phase: each message
// gets a goroutine bounded by the configured worker semaphore. Results are
// collected for the sequential Ack/Nack pass that runs after the flush.
func (s *EventService) runTriggerBatchPhase1(messages []*natsModel.Message) []messageResult {
	results := make([]messageResult, len(messages))
	sem := make(chan struct{}, s.workers)
	var wg sync.WaitGroup
	for i, msg := range messages {
		sem <- struct{}{}
		wg.Add(1)
		go func(idx int, m *natsModel.Message) {
			defer wg.Done()
			defer func() { <-sem }()
			results[idx] = s.processOneMessage(m)
		}(i, msg)
	}
	wg.Wait()
	return results
}

// flushBatchPublishes performs the single CorePublisher flush at the end of
// Phase 1 so all fire-and-forget publishes hit the wire in one TCP round.
// Shared between trigger and workflow batch flows.
func (s *EventService) flushBatchPublishes() {
	if err := s.deps.NatsBus.FlushConnection(); err != nil {
		logger.Error(err, "[SERVICE:Event] Batch flush failed")
	}
}

// runTriggerBatchPhase3Ack walks collected results and acks/nacks/rejects
// each message accordingly. Sequential because NATS protocol calls are
// fast enough that the parallelism budget belongs in Phase 1.
func (s *EventService) runTriggerBatchPhase3Ack(results []messageResult) {
	for _, r := range results {
		switch r.action {
		case "ack":
			r.msg.Ack()
			s.deps.Metrics.MessagesTotal.WithLabelValues("ack").Inc()
		case "nack":
			r.msg.Nack(r.nackErr)
			if r.isDLQ {
				s.deps.Metrics.MessagesTotal.WithLabelValues("dlq").Inc()
			} else {
				s.deps.Metrics.MessagesTotal.WithLabelValues("nack").Inc()
			}
		case "reject":
			r.msg.Reject(r.rejectReason)
			s.deps.Metrics.MessagesTotal.WithLabelValues("reject").Inc()
		}
	}
}

// parseTriggerExecuteEvent decodes the inbound JSON payload into the typed
// TriggerExecuteEvent used by the single-message orchestration. Wraps the
// JSON error so callers can surface a stable parse-failure message.
func (s *EventService) parseTriggerExecuteEvent(data []byte) (*triggerDtos.TriggerExecuteEvent, error) {
	var event triggerDtos.TriggerExecuteEvent
	if err := json.Unmarshal(data, &event); err != nil {
		logger.Error(err, "[SERVICE:Event] Failed to deserialize TriggerExecuteEvent")
		return nil, fmt.Errorf("failed to deserialize event: %w", err)
	}
	logger.Info(fmt.Sprintf("[SERVICE:Event] Processing trigger execution: triggerId=%s, executionId=%s",
		event.TriggerID, event.ExecutionID))
	return &event, nil
}

// fetchTriggerForExecution loads the trigger configuration and validates the
// enabled flag. Returns skip=true (with err=nil) when the trigger exists but
// is disabled — caller treats that as a successful Ack with no action taken.
func (s *EventService) fetchTriggerForExecution(event *triggerDtos.TriggerExecuteEvent) (*triggerDtos.TriggerResponse, bool, error) {
	ctx := context.Background()
	trigger, err := s.deps.TriggerService.GetTriggerById(ctx, &event.TriggerID)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to fetch trigger: %s", event.TriggerID))
		return nil, false, fmt.Errorf("failed to fetch trigger: %w", err)
	}
	if trigger.Enabled == nil || !*trigger.Enabled {
		status := "nil"
		if trigger.Enabled != nil {
			status = fmt.Sprintf("%t", *trigger.Enabled)
		}
		logger.Warn(fmt.Sprintf("[SERVICE:Event] Trigger is not active: %s (enabled: %s)", event.TriggerID, status))
		return nil, true, nil
	}
	return trigger, false, nil
}

// resolveTriggerConfig converts the trigger's typed config to a generic map
// and resolves any {{event.*}} placeholders against the inbound payload.
// Returns the resolved config + triggerType (used by the caller to pick an
// executor and label observability).
func (s *EventService) resolveTriggerConfig(trigger *triggerDtos.TriggerResponse, event *triggerDtos.TriggerExecuteEvent) (map[string]interface{}, string, error) {
	triggerName, triggerType, category := triggerLabels(trigger)
	logger.Info(fmt.Sprintf("[SERVICE:Event] Trigger found: name=%s, type=%s, category=%s", triggerName, triggerType, category))
	configMap, err := TriggerConfigToMap(trigger.Config)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to convert trigger config: %s", event.TriggerID))
		return nil, "", fmt.Errorf("failed to convert trigger config: %w", err)
	}
	resolved, err := ResolvePlaceholdersInMap(configMap, event.Payload)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to resolve placeholders for trigger: %s", event.TriggerID))
		return nil, "", fmt.Errorf("failed to resolve placeholders: %w", err)
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Event] Placeholders resolved for trigger: %s", event.TriggerID))
	return resolved, triggerType, nil
}

// executeResolvedTrigger picks the executor for the given trigger type and
// runs it against the resolved config. Used by the legacy single-message
// path; the batch path uses processOneMessage which carries its own metrics.
func (s *EventService) executeResolvedTrigger(triggerType string, resolvedConfig map[string]interface{}, event *triggerDtos.TriggerExecuteEvent) error {
	ctx := context.Background()
	executor, exists := s.deps.ExecutorRegistry.GetExecutor(triggerType)
	if !exists {
		logger.Error(nil, fmt.Sprintf("[SERVICE:Event] No executor found for trigger type: %s", triggerType))
		return fmt.Errorf("no executor found for trigger type: %s", triggerType)
	}
	logger.Info(fmt.Sprintf("[SERVICE:Event] Executing trigger: %s (type: %s)", event.TriggerID, triggerType))
	if err := executor.Execute(ctx, resolvedConfig); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Trigger execution failed: %s", event.TriggerID))
		return fmt.Errorf("trigger execution failed: %w", err)
	}
	logger.Info(fmt.Sprintf("[SERVICE:Event] Trigger executed successfully: triggerId=%s, executionId=%s",
		event.TriggerID, event.ExecutionID))
	return nil
}

// triggerLabels unwraps the optional Name/TriggerType/Category fields into
// safe defaults so log lines and metric labels stay consistent.
func triggerLabels(trigger *triggerDtos.TriggerResponse) (string, string, string) {
	name, ttype, category := "", "", ""
	if trigger.Name != nil {
		name = *trigger.Name
	}
	if trigger.TriggerType != nil {
		ttype = *trigger.TriggerType
	}
	if trigger.Category != nil {
		category = *trigger.Category
	}
	return name, ttype, category
}

// processOneMessage handles a single message in the parallel Phase 1.
// Thread-safe — every dependency is safe for concurrent use. Returns a
// messageResult that the sequential Phase 3 acts on.
func (s *EventService) processOneMessage(msg *natsModel.Message) messageResult {
	ctx := context.Background()
	var event triggerDtos.TriggerExecuteEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Error(err, "[SERVICE:Event] Failed to deserialize TriggerExecuteEvent")
		s.deps.Metrics.TriggersProcessed.WithLabelValues("error").Inc()
		return messageResult{msg: msg, action: "reject", rejectReason: fmt.Sprintf("invalid JSON: %s", err.Error())}
	}
	msg.OrgId = event.OrgID
	msg.PathKey = event.PathKey
	msg.EventTrackerId = event.EventTrackerId
	logger.Info(fmt.Sprintf("[SERVICE:Event] Processing trigger: triggerId=%s, executionId=%s, source=%s, eventTrackerId=%s",
		event.TriggerID, event.ExecutionID, event.Source, event.EventTrackerId))
	startTime := time.Now()
	var cacheMetrics common.CacheMetrics
	trigger, err := s.deps.TriggerService.GetTriggerById(ctx, &event.TriggerID, &cacheMetrics)
	if cacheMetrics.Hit {
		s.deps.Metrics.TriggerCacheTotal.WithLabelValues("hit").Inc()
	} else {
		s.deps.Metrics.TriggerCacheTotal.WithLabelValues("miss").Inc()
	}
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to fetch trigger: %s", event.TriggerID))
		s.deps.Metrics.TriggersProcessed.WithLabelValues("error").Inc()
		s.deps.Metrics.TriggerProcessingDuration.Observe(time.Since(startTime).Seconds())
		isDLQ := false
		if attempt, max, _ := msg.GetRetryInfo(); max > 0 && attempt > max {
			isDLQ = true
		}
		return messageResult{msg: msg, action: "nack", nackErr: fmt.Errorf("failed to fetch trigger: %w", err), isDLQ: isDLQ}
	}
	if trigger.Enabled == nil || !*trigger.Enabled {
		status := "nil"
		if trigger.Enabled != nil {
			status = fmt.Sprintf("%t", *trigger.Enabled)
		}
		logger.Warn(fmt.Sprintf("[SERVICE:Event] Trigger is not active: %s (enabled: %s)", event.TriggerID, status))
		s.deps.Metrics.TriggersProcessed.WithLabelValues("disabled").Inc()
		s.deps.Metrics.TriggerProcessingDuration.Observe(time.Since(startTime).Seconds())
		return messageResult{msg: msg, action: "ack"}
	}
	triggerName, triggerType, category := triggerLabels(trigger)
	logger.Info(fmt.Sprintf("[SERVICE:Event] Trigger found: name=%s, type=%s, category=%s",
		triggerName, triggerType, category))
	configMap, err := TriggerConfigToMap(trigger.Config)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to convert trigger config: %s", event.TriggerID))
		s.publishTriggerEvent(ctx, &event, triggerName, triggerType, category, startTime, false, err.Error(), nil)
		s.deps.Metrics.TriggersProcessed.WithLabelValues("error").Inc()
		s.deps.Metrics.PlaceholderResolutions.WithLabelValues("error").Inc()
		s.deps.Metrics.TriggerProcessingDuration.Observe(time.Since(startTime).Seconds())
		return messageResult{msg: msg, action: "reject", rejectReason: fmt.Sprintf("failed to convert trigger config: %s", err.Error())}
	}
	resolvedConfig, err := ResolvePlaceholdersInMap(configMap, event.Payload)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to resolve placeholders for trigger: %s", event.TriggerID))
		s.publishTriggerEvent(ctx, &event, triggerName, triggerType, category, startTime, false, err.Error(), nil)
		s.deps.Metrics.TriggersProcessed.WithLabelValues("error").Inc()
		s.deps.Metrics.PlaceholderResolutions.WithLabelValues("error").Inc()
		s.deps.Metrics.TriggerProcessingDuration.Observe(time.Since(startTime).Seconds())
		return messageResult{msg: msg, action: "reject", rejectReason: fmt.Sprintf("failed to resolve placeholders: %s", err.Error())}
	}
	s.deps.Metrics.PlaceholderResolutions.WithLabelValues("success").Inc()
	logger.Debug(fmt.Sprintf("[SERVICE:Event] Placeholders resolved for trigger: %s", event.TriggerID))
	executor, exists := s.deps.ExecutorRegistry.GetExecutor(triggerType)
	if !exists {
		errMsg := fmt.Sprintf("no executor found for trigger type: %s", triggerType)
		logger.Error(nil, fmt.Sprintf("[SERVICE:Event] %s", errMsg))
		s.publishTriggerEvent(ctx, &event, triggerName, triggerType, category, startTime, false, errMsg, resolvedConfig)
		s.deps.Metrics.TriggersProcessed.WithLabelValues("no_executor").Inc()
		s.deps.Metrics.TriggerProcessingDuration.Observe(time.Since(startTime).Seconds())
		return messageResult{msg: msg, action: "reject", rejectReason: errMsg}
	}
	logger.Info(fmt.Sprintf("[SERVICE:Event] Executing trigger: %s (type: %s)", event.TriggerID, triggerType))
	execStart := time.Now()
	if err := executor.Execute(ctx, resolvedConfig); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Trigger execution failed: %s", event.TriggerID))
		s.deps.Metrics.ExecutorDuration.WithLabelValues(triggerType).Observe(time.Since(execStart).Seconds())
		s.deps.Metrics.ExecutorTotal.WithLabelValues(triggerType, "error").Inc()
		s.publishTriggerEvent(ctx, &event, triggerName, triggerType, category, startTime, false, err.Error(), resolvedConfig)
		s.deps.Metrics.TriggersProcessed.WithLabelValues("error").Inc()
		s.deps.Metrics.TriggerProcessingDuration.Observe(time.Since(startTime).Seconds())
		isDLQ := false
		if attempt, max, _ := msg.GetRetryInfo(); max > 0 && attempt > max {
			isDLQ = true
		}
		return messageResult{msg: msg, action: "nack", nackErr: fmt.Errorf("trigger execution failed: %w", err), isDLQ: isDLQ}
	}
	s.deps.Metrics.ExecutorDuration.WithLabelValues(triggerType).Observe(time.Since(execStart).Seconds())
	s.deps.Metrics.ExecutorTotal.WithLabelValues(triggerType, "success").Inc()
	s.publishTriggerEvent(ctx, &event, triggerName, triggerType, category, startTime, true, "", resolvedConfig)
	s.deps.Metrics.TriggersProcessed.WithLabelValues("success").Inc()
	s.deps.Metrics.TriggerProcessingDuration.Observe(time.Since(startTime).Seconds())
	logger.Info(fmt.Sprintf("[SERVICE:Event] Trigger executed successfully: triggerId=%s, executionId=%s",
		event.TriggerID, event.ExecutionID))
	return messageResult{msg: msg, action: "ack"}
}
