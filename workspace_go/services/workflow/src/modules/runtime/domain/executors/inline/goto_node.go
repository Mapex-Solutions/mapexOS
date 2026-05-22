package inline

import (
	"context"
	"fmt"
	"time"

	"workflow/src/modules/runtime/domain/entities"
)

// GotoExecutor handles GoTo portal nodes (sender and receiver roles).
// Sender nodes are routed to their matching receiver via GraphBuilder adjacency injection.
// Receiver nodes continue to the next user-drawn edge.
type GotoExecutor struct{}

// NewGotoExecutor creates a new GotoExecutor.
func NewGotoExecutor() entities.NodeExecutor {
	return &GotoExecutor{}
}

// NodeType returns the node type identifier for GoTo nodes.
func (e *GotoExecutor) NodeType() string {
	return "core/goto"
}

// Execute processes a GoTo node.
// Sender: validates that GraphBuilder injected an edge to the matching receiver,
// returning an error if no receiver is reachable (orphaned sender).
// Receiver: passthrough — continues to the next node via user-drawn edges.
func (e *GotoExecutor) Execute(_ context.Context, ctx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := ctx.ParsedConfig.(*entities.GotoNodeConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config for goto node %s: expected *GotoNodeConfig", ctx.NodeID)
	}

	// Sender: validate that GraphBuilder injected an edge to receiver
	if cfg.Role == "sender" {
		if !ctx.Graph.HasEdge(ctx.NodeID, "out") {
			return &entities.NodeExecutionResult{
				Error: &entities.ExecutionError{
					Code:      "GOTO_NO_RECEIVER",
					Message:   fmt.Sprintf("goto sender '%s' has no matching receiver (pairLabel: %s)", ctx.Label, cfg.PairLabel),
					NodeID:    ctx.NodeID,
					NodeType:  ctx.NodeType,
					Timestamp: time.Now(),
				},
			}, nil
		}
	}

	return &entities.NodeExecutionResult{
		OutputHandles: []string{"out"},
	}, nil
}
