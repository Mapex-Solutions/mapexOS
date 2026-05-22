package message

import (
	runtimeContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/runtime"
	sharedTypes "workflow/src/shared/types"
)

/*
 * NATS MESSAGE TYPES
 * Serialized as JSON when published/consumed via NATS JetStream streams.
 *
 * Shared types (StateEvent, ResumeMessage, ExecutionError) live in src/shared/types/.
 * Cross-service payloads (CodeExecutionRequest, WorkflowTriggerRequest) live in
 * packages/contracts/services/workflow/runtime/.
 * This file exposes them to the consumer layer via aliases.
 */

// StateEvent is an alias to the shared type. Published to WORKFLOW-STATE stream.
type StateEvent = sharedTypes.StateEvent

// ResumeMessage is an alias to the shared type. Published to WORKFLOW-RESUME stream.
type ResumeMessage = sharedTypes.ResumeMessage

// CodeExecutionRequest is the cross-service payload consumed by js-workflow-executor.
type CodeExecutionRequest = runtimeContract.CodeExecutionRequest

// WorkflowTriggerRequest is the cross-service payload consumed by the triggers service.
type WorkflowTriggerRequest = runtimeContract.WorkflowTriggerRequest
