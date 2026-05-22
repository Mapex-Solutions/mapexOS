package services

import (
	"context"
	"fmt"
	"math"
	"time"

	defPorts "workflow/src/modules/definitions/application/ports"
	appConstants "workflow/src/modules/runtime/application/constants"
	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// handleNodeError checks the node's ErrorHandler config and decides whether to retry
// or let the error flow to the "error" output handle.
// Returns (true, nil) if the error was handled (retry suspended or continued via error handle).
// Returns (false, nil) if the caller should failExecution.
func (s *RuntimeService) handleNodeError(
	ctx context.Context,
	execution *entities.WorkflowExecution,
	graph *entities.ExecutionGraph,
	nodeID string,
	execErr *entities.ExecutionError,
) (bool, error) {
	node, ok := graph.GetNode(nodeID)
	if !ok {
		return false, nil
	}

	// Check if node has error output handle — if not, fail execution (sync nodes)
	if !graphHasErrorHandle(graph, nodeID) {
		return false, nil
	}

	// Check retry config
	if node.ErrorHandler != nil && node.ErrorHandler.Enabled {
		attempt := getRetryAttempt(execution.NodeStates, nodeID)
		maxAttempts := node.ErrorHandler.MaxAttempts
		if maxAttempts > constants.MaxRetryAttempts {
			maxAttempts = constants.MaxRetryAttempts
		}

		if attempt < maxAttempts {
			logger.Info(fmt.Sprintf("[SERVICE:Runtime] Retrying node %s (attempt %d/%d) in execution %s",
				nodeID, attempt+1, maxAttempts, execution.WorkflowUUID))
			return true, s.suspendForRetry(execution, nodeID, node.Type, node.ErrorHandler, attempt)
		}

		logger.Info(fmt.Sprintf("[SERVICE:Runtime] Retry exhausted for node %s (%d/%d) in execution %s, following error handle",
			nodeID, attempt, maxAttempts, execution.WorkflowUUID))
	}

	// No retry or retries exhausted → follow "error" output handle
	return s.continueViaErrorHandle(ctx, execution, graph, nodeID, execErr)
}

// handleResumeError checks the node's ErrorHandler before failing on async callback errors.
// Same logic as handleNodeError but for resume path.
func (s *RuntimeService) handleResumeError(
	ctx context.Context,
	execution *entities.WorkflowExecution,
	graph *entities.ExecutionGraph,
	nodeID string,
	execErr *entities.ExecutionError,
) (bool, error) {
	return s.handleNodeError(ctx, execution, graph, nodeID, execErr)
}

// suspendForRetry stores retry state in NodeStates and suspends the execution with a timer.
// NATS Schedule delivers the timer resume to WORKFLOW-RESUME at expiry.
func (s *RuntimeService) suspendForRetry(
	execution *entities.WorkflowExecution,
	nodeID string,
	nodeType string,
	eh *defPorts.ErrorHandlerConfig,
	attempt int,
) error {
	delay := calculateRetryDelay(eh, attempt)
	expiresAt := time.Now().Add(delay)

	retryState := map[string]interface{}{
		appConstants.NodeStateKeyWaitType:     constants.WaitTypeRetryTimer,
		appConstants.NodeStateKeyRetryAttempt: attempt + 1,
		appConstants.NodeStateKeyExpiresAt:    expiresAt,
		appConstants.NodeStateKeyNodeType:     nodeType,
	}

	return s.suspendExecution(execution, nodeID, constants.NodeTypeRetry, retryState)
}

// continueViaErrorHandle follows the "error" output handle from the failed node.
// Returns (true, nil) if the error handle exists and execution continues.
// Returns (false, nil) if no error handle exists — caller should failExecution.
func (s *RuntimeService) continueViaErrorHandle(
	ctx context.Context,
	execution *entities.WorkflowExecution,
	graph *entities.ExecutionGraph,
	nodeID string,
	execErr *entities.ExecutionError,
) (bool, error) {
	nextNodes := graph.ResolveNextNodes(nodeID, []string{constants.OutputHandleError})
	if len(nextNodes) == 0 {
		return false, nil
	}

	// Clear any retry state for this node
	delete(execution.NodeStates, nodeID)

	// Continue execution from the error handle target
	execution.ActiveNodeIDs = []string{nextNodes[0]}
	if err := s.checkpoint(execution); err != nil {
		return false, err
	}

	if err := s.execute(ctx, execution, graph, nextNodes[0]); err != nil {
		return true, err
	}
	return true, nil
}

// calculateRetryDelay computes the backoff delay for a given attempt.
// delay = initialInterval * (backoffMultiplier ^ attempt)
func calculateRetryDelay(eh *defPorts.ErrorHandlerConfig, attempt int) time.Duration {
	interval := float64(eh.InitialInterval)
	multiplier := eh.BackoffMultiplier
	if multiplier < 1 {
		multiplier = 1
	}

	delaySeconds := interval * math.Pow(multiplier, float64(attempt))

	// Convert to the configured unit
	switch eh.IntervalUnit {
	case constants.IntervalUnitMinutes:
		delaySeconds *= constants.SecondsPerMinute
	case constants.IntervalUnitHours:
		delaySeconds *= constants.SecondsPerHour
	}

	// Cap at MaxRetryDelaySeconds to keep schedule timers bounded.
	if delaySeconds > constants.MaxRetryDelaySeconds {
		delaySeconds = constants.MaxRetryDelaySeconds
	}

	return time.Duration(delaySeconds) * time.Second
}

// getRetryAttempt reads the current retry attempt from NodeStates.
func getRetryAttempt(nodeStates map[string]map[string]interface{}, nodeID string) int {
	ns := nodeStates[nodeID]
	if ns == nil {
		return 0
	}
	if v, ok := ns[appConstants.NodeStateKeyInternalRetry]; ok {
		switch t := v.(type) {
		case int:
			return t
		case float64:
			return int(t)
		}
	}
	return 0
}

// graphHasErrorHandle checks if a node has an "error" output edge in the graph.
func graphHasErrorHandle(graph *entities.ExecutionGraph, nodeID string) bool {
	return graph.HasEdge(nodeID, constants.OutputHandleError)
}
