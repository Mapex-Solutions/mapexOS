package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	defPorts "workflow/src/modules/definitions/application/ports"
	instancePorts "workflow/src/modules/instances/application/ports"
	appConstants "workflow/src/modules/runtime/application/constants"
	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"
	domainServices "workflow/src/modules/runtime/domain/services"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// prepareWorkflow loads a definition, validates it, builds the execution graph,
// and locates the start node. Returns permanentError for validation failures
// (Reject) and regular errors for transient failures (Nack).
func (s *RuntimeService) prepareWorkflow(ctx context.Context, definitionID string) (*preparedWorkflow, error) {
	def, err := s.deps.DefinitionLoader.GetDefinition(ctx, definitionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get definition %s: %w", definitionID, err)
	}
	if def == nil || !def.Enabled {
		return nil, &permanentError{fmt.Sprintf("definition %s not found or disabled", definitionID)}
	}
	if def.Status != "" && def.Status != string(defPorts.StatusValid) {
		return nil, &permanentError{fmt.Sprintf("definition %s has status %q (missing plugins: %v)", definitionID, def.Status, def.MissingPlugins)}
	}

	graph := domainServices.BuildGraph(def)

	startNodeID := ""
	for _, node := range graph.Nodes {
		if node.Type == constants.NodeTypeStart {
			startNodeID = node.ID
			break
		}
	}
	if startNodeID == "" {
		return nil, &permanentError{fmt.Sprintf("no start node found in definition %s", definitionID)}
	}

	return &preparedWorkflow{def: def, graph: graph, startNodeID: startNodeID}, nil
}

// newBaseExecution creates a WorkflowExecution with all common fields populated
// from the prepared workflow. Callers set mode-specific fields (UUID, instance,
// parent, event payload, etc.) after this call.
func newBaseExecution(prepared *preparedWorkflow) *entities.WorkflowExecution {
	now := time.Now()
	return &entities.WorkflowExecution{
		ID:             model.NewObjectID(),
		DefinitionID:   prepared.def.ID,
		DefinitionName: prepared.def.Name,
		OrgID:          prepared.def.OrgID,
		PathKey:        prepared.def.PathKey,
		Version:        1,
		Status:         entities.ExecStatusRunning,
		ActiveNodeIDs:  []string{prepared.startNodeID},
		NodeStates:     make(map[string]map[string]interface{}),
		State:          instancePorts.InitializeState(prepared.def.States),
		ExecutionPath:  []entities.PathEntry{},
		NodeOutputs:    make(map[string]interface{}),
		StartedAt:      &now,
		Created:        now,
		Updated:        now,
	}
}

// createAndExecute persists the execution to KV, notifies the Archiver,
// and starts the DAG walker.
func (s *RuntimeService) createAndExecute(
	ctx context.Context,
	execution *entities.WorkflowExecution,
	graph *entities.ExecutionGraph,
	startNodeID string,
) error {
	if err := s.deps.ExecutionStateRepo.Create(execution); err != nil {
		return fmt.Errorf("failed to create execution state: %w", err)
	}

	s.deps.Metrics.ExecutionStartedTotal.WithLabelValues(execution.TriggerSource).Inc()

	if err := s.deps.RuntimePublisher.PublishStateEvent(execution, constants.StatusCreated); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Runtime] Failed to publish state event 'created' for execution %s", execution.WorkflowUUID))
	}

	return s.execute(ctx, execution, graph, startNodeID)
}

// handlePrepareError routes errors from prepareWorkflow to Reject (permanent) or Nack (transient).
func handlePrepareError(msg *natsModel.Message, err error) {
	var perm *permanentError
	if errors.As(err, &perm) {
		msg.Reject(perm.reason)
	} else {
		msg.Nack(err)
	}
}

// generateExecutionToken creates a deterministic execution token from workflowUUID, nodeID, and attempt.
// Used to validate that a callback belongs to the current dispatch (prevents stale/duplicate callbacks).
// Returns a 32-char hex string (16 bytes of SHA256).
func generateExecutionToken(workflowUUID, nodeID string, attempt int) string {
	raw := workflowUUID + ":" + nodeID + ":" + strconv.Itoa(attempt)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:16])
}

// buildMsgId creates a deterministic NATS Msg-Id for dispatch deduplication.
// Format: {workflowUUID}:{nodeID}:{attempt}:{loopIndex}
func buildMsgId(workflowUUID, nodeID string, attempt int, state map[string]interface{}) string {
	loopIndex := 0
	if v, ok := state[appConstants.NodeStateKeyLoopIndex].(float64); ok {
		loopIndex = int(v)
	} else if v, ok := state[appConstants.NodeStateKeyLoopIndex].(int); ok {
		loopIndex = v
	}
	return workflowUUID + ":" + nodeID + ":" + strconv.Itoa(attempt) + ":" + strconv.Itoa(loopIndex)
}

// generateSubworkflowUUID creates a deterministic UUIDv5 for a child subworkflow execution.
// Same parent + same node + same attempt always produces the same UUID.
// Uses UUIDv5 (SHA1-based, RFC 4122) with the parent UUID as namespace.
func generateSubworkflowUUID(parentUUID, nodeID string, attempt int) string {
	ns, err := uuid.Parse(parentUUID)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Runtime] Failed to parse parent UUID %q for subworkflow UUIDv5, falling back to random UUID", parentUUID))
		return uuid.New().String()
	}
	name := nodeID + ":" + strconv.Itoa(attempt)
	return uuid.NewSHA1(ns, []byte(name)).String()
}
