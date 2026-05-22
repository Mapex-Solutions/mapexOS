package types

import "workflow/src/modules/runtime/domain/entities"

// Branch status constants used by fanout execution.
const (
	BranchCompleted = "completed"
	BranchFailed    = "failed"
	BranchWaiting   = "waiting"
)

// BranchResult holds the outcome of a single fanout branch execution.
// Used by the RuntimeService to collect results from parallel goroutines
// and merge state patches after all branches complete.
type BranchResult struct {
	// Status is the final status of the branch ("completed", "failed", "waiting").
	Status string

	// StatePatch is the delta to merge into the instance state after branch completes.
	StatePatch map[string]interface{}

	// NodeOutputs are the outputs produced by nodes in this branch.
	NodeOutputs map[string]interface{}

	// ExecPath is the execution path entries recorded during branch execution.
	ExecPath []entities.PathEntry

	// NodeState is set when the branch suspends on an async node (contains waitType).
	NodeState map[string]interface{}

	// WaitNodeID is the node ID where the branch suspended.
	WaitNodeID string

	// MergeNodeID is the merge node where the branch converged.
	MergeNodeID string

	// Err is the execution error if the branch failed.
	Err *entities.ExecutionError
}
