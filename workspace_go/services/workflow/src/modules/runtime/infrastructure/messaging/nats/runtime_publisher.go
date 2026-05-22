package nats

import (
	"fmt"
	"time"

	"workflow/src/modules/runtime/application/constants"
	"workflow/src/modules/runtime/application/ports"
	"workflow/src/modules/runtime/domain/entities"
	sharedTypes "workflow/src/shared/types"

	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/executions"
	runtimeContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/runtime"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time check
var _ ports.RuntimePublisherPort = (*RuntimePublisher)(nil)

// NewRuntimePublisher creates a new RuntimePublisher with DIG-injected dependencies.
func NewRuntimePublisher(params RuntimePublisherParams) ports.RuntimePublisherPort {
	return &RuntimePublisher{
		publisher:       params.Publisher,
		scheduleManager: params.ScheduleManager,
	}
}

// PublishStateEvent publishes a workflow state change event to the NATS subject workflow.state.{status}.
func (p *RuntimePublisher) PublishStateEvent(execution *entities.WorkflowExecution, status string) error {
	orgID := ""
	if execution.OrgID != nil {
		orgID = execution.OrgID.Hex()
	}
	event := sharedTypes.StateEvent{
		InstanceID:     execution.WorkflowUUID,
		ExecutionId:    execution.ID.Hex(),
		WorkflowID:     execution.DefinitionID.Hex(),
		OrgID:          orgID,
		WorkflowName:   execution.WorkflowName,
		InstanceName:   execution.InstanceName,
		DefinitionName: execution.DefinitionName,
		InstanceObjID:  execution.InstanceID.Hex(),
		Status:         string(execution.Status),
		ActiveNodeIDs:  execution.ActiveNodeIDs,
		Version:        execution.Version,
		TriggerSource:  execution.TriggerSource,
	}
	subject := fmt.Sprintf(constants.StatePatternSubject, status)
	if err := p.publisher.Publish(natsModel.PublishConfig{
		Subject: subject,
		Data:    event,
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:RuntimePublisher] Failed to publish state event %s for %s", status, execution.WorkflowUUID))
		return err
	}
	return nil
}

// PublishResumeMessage publishes a resume/re-enqueue message for a suspended execution.
func (p *RuntimePublisher) PublishResumeMessage(executionID string, nodeID string, status string) error {
	msg := sharedTypes.ResumeMessage{
		InstanceID: executionID,
		NodeID:     nodeID,
		Status:     status,
	}
	subject := fmt.Sprintf(constants.ResumeReenqueuePatternSubject, executionID)
	if err := p.publisher.Publish(natsModel.PublishConfig{
		Subject: subject,
		Data:    msg,
	}); err != nil {
		return fmt.Errorf("[INFRA:RuntimePublisher] failed to publish resume for %s: %w", executionID, err)
	}
	return nil
}

// PublishResumeTimer re-publishes a raw fired-schedule body to WORKFLOW-RESUME on
// the per-instance timer subject workflow.resume.timer.{instanceId}.
func (p *RuntimePublisher) PublishResumeTimer(instanceID string, body map[string]interface{}) error {
	subject := fmt.Sprintf(constants.ResumeTimerPatternSubject, instanceID)
	return p.publisher.Publish(natsModel.PublishConfig{
		Subject: subject,
		Data:    body,
	})
}

// DispatchCodeExecution publishes a code execution request to the WORKFLOW-JS-CODE stream.
func (p *RuntimePublisher) DispatchCodeExecution(execution *entities.WorkflowExecution, nodeID string, nodeState map[string]interface{}, executionToken, msgId string) error {
	orgID := ""
	if execution.OrgID != nil {
		orgID = execution.OrgID.Hex()
	}

	callbackSubject := fmt.Sprintf(constants.ResumeCallbackPatternSubject, execution.WorkflowUUID)

	req := runtimeContract.CodeExecutionRequest{
		OrgID:           orgID,
		PathKey:         execution.PathKey,
		WorkflowID:      execution.DefinitionID.Hex(),
		NodeID:          nodeID,
		InstanceID:      execution.WorkflowUUID,
		CallbackSubject: callbackSubject,
		ExecutionToken:  executionToken,
		Timeout:         model.MapGetInt(nodeState, "timeout"),
		EventPayload:    execution.EventPayload,
		State:           execution.State,
		Inputs:          execution.ExternalInputs,
		NodeOutputs:     execution.NodeOutputs,
	}
	if err := p.publisher.Publish(natsModel.PublishConfig{
		Subject: constants.CodeSubject,
		Data:    req,
		Headers: map[string]string{"Nats-Msg-Id": msgId},
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:RuntimePublisher] Failed to publish code request for %s", execution.WorkflowUUID))
		return err
	}
	return nil
}

// DispatchSubworkflowExecution publishes a subworkflow execution request to WORKFLOW-EXECUTION stream.
func (p *RuntimePublisher) DispatchSubworkflowExecution(execution *entities.WorkflowExecution, nodeID string, nodeState map[string]interface{}, executionToken, msgId string) error {
	depth := model.MapGetInt(nodeState, "depth")
	wfID := fmt.Sprintf("%v", nodeState["workflowId"])
	inputData := model.MapGetMap(nodeState, "inputData")
	callbackSubject := fmt.Sprintf(constants.ResumeCallbackPatternSubject, execution.WorkflowUUID)

	msg := v1.WorkflowExecutionMessage{
		Mode: "subworkflow",
		Data: map[string]interface{}{
			"definitionId":    wfID,
			"parentInstanceId": execution.WorkflowUUID,
			"parentNodeId":    nodeID,
			"callbackSubject":  callbackSubject,
			"executionToken":   executionToken,
			"depth":            depth,
			"inputData":        inputData,
		},
	}
	subject := fmt.Sprintf(constants.ExecutionSubworkflowPatternSubject, execution.WorkflowUUID)
	if err := p.publisher.Publish(natsModel.PublishConfig{
		Subject: subject,
		Data:    msg,
		Headers: map[string]string{"Nats-Msg-Id": msgId},
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:RuntimePublisher] Failed to publish subworkflow execution for %s", execution.WorkflowUUID))
		return err
	}
	return nil
}

// DispatchWorkflowTrigger publishes a workflow trigger request to the TRIGGERS stream.
// Subject: trigger.WORKFLOW.execute
// mode "trigger": uses a registered trigger entity (triggerId in data)
// mode "plugin": uses a fully resolved plugin action pipeline (action + hooks in data)
func (p *RuntimePublisher) DispatchWorkflowTrigger(execution *entities.WorkflowExecution, nodeID string, mode string, data map[string]interface{}, executionToken, msgId string) error {
	orgID := ""
	if execution.OrgID != nil {
		orgID = execution.OrgID.Hex()
	}

	callbackSubject := fmt.Sprintf(constants.ResumeCallbackPatternSubject, execution.WorkflowUUID)

	req := runtimeContract.WorkflowTriggerRequest{
		Mode:            mode,
		OrgID:           orgID,
		PathKey:         execution.PathKey,
		WorkflowID:      execution.DefinitionID.Hex(),
		InstanceID:      execution.WorkflowUUID,
		NodeID:          nodeID,
		CallbackSubject: callbackSubject,
		ExecutionToken:  executionToken,
		Data:            data,
	}

	if err := p.publisher.Publish(natsModel.PublishConfig{
		Subject: constants.TriggerWorkflowExecuteSubject,
		Data:    req,
		Headers: map[string]string{"Nats-Msg-Id": msgId},
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:RuntimePublisher] Failed to publish workflow trigger (mode=%s) for %s node %s", mode, execution.WorkflowUUID, nodeID))
		return err
	}
	return nil
}

// PublishSignalResume publishes a resume message with signal data to WORKFLOW-RESUME.
func (p *RuntimePublisher) PublishSignalResume(executionID string, nodeID string, signalData map[string]interface{}) error {
	msg := sharedTypes.ResumeMessage{
		InstanceID: executionID,
		NodeID:     nodeID,
		Status:     "success",
		SignalData: signalData,
	}
	subject := fmt.Sprintf(constants.ResumeSignalPatternSubject, executionID)
	if err := p.publisher.Publish(natsModel.PublishConfig{
		Subject: subject,
		Data:    msg,
	}); err != nil {
		return fmt.Errorf("[INFRA:RuntimePublisher] failed to publish signal resume for %s: %w", executionID, err)
	}
	return nil
}

// PublishCallbackResume publishes a resume message to a specific callback subject.
// Used by subworkflow child to notify parent on terminal state.
func (p *RuntimePublisher) PublishCallbackResume(subject string, resume sharedTypes.ResumeMessage) error {
	if err := p.publisher.Publish(natsModel.PublishConfig{
		Subject: subject,
		Data:    resume,
	}); err != nil {
		return fmt.Errorf("[INFRA:RuntimePublisher] failed to publish callback resume to %s: %w", subject, err)
	}
	return nil
}

// PublishSchedule publishes a NATS scheduled message for a timer.
// The message body varies by waitType:
//   - "timer" (delay): normal resume — schedule IS the completion
//   - "retryTimer": retry trigger — HandleResume detects retryTimer in NodeState
//   - default (callback/signal/condition): timeout safety net
func (p *RuntimePublisher) PublishSchedule(wfUUID, nodeID string, expiresAt time.Time, waitType string, enableOutput bool) error {
	var data map[string]interface{}

	switch waitType {
	case "timer":
		data = map[string]interface{}{
			"instanceId": wfUUID,
			"nodeId":     nodeID,
			"status":     "success",
		}
	case "retryTimer":
		data = map[string]interface{}{
			"instanceId": wfUUID,
			"nodeId":     nodeID,
			"isTimeout":  true,
		}
	default:
		data = map[string]interface{}{
			"instanceId":   wfUUID,
			"nodeId":       nodeID,
			"isTimeout":    true,
			"enableOutput": enableOutput,
		}
	}

	subject := fmt.Sprintf("%s.%s.%s", constants.ScheduleSubjectPrefix, wfUUID, nodeID)
	if err := p.scheduleManager.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       subject,
		TargetSubject: constants.ScheduleTargetSubject,
		ScheduleAt:    expiresAt,
		Data:          data,
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:RuntimePublisher] Failed to publish schedule for %s node %s (waitType=%s)", wfUUID, nodeID, waitType))
		return err
	}

	logger.Debug(fmt.Sprintf("[INFRA:RuntimePublisher] Schedule published: %s node %s waitType=%s expiresAt=%s",
		wfUUID, nodeID, waitType, expiresAt.UTC().Format(time.RFC3339)))
	return nil
}

// PurgeSchedule cancels a pending schedule for a specific node.
// Idempotent: returns nil if the schedule already fired or was never published.
func (p *RuntimePublisher) PurgeSchedule(wfUUID, nodeID string) error {
	subject := fmt.Sprintf("%s.%s.%s", constants.ScheduleSubjectPrefix, wfUUID, nodeID)
	return p.scheduleManager.PurgeStreamSubject(constants.ScheduleStreamName, subject)
}

// PurgeAllSchedules cancels all pending schedules for an entire workflow execution.
// Called by failExecution and completeExecution for explicit cleanup.
func (p *RuntimePublisher) PurgeAllSchedules(wfUUID string) error {
	subject := fmt.Sprintf("%s.%s.>", constants.ScheduleSubjectPrefix, wfUUID)
	return p.scheduleManager.PurgeStreamSubject(constants.ScheduleStreamName, subject)
}

