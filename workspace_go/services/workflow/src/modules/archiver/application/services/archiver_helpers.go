package services

import (
	"encoding/json"
	"fmt"
	"time"

	archiverConstants "workflow/src/modules/archiver/application/constants"
	dtos "workflow/src/modules/archiver/application/dtos"
	"workflow/src/modules/archiver/application/ports"
	archiverTypes "workflow/src/modules/archiver/application/types"
	"workflow/src/modules/archiver/domain/repositories"
	runtimeConstants "workflow/src/modules/runtime/application/constants"
	runtimePorts "workflow/src/modules/runtime/application/ports"

	common "github.com/Mapex-Solutions/MapexOS/contracts/common"
	execDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/executions"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// State-event subjects are built at package init from the same
// StatePatternSubject constant the runtime publisher uses, so the
// archiver classifies on the env-prefixed subject (e.g. when
// GO_ENV=dev the full subject is dev.mapexos.workflow.state.created).
// Comparing against an unprefixed literal here previously rejected
// every state event as "unknown" and silently dropped the workflow
// state history to the DLQ.
var (
	stateSubjectCreated   = fmt.Sprintf(runtimeConstants.StatePatternSubject, "created")
	stateSubjectWaiting   = fmt.Sprintf(runtimeConstants.StatePatternSubject, "waiting")
	stateSubjectResumed   = fmt.Sprintf(runtimeConstants.StatePatternSubject, "resumed")
	stateSubjectCompleted = fmt.Sprintf(runtimeConstants.StatePatternSubject, "completed")
	stateSubjectFailed    = fmt.Sprintf(runtimeConstants.StatePatternSubject, "failed")
	stateSubjectCancelled = fmt.Sprintf(runtimeConstants.StatePatternSubject, "cancelled")
)

/*
 * HTTP QUERY HELPERS
 */

// mapExecutionToDTO converts a WorkflowExecution entity to the contract DTO.
func mapExecutionToDTO(exec *runtimePorts.WorkflowExecution) dtos.ExecutionResponseDTO {
	id := exec.ID.Hex()
	wfUUID := exec.WorkflowUUID
	name := exec.WorkflowName
	instName := exec.InstanceName
	defName := exec.DefinitionName
	version := exec.Version
	status := execDtos.ExecutionStatus(exec.Status)
	depth := exec.Depth

	dto := dtos.ExecutionResponseDTO{
		ID:             &id,
		WorkflowUUID:   &wfUUID,
		InstanceID:     &exec.InstanceID,
		DefinitionID:   &exec.DefinitionID,
		WorkflowName:   &name,
		InstanceName:   &instName,
		DefinitionName: &defName,
		OrgID:          exec.OrgID,
		Version:        &version,
		Status:         &status,
		ActiveNodeIDs:  exec.ActiveNodeIDs,
		State:          exec.State,
		EventPayload:   exec.EventPayload,
		NodeOutputs:    exec.NodeOutputs,
		Depth:          &depth,
	}

	// ExecutionPath
	if len(exec.ExecutionPath) > 0 {
		pathEntries := make([]execDtos.PathEntryResponse, 0, len(exec.ExecutionPath))
		for _, pe := range exec.ExecutionPath {
			entry := execDtos.PathEntryResponse{
				NodeID:       pe.NodeID,
				NodeType:     pe.NodeType,
				Status:       pe.Status,
				DurationMs:   pe.DurationMs,
				OutputHandle: pe.OutputHandle,
				Error:        pe.Error,
			}
			pathEntries = append(pathEntries, entry)
		}
		dto.ExecutionPath = pathEntries
	}

	// ErrorInfo
	if exec.ErrorInfo != nil {
		dto.ErrorInfo = &execDtos.ErrorInfoResponse{
			Code:       exec.ErrorInfo.Code,
			Message:    exec.ErrorInfo.Message,
			NodeID:     exec.ErrorInfo.NodeID,
			NodeType:   exec.ErrorInfo.NodeType,
			StackTrace: exec.ErrorInfo.StackTrace,
		}
	}

	// TriggerSource
	if exec.TriggerSource != "" {
		ts := exec.TriggerSource
		dto.TriggerSource = &ts
	}

	// ParentExecutionID
	if exec.ParentExecutionID != nil {
		dto.ParentExecutionID = exec.ParentExecutionID
	}

	// Timestamps
	dto.Created = &common.NullTime{Time: exec.Created}
	dto.Updated = &common.NullTime{Time: exec.Updated}
	if exec.StartedAt != nil {
		dto.StartedAt = &common.NullTime{Time: *exec.StartedAt}
	}
	if exec.CompletedAt != nil {
		dto.CompletedAt = &common.NullTime{Time: *exec.CompletedAt}
	}

	return dto
}

/*
 * HELPERS
 */

// fetchFullState retrieves the complete execution from NATS KV for terminal archiving.
func (s *ArchiverService) fetchFullState(executionID string) (*runtimePorts.WorkflowExecution, string, error) {
	kvKey := fmt.Sprintf("exec.%s", executionID)
	entry, err := s.deps.KVStore.Get(kvKey)
	if err != nil {
		return nil, kvKey, fmt.Errorf("KV Get failed: %w", err)
	}

	var execution runtimePorts.WorkflowExecution
	if err := json.Unmarshal(entry.Value, &execution); err != nil {
		return nil, kvKey, fmt.Errorf("unmarshal failed: %w", err)
	}

	return &execution, kvKey, nil
}

// publishWorkflowEvent publishes a terminal execution to EVENTS-WORKFLOW for ClickHouse cold storage.
// JSON fields (executionPath, nodeOutputs, errorInfo, eventPayload) are serialized as strings.
func (s *ArchiverService) publishWorkflowEvent(exec *runtimePorts.WorkflowExecution) {
	// Calculate duration
	var durationMs int64
	if exec.CompletedAt != nil {
		durationMs = exec.CompletedAt.Sub(exec.Created).Milliseconds()
	}

	// Determine success
	success := exec.Status == runtimePorts.ExecStatusCompleted

	// Extract error message
	errorMessage := ""
	if exec.ErrorInfo != nil {
		errorMessage = exec.ErrorInfo.Message
	}

	// Serialize JSON fields using response DTOs (camelCase json tags)
	pathDTOs := make([]execDtos.PathEntryResponse, 0, len(exec.ExecutionPath))
	for _, pe := range exec.ExecutionPath {
		entry := execDtos.PathEntryResponse{
			NodeID:       pe.NodeID,
			NodeType:     pe.NodeType,
			Status:       pe.Status,
			DurationMs:   pe.DurationMs,
			OutputHandle: pe.OutputHandle,
			Error:        pe.Error,
		}
		if !pe.EnteredAt.IsZero() {
			entry.EnteredAt = &common.NullTime{Time: pe.EnteredAt}
		}
		if pe.ExitedAt != nil && !pe.ExitedAt.IsZero() {
			entry.ExitedAt = &common.NullTime{Time: *pe.ExitedAt}
		}
		pathDTOs = append(pathDTOs, entry)
	}
	executionPathJSON, _ := json.Marshal(pathDTOs)
	nodeOutputsJSON, _ := json.Marshal(exec.NodeOutputs)
	errorInfoJSON := ""
	if exec.ErrorInfo != nil {
		errDTO := execDtos.ErrorInfoResponse{
			Code:       exec.ErrorInfo.Code,
			Message:    exec.ErrorInfo.Message,
			NodeID:     exec.ErrorInfo.NodeID,
			NodeType:   exec.ErrorInfo.NodeType,
			StackTrace: exec.ErrorInfo.StackTrace,
		}
		b, _ := json.Marshal(errDTO)
		errorInfoJSON = string(b)
	}
	eventPayloadJSON, _ := json.Marshal(exec.EventPayload)
	stateJSON, _ := json.Marshal(exec.State)
	externalInputsJSON, _ := json.Marshal(exec.ExternalInputs)

	// Build org/path info
	orgId := ""
	if exec.OrgID != nil {
		orgId = exec.OrgID.Hex()
	}

	// Finished timestamp
	finished := time.Now()
	if exec.CompletedAt != nil {
		finished = *exec.CompletedAt
	}

	// Parent execution ID
	parentExecutionId := ""
	if exec.ParentExecutionID != nil {
		parentExecutionId = *exec.ParentExecutionID
	}

	payload := map[string]interface{}{
		"created":           exec.Created,
		"finished":          finished,
		"eventTrackerId":    exec.EventTrackerId,
		"executionId":       exec.ID.Hex(),
		"orgId":             orgId,
		"pathKey":           exec.PathKey,
		"workflowUUID":      exec.WorkflowUUID,
		"instanceId":        exec.InstanceID.Hex(),
		"definitionId":      exec.DefinitionID.Hex(),
		"workflowName":      exec.WorkflowName,
		"instanceName":      exec.InstanceName,
		"definitionName":    exec.DefinitionName,
		"status":            string(exec.Status),
		"success":           success,
		"durationMs":        durationMs,
		"errorMessage":      errorMessage,
		"executionPath":     string(executionPathJSON),
		"nodeOutputs":       string(nodeOutputsJSON),
		"errorInfo":         errorInfoJSON,
		"eventPayload":      string(eventPayloadJSON),
		"triggerSource":     exec.TriggerSource,
		"parentExecutionId": parentExecutionId,
		"depth":             exec.Depth,
		"retentionDays":     archiveRetentionDays(exec.RetentionDays),
		"state":             string(stateJSON),
		"externalInputs":    string(externalInputsJSON),
	}

	if err := s.deps.Publisher.Publish(natsModel.PublishConfig{
		Subject: archiverConstants.EventsWorkflowSubject,
		Data:    payload,
	}); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Archiver] Failed to publish workflow event for %s: %s", exec.WorkflowUUID, err))
	}
}

// archiveRetentionDays maps the execution's RetentionDays (uint16, 0=org default)
// to the events contract convention (int16: -1=org default, 0=forever, >0=custom).
func archiveRetentionDays(execRetention uint16) int16 {
	if execRetention == 0 {
		return -1 // org default
	}
	return int16(execRetention)
}

func buildLightweightStub(event ports.StateEvent) (repositories.LightweightExecution, error) {
	definitionID, err := model.ToObjectID(event.WorkflowID)
	if err != nil {
		return repositories.LightweightExecution{}, fmt.Errorf("invalid workflowId: %w", err)
	}

	// Parse pre-generated executionId (MongoDB ObjectId hex from runtime)
	var executionObjId model.ObjectId
	if event.ExecutionId != "" {
		oid, oidErr := model.ToObjectID(event.ExecutionId)
		if oidErr == nil {
			executionObjId = oid
		}
	}

	// Parse instanceId (MongoDB ObjectId hex)
	var instanceObjId model.ObjectId
	if event.InstanceObjID != "" {
		oid, oidErr := model.ToObjectID(event.InstanceObjID)
		if oidErr == nil {
			instanceObjId = oid
		}
	}

	now := time.Now()
	stub := repositories.LightweightExecution{
		ID:             executionObjId,
		WorkflowUUID:   event.InstanceID,
		InstanceID:     instanceObjId,
		DefinitionID:   definitionID,
		WorkflowName:   event.WorkflowName,
		InstanceName:   event.InstanceName,
		DefinitionName: event.DefinitionName,
		Version:        event.Version,
		Status:         event.Status,
		ActiveNodeIDs:  event.ActiveNodeIDs,
		TriggerSource:  event.TriggerSource,
		StartedAt:      now,
		Created:        now,
		Updated:        now,
	}

	if event.OrgID != "" {
		orgID, err := model.ToObjectID(event.OrgID)
		if err == nil {
			stub.OrgID = &orgID
		}
	}

	return stub, nil
}

func buildWaitingUpdate(event ports.StateEvent) (repositories.WaitingUpdate, error) {
	if event.InstanceID == "" {
		return repositories.WaitingUpdate{}, fmt.Errorf("missing instanceId")
	}

	return repositories.WaitingUpdate{
		WorkflowUUID:  event.InstanceID,
		Status:        event.Status,
		ActiveNodeIDs: event.ActiveNodeIDs,
		Updated:       time.Now(),
	}, nil
}

func isCreatedEvent(subject string) bool {
	return subject == stateSubjectCreated
}

func isWaitingEvent(subject string) bool {
	return subject == stateSubjectWaiting
}

func isResumedEvent(subject string) bool {
	return subject == stateSubjectResumed
}

func isTerminalEvent(subject string) bool {
	return subject == stateSubjectCompleted ||
		subject == stateSubjectFailed ||
		subject == stateSubjectCancelled
}

func ackBatch(refs []archiverTypes.MsgRef, batchType string) {
	for _, ref := range refs {
		if ref.Batch == batchType {
			ref.Msg.Ack()
		}
	}
}

func nackBatch(refs []archiverTypes.MsgRef, batchType string) {
	for _, ref := range refs {
		if ref.Batch == batchType {
			ref.Msg.Nack(fmt.Errorf("batch %s write failed", batchType))
		}
	}
}
