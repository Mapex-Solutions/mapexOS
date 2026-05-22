package types

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

/**
 * NodeContext holds the working state maps for step execution.
 * Main walker populates from instance; branch walker from isolated copies.
 */
type NodeContext struct {
	// InstanceID is the ID of the workflow instance being executed.
	InstanceID model.ObjectId

	// State is the mutable workflow state map.
	State map[string]interface{}

	// EventPayload is the immutable event data that triggered the workflow.
	EventPayload map[string]interface{}

	// NodeOutputs holds outputs produced by previously executed nodes.
	NodeOutputs map[string]interface{}

	// NodeStates holds per-node internal state (loop counter, wait info, merge branches).
	NodeStates map[string]map[string]interface{}

	// ExternalInputs holds external input values provided at trigger time.
	ExternalInputs map[string]interface{}

	// Depth is the subworkflow recursion depth.
	Depth int
}
