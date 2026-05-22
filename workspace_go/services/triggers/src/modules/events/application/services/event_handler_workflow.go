package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// runWorkflowBatchPhase1 launches the parallel-execution phase for the
// workflow batch. Each message gets its own goroutine; results are
// collected for the sequential Ack/Nack pass.
func (s *EventService) runWorkflowBatchPhase1(messages []*natsModel.Message) []messageResult {
	results := make([]messageResult, len(messages))
	var wg sync.WaitGroup
	for i, msg := range messages {
		wg.Add(1)
		go func(idx int, m *natsModel.Message) {
			defer wg.Done()
			results[idx] = s.processWorkflowMessage(m)
		}(i, msg)
	}
	wg.Wait()
	return results
}

// runWorkflowBatchPhase3Ack walks collected workflow results and emits the
// matching Ack/Nack/Reject + workflow-specific metric counter.
func (s *EventService) runWorkflowBatchPhase3Ack(results []messageResult) {
	for _, r := range results {
		switch r.action {
		case "ack":
			r.msg.Ack()
			s.deps.Metrics.WorkflowMessagesTotal.WithLabelValues("ack").Inc()
		case "nack":
			r.msg.Nack(r.nackErr)
			s.deps.Metrics.WorkflowMessagesTotal.WithLabelValues("nack").Inc()
		case "reject":
			r.msg.Reject(r.rejectReason)
			s.deps.Metrics.WorkflowMessagesTotal.WithLabelValues("reject").Inc()
		}
	}
}

// processWorkflowMessage handles a single workflow execution request.
// Deserializes the WorkflowTriggerRequest, routes by mode (plugin or
// trigger), executes, and publishes the resume callback to the workflow
// service. Thread-safe — used by the parallel Phase 1.
func (s *EventService) processWorkflowMessage(msg *natsModel.Message) messageResult {
	startTime := time.Now()
	var req struct {
		Mode            string                 `json:"mode"`
		OrgID           string                 `json:"orgId"`
		PathKey         string                 `json:"pathKey"`
		WorkflowID      string                 `json:"workflowId"`
		InstanceID      string                 `json:"instanceId"`
		NodeID          string                 `json:"nodeId"`
		CallbackSubject string                 `json:"callbackSubject"`
		ExecutionToken  string                 `json:"executionToken,omitempty"`
		Data            map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		logger.Error(err, "[SERVICE:Event] Failed to deserialize WorkflowTriggerRequest")
		s.deps.Metrics.WorkflowProcessed.WithLabelValues("unknown", "error").Inc()
		return messageResult{msg: msg, action: "reject", rejectReason: fmt.Sprintf("invalid JSON: %s", err.Error())}
	}
	msg.OrgId = req.OrgID
	msg.PathKey = req.PathKey
	logger.Info(fmt.Sprintf("[SERVICE:Event] Processing workflow request: mode=%s, nodeId=%s, instanceId=%s",
		req.Mode, req.NodeID, req.InstanceID))
	var output interface{}
	var execErr error
	execStart := time.Now()
	switch req.Mode {
	case "plugin":
		_, output, execErr = s.executePluginAction(req.Data)
		actionType := "http"
		if action, _ := req.Data["action"].(map[string]interface{}); action != nil {
			if t, _ := action["type"].(string); t != "" {
				actionType = t
			}
		}
		s.deps.Metrics.WorkflowExecutorDuration.WithLabelValues(actionType).Observe(time.Since(execStart).Seconds())
		if execErr != nil {
			s.deps.Metrics.WorkflowExecutorTotal.WithLabelValues(actionType, "error").Inc()
		} else {
			s.deps.Metrics.WorkflowExecutorTotal.WithLabelValues(actionType, "success").Inc()
		}
	case "trigger":
		_, output, execErr = s.executeTriggerEntity(req.Data)
	default:
		s.deps.Metrics.WorkflowProcessed.WithLabelValues(req.Mode, "error").Inc()
		return messageResult{msg: msg, action: "reject", rejectReason: fmt.Sprintf("unknown mode: %s", req.Mode)}
	}
	resume := map[string]interface{}{
		"instanceId": req.InstanceID,
		"nodeId":     req.NodeID,
		"status":     "success",
	}
	if req.ExecutionToken != "" {
		resume["executionToken"] = req.ExecutionToken
	}
	if execErr != nil {
		logger.Debug(fmt.Sprintf("[SERVICE:Event] Workflow plugin execution failed: mode=%s, nodeId=%s, error=%s",
			req.Mode, req.NodeID, execErr.Error()))
		resume["status"] = "error"
		resume["error"] = map[string]interface{}{
			"code":    "PLUGIN_EXECUTION_ERROR",
			"message": execErr.Error(),
			"nodeId":  req.NodeID,
		}
		s.deps.Metrics.WorkflowProcessed.WithLabelValues(req.Mode, "error").Inc()
	} else {
		logger.Debug(fmt.Sprintf("[SERVICE:Event] Workflow plugin execution success: mode=%s, nodeId=%s",
			req.Mode, req.NodeID))
		if output != nil {
			resume["output"] = output
		}
		s.deps.Metrics.WorkflowProcessed.WithLabelValues(req.Mode, "success").Inc()
	}
	s.deps.Metrics.WorkflowProcessingDuration.WithLabelValues(req.Mode).Observe(time.Since(startTime).Seconds())
	if req.CallbackSubject != "" {
		logger.Debug(fmt.Sprintf("[SERVICE:Event] Publishing resume callback: subject=%s, status=%s",
			req.CallbackSubject, resume["status"]))
		if err := s.deps.NatsBus.PublishCore(natsModel.PublishCoreConfig{
			Subject: req.CallbackSubject,
			Data:    resume,
		}); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to publish resume callback for %s", req.InstanceID))
			s.deps.Metrics.WorkflowResumePublished.WithLabelValues("error").Inc()
			return messageResult{msg: msg, action: "nack", nackErr: err}
		}
		s.deps.Metrics.WorkflowResumePublished.WithLabelValues("ok").Inc()
	}
	return messageResult{msg: msg, action: "ack"}
}
