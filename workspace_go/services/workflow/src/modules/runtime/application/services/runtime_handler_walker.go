package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	appConstants "workflow/src/modules/runtime/application/constants"
	runtimeTypes "workflow/src/modules/runtime/application/types"
	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"

	"github.com/Mapex-Solutions/mapexGoKit/utils/deepcopy"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// execute is the main DAG walker loop. Per-step checkpoint: after each node completes,
// the result is persisted before advancing to the next node.
func (s *RuntimeService) execute(
	ctx context.Context,
	execution *entities.WorkflowExecution,
	graph *entities.ExecutionGraph,
	startNodeID string,
) error {
	s.activeWalks.Add(1)
	defer s.activeWalks.Done()

	s.deps.Metrics.ExecutionActive.Inc()
	defer s.deps.Metrics.ExecutionActive.Dec()

	stepCount := 0
	defer func() {
		s.deps.Metrics.ExecutionSteps.Observe(float64(stepCount))
	}()

	nc := &runtimeTypes.NodeContext{
		State:          execution.State,
		EventPayload:   execution.EventPayload,
		NodeOutputs:    execution.NodeOutputs,
		NodeStates:     execution.NodeStates,
		ExternalInputs: execution.ExternalInputs,
		Depth:          execution.Depth,
	}

	currentNodeID := startNodeID

	for step := 0; step < constants.MaxInlineSteps; step++ {
		stepCount++
		result, pathEntry, execErr := s.executeStep(ctx, currentNodeID, graph, nc)
		if execErr != nil {
			logger.Debug(fmt.Sprintf("[SERVICE:Runtime] Node %s (%s) failed: code=%s message=%s",
				execErr.NodeID, execErr.NodeType, execErr.Code, execErr.Message))
			execution.ExecutionPath = append(execution.ExecutionPath, pathEntry)

			// Check error handler (retry or error output handle) before failing
			handled, err := s.handleNodeError(ctx, execution, graph, currentNodeID, execErr)
			if err != nil {
				return err
			}
			if handled {
				return nil
			}
			return s.failExecution(execution, execErr)
		}

		// Handle wait → checkpoint + dispatch by node type + publish waiting event
		if result.NodeState != nil && result.NodeState[appConstants.NodeStateKeyWaitType] != nil {
			execution.ExecutionPath = append(execution.ExecutionPath, pathEntry)
			node, ok := graph.GetNode(currentNodeID)
			if !ok {
				return s.failExecution(execution, &entities.ExecutionError{
					Code:      constants.ErrCodeNodeNotFound,
					Message:   fmt.Sprintf("node %s not found in graph", currentNodeID),
					NodeID:    currentNodeID,
					Timestamp: time.Now(),
				})
			}
			if err := s.suspendExecution(execution, currentNodeID, node.Type, result.NodeState); err != nil {
				return err
			}
			return nil
		}

		execution.ExecutionPath = append(execution.ExecutionPath, pathEntry)

		// Apply NodeState (non-wait, e.g., loop counter, merge count, sequence step)
		if result.NodeState != nil {
			execution.NodeStates[currentNodeID] = result.NodeState
		}

		// When a stateful node emits "done", clear its NodeState so it can be reused
		// in future iterations (e.g., sequence or loop inside a parent loop).
		if len(result.OutputHandles) > 0 && result.OutputHandles[0] == constants.OutputHandleDone {
			delete(execution.NodeStates, currentNodeID)
		}

		// Track loop context:
		// - "body": push only on first iteration (skip if already on top of stack)
		// - "done": pop from stack so the walker can find the parent loop
		if node, ok := graph.GetNode(currentNodeID); ok && node.Type == constants.NodeTypeLoop {
			if len(result.OutputHandles) > 0 && result.OutputHandles[0] == constants.OutputHandleBody {
				if peekLoopStack(execution.NodeStates) != currentNodeID {
					pushLoopStack(execution.NodeStates, currentNodeID)
				}
			} else if result.OutputHandles[0] == constants.OutputHandleDone {
				popLoopStack(execution.NodeStates)
			}
		}

		// Apply state patch (nil sentinel = delete key)
		if result.StatePatch != nil {
			for k, v := range result.StatePatch {
				if v == nil {
					delete(execution.State, k)
				} else {
					execution.State[k] = v
				}
			}
		}

		// Apply node output
		if result.NodeOutput != nil {
			execution.NodeOutputs[currentNodeID] = result.NodeOutput
		}

		// No output handles → check loop stack or complete
		if len(result.OutputHandles) == 0 {
			if loopNodeID := peekLoopStack(execution.NodeStates); loopNodeID != "" {
				popLoopStack(execution.NodeStates)
				currentNodeID = loopNodeID
				execution.ActiveNodeIDs = []string{currentNodeID}
				if err := s.checkpoint(execution); err != nil {
					return err
				}
				if s.deps.ShutdownManager.IsShuttingDown() {
					return nil
				}
				continue
			}
			return s.completeOrResuspend(execution)
		}

		// Resolve next nodes
		nextNodes := graph.ResolveNextNodes(currentNodeID, result.OutputHandles)
		if len(nextNodes) == 0 {
			if loopNodeID := peekLoopStack(execution.NodeStates); loopNodeID != "" {
				popLoopStack(execution.NodeStates)
				currentNodeID = loopNodeID
				execution.ActiveNodeIDs = []string{currentNodeID}
				if err := s.checkpoint(execution); err != nil {
					return err
				}
				if s.deps.ShutdownManager.IsShuttingDown() {
					return nil
				}
				continue
			}
			return s.completeOrResuspend(execution)
		}

		// Single next node → advance and continue
		if len(nextNodes) == 1 {
			currentNodeID = nextNodes[0]
			execution.ActiveNodeIDs = []string{currentNodeID}
			if err := s.checkpoint(execution); err != nil {
				return err
			}
			if s.deps.ShutdownManager.IsShuttingDown() {
				s.deps.Metrics.ShutdownStopsTotal.Inc()
				return nil
			}
			continue
		}

		// Multiple outputs — route by node type
		node, ok := graph.GetNode(currentNodeID)
		if !ok {
			return s.failExecution(execution, &entities.ExecutionError{
				Code:      constants.ErrCodeNodeNotFound,
				Message:   fmt.Sprintf("node %s not found in graph", currentNodeID),
				NodeID:    currentNodeID,
				Timestamp: time.Now(),
			})
		}
		switch node.Type {
		case constants.NodeTypeSwitch:
			if err := s.executeSwitchBranches(ctx, execution, graph, nextNodes); err != nil {
				return err
			}
			return nil

		default:
			fanoutMode := constants.FanoutModeWaitAll
			if cfg, ok := graph.ParsedConfigs[currentNodeID].(*entities.FanoutNodeConfig); ok && cfg.Mode != "" {
				fanoutMode = cfg.Mode
			}
			mergeNodeID, err := s.executeFanout(ctx, execution, graph, nextNodes, fanoutMode)
			if err != nil {
				return err
			}
			if execution.Status != entities.ExecStatusRunning {
				return nil
			}
			if mergeNodeID != "" {
				currentNodeID = mergeNodeID
				execution.ActiveNodeIDs = []string{currentNodeID}
				continue
			}
			return s.completeExecution(execution)
		}
	}

	// MaxInlineSteps exceeded → re-enqueue to WORKFLOW-RESUME
	logger.Warn(fmt.Sprintf("[SERVICE:Runtime] MaxInlineSteps exceeded for execution %s, re-enqueueing", execution.WorkflowUUID))
	if err := s.checkpoint(execution); err != nil {
		return err
	}
	if s.deps.ShutdownManager.IsShuttingDown() {
		return nil
	}
	activeNode := ""
	if len(execution.ActiveNodeIDs) > 0 {
		activeNode = execution.ActiveNodeIDs[0]
	}
	return s.deps.RuntimePublisher.PublishResumeMessage(execution.WorkflowUUID, activeNode, constants.ResumeTypeReenqueue)
}

// executeStep resolves and executes a single node. Returns (result, pathEntry, nil) on success
// or (nil, pathEntry, execError) on failure. PathEntry is ALWAYS populated for tracking.
func (s *RuntimeService) executeStep(
	ctx context.Context,
	nodeID string,
	graph *entities.ExecutionGraph,
	nc *runtimeTypes.NodeContext,
) (*entities.NodeExecutionResult, entities.PathEntry, *entities.ExecutionError) {
	emptyPath := entities.PathEntry{NodeID: nodeID, Status: constants.StatusError, EnteredAt: time.Now()}

	node, ok := graph.GetNode(nodeID)
	if !ok {
		return nil, emptyPath, &entities.ExecutionError{
			Code: constants.ErrCodeNodeNotFound, Message: fmt.Sprintf("node %s not found in graph", nodeID),
			NodeID: nodeID, Timestamp: time.Now(),
		}
	}
	emptyPath.NodeType = node.Type

	executor, err := s.registry.Get(node.Type)
	if err != nil {
		return nil, emptyPath, &entities.ExecutionError{
			Code: constants.ErrCodeExecutorNotFound, Message: err.Error(),
			NodeID: nodeID, NodeType: node.Type, Timestamp: time.Now(),
		}
	}

	var nodeTimeout *entities.TimeoutConfig
	if node.Timeout != nil {
		nodeTimeout = &entities.TimeoutConfig{
			Duration:     node.Timeout.Duration,
			Unit:         node.Timeout.Unit,
			EnableOutput: node.Timeout.EnableOutput,
		}
	}

	execCtx := &entities.NodeExecutionContext{
		InstanceID:     nc.InstanceID,
		State:          nc.State,
		EventPayload:   nc.EventPayload,
		NodeOutputs:    nc.NodeOutputs,
		NodeStates:     nc.NodeStates,
		ExternalInputs: nc.ExternalInputs,
		Depth:          nc.Depth,
		NodeID:         nodeID,
		NodeType:       node.Type,
		ParsedConfig:   graph.ParsedConfigs[nodeID],
		Label:          node.Label,
		Timeout:        nodeTimeout,
		Graph:          graph,
		Timezone:       graph.Timezone,
	}

	enteredAt := time.Now()
	result, err := executor.Execute(ctx, execCtx)
	exitedAt := time.Now()

	pathEntry := entities.PathEntry{
		NodeID:     nodeID,
		NodeType:   node.Type,
		Status:     constants.StatusCompleted,
		EnteredAt:  enteredAt,
		ExitedAt:   &exitedAt,
		DurationMs: exitedAt.Sub(enteredAt).Milliseconds(),
	}

	if err != nil {
		pathEntry.Status = constants.StatusError
		errMsg := err.Error()
		pathEntry.Error = &errMsg
		return nil, pathEntry, &entities.ExecutionError{
			Code: constants.ErrCodeExecutionError, Message: err.Error(),
			NodeID: nodeID, NodeType: node.Type, Timestamp: time.Now(),
		}
	}

	if len(result.OutputHandles) > 0 {
		pathEntry.OutputHandle = result.OutputHandles[0]
	}

	if result.Error != nil {
		pathEntry.Status = constants.StatusError
		return nil, pathEntry, result.Error
	}

	if result.NodeState != nil && result.NodeState[appConstants.NodeStateKeyWaitType] != nil {
		pathEntry.Status = constants.StatusWaiting
	}

	return result, pathEntry, nil
}

// executeSwitchBranches executes each branch in sequence (not parallel).
// Used when core/switch matchMode="all" returns multiple output handles.
func (s *RuntimeService) executeSwitchBranches(
	ctx context.Context,
	execution *entities.WorkflowExecution,
	graph *entities.ExecutionGraph,
	branchStartNodes []string,
) error {
	for _, startNodeID := range branchStartNodes {
		if err := s.execute(ctx, execution, graph, startNodeID); err != nil {
			return err
		}
		if execution.Status != entities.ExecStatusRunning {
			return nil
		}
	}
	return nil
}

// executeFanout spawns goroutines per branch with isolated state copies.
func (s *RuntimeService) executeFanout(
	ctx context.Context,
	execution *entities.WorkflowExecution,
	graph *entities.ExecutionGraph,
	branchStartNodes []string,
	fanoutMode string,
) (string, error) {
	n := len(branchStartNodes)
	if n > constants.MaxFanoutBranches {
		return "", s.failExecution(execution, &entities.ExecutionError{
			Code:      constants.ErrCodeMaxFanoutExceeded,
			Message:   fmt.Sprintf("fanout has %d branches (max %d)", n, constants.MaxFanoutBranches),
			NodeID:    execution.ActiveNodeIDs[0],
			Timestamp: time.Now(),
		})
	}

	baseState := deepcopy.Map(execution.State)
	baseOutputs := deepcopy.Map(execution.NodeOutputs)
	baseNodeStates := deepcopy.MapOfMaps(execution.NodeStates)

	results := make([]runtimeTypes.BranchResult, n)
	var wg sync.WaitGroup

	for i, startNode := range branchStartNodes {
		wg.Add(1)
		go func(idx int, startNodeID string) {
			defer wg.Done()
			nc := &runtimeTypes.NodeContext{
				State:          deepcopy.Map(baseState),
				EventPayload:   execution.EventPayload,
				NodeOutputs:    deepcopy.Map(baseOutputs),
				NodeStates:     deepcopy.MapOfMaps(baseNodeStates),
				ExternalInputs: execution.ExternalInputs,
				Depth:          execution.Depth,
			}
			results[idx] = s.executeBranch(ctx, nc, graph, startNodeID)
		}(i, startNode)
	}
	wg.Wait()

	mergeNodeID := ""
	for _, br := range results {
		execution.ExecutionPath = append(execution.ExecutionPath, br.ExecPath...)
		if br.MergeNodeID != "" {
			mergeNodeID = br.MergeNodeID
		}
	}

	allCompleted := true
	for _, br := range results {
		if br.Status == runtimeTypes.BranchFailed {
			return "", s.failExecution(execution, br.Err)
		}
		if br.Status != runtimeTypes.BranchCompleted {
			allCompleted = false
		}
	}

	if !allCompleted {
		waitingNodes := []string{}
		for _, br := range results {
			if br.Status == runtimeTypes.BranchWaiting {
				execution.NodeStates[br.WaitNodeID] = br.NodeState
				waitingNodes = append(waitingNodes, br.WaitNodeID)
			}
		}
		if len(waitingNodes) == 0 {
			return "", s.failExecution(execution, &entities.ExecutionError{
				Code: constants.ErrCodeFanoutInconsistent, Message: "fanout branches neither completed nor waiting",
				Timestamp: time.Now(),
			})
		}
		execution.Status = entities.ExecStatusWaiting
		execution.ActiveNodeIDs = waitingNodes
		for _, br := range results {
			if br.Status == runtimeTypes.BranchCompleted {
				for k, v := range br.StatePatch {
					execution.State[k] = v
				}
				for k, v := range br.NodeOutputs {
					execution.NodeOutputs[k] = v
				}
			}
		}
		if fanoutMode == constants.FanoutModeFirstCompleted {
			execution.NodeStates[appConstants.NodeStateKeyFanoutMeta] = map[string]interface{}{appConstants.NodeStateKeyMode: constants.FanoutModeFirstCompleted}
		}
		if err := s.checkpoint(execution); err != nil {
			return "", err
		}
		s.suspendFanoutExecution(execution)
		return "", nil
	}

	for _, br := range results {
		for k, v := range br.StatePatch {
			execution.State[k] = v
		}
		for k, v := range br.NodeOutputs {
			execution.NodeOutputs[k] = v
		}
	}

	if mergeNodeID != "" {
		if execution.NodeStates == nil {
			execution.NodeStates = make(map[string]map[string]interface{})
		}
		execution.NodeStates[mergeNodeID] = map[string]interface{}{constants.BranchCountKey: n}
	}

	if err := s.checkpoint(execution); err != nil {
		return "", err
	}

	return mergeNodeID, nil
}

// executeBranch runs a single fanout branch with isolated state.
func (s *RuntimeService) executeBranch(
	ctx context.Context,
	nc *runtimeTypes.NodeContext,
	graph *entities.ExecutionGraph,
	startNodeID string,
) runtimeTypes.BranchResult {
	result := runtimeTypes.BranchResult{
		Status:      runtimeTypes.BranchCompleted,
		StatePatch:  make(map[string]interface{}),
		NodeOutputs: make(map[string]interface{}),
	}

	currentNodeID := startNodeID

	for step := 0; step < constants.MaxInlineSteps; step++ {
		node, ok := graph.GetNode(currentNodeID)
		if !ok {
			result.Status = runtimeTypes.BranchFailed
			result.Err = &entities.ExecutionError{
				Code: constants.ErrCodeNodeNotFound, Message: fmt.Sprintf("node %s not found in graph", currentNodeID),
				NodeID: currentNodeID, Timestamp: time.Now(),
			}
			return result
		}
		if node.Type == constants.NodeTypeMerge {
			result.MergeNodeID = currentNodeID
			return result
		}

		execResult, pathEntry, execErr := s.executeStep(ctx, currentNodeID, graph, nc)
		if execErr != nil {
			result.ExecPath = append(result.ExecPath, pathEntry)
			result.Status = runtimeTypes.BranchFailed
			result.Err = execErr
			return result
		}

		if execResult.NodeState != nil && execResult.NodeState[appConstants.NodeStateKeyWaitType] != nil {
			result.ExecPath = append(result.ExecPath, pathEntry)
			result.Status = runtimeTypes.BranchWaiting
			execResult.NodeState[appConstants.NodeStateKeyNodeType] = node.Type
			result.NodeState = execResult.NodeState
			result.WaitNodeID = currentNodeID
			return result
		}

		result.ExecPath = append(result.ExecPath, pathEntry)

		if execResult.NodeState != nil {
			nc.NodeStates[currentNodeID] = execResult.NodeState
		}

		if len(execResult.OutputHandles) > 0 && execResult.OutputHandles[0] == constants.OutputHandleDone {
			delete(nc.NodeStates, currentNodeID)
		}

		if node.Type == constants.NodeTypeLoop && len(execResult.OutputHandles) > 0 {
			if execResult.OutputHandles[0] == constants.OutputHandleBody {
				if peekLoopStack(nc.NodeStates) != currentNodeID {
					pushLoopStack(nc.NodeStates, currentNodeID)
				}
			} else if execResult.OutputHandles[0] == constants.OutputHandleDone {
				popLoopStack(nc.NodeStates)
			}
		}

		if execResult.StatePatch != nil {
			for k, v := range execResult.StatePatch {
				nc.State[k] = v
				result.StatePatch[k] = v
			}
		}

		if execResult.NodeOutput != nil {
			nc.NodeOutputs[currentNodeID] = execResult.NodeOutput
			result.NodeOutputs[currentNodeID] = execResult.NodeOutput
		}

		if len(execResult.OutputHandles) == 0 {
			if loopNodeID := peekLoopStack(nc.NodeStates); loopNodeID != "" {
				popLoopStack(nc.NodeStates)
				currentNodeID = loopNodeID
				continue
			}
			return result
		}

		nextNodes := graph.ResolveNextNodes(currentNodeID, execResult.OutputHandles)
		if len(nextNodes) == 0 {
			if loopNodeID := peekLoopStack(nc.NodeStates); loopNodeID != "" {
				popLoopStack(nc.NodeStates)
				currentNodeID = loopNodeID
				continue
			}
			return result
		}

		currentNodeID = nextNodes[0]
	}

	result.Status = runtimeTypes.BranchFailed
	result.Err = &entities.ExecutionError{
		Code: constants.ErrCodeBranchMaxSteps, Message: fmt.Sprintf("branch exceeded %d steps", constants.MaxInlineSteps),
		NodeID: currentNodeID, Timestamp: time.Now(),
	}
	return result
}
