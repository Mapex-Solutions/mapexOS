package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"workflow/src/modules/definitions/application/constants"
	"workflow/src/modules/definitions/application/ports"
	domainConstants "workflow/src/modules/definitions/domain/constants"
	"workflow/src/modules/definitions/domain/entities"

	contractDefs "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/definitions"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/* PRIVATE METHODS — MinIO L2 CACHE MANAGEMENT */

// writeNodeScripts writes all code node scripts to MinIO (L2 cache).
// Called on definition create to populate L2 for js-workflow-executor.
// Returns error on first write failure — scripts are critical for code node execution.
func (s *DefinitionService) writeNodeScripts(ctx context.Context, def *entities.WorkflowDefinition) error {
	orgId := extractOrgId(def)
	defId := def.ID.Hex()

	for _, node := range def.Nodes {
		if node.Type != domainConstants.NodeTypeCode {
			continue
		}
		script := getNodeScript(node)
		if script == "" {
			continue
		}
		if err := s.deps.DefinitionStoragePort.WriteNodeScript(ctx, orgId, defId, node.ID, []byte(script)); err != nil {
			return fmt.Errorf("failed to write script for node %s in definition %s: %w", node.ID, defId, err)
		}
	}
	return nil
}

// deleteAllNodeData removes all scripts and bytecodes for a definition's code nodes.
func (s *DefinitionService) deleteAllNodeData(ctx context.Context, def *entities.WorkflowDefinition, nodeIds []string) {
	if len(nodeIds) == 0 {
		return
	}
	orgId := extractOrgId(def)
	if err := s.deps.DefinitionStoragePort.DeleteAllNodeData(ctx, orgId, def.ID.Hex(), nodeIds); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Definition] Failed to delete node data for %s: %v", def.ID.Hex(), err))
	}
}

// publishDefinitionInvalidate publishes a FANOUT message for cache invalidation.
// Consumers (js-workflow-executor) invalidate L0 + L1 for the specified nodeIds.
func (s *DefinitionService) publishDefinitionInvalidate(ctx context.Context, def *entities.WorkflowDefinition, nodeIds []string) {
	if def == nil || len(nodeIds) == 0 {
		return
	}

	orgId := extractOrgId(def)

	payload := ports.DefinitionInvalidatePayload{
		OrgId:        orgId,
		DefinitionId: def.ID.Hex(),
		NodeIds:      nodeIds,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Definition] Failed to marshal FANOUT payload for %s: %v", def.ID.Hex(), err))
		return
	}

	if err := s.deps.NatsBus.PublishFanout(ctx, constants.FanoutDefinitionSubject, data); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Definition] Failed to publish FANOUT for %s: %v", def.ID.Hex(), err))
	}
}

/* PRIVATE HELPERS — CODE NODE DIFF */

// extractOrgId returns the hex string of the definition's OrgID, or empty string.
func extractOrgId(def *entities.WorkflowDefinition) string {
	if def.OrgID != nil {
		return def.OrgID.Hex()
	}
	return ""
}

// getCodeNodeIds returns the IDs of all code nodes in a definition.
func getCodeNodeIds(def *entities.WorkflowDefinition) []string {
	var ids []string
	for _, node := range def.Nodes {
		if node.Type == domainConstants.NodeTypeCode {
			ids = append(ids, node.ID)
		}
	}
	return ids
}

// diffCodeNodes compares the code-node sets between two workflow
// definitions and returns three slices:
//   - added:    code-node ids present in `after` but not in `before`
//   - removed:  code-node ids present in `before` but not in `after`
//   - modified: code-node ids present in both whose script changed
//
// Non-code nodes are ignored. Used by Update to drive bytecode cache
// invalidation + re-compile scheduling.
func diffCodeNodes(before, after *entities.WorkflowDefinition) (added, removed, modified []string) {
	beforeMap := indexCodeNodesById(before)
	afterMap := indexCodeNodesById(after)

	for id := range afterMap {
		if _, existed := beforeMap[id]; !existed {
			added = append(added, id)
		}
	}
	for id := range beforeMap {
		if _, stillThere := afterMap[id]; !stillThere {
			removed = append(removed, id)
		}
	}
	for id, beforeNode := range beforeMap {
		afterNode, stillThere := afterMap[id]
		if !stillThere {
			continue
		}
		if getNodeScript(beforeNode) != getNodeScript(afterNode) {
			modified = append(modified, id)
		}
	}
	return added, removed, modified
}

// indexCodeNodesById builds an id -> node map for all code-typed nodes
// in a definition. Non-code nodes are skipped.
func indexCodeNodesById(def *entities.WorkflowDefinition) map[string]entities.WorkflowNode {
	out := make(map[string]entities.WorkflowNode)
	if def == nil {
		return out
	}
	for _, node := range def.Nodes {
		if node.Type == domainConstants.NodeTypeCode {
			out[node.ID] = node
		}
	}
	return out
}

// getNodeScript extracts the script string from a code node's config map.
func getNodeScript(node entities.WorkflowNode) string {
	return model.MapGetString(node.Config, domainConstants.NodeConfigKeyScript)
}

// findNode finds a node by ID in a node slice.
func findNode(nodes []entities.WorkflowNode, nodeId string) (entities.WorkflowNode, bool) {
	for _, node := range nodes {
		if node.ID == nodeId {
			return node, true
		}
	}
	return entities.WorkflowNode{}, false
}

/* PRIVATE HELPERS — CONTRACT → ENTITY CONVERSION */

// contractNodesToEntity converts contract WorkflowNode slice to entity WorkflowNode slice.
func contractNodesToEntity(nodes []contractDefs.WorkflowNode) []entities.WorkflowNode {
	result := make([]entities.WorkflowNode, len(nodes))
	for i, n := range nodes {
		result[i] = entities.WorkflowNode{
			ID:           n.ID,
			Type:         n.Type,
			Label:        n.Label,
			Position:     entities.Position{X: n.Position.X, Y: n.Position.Y},
			Config:       n.Config,
			ParentNodeID: n.ParentNodeID,
		}
	}
	return result
}

// contractEdgesToEntity converts contract WorkflowEdge slice to entity WorkflowEdge slice.
func contractEdgesToEntity(edges []contractDefs.WorkflowEdge) []entities.WorkflowEdge {
	result := make([]entities.WorkflowEdge, len(edges))
	for i, e := range edges {
		result[i] = entities.WorkflowEdge{
			ID:           e.ID,
			Source:       e.Source,
			SourceHandle: e.SourceHandle,
			Target:       e.Target,
			TargetHandle: e.TargetHandle,
			Label:        e.Label,
			PathOffsetX:  e.PathOffsetX,
			PathOffsetY:  e.PathOffsetY,
		}
	}
	return result
}

// trackDefinitionMetrics is the shared metric-recording shim used by every
// definition operation: increments the per-(op, outcome) counter and
// observes the duration histogram. Centralising it keeps every public method
// readable and prevents drift in metric label values.
func (s *DefinitionService) trackDefinitionMetrics(op, outcome string, start time.Time) {
	s.deps.Metrics.DefinitionOperations.WithLabelValues(op, outcome).Inc()
	s.deps.Metrics.DefinitionOperationDuration.WithLabelValues(op).Observe(time.Since(start).Seconds())
}

// notFoundDefinition wraps the canonical 404 envelope used by every
// definition endpoint, keeping the message consistent across CRUD operations.
func notFoundDefinition() error {
	return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Workflow definition not found"}}
}

// scriptUploadFailedError wraps a script upload failure with a stable
// message used by both Create and Update so HTTP clients can match it.
func scriptUploadFailedError(op string, err error) error {
	return fmt.Errorf("definition %sd but script upload failed: %w", op, err)
}
