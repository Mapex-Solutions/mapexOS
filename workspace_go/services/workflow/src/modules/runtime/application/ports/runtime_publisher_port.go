package ports

import (
	"time"

	"workflow/src/modules/runtime/domain/entities"
	sharedTypes "workflow/src/shared/types"
)

// RuntimePublisherPort abstracts all NATS JetStream publishing operations
// for the runtime module. Implementations handle message construction,
// subject formatting, and the actual publish call.
type RuntimePublisherPort interface {
	// PublishStateEvent notifies the Archiver about execution lifecycle transitions.
	// Status values: "created", "waiting", "resumed", "completed", "failed", "cancelled".
	PublishStateEvent(execution *entities.WorkflowExecution, status string) error

	// PublishResumeMessage publishes a resume/re-enqueue message to WORKFLOW-RESUME.
	PublishResumeMessage(executionID string, nodeID string, status string) error

	// PublishResumeTimer re-publishes a raw fired-schedule body to WORKFLOW-RESUME on
	// the per-instance timer subject workflow.resume.timer.{instanceId}. Used by the
	// schedule-fire handler to forward NATS Schedule deliveries to HandleResume.
	PublishResumeTimer(instanceID string, body map[string]interface{}) error

	// DispatchCodeExecution sends a code execution request to WORKFLOW-JS-CODE stream.
	// executionToken is included in payload for callback validation; msgId is the Nats-Msg-Id header for dedup.
	DispatchCodeExecution(execution *entities.WorkflowExecution, nodeID string, nodeState map[string]interface{}, executionToken, msgId string) error

	// DispatchSubworkflowExecution publishes a subworkflow execution request to WORKFLOW-EXECUTION stream
	// with mode=subworkflow. The child workflow is processed by HandleExecution → handleSubworkflow.
	DispatchSubworkflowExecution(execution *entities.WorkflowExecution, nodeID string, nodeState map[string]interface{}, executionToken, msgId string) error

	// DispatchWorkflowTrigger sends a workflow trigger request to the Triggers Service.
	// Subject: trigger.WORKFLOW.execute
	// mode "trigger": uses a registered trigger entity (triggerId in data)
	// mode "plugin": uses a fully resolved plugin action pipeline (action + hooks in data)
	DispatchWorkflowTrigger(execution *entities.WorkflowExecution, nodeID string, mode string, data map[string]interface{}, executionToken, msgId string) error

	// PublishSignalResume publishes a resume message with signal data to WORKFLOW-RESUME.
	// Used when a signal is delivered to a waiting execution.
	PublishSignalResume(executionID string, nodeID string, signalData map[string]interface{}) error

	// PublishCallbackResume publishes a resume message to a specific callback subject.
	// Used by subworkflow child to notify parent on terminal state.
	PublishCallbackResume(subject string, resume sharedTypes.ResumeMessage) error

	// PublishSchedule publishes a NATS scheduled message for a timer.
	// The message body varies by waitType:
	//   - "timer" (delay): normal resume {instanceId, nodeId, status: "completed", outputHandle: "out"}
	//   - "retryTimer": {instanceId, nodeId, isTimeout: true}
	//   - "callback"/"signal"/"condition": {instanceId, nodeId, isTimeout: true, enableOutput: <bool>}
	// Subject: workflow.schedule.{wfUUID}.{nodeID} (one schedule per node per execution).
	PublishSchedule(wfUUID, nodeID string, expiresAt time.Time, waitType string, enableOutput bool) error

	// PurgeSchedule cancels a pending schedule for a specific node.
	// Idempotent: returns nil if the schedule already fired or was never published.
	PurgeSchedule(wfUUID, nodeID string) error

	// PurgeAllSchedules cancels all pending schedules for an entire workflow execution.
	// Called by failExecution and completeExecution for explicit cleanup.
	PurgeAllSchedules(wfUUID string) error
}
