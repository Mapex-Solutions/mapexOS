package services

import (
	"sync"

	defPorts "workflow/src/modules/definitions/application/ports"
	"workflow/src/modules/runtime/application/di"
	runtimeEntities "workflow/src/modules/runtime/domain/entities"
	"workflow/src/modules/runtime/domain/executors"
)

// RuntimeService orchestrates workflow execution: creates executions, runs the DAG walker,
// checkpoints state via ExecutionStateRepository, and handles resume signals.
type RuntimeService struct {
	deps        di.RuntimeServiceDependenciesInjection
	registry    *executors.ExecutorRegistry
	activeWalks sync.WaitGroup
}

// preparedWorkflow is the intermediate artifact produced by prepareWorkflow:
// a validated definition + its execution graph + the identified start node.
// Returned to handlers before a WorkflowExecution is constructed.
type preparedWorkflow struct {
	def         *defPorts.WorkflowDefinition
	graph       *runtimeEntities.ExecutionGraph
	startNodeID string
}

// permanentError marks validation failures that must NOT be retried.
// Handlers see this via errors.As and Reject the NATS message (instead of Nack).
type permanentError struct {
	reason string
}

func (e *permanentError) Error() string { return e.reason }
