package entities

import "context"

// NodeExecutor is the domain contract for all workflow node executors.
// Lives in entities to avoid circular imports between executors/ and its sub-packages.
type NodeExecutor interface {
	// Execute runs the node logic and returns an execution result.
	Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error)

	// NodeType returns the node type identifier (e.g., "core/condition").
	NodeType() string
}
