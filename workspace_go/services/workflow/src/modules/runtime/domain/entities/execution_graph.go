package entities

import (
	defPorts "workflow/src/modules/definitions/application/ports"
)

// ExecutionGraph is built from a WorkflowDefinition and cached in TieredCache.
// Provides adjacency-based node resolution for the DAG walker.
type ExecutionGraph struct {
	// Adjacency maps nodeId → handleId → targetNodeId
	Adjacency map[string]map[string]string

	// Nodes provides direct lookup by node ID (excludes visual-only nodes)
	Nodes map[string]*defPorts.WorkflowNode

	// GotoPairs maps pairLabel → receiverNodeId for goto resolution
	GotoPairs map[string]string

	// ParsedConfigs maps nodeId → typed config struct (parsed once by GraphBuilder).
	// Executors receive the already-typed config via NodeExecutionContext.ParsedConfig.
	ParsedConfigs map[string]interface{}

	// Timezone resolved from WorkflowDefinition (literal value).
	Timezone string
}

// ResolveNextNodes returns the target node IDs for the given output handles.
// Used by the DAG walker to determine where to go after a node completes.
func (g *ExecutionGraph) ResolveNextNodes(nodeID string, outputHandles []string) []string {
	targets := make([]string, 0, len(outputHandles))
	edges, ok := g.Adjacency[nodeID]
	if !ok {
		return targets
	}
	for _, handle := range outputHandles {
		if target, ok := edges[handle]; ok {
			targets = append(targets, target)
		}
	}
	return targets
}

// GetNode retrieves a node by ID from the graph.
func (g *ExecutionGraph) GetNode(nodeID string) (*defPorts.WorkflowNode, bool) {
	node, ok := g.Nodes[nodeID]
	return node, ok
}

// ResolveDefaultHandle returns the primary success output handle for a node.
// Scans the adjacency map and returns the first handle that is NOT "error" or "timeout".
// Falls back to "out" if no edges exist (safety net for edge-less terminal paths).
func (g *ExecutionGraph) ResolveDefaultHandle(nodeID string) string {
	edges, ok := g.Adjacency[nodeID]
	if !ok || len(edges) == 0 {
		return "out"
	}
	for handle := range edges {
		if handle != "error" && handle != "timeout" {
			return handle
		}
	}
	return "out"
}

// HasEdge checks if a specific edge (node + handle) exists in the adjacency map.
func (g *ExecutionGraph) HasEdge(nodeID string, handle string) bool {
	edges, ok := g.Adjacency[nodeID]
	if !ok {
		return false
	}
	_, ok = edges[handle]
	return ok
}
