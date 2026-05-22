package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"workflow/src/modules/runtime/application/constants"
	"workflow/src/modules/runtime/application/ports"
	domainConstants "workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"
	domainServices "workflow/src/modules/runtime/domain/services"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// parseResumeMessage decodes the WORKFLOW-RESUME payload. A bad payload is
// permanently rejected — there is no useful retry for malformed JSON.
func (s *RuntimeService) parseResumeMessage(msg *natsModel.Message) (ports.ResumeMessage, bool) {
	var resume ports.ResumeMessage
	if err := json.Unmarshal(msg.Data, &resume); err != nil {
		msg.Reject(fmt.Sprintf("invalid resume message: %s", err))
		return resume, false
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] Resume received: instanceId=%s, nodeId=%s, status=%s, outputHandle=%s, isTimeout=%t, hasOutput=%t, hasError=%t",
		resume.InstanceID, resume.NodeID, resume.Status, resume.OutputHandle, resume.IsTimeout, resume.Output != nil, resume.Error != nil))
	return resume, true
}

// loadResumeExecution fetches the current KV state for the addressed
// execution, sets DLQ context for multi-tenant routing, and validates the
// status is resumable. Returns (nil, 0, false) and acks the message when the
// execution has already terminated or vanished — both are no-op outcomes for
// a resume signal, not error conditions.
func (s *RuntimeService) loadResumeExecution(msg *natsModel.Message, resume *ports.ResumeMessage) (*entities.WorkflowExecution, uint64, bool) {
	execution, revision, err := s.deps.ExecutionStateRepo.GetWithRevision(resume.InstanceID)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Runtime] Resume skipped for %s: execution already finished or cancelled, no longer in KV — resume event discarded", resume.InstanceID))
		msg.Ack()
		return nil, 0, false
	}
	if execution.OrgID != nil {
		msg.OrgId = execution.OrgID.Hex()
	}
	msg.PathKey = execution.PathKey

	if execution.Status != entities.ExecStatusWaiting && execution.Status != entities.ExecStatusRunning {
		logger.Warn(fmt.Sprintf("[SERVICE:Runtime] Resume skipped for %s: execution already ended with status %q — cannot resume a finished workflow", resume.InstanceID, execution.Status))
		msg.Ack()
		return nil, 0, false
	}
	return execution, revision, true
}

// validateResumeToken enforces G5: discard stale/duplicate callbacks whose
// execution token does not match the one stored on the waiting node. Returns
// false (and acks the msg) when the callback is rejected as stale.
func (s *RuntimeService) validateResumeToken(msg *natsModel.Message, execution *entities.WorkflowExecution, resume *ports.ResumeMessage) bool {
	if resume.ExecutionToken == "" || resume.NodeID == "" {
		return true
	}
	ns := execution.NodeStates[resume.NodeID]
	if ns == nil {
		return true
	}
	expectedToken, ok := ns[constants.NodeStateKeyExecutionToken].(string)
	if !ok || expectedToken == "" {
		return true
	}
	if resume.ExecutionToken == expectedToken {
		return true
	}
	logger.Warn(fmt.Sprintf("[SERVICE:Runtime] Stale callback rejected for %s node %s: token mismatch (expected %s, got %s)",
		resume.InstanceID, resume.NodeID, expectedToken, resume.ExecutionToken))
	s.deps.Metrics.TokenRejectionsTotal.Inc()
	msg.Ack()
	return false
}

// maybePurgeNonTimeoutSchedule cancels the pending NATS schedule when the
// resume arrived from a callback/signal (i.e. before the timer fired). Best-
// effort: a purge failure is logged but does not block the resume.
func (s *RuntimeService) maybePurgeNonTimeoutSchedule(resume *ports.ResumeMessage) {
	if resume.IsTimeout || resume.NodeID == "" {
		return
	}
	if err := s.deps.RuntimePublisher.PurgeSchedule(resume.InstanceID, resume.NodeID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Runtime] Failed to purge schedule for %s node %s: %s",
			resume.InstanceID, resume.NodeID, err))
	}
}

// tryEarlyResumeBranches runs the three short-circuit handlers in order
// (retry-timer, timeout-without-output, error). Returns true when any of
// them claims the message so HandleResume bails before the normal CAS
// resume path.
func (s *RuntimeService) tryEarlyResumeBranches(ctx context.Context, msg *natsModel.Message, execution *entities.WorkflowExecution, resume *ports.ResumeMessage) bool {
	if s.tryHandleRetryTimerResume(ctx, msg, execution, resume) {
		return true
	}
	if s.tryHandleTimeoutResume(msg, execution, resume) {
		return true
	}
	if s.tryHandleErrorResume(ctx, msg, execution, resume) {
		return true
	}
	return false
}

// tryHandleRetryTimerResume detects a retry-timer fire (IsTimeout + node wait
// type = retry timer) and re-runs the same node from scratch. Returns true
// when the resume was fully handled here (msg acked/nacked); false to let the
// next handler take it.
func (s *RuntimeService) tryHandleRetryTimerResume(ctx context.Context, msg *natsModel.Message, execution *entities.WorkflowExecution, resume *ports.ResumeMessage) bool {
	if !resume.IsTimeout {
		return false
	}
	ns := execution.NodeStates[resume.NodeID]
	if ns == nil || ns[constants.NodeStateKeyWaitType] != domainConstants.WaitTypeRetryTimer {
		return false
	}

	retryAttempt := readRetryAttempt(ns)
	logger.Info(fmt.Sprintf("[SERVICE:Runtime] Retry timer fired for %s node %s (attempt %d)",
		resume.InstanceID, resume.NodeID, retryAttempt))

	s.applyRetryTimerCheckpoint(execution, resume.NodeID, retryAttempt)
	if err := s.checkpoint(execution); err != nil {
		msg.Nack(err)
		return true
	}
	if err := s.publishResumedStateEvent(execution); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish resumed state event for retry %s", resume.InstanceID))
	}

	graph, err := s.loadGraphForResume(ctx, execution)
	if err != nil {
		msg.Nack(fmt.Errorf("failed to get definition for retry: %w", err))
		return true
	}
	if err := s.execute(ctx, execution, graph, resume.NodeID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Retry execution failed for %s", resume.InstanceID))
		msg.Nack(err)
		return true
	}
	msg.Ack()
	return true
}

// tryHandleTimeoutResume routes a NATS-Schedule timeout. Two branches:
//   - EnableOutput=true → mutate resume so the normal flow routes through the
//     "timeout" handle; return false so the orchestration continues.
//   - EnableOutput=false → mark the path entry as timed out and fail the
//     whole execution. Returns true (msg acked/nacked).
//
// Retry-timer fires never reach here — they were absorbed by the previous step.
func (s *RuntimeService) tryHandleTimeoutResume(msg *natsModel.Message, execution *entities.WorkflowExecution, resume *ports.ResumeMessage) bool {
	if !resume.IsTimeout {
		return false
	}
	if resume.EnableOutput {
		resume.OutputHandle = domainConstants.OutputHandleTimeout
		resume.Status = domainConstants.StatusCompleted
		logger.Info(fmt.Sprintf("[SERVICE:Runtime] Timeout with output for %s node %s — routing to timeout handle", resume.InstanceID, resume.NodeID))
		return false
	}

	timeoutMsg := fmt.Sprintf("TIMEOUT_EXCEEDED: node %s timed out", resume.NodeID)
	markPathStatusOnNode(execution, resume.NodeID, domainConstants.StatusTimeout, "", &timeoutMsg)
	if err := s.failExecution(execution, &entities.ExecutionError{
		Code:      domainConstants.ErrCodeTimeoutExceeded,
		Message:   fmt.Sprintf("node %s timed out waiting for callback/signal", resume.NodeID),
		NodeID:    resume.NodeID,
		Timestamp: time.Now(),
	}); err != nil {
		msg.Nack(err)
		return true
	}
	msg.Ack()
	return true
}

// tryHandleErrorResume reacts to a callback that reported failure. The
// dispatch consults the per-node error handler (retry / error-handle) before
// falling through to a hard failExecution. Returns true once a terminal
// decision has been made (msg acked/nacked).
func (s *RuntimeService) tryHandleErrorResume(ctx context.Context, msg *natsModel.Message, execution *entities.WorkflowExecution, resume *ports.ResumeMessage) bool {
	if resume.Error == nil && resume.Status != domainConstants.StatusError {
		return false
	}
	s.normaliseResumeError(execution, resume)
	logger.Warn(fmt.Sprintf("[SERVICE:Runtime] Resume error for %s: nodeId=%s, nodeType=%s, code=%s, message=%s",
		resume.InstanceID, resume.Error.NodeID, resume.Error.NodeType, resume.Error.Code, resume.Error.Message))

	errMsg := fmt.Sprintf("%s: %s", resume.Error.Code, resume.Error.Message)
	markPathStatusOnNode(execution, resume.NodeID, domainConstants.StatusError, domainConstants.OutputHandleError, &errMsg)

	if handled, handleErr := s.routeResumeErrorToHandler(ctx, execution, resume); handleErr != nil {
		msg.Nack(handleErr)
		return true
	} else if handled {
		msg.Ack()
		return true
	}

	if err := s.failExecution(execution, resume.Error); err != nil {
		msg.Nack(err)
		return true
	}
	msg.Ack()
	return true
}

// runNormalResume is the success path: apply patches/output to the execution,
// pick the next node from the graph (or pop a pending loop), and continue the
// DAG walker. The CAS-retry loop guards against concurrent fanout callbacks
// stomping on each other.
func (s *RuntimeService) runNormalResume(ctx context.Context, msg *natsModel.Message, execution *entities.WorkflowExecution, revision uint64, resume *ports.ResumeMessage) {
	graph, err := s.loadGraphForResume(ctx, execution)
	if err != nil {
		msg.Nack(fmt.Errorf("failed to get definition: %w", err))
		return
	}

	startNodeID, ok := s.resolveResumeStartNode(msg, execution, revision, resume, graph)
	if !ok {
		return
	}
	if startNodeID == "" {
		msg.Reject(fmt.Sprintf("execution %s has no active node to resume from", resume.InstanceID))
		return
	}
	if err := s.execute(ctx, execution, graph, startNodeID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Resume execution failed for execution %s", resume.InstanceID))
		msg.Nack(err)
		return
	}
	msg.Ack()
}

// resolveResumeStartNode picks the node where the walker continues. For an
// execution already in Running status it just re-enqueues the first active
// node. For a Waiting execution it CAS-applies the resume payload, resolves
// the next node by graph, and falls back to a popped loop frame when the
// current branch is exhausted. Returns ("", false) when the orchestration
// already finished by a side-effect (msg acked/nacked here) so the caller
// must not continue.
func (s *RuntimeService) resolveResumeStartNode(msg *natsModel.Message, execution *entities.WorkflowExecution, revision uint64, resume *ports.ResumeMessage, graph *entities.ExecutionGraph) (string, bool) {
	if execution.Status == entities.ExecStatusRunning {
		if len(execution.ActiveNodeIDs) > 0 {
			return execution.ActiveNodeIDs[0], true
		}
		return "", true
	}

	if !s.applyResumePatchesWithCAS(msg, execution, revision, resume) {
		return "", false
	}
	if err := s.publishResumedStateEvent(execution); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish resumed state event for execution %s", resume.InstanceID))
	}

	outputHandle := graph.ResolveDefaultHandle(resume.NodeID)
	nextNodes := graph.ResolveNextNodes(resume.NodeID, []string{outputHandle})
	if len(nextNodes) > 0 {
		return nextNodes[0], true
	}
	if loopNodeID := peekLoopStack(execution.NodeStates); loopNodeID != "" {
		popLoopStack(execution.NodeStates)
		return loopNodeID, true
	}
	if err := s.completeOrResuspend(execution); err != nil {
		msg.Nack(err)
		return "", false
	}
	msg.Ack()
	return "", false
}

// applyResumePatchesWithCAS runs the bounded retry loop that merges the
// callback's state patch + output into the execution and persists it under
// optimistic concurrency control. On exhausting retries it nacks the msg and
// returns false so the caller stops. When the execution has been removed or
// has changed status mid-flight (concurrent terminal resume), it acks and
// returns false — there is nothing left to advance.
func (s *RuntimeService) applyResumePatchesWithCAS(msg *natsModel.Message, execution *entities.WorkflowExecution, revision uint64, resume *ports.ResumeMessage) bool {
	currentExecution := execution
	currentRevision := revision
	for casAttempt := 0; casAttempt < constants.MaxCASRetries; casAttempt++ {
		if casAttempt > 0 {
			reloaded, reloadedRev, err := s.deps.ExecutionStateRepo.GetWithRevision(resume.InstanceID)
			if err != nil {
				logger.Warn(fmt.Sprintf("[SERVICE:Runtime] CAS retry %d: execution %s no longer in KV", casAttempt, resume.InstanceID))
				msg.Ack()
				return false
			}
			if reloaded.Status != entities.ExecStatusWaiting {
				logger.Warn(fmt.Sprintf("[SERVICE:Runtime] CAS retry %d: execution %s status changed to %q", casAttempt, resume.InstanceID, reloaded.Status))
				msg.Ack()
				return false
			}
			currentExecution = reloaded
			currentRevision = reloadedRev
			*execution = *reloaded
		}

		s.mergeResumePayloadIntoExecution(currentExecution, resume)
		currentExecution.Updated = time.Now()
		if err := s.deps.ExecutionStateRepo.SaveWithRevision(currentExecution, currentRevision); err != nil {
			logger.Debug(fmt.Sprintf("[SERVICE:Runtime] CAS conflict %d/%d for %s: %s", casAttempt+1, constants.MaxCASRetries, resume.InstanceID, err))
			s.deps.Metrics.CASRetriesTotal.Inc()
			if casAttempt == constants.MaxCASRetries-1 {
				msg.Nack(fmt.Errorf("[SERVICE:Runtime] CAS exhausted after %d retries for %s", constants.MaxCASRetries, resume.InstanceID))
				return false
			}
			continue
		}
		*execution = *currentExecution
		return true
	}
	return false
}

// mergeResumePayloadIntoExecution applies the patches from a normal resume:
// state map merges, node output, signal data fallback, node-state cleanup,
// active-node removal, status flip, and path-entry status update. Pure
// in-memory mutation — persistence happens in the CAS loop above.
func (s *RuntimeService) mergeResumePayloadIntoExecution(execution *entities.WorkflowExecution, resume *ports.ResumeMessage) {
	if resume.StatePatch != nil {
		for k, v := range resume.StatePatch {
			execution.State[k] = v
		}
	}
	if resume.Output != nil {
		execution.NodeOutputs[resume.NodeID] = resume.Output
	}
	if resume.SignalData != nil && resume.Output == nil {
		execution.NodeOutputs[resume.NodeID] = resume.SignalData
	}

	delete(execution.NodeStates, resume.NodeID)
	remaining := make([]string, 0, len(execution.ActiveNodeIDs))
	for _, id := range execution.ActiveNodeIDs {
		if id != resume.NodeID {
			remaining = append(remaining, id)
		}
	}
	execution.ActiveNodeIDs = remaining
	execution.Status = entities.ExecStatusRunning

	pathStatus := domainConstants.StatusCompleted
	if resume.IsTimeout {
		pathStatus = domainConstants.StatusTimeout
	}
	for i := len(execution.ExecutionPath) - 1; i >= 0; i-- {
		if execution.ExecutionPath[i].NodeID == resume.NodeID && execution.ExecutionPath[i].Status == domainConstants.StatusWaiting {
			execution.ExecutionPath[i].Status = pathStatus
			now := time.Now()
			execution.ExecutionPath[i].ExitedAt = &now
			break
		}
	}
}

// applyRetryTimerCheckpoint resets the node state for a retry-timer fire so
// the node executes again, marks the execution as Running with the retried
// node active, and updates the path entry status from Waiting → Retrying.
func (s *RuntimeService) applyRetryTimerCheckpoint(execution *entities.WorkflowExecution, nodeID string, retryAttempt int) {
	delete(execution.NodeStates, nodeID)
	execution.NodeStates[nodeID] = map[string]interface{}{
		constants.NodeStateKeyInternalRetry: retryAttempt,
	}
	execution.Status = entities.ExecStatusRunning
	execution.ActiveNodeIDs = []string{nodeID}

	for i := len(execution.ExecutionPath) - 1; i >= 0; i-- {
		if execution.ExecutionPath[i].NodeID == nodeID && execution.ExecutionPath[i].Status == domainConstants.StatusWaiting {
			execution.ExecutionPath[i].Status = domainConstants.StatusRetrying
			now := time.Now()
			execution.ExecutionPath[i].ExitedAt = &now
			break
		}
	}
}

// loadGraphForResume rebuilds the execution graph for the in-flight execution.
// The TieredCache fronting the definition loader makes this a hit on the hot
// path so the cost is bounded even though it is called per resume.
func (s *RuntimeService) loadGraphForResume(ctx context.Context, execution *entities.WorkflowExecution) (*entities.ExecutionGraph, error) {
	defID := execution.DefinitionID.Hex()
	def, err := s.deps.DefinitionLoader.GetDefinition(ctx, defID, execution.OrgID)
	if err != nil {
		return nil, err
	}
	return domainServices.BuildGraph(def), nil
}

// normaliseResumeError fills the ExecutionError with the contextual fields
// (code, nodeID, timestamp, nodeType) when the callback omitted them — so
// downstream archival, error handlers, and parent callbacks see a complete
// error envelope regardless of which client published the resume.
func (s *RuntimeService) normaliseResumeError(execution *entities.WorkflowExecution, resume *ports.ResumeMessage) {
	if resume.Error == nil {
		resume.Error = &entities.ExecutionError{
			Code:    domainConstants.ErrCodeExternalError,
			Message: "external service reported error status",
		}
	}
	if resume.Error.NodeID == "" {
		resume.Error.NodeID = resume.NodeID
	}
	if resume.Error.Timestamp.IsZero() {
		resume.Error.Timestamp = time.Now()
	}
	if resume.Error.NodeType != "" {
		return
	}
	for i := len(execution.ExecutionPath) - 1; i >= 0; i-- {
		if execution.ExecutionPath[i].NodeID == resume.NodeID {
			resume.Error.NodeType = execution.ExecutionPath[i].NodeType
			return
		}
	}
}

// routeResumeErrorToHandler asks the per-node error policy whether to retry,
// route via an error handle, or surface the failure. Returns (true, nil) when
// the handler absorbed the error (retry scheduled or branch chosen); (false,
// nil) when the orchestration should fall through to failExecution; or an
// error from the handler itself.
func (s *RuntimeService) routeResumeErrorToHandler(ctx context.Context, execution *entities.WorkflowExecution, resume *ports.ResumeMessage) (bool, error) {
	graph, err := s.loadGraphForResume(ctx, execution)
	if err != nil || graph == nil {
		return false, nil
	}
	return s.handleResumeError(ctx, execution, graph, resume.NodeID, resume.Error)
}

// markPathStatusOnNode rewrites the most-recent path entry for nodeID from
// Waiting to the supplied terminal status, recording exit time, optional
// output handle, and optional error message. Mirrors the path-entry update
// logic used by both retry-timer, timeout, and error resume paths.
func markPathStatusOnNode(execution *entities.WorkflowExecution, nodeID, status, outputHandle string, errMsg *string) {
	now := time.Now()
	for i := len(execution.ExecutionPath) - 1; i >= 0; i-- {
		entry := &execution.ExecutionPath[i]
		if entry.NodeID != nodeID || entry.Status != domainConstants.StatusWaiting {
			continue
		}
		entry.Status = status
		if outputHandle != "" {
			entry.OutputHandle = outputHandle
		}
		entry.ExitedAt = &now
		if errMsg != nil {
			entry.Error = errMsg
		}
		return
	}
}

// readRetryAttempt extracts the numeric retry attempt from a node-state map,
// tolerating both float64 (JSON-deserialised) and int (Go-native) shapes.
func readRetryAttempt(ns map[string]interface{}) int {
	if v, ok := ns[constants.NodeStateKeyRetryAttempt].(float64); ok {
		return int(v)
	}
	if v, ok := ns[constants.NodeStateKeyRetryAttempt].(int); ok {
		return v
	}
	return 0
}
