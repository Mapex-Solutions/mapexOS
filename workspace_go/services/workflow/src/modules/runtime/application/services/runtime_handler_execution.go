package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	instancePorts "workflow/src/modules/instances/application/ports"
	appConstants "workflow/src/modules/runtime/application/constants"
	"workflow/src/modules/runtime/application/ports"
	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"

	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/executions"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"github.com/google/uuid"
)

// parseExecutionMessage decodes the WORKFLOW-EXECUTION payload and validates
// that the dispatch mode is present. A missing/invalid payload is rejected
// (permanent) because retrying cannot fix a malformed message.
func (s *RuntimeService) parseExecutionMessage(msg *natsModel.Message) (*v1.WorkflowExecutionMessage, bool) {
	var execMsg v1.WorkflowExecutionMessage
	if err := json.Unmarshal(msg.Data, &execMsg); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] parse FAILED — rejecting msg subject=%s body=%s", msg.Subject, string(msg.Data)))
		msg.Reject(fmt.Sprintf("invalid execution message: %s", err))
		return nil, false
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] parse OK mode=%s dataKeys=%v eventNil=%t", execMsg.Mode, mapKeys(execMsg.Data), execMsg.Event == nil))
	if execMsg.Mode == "" {
		logger.Error(fmt.Errorf("mode missing"), fmt.Sprintf("[SERVICE:Runtime] mode empty — rejecting body=%s", string(msg.Data)))
		msg.Reject("mode is required in execution message")
		return nil, false
	}
	return &execMsg, true
}

// mapKeys returns the keys of a generic map for log-line cardinality —
// shows what fields the dispatcher will see without leaking values.
func mapKeys(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// dispatchExecutionMode routes the execution message to the per-mode handler.
// Routing is by the "mode" field — never by NATS subject — so the same stream
// can carry newInstance, signal, signalOrStart, and subworkflow envelopes.
func (s *RuntimeService) dispatchExecutionMode(msg *natsModel.Message, execMsg *v1.WorkflowExecutionMessage) {
	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] dispatch mode=%s", execMsg.Mode))
	switch execMsg.Mode {
	case constants.ExecutionModeNewInstance:
		s.handleNewInstance(msg, execMsg)
	case constants.ExecutionModeSignal:
		s.handleSignalMode(msg, execMsg)
	case constants.ExecutionModeSignalOrStart:
		s.handleSignalOrStart(msg, execMsg)
	case constants.ExecutionModeSubworkflow:
		s.handleSubworkflow(msg, execMsg)
	default:
		logger.Error(fmt.Errorf("unknown mode"), fmt.Sprintf("[SERVICE:Runtime] unknown mode=%s — rejecting", execMsg.Mode))
		msg.Reject(fmt.Sprintf("unknown execution mode: %s", execMsg.Mode))
	}
}

// handleNewInstance loads instance + definition, creates a new execution, and runs the DAG walker.
func (s *RuntimeService) handleNewInstance(msg *natsModel.Message, execMsg *v1.WorkflowExecutionMessage) {
	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] handleNewInstance entered dataKeys=%v", mapKeys(execMsg.Data)))
	instanceID, _ := execMsg.Data[appConstants.ExecDataKeyInstanceID].(string)
	if instanceID == "" {
		logger.Error(fmt.Errorf("instanceId missing"), fmt.Sprintf("[SERVICE:Runtime] handleNewInstance: data.instanceId empty data=%+v", execMsg.Data))
		msg.Reject("data.instanceId is required for mode 'newInstance'")
		return
	}

	workflowUUID, _ := execMsg.Data[appConstants.ExecDataKeyWorkflowUUID].(string)
	eventTrackerId, _ := execMsg.Data[appConstants.ExecDataKeyEventTrackerID].(string)
	ctx := context.Background()

	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] handleNewInstance loading instance instanceId=%s", instanceID))
	instance, err := s.deps.InstanceLoader.GetInstance(ctx, instanceID)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] handleNewInstance GetInstance err instanceId=%s", instanceID))
		msg.Nack(fmt.Errorf("failed to load instance %s: %w", instanceID, err))
		return
	}
	if instance == nil {
		logger.Error(fmt.Errorf("instance not found"), fmt.Sprintf("[SERVICE:Runtime] handleNewInstance instance nil instanceId=%s", instanceID))
		msg.Reject(fmt.Sprintf("instance %s not found", instanceID))
		return
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] handleNewInstance instance loaded id=%s defId=%s defVer=%d enabled=%t", instanceID, instance.DefinitionID.Hex(), instance.DefinitionVersion, instance.Enabled))

	if instance.OrgID != nil {
		msg.OrgId = instance.OrgID.Hex()
	}
	msg.PathKey = instance.PathKey

	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] handleNewInstance prepareWorkflow defId=%s", instance.DefinitionID.Hex()))
	prepared, err := s.prepareWorkflow(ctx, instance.DefinitionID.Hex())
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] handleNewInstance prepareWorkflow FAILED defId=%s", instance.DefinitionID.Hex()))
		handlePrepareError(msg, err)
		return
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] handleNewInstance prepareWorkflow OK startNode=%s", prepared.startNodeID))

	executionUUID := s.resolveExecutionUUID(instance.UniqueExecution, instance.WorkflowUUID, workflowUUID)

	instObjId, _ := model.ToObjectID(instanceID)
	execution := newBaseExecution(prepared)
	execution.WorkflowUUID = executionUUID
	execution.InstanceID = instObjId
	execution.InstanceName = instance.Name
	execution.WorkflowName = instance.Name
	execution.EventTrackerId = eventTrackerId
	execution.TriggerSource = constants.TriggerSourceWorkflow
	execution.ExternalInputs = instancePorts.InitializeExternalInputs(prepared.def.ExternalInputs, nil)
	execution.EventPayload = execMsg.Event

	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] handleNewInstance createAndExecute uuid=%s startNode=%s", execution.WorkflowUUID, prepared.startNodeID))
	if err := s.createAndExecute(ctx, execution, prepared.graph, prepared.startNodeID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] handleNewInstance createAndExecute FAILED uuid=%s instanceId=%s", execution.WorkflowUUID, instanceID))
		msg.Nack(err)
		return
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Runtime] handleNewInstance createAndExecute OK uuid=%s status=%s", execution.WorkflowUUID, execution.Status))

	logger.Info(fmt.Sprintf("[SERVICE:Runtime] Execution newInstance → instanceId=%s defId=%s uuid=%s status=%s",
		instanceID, instance.DefinitionID.Hex(), executionUUID, execution.Status))
	msg.Ack()
}

// handleSubworkflow creates a child execution from parent context and runs the DAG walker.
func (s *RuntimeService) handleSubworkflow(msg *natsModel.Message, execMsg *v1.WorkflowExecutionMessage) {
	definitionID, _ := execMsg.Data[appConstants.ExecDataKeyDefinitionID].(string)
	parentInstanceID, _ := execMsg.Data[appConstants.ExecDataKeyParentInstanceID].(string)
	parentNodeID, _ := execMsg.Data[appConstants.ExecDataKeyParentNodeID].(string)
	depth := 0

	if d, ok := execMsg.Data[appConstants.ExecDataKeyDepth].(float64); ok {
		depth = int(d)
	}

	inputData, _ := execMsg.Data[appConstants.ExecDataKeyInputData].(map[string]interface{})
	callbackSubject, _ := execMsg.Data[appConstants.ExecDataKeyCallbackSubject].(string)
	executionToken, _ := execMsg.Data[appConstants.ExecDataKeyExecutionToken].(string)

	if definitionID == "" {
		msg.Reject("data.definitionId is required for mode 'subworkflow'")
		return
	}

	ctx := context.Background()

	prepared, err := s.prepareWorkflow(ctx, definitionID)
	if err != nil {
		handlePrepareError(msg, err)
		return
	}

	if prepared.def.OrgID != nil {
		msg.OrgId = prepared.def.OrgID.Hex()
	}
	msg.PathKey = prepared.def.PathKey

	execution := newBaseExecution(prepared)
	execution.WorkflowName = prepared.def.Name
	execution.InstanceName = prepared.def.Name

	retryAttempt := 0
	if a, ok := execMsg.Data[appConstants.ExecDataKeyRetryAttempt].(float64); ok {
		retryAttempt = int(a)
	}
	execution.WorkflowUUID = generateSubworkflowUUID(parentInstanceID, parentNodeID, retryAttempt)
	execution.TriggerSource = constants.TriggerSourceSubworkflow
	execution.ExternalInputs = make(map[string]interface{})
	execution.EventPayload = inputData
	execution.ParentNodeID = parentNodeID
	execution.Depth = depth
	execution.CallbackSubject = callbackSubject
	execution.ParentExecutionToken = executionToken

	if parentInstanceID != "" {
		execution.ParentExecutionID = &parentInstanceID
	}

	if err := s.createAndExecute(ctx, execution, prepared.graph, prepared.startNodeID); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			s.handleSubworkflowDedup(execution, callbackSubject, executionToken)
			msg.Ack()
			return
		}
		msg.Nack(err)
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:Runtime] Execution subworkflow → defId=%s uuid=%s depth=%d parent=%s status=%s",
		definitionID, execution.WorkflowUUID, depth, parentInstanceID, execution.Status))
	msg.Ack()
}

// handleSubworkflowDedup handles the case where a child subworkflow already exists in KV.
// On redelivery, the deterministic UUID produces the same key → KV.Create fails.
// If the child completed/failed, re-publish callback to parent. If running/waiting, do nothing.
func (s *RuntimeService) handleSubworkflowDedup(execution *entities.WorkflowExecution, callbackSubject, executionToken string) {
	existing, err := s.deps.ExecutionStateRepo.Get(execution.WorkflowUUID)
	if err != nil {
		logger.Debug(fmt.Sprintf("[SERVICE:Runtime] Subworkflow dedup: child %s not in KV (already archived), skipping", execution.WorkflowUUID))
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:Runtime] Subworkflow dedup detected: child %s status=%s", execution.WorkflowUUID, existing.Status))

	if existing.Status.IsTerminal() && callbackSubject != "" {
		existing.CallbackSubject = callbackSubject
		existing.ParentExecutionToken = executionToken
		s.publishParentCallback(existing)
	}
}

// resolveExecutionUUID picks the workflow UUID for a new execution using the
// precedence: instance-pinned (UniqueExecution) > caller-supplied > generated.
// Centralised so newInstance, subworkflow, and HTTP-trigger paths share rules.
func (s *RuntimeService) resolveExecutionUUID(uniqueExecution bool, instanceUUID, requestedUUID string) string {
	if uniqueExecution && instanceUUID != "" {
		return instanceUUID
	}
	if requestedUUID != "" {
		return requestedUUID
	}
	return uuid.New().String()
}

// prepareInstanceExecution loads the workflow instance + validates it +
// prepares the executable graph. Used by ExecuteByInstanceID as the
// pre-flight bundle so the public method stays in the §3 line budget.
// Errors are wrapped with operation context and bubble straight to the
// HTTP caller.
func (s *RuntimeService) prepareInstanceExecution(ctx context.Context, instanceID string) (*instancePorts.WorkflowInstance, *preparedWorkflow, error) {
	instance, err := s.deps.InstanceLoader.GetInstance(ctx, instanceID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load instance %s: %w", instanceID, err)
	}
	if err := s.validateInstanceForExecution(instance, instanceID); err != nil {
		return nil, nil, err
	}
	prepared, err := s.prepareWorkflow(ctx, instance.DefinitionID.Hex())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to prepare workflow: %w", err)
	}
	return instance, prepared, nil
}

// validateInstanceForExecution returns a permanent error if the instance is
// missing or disabled. Used by the HTTP-triggered ExecuteByInstanceID path
// where the caller (HTTP handler) must surface the failure synchronously.
func (s *RuntimeService) validateInstanceForExecution(instance *instancePorts.WorkflowInstance, instanceID string) error {
	if instance == nil {
		return fmt.Errorf("instance %s not found", instanceID)
	}
	if !instance.Enabled {
		return fmt.Errorf("instance %s is disabled", instanceID)
	}
	return nil
}

// buildHttpExecution constructs the WorkflowExecution entity for an HTTP-driven
// run from the resolved instance + prepared definition. Trigger source is
// fixed to HTTP so the archived execution preserves provenance.
func (s *RuntimeService) buildHttpExecution(instance *instancePorts.WorkflowInstance, prepared *preparedWorkflow, executionUUID string, eventPayload map[string]interface{}) *entities.WorkflowExecution {
	instObjId, _ := model.ToObjectID(instance.ID.Hex())
	execution := newBaseExecution(prepared)
	execution.WorkflowUUID = executionUUID
	execution.InstanceID = instObjId
	execution.InstanceName = instance.Name
	execution.WorkflowName = instance.Name
	execution.TriggerSource = constants.TriggerSourceHTTP
	execution.ExternalInputs = instancePorts.InitializeExternalInputs(prepared.def.ExternalInputs, nil)
	execution.EventPayload = eventPayload
	return execution
}

// buildExecuteResult shapes the post-execution snapshot returned to HTTP
// callers. Mirrors the fields the gateway needs (uuid, status, error info)
// without leaking the full domain entity across the port boundary.
func (s *RuntimeService) buildExecuteResult(execution *entities.WorkflowExecution) *ports.ExecuteResult {
	result := &ports.ExecuteResult{
		WorkflowUUID: execution.WorkflowUUID,
		Status:       string(execution.Status),
	}
	if execution.ErrorInfo != nil {
		result.ErrorInfo = execution.ErrorInfo
	}
	return result
}
