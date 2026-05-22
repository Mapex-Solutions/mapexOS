package entities

import "errors"

// Sentinel errors for execution limit violations.
// Used by executors and the DAG walker to enforce safety boundaries.
var (
	// ErrExecutorNotFound is returned when no executor is registered for a node type.
	ErrExecutorNotFound = errors.New("executor not found for node type")
	// ErrMaxInlineSteps is returned when the DAG walker exceeds the maximum inline steps.
	ErrMaxInlineSteps = errors.New("max inline steps exceeded")
	// ErrMaxLoopIterations is returned when a loop node exceeds the maximum iterations.
	ErrMaxLoopIterations = errors.New("max loop iterations exceeded")
	// ErrMaxFanoutBranches is returned when a fanout node exceeds the maximum branches.
	ErrMaxFanoutBranches = errors.New("max fanout branches exceeded")
	// ErrMaxSubworkflowDepth is returned when subworkflow nesting exceeds the maximum depth.
	ErrMaxSubworkflowDepth = errors.New("max subworkflow depth exceeded")
)
