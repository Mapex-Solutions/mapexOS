package services

import (
	"strings"

	defPorts "workflow/src/modules/definitions/application/ports"
	"workflow/src/modules/runtime/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// BuildGraph creates an ExecutionGraph from a WorkflowDefinition.
// The graph is cached in TieredCache and reused across instance executions.
//
// Algorithm:
//
//	Filter visual-only nodes (text_note, group_frame)
//	Index nodes by ID
//	Resolve goto pairs (receiver by pairLabel)
//	Build adjacency from edges (skip visual handles with __ prefix)
//	Inject goto sender→receiver as logical edges
//	Parse node configs into typed structs (cost paid once, cached in graph)
//	Resolve timezone from definition
func BuildGraph(def *defPorts.WorkflowDefinition) *entities.ExecutionGraph {
	graph := &entities.ExecutionGraph{
		Adjacency:     make(map[string]map[string]string),
		Nodes:         make(map[string]*defPorts.WorkflowNode),
		GotoPairs:     make(map[string]string),
		ParsedConfigs: make(map[string]interface{}),
	}

	// Index nodes — filter visual-only types
	for i := range def.Nodes {
		node := &def.Nodes[i]
		if isVisualOnly(node.Type) {
			continue
		}
		graph.Nodes[node.ID] = node
	}

	// Resolve goto pairs (receiver by pairLabel)
	for _, node := range graph.Nodes {
		if node.Type != "core/goto" {
			continue
		}
		cfg := parseGotoConfig(node.Config)
		if cfg.role == "receiver" && cfg.pairLabel != "" {
			graph.GotoPairs[cfg.pairLabel] = node.ID
		}
	}

	// Build adjacency from edges
	for i := range def.Edges {
		edge := &def.Edges[i]

		// Skip edges involving visual-only nodes
		if _, ok := graph.Nodes[edge.Source]; !ok {
			continue
		}
		if _, ok := graph.Nodes[edge.Target]; !ok {
			continue
		}

		handle := edge.SourceHandle
		if handle == "" {
			handle = "out"
		}

		// Skip visual handles (prefixed with __)
		if strings.HasPrefix(handle, "__") {
			continue
		}

		if graph.Adjacency[edge.Source] == nil {
			graph.Adjacency[edge.Source] = make(map[string]string)
		}
		graph.Adjacency[edge.Source][handle] = edge.Target
	}

	// Inject goto sender → receiver as logical edges
	for _, node := range graph.Nodes {
		if node.Type != "core/goto" {
			continue
		}
		cfg := parseGotoConfig(node.Config)
		if cfg.role != "sender" {
			continue
		}
		receiverID, ok := graph.GotoPairs[cfg.pairLabel]
		if !ok {
			continue
		}
		if graph.Adjacency[node.ID] == nil {
			graph.Adjacency[node.ID] = make(map[string]string)
		}
		graph.Adjacency[node.ID]["out"] = receiverID
	}

	// Parse node configs into typed structs (cost paid once, cached in graph)
	for _, node := range graph.Nodes {
		if parsed := parseNodeConfig(node.Type, node.Config); parsed != nil {
			graph.ParsedConfigs[node.ID] = parsed
		}
	}

	// Resolve timezone from definition (literal only at graph-build time)
	if def.Timezone.Type == defPorts.FieldValueLiteral && def.Timezone.Value != "" {
		graph.Timezone = def.Timezone.Value
	}

	return graph
}

// parseGotoConfig extracts role and pairLabel from a goto node's raw config map.
func parseGotoConfig(config map[string]interface{}) gotoConfig {
	return gotoConfig{
		role:      model.MapGetString(config, "role"),
		pairLabel: model.MapGetString(config, "pairLabel"),
	}
}

// isVisualOnly returns true for node types that are visual-only (not part of execution).
func isVisualOnly(nodeType string) bool {
	return nodeType == "core/text_note" || nodeType == "core/group_frame"
}
