package services

import (
	"fmt"
	"time"

	appConstants "workflow/src/modules/runtime/application/constants"
	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	sharedTypes "workflow/src/shared/types"
)

// State Transitions

// checkpoint persists the current execution state via the repository.
func (s *RuntimeService) checkpoint(execution *entities.WorkflowExecution) error {
	start := time.Now()
	execution.Updated = start
	err := s.deps.ExecutionStateRepo.Save(execution)
	s.deps.Metrics.CheckpointDuration.Observe(time.Since(start).Seconds())
	return err
}

// failExecution marks the execution as failed, checkpoints, and publishes the state event.
func (s *RuntimeService) failExecution(execution *entities.WorkflowExecution, errInfo *entities.ExecutionError) error {
	execution.Status = entities.ExecStatusFailed
	execution.ErrorInfo = errInfo
	execution.ActiveNodeIDs = nil
	// Purge all pending schedules for this execution
	if err := s.deps.RuntimePublisher.PurgeAllSchedules(execution.WorkflowUUID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Runtime] Failed to purge schedules on fail for %s: %s", execution.WorkflowUUID, err))
	}
	if err := s.checkpoint(execution); err != nil {
		return err
	}
	if err := s.deps.RuntimePublisher.PublishStateEvent(execution, constants.StatusFailed); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish state event 'failed' for execution %s", execution.WorkflowUUID))
	}
	if execution.CallbackSubject != "" {
		s.publishParentCallback(execution)
	}
	s.deps.Metrics.ExecutionFailedTotal.WithLabelValues(execution.TriggerSource).Inc()
	if execution.StartedAt != nil {
		s.deps.Metrics.ExecutionDuration.WithLabelValues(execution.TriggerSource).Observe(time.Since(*execution.StartedAt).Seconds())
	}
	return nil
}

// completeExecution marks the execution as completed, checkpoints, and publishes the state event.
func (s *RuntimeService) completeExecution(execution *entities.WorkflowExecution) error {
	now := time.Now()
	execution.Status = entities.ExecStatusCompleted
	execution.CompletedAt = &now
	execution.ActiveNodeIDs = nil
	// Purge all pending schedules for this execution
	if err := s.deps.RuntimePublisher.PurgeAllSchedules(execution.WorkflowUUID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Runtime] Failed to purge schedules on complete for %s: %s", execution.WorkflowUUID, err))
	}
	if err := s.checkpoint(execution); err != nil {
		return err
	}
	if err := s.deps.RuntimePublisher.PublishStateEvent(execution, constants.StatusCompleted); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish state event 'completed' for execution %s", execution.WorkflowUUID))
	}
	if execution.CallbackSubject != "" {
		s.publishParentCallback(execution)
	}
	s.deps.Metrics.ExecutionCompletedTotal.WithLabelValues(execution.TriggerSource).Inc()
	if execution.StartedAt != nil {
		s.deps.Metrics.ExecutionDuration.WithLabelValues(execution.TriggerSource).Observe(time.Since(*execution.StartedAt).Seconds())
	}
	return nil
}

// completeOrResuspend checks if other nodes are still waiting (e.g. from fanout branches).
// If so, returns to waiting state. Otherwise, completes the execution.
func (s *RuntimeService) completeOrResuspend(execution *entities.WorkflowExecution) error {
	if len(execution.ActiveNodeIDs) > 0 && len(execution.NodeStates) > 0 {
		// Check fanout mode — "firstCompleted" skips re-suspend
		if meta, ok := execution.NodeStates[appConstants.NodeStateKeyFanoutMeta]; ok {
			if mode, _ := meta[appConstants.NodeStateKeyMode].(string); mode == constants.FanoutModeFirstCompleted {
				logger.Info(fmt.Sprintf("[SERVICE:Runtime] Execution %s completing (firstCompleted mode): cancelling %d waiting node(s)",
					execution.WorkflowUUID, len(execution.ActiveNodeIDs)))
				now := time.Now()
				for _, id := range execution.ActiveNodeIDs {
					for i := len(execution.ExecutionPath) - 1; i >= 0; i-- {
						if execution.ExecutionPath[i].NodeID == id && execution.ExecutionPath[i].Status == constants.StatusWaiting {
							execution.ExecutionPath[i].Status = constants.StatusCancelled
							execution.ExecutionPath[i].ExitedAt = &now
							break
						}
					}
					delete(execution.NodeStates, id)
				}
				delete(execution.NodeStates, appConstants.NodeStateKeyFanoutMeta)
				execution.ActiveNodeIDs = nil
				return s.completeExecution(execution)
			}
		}

		stillWaiting := make([]string, 0)
		for _, id := range execution.ActiveNodeIDs {
			if ns, ok := execution.NodeStates[id]; ok && ns[appConstants.NodeStateKeyWaitType] != nil {
				stillWaiting = append(stillWaiting, id)
			}
		}
		if len(stillWaiting) > 0 {
			execution.Status = entities.ExecStatusWaiting
			execution.ActiveNodeIDs = stillWaiting
			logger.Info(fmt.Sprintf("[SERVICE:Runtime] Execution %s re-suspended: %d node(s) still waiting %v",
				execution.WorkflowUUID, len(stillWaiting), stillWaiting))
			if err := s.checkpoint(execution); err != nil {
				return err
			}
			if err := s.deps.RuntimePublisher.PublishStateEvent(execution, constants.StatusWaiting); err != nil {
				logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish re-suspend state event for execution %s", execution.WorkflowUUID))
			}
			return nil
		}
	}
	return s.completeExecution(execution)
}

// Suspension + Dispatch

// suspendExecution marks the execution as waiting, generates execution token,
// dispatches to external service FIRST, then checkpoints and publishes state event.
// If dispatch fails, no checkpoint is written — the original NATS message is not ACKed
// and NATS redelivers. NATS Msg-Id dedup prevents duplicate dispatch on redelivery.
func (s *RuntimeService) suspendExecution(
	execution *entities.WorkflowExecution,
	nodeID string,
	nodeType string,
	nodeState map[string]interface{},
) error {
	execution.Status = entities.ExecStatusWaiting
	execution.ActiveNodeIDs = []string{nodeID}

	// Preserve __retryAttempt across suspension cycles (retry timer → callback → retry timer)
	if prev := execution.NodeStates[nodeID]; prev != nil {
		if attempt, ok := prev[appConstants.NodeStateKeyInternalRetry]; ok {
			nodeState[appConstants.NodeStateKeyInternalRetry] = attempt
		}
	}
	execution.NodeStates[nodeID] = nodeState

	// Generate execution token + Msg-Id BEFORE dispatch
	attempt := getRetryAttempt(execution.NodeStates, nodeID)
	token := generateExecutionToken(execution.WorkflowUUID, nodeID, attempt)
	nodeState[appConstants.NodeStateKeyExecutionToken] = token
	msgId := buildMsgId(execution.WorkflowUUID, nodeID, attempt, execution.State)

	// Dispatch FIRST — if this fails, no checkpoint, NATS redelivers
	if err := s.dispatchByNodeType(execution, nodeID, nodeType, nodeState, token, msgId); err != nil {
		return err
	}

	// Checkpoint AFTER dispatch succeeded
	if err := s.checkpoint(execution); err != nil {
		return err
	}

	// Publish NATS schedule for timer-based resume
	if expiresAt, ok := nodeState[appConstants.NodeStateKeyExpiresAt].(time.Time); ok {
		waitType, _ := nodeState[appConstants.NodeStateKeyWaitType].(string)
		enableOutput, _ := nodeState[appConstants.NodeStateKeyEnableOutput].(bool)
		if err := s.deps.RuntimePublisher.PublishSchedule(execution.WorkflowUUID, nodeID, expiresAt, waitType, enableOutput); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish schedule for %s node %s", execution.WorkflowUUID, nodeID))
		}
	}

	if err := s.publishWaitingStateEvent(execution); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish waiting state event for execution %s", execution.WorkflowUUID))
	}

	return nil
}

// suspendFanoutExecution generates tokens, dispatches all waiting nodes, then publishes state event.
// The checkpoint is handled by the caller (walker fanout path) AFTER this function returns.
func (s *RuntimeService) suspendFanoutExecution(execution *entities.WorkflowExecution) {
	// Dispatch ALL nodes FIRST (token + Msg-Id per node)
	for _, nodeID := range execution.ActiveNodeIDs {
		ns := execution.NodeStates[nodeID]
		if ns == nil {
			continue
		}
		nodeType, _ := ns[appConstants.NodeStateKeyNodeType].(string)
		attempt := getRetryAttempt(execution.NodeStates, nodeID)
		token := generateExecutionToken(execution.WorkflowUUID, nodeID, attempt)
		ns[appConstants.NodeStateKeyExecutionToken] = token
		msgId := buildMsgId(execution.WorkflowUUID, nodeID, attempt, execution.State)
		if err := s.dispatchByNodeType(execution, nodeID, nodeType, ns, token, msgId); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to dispatch fanout waiting node %s (type=%s) in execution %s", nodeID, nodeType, execution.WorkflowUUID))
		}
	}

	// Publish NATS schedules for timer-based resume (one per node with expiresAt)
	for _, nodeID := range execution.ActiveNodeIDs {
		ns := execution.NodeStates[nodeID]
		if ns == nil {
			continue
		}
		if expiresAt, ok := ns[appConstants.NodeStateKeyExpiresAt].(time.Time); ok {
			waitType, _ := ns[appConstants.NodeStateKeyWaitType].(string)
			enableOutput, _ := ns[appConstants.NodeStateKeyEnableOutput].(bool)
			if err := s.deps.RuntimePublisher.PublishSchedule(execution.WorkflowUUID, nodeID, expiresAt, waitType, enableOutput); err != nil {
				logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish fanout schedule for %s node %s", execution.WorkflowUUID, nodeID))
			}
		}
	}

	// State event AFTER dispatches
	if err := s.publishWaitingStateEvent(execution); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish waiting state event (fanout) for execution %s", execution.WorkflowUUID))
	}
}

// dispatchByNodeType routes a suspended node to the appropriate external service publisher.
func (s *RuntimeService) dispatchByNodeType(
	execution *entities.WorkflowExecution,
	nodeID string,
	nodeType string,
	nodeState map[string]interface{},
	executionToken, msgId string,
) error {
	var err error
	var dispatchType string

	switch nodeType {
	case constants.NodeTypeDelay, constants.NodeTypeRetry:
		return nil // Timer-based: NATS Schedule handles delivery, no external dispatch needed
	case constants.NodeTypeCode:
		dispatchType = constants.DispatchTypeCode
		err = s.deps.RuntimePublisher.DispatchCodeExecution(execution, nodeID, nodeState, executionToken, msgId)
	case constants.NodeTypeSubworkflow:
		dispatchType = constants.DispatchTypeSubworkflow
		err = s.deps.RuntimePublisher.DispatchSubworkflowExecution(execution, nodeID, nodeState, executionToken, msgId)
	case constants.NodeTypeTriggerEvent:
		dispatchType = constants.DispatchTypeTrigger
		triggerData := map[string]interface{}{
			appConstants.NodeStateKeyTriggerID: nodeState[appConstants.NodeStateKeyTriggerID],
			appConstants.NodeStateKeyPayload:   nodeState[appConstants.NodeStateKeyPayload],
		}
		logger.Debug(fmt.Sprintf("[SERVICE:Runtime] Dispatching trigger entity: instanceId=%s, nodeId=%s, triggerId=%v",
			execution.WorkflowUUID, nodeID, nodeState[appConstants.NodeStateKeyTriggerID]))
		err = s.deps.RuntimePublisher.DispatchWorkflowTrigger(execution, nodeID, constants.DispatchTypeTrigger, triggerData, executionToken, msgId)
	case constants.NodeTypeWaitSignal, constants.NodeTypeWaitFor:
		return nil
	default:
		dispatchType = constants.DispatchTypePlugin
		err = s.dispatchPluginByActionType(execution, nodeID, nodeState, executionToken, msgId)
	}

	if err != nil {
		s.deps.Metrics.DispatchTotal.WithLabelValues(dispatchType, constants.DispatchOutcomeError).Inc()
	} else {
		s.deps.Metrics.DispatchTotal.WithLabelValues(dispatchType, constants.DispatchOutcomeSuccess).Inc()
	}
	return err
}

// dispatchPluginByActionType routes a plugin node's execution based on the action type.
//
// Routing:
//   - "http", "mqtt", "nats", "email", "rabbitmq", "websocket" → Triggers Service (mode "plugin")
//   - "script" → JS Workflow Executor (V8)
func (s *RuntimeService) dispatchPluginByActionType(
	execution *entities.WorkflowExecution,
	nodeID string,
	nodeState map[string]interface{},
	executionToken, msgId string,
) error {
	action, _ := nodeState[appConstants.NodeStateKeyAction].(map[string]interface{})
	actionType := ""
	if action != nil {
		actionType, _ = action[appConstants.ActionKeyType].(string)
	}

	switch actionType {
	case constants.ActionTypeHTTP, constants.ActionTypeMQTT, constants.ActionTypeNATS, constants.ActionTypeEmail, constants.ActionTypeRabbitMQ, constants.ActionTypeWebsocket:
		pluginData := map[string]interface{}{
			appConstants.NodeStateKeyPluginID:  nodeState[appConstants.NodeStateKeyPluginID],
			appConstants.NodeStateKeyNodeType:  nodeState[appConstants.NodeStateKeyNodeType],
			appConstants.NodeStateKeyOperation: nodeState[appConstants.NodeStateKeyOperation],
			appConstants.NodeStateKeyAction:    nodeState[appConstants.NodeStateKeyAction],
		}
		if hooks, ok := nodeState[appConstants.NodeStateKeyHooks]; ok {
			pluginData[appConstants.NodeStateKeyHooks] = hooks
		}
		logger.Debug(fmt.Sprintf("[SERVICE:Runtime] Dispatching plugin action: instanceId=%s, nodeId=%s, pluginId=%v, operation=%v, actionType=%s",
			execution.WorkflowUUID, nodeID, nodeState[appConstants.NodeStateKeyPluginID], nodeState[appConstants.NodeStateKeyOperation], actionType))
		return s.deps.RuntimePublisher.DispatchWorkflowTrigger(execution, nodeID, constants.DispatchTypePlugin, pluginData, executionToken, msgId)
	case constants.ActionTypeScript:
		return s.deps.RuntimePublisher.DispatchCodeExecution(execution, nodeID, nodeState, executionToken, msgId)
	default:
		return fmt.Errorf("[SERVICE:Runtime] unsupported plugin action type '%s' for node %s", actionType, nodeID)
	}
}

// State Event Helpers

func (s *RuntimeService) publishWaitingStateEvent(execution *entities.WorkflowExecution) error {
	return s.deps.RuntimePublisher.PublishStateEvent(execution, constants.StatusWaiting)
}

func (s *RuntimeService) publishResumedStateEvent(execution *entities.WorkflowExecution) error {
	return s.deps.RuntimePublisher.PublishStateEvent(execution, constants.StatusResumed)
}

// publishParentCallback sends a ResumeMessage to the parent execution's callback subject.
// Called when a subworkflow child reaches terminal state (completed or failed).
func (s *RuntimeService) publishParentCallback(execution *entities.WorkflowExecution) {
	if execution.ParentExecutionID == nil {
		return
	}

	resume := sharedTypes.ResumeMessage{
		InstanceID:     *execution.ParentExecutionID,
		NodeID:         execution.ParentNodeID,
		ExecutionToken: execution.ParentExecutionToken,
	}

	if execution.Status == entities.ExecStatusCompleted {
		resume.Status = constants.StatusCompleted
		resume.Output = execution.State
	} else {
		resume.Status = constants.StatusError
		if execution.ErrorInfo != nil {
			resume.Error = execution.ErrorInfo
		}
	}

	if err := s.deps.RuntimePublisher.PublishCallbackResume(execution.CallbackSubject, resume); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish parent callback for child %s to %s", execution.WorkflowUUID, execution.CallbackSubject))
	}
}

// Loop Stack Helpers

// pushLoopStack adds a loop nodeId to the stack in NodeStates.
func pushLoopStack(nodeStates map[string]map[string]interface{}, loopNodeID string) {
	stack := getLoopStack(nodeStates)
	stack = append(stack, loopNodeID)
	if nodeStates[constants.LoopStackKey] == nil {
		nodeStates[constants.LoopStackKey] = make(map[string]interface{})
	}
	nodeStates[constants.LoopStackKey][appConstants.NodeStateKeyStack] = stack
}

// popLoopStack removes and returns the top loop nodeId from the stack.
func popLoopStack(nodeStates map[string]map[string]interface{}) string {
	stack := getLoopStack(nodeStates)
	if len(stack) == 0 {
		return ""
	}
	top := stack[len(stack)-1]
	topStr, _ := top.(string)
	stack = stack[:len(stack)-1]
	if len(stack) == 0 {
		delete(nodeStates, constants.LoopStackKey)
	} else {
		nodeStates[constants.LoopStackKey][appConstants.NodeStateKeyStack] = stack
	}
	return topStr
}

// peekLoopStack returns the top loop nodeId without removing it.
func peekLoopStack(nodeStates map[string]map[string]interface{}) string {
	stack := getLoopStack(nodeStates)
	if len(stack) == 0 {
		return ""
	}
	topStr, _ := stack[len(stack)-1].(string)
	return topStr
}

func getLoopStack(nodeStates map[string]map[string]interface{}) []interface{} {
	if nodeStates[constants.LoopStackKey] == nil {
		return nil
	}
	stack, _ := nodeStates[constants.LoopStackKey][appConstants.NodeStateKeyStack].([]interface{})
	return stack
}
