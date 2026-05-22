package services

import (
	"context"
	"fmt"
	"strings"

	"workflow/src/modules/runtime/application/di"
	"workflow/src/modules/runtime/application/ports"
	"workflow/src/modules/runtime/domain/executors"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time check
var _ ports.RuntimeServicePort = (*RuntimeService)(nil)

// New creates a new RuntimeService with all dependencies injected. The
// constructor builds the per-node executor registry from the injected
// condition evaluator, value resolver, vault, and plugin dependencies.
func New(deps di.RuntimeServiceDependenciesInjection) ports.RuntimeServicePort {
	registry := executors.BuildRegistry(deps.ConditionEvaluator, deps.ValueResolver, deps.VaultService, deps.PluginRepo)
	return &RuntimeService{
		deps:     deps,
		registry: registry,
	}
}

// HandleExecution processes a WORKFLOW-EXECUTION message with multi-mode
// dispatch (newInstance | signal | signalOrStart | subworkflow). The mode is
// read from the payload — never from the subject — so a single stream can
// carry every execution kind.
func (s *RuntimeService) HandleExecution(msg *natsModel.Message) {
	execMsg, ok := s.parseExecutionMessage(msg)
	if !ok {
		return
	}
	s.dispatchExecutionMode(msg, execMsg)
}

// HandleResume continues a WAITING execution: applies the callback's state
// patch / output / signal data, runs error or retry handlers when the
// callback reports failure, and re-enters the DAG walker at the next node.
// Special branches (retry-timer, timeout-without-output, error) terminate
// early; everything else falls through to the normal CAS resume.
func (s *RuntimeService) HandleResume(msg *natsModel.Message) {
	ctx := context.Background()
	resume, ok := s.parseResumeMessage(msg)
	if !ok {
		return
	}
	execution, revision, ok := s.loadResumeExecution(msg, &resume)
	if !ok {
		return
	}
	if !s.validateResumeToken(msg, execution, &resume) {
		return
	}
	s.maybePurgeNonTimeoutSchedule(&resume)
	if s.tryEarlyResumeBranches(ctx, msg, execution, &resume) {
		return
	}
	s.runNormalResume(ctx, msg, execution, revision, &resume)
}

// HandleScheduleFire is a thin re-publisher: a fired NATS schedule already
// carries the resume payload, so the consumer extracts the instanceId and
// forwards the body to mapexos.workflow.resume.timer.{instanceId} where
// HandleResume picks it up.
func (s *RuntimeService) HandleScheduleFire(msg *natsModel.Message) {
	body, instanceId, ok := s.parseScheduleFireBody(msg)
	if !ok {
		return
	}
	s.republishScheduleResume(msg, instanceId, body)
}

// ExecuteByInstanceID starts a workflow execution for the given instance.
// Same lifecycle as a NATS-driven newInstance, but synchronous: the HTTP
// caller receives the result (uuid, status, error info) instead of an Ack.
// Errors are returned to the caller, never logged-and-swallowed.
func (s *RuntimeService) ExecuteByInstanceID(ctx context.Context, instanceID string, eventPayload map[string]interface{}, workflowUUID string) (*ports.ExecuteResult, error) {
	instance, prepared, err := s.prepareInstanceExecution(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	executionUUID := s.resolveExecutionUUID(instance.UniqueExecution, instance.WorkflowUUID, workflowUUID)
	execution := s.buildHttpExecution(instance, prepared, executionUUID, eventPayload)
	if err := s.createAndExecute(ctx, execution, prepared.graph, prepared.startNodeID); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, fmt.Errorf("execution already running with UUID %s", executionUUID)
		}
		return nil, fmt.Errorf("execution failed: %w", err)
	}
	logger.Info(fmt.Sprintf("[SERVICE:Runtime] HTTP Execute → instanceId=%s uuid=%s status=%s",
		instanceID, executionUUID, execution.Status))
	return s.buildExecuteResult(execution), nil
}

// WaitForActiveWalkers blocks until all active execute() calls have returned.
// Called by the shutdown drainer (WalkerDrainer) so in-flight walkers commit
// before the process exits.
func (s *RuntimeService) WaitForActiveWalkers() {
	s.activeWalks.Wait()
}
