package entities

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

// NodeExecutionContext is the input for all node executors.
// Carries instance state, node info, pre-parsed config, and graph reference.
// Service dependencies (evaluator, resolver, publisher) are injected separately
// by the DAG walker, not embedded in this struct.
type NodeExecutionContext struct {
	// Instance state
	InstanceID     model.ObjectId
	State          map[string]interface{}
	EventPayload   map[string]interface{}
	NodeOutputs    map[string]interface{}
	NodeStates     map[string]map[string]interface{}
	ExternalInputs map[string]interface{}
	Depth          int

	// Node info
	NodeID       string
	NodeType     string
	ParsedConfig interface{} // typed config struct (parsed once by GraphBuilder, e.g., *ConditionNodeConfig)
	Label        string

	// Timeout configuration for this node (nil = use executor default)
	Timeout *TimeoutConfig

	// Graph (for resolving output handles and adjacency)
	Graph *ExecutionGraph

	// Timezone (resolved from definition)
	Timezone string
}

// NodeExecutionResult is the output from all node executors.
// Carries state patches, outputs, control flow handles, and optional suspension or error.
type NodeExecutionResult struct {
	// OutputHandles are the handles to follow in the graph (e.g., "out", "true", "false")
	OutputHandles []string

	// StatePatch is the delta to merge into the instance state (user variables)
	StatePatch map[string]interface{}

	// NodeState is the internal state of this node (loop counter, wait info, merge branches).
	// Stored in instance.NodeStates[nodeId]. If waitType key is present, runtime suspends.
	NodeState map[string]interface{}

	// NodeOutput is the output of the node (code result, subworkflow result)
	NodeOutput interface{}

	// LogEntries are step logs emitted by the node
	LogEntries []LogEntry

	// Error is a structured execution error
	Error *ExecutionError
}
