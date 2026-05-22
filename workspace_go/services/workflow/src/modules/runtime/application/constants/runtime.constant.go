package constants

import (
	"strings"

	runtimeContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/runtime"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/*
 * NATS KV CONFIGURATION
 */

const (
	// KVBucketName is the NATS KV bucket for workflow instance hot state.
	// Stores the full instance state during execution (~1-5KB per instance).
	// Deleted by the Archiver after workflow completion. NATS KV buckets are
	// independent of JetStream streams and stay as a const literal.
	KVBucketName = "WORKFLOW-INSTANCES"

	// KVBucketDescription is a human-readable description for the KV bucket.
	KVBucketDescription = "Workflow instance hot state (execution in progress)"

	// KVReplicas is the number of replicas for the KV bucket.
	// Production: 3 (cluster HA). Standalone: 1.
	KVReplicas = 1
)

/*
 * NATS STREAM CONFIGURATION
 *
 * Stream and subject names resolve at package init from GO_ENV via the
 * mapexGoKit config helpers. Cross-service values are re-exported from
 * packages/contracts/services/workflow/runtime; intra-service values are
 * built directly with config.StreamName / config.Subject.
 */

// workflowSubjectBase is the env-prefixed base used to build per-instance
// pattern subjects (state, resume, signal, etc.) at runtime via fmt.Sprintf.
// Resolved at package init — e.g. "dev.mapexos.workflow".
var workflowSubjectBase = strings.ToLower(config.GetEnv()) + ".mapexos.workflow"

var (
	// ResumeStreamName is the JetStream stream for resume signals.
	// Cross-service contract — re-exported from packages/contracts/services/workflow/runtime.
	ResumeStreamName = runtimeContract.StreamWorkflowResume

	// ResumeSubject is the subject pattern for resume messages.
	// Cross-service contract — re-exported from packages/contracts/services/workflow/runtime.
	ResumeSubject = runtimeContract.SubjectWorkflowResume

	// StateStreamName is the JetStream stream for workflow state lifecycle events.
	// Consumed by the Archiver module to persist state snapshots to MongoDB.
	// Resolved at package init — e.g. "DEV-MAPEXOS-WORKFLOW-STATE".
	StateStreamName = config.StreamName("WORKFLOW", "STATE")

	// StateSubject is the subject pattern for state lifecycle events.
	// Resolved at package init — e.g. "dev.mapexos.workflow.state.>".
	StateSubject = config.Subject("workflow", "state") + ".>"

	// LogsStreamName is the JetStream stream for per-step execution logs.
	// Consumed by the Archiver module to persist logs to ClickHouse. Resolved
	// at package init — e.g. "DEV-MAPEXOS-WORKFLOW-LOGS".
	LogsStreamName = config.StreamName("WORKFLOW", "LOGS")

	// LogsSubject is the subject pattern for execution log events. Resolved
	// at package init — e.g. "dev.mapexos.workflow.logs.>".
	LogsSubject = config.Subject("workflow", "logs") + ".>"

	// ScheduleStreamName is the JetStream stream for NATS message scheduling.
	// Stores scheduled messages that NATS delivers to TargetSubject at the
	// specified time. File storage: schedules survive NATS restarts — no
	// reconciliation needed. Resolved at package init —
	// e.g. "DEV-MAPEXOS-WORKFLOW-SCHEDULE".
	ScheduleStreamName = config.StreamName("WORKFLOW", "SCHEDULE")

	// ScheduleSubjectPrefix is the base subject for schedule messages.
	// Resolved at package init — e.g. "dev.mapexos.workflow.schedule".
	ScheduleSubjectPrefix = config.Subject("workflow", "schedule")

	// ScheduleSubjectPattern is the wildcard pattern for the schedule stream subjects.
	ScheduleSubjectPattern = config.Subject("workflow", "schedule") + ".>"

	// ScheduleTargetSubject is the subject where NATS delivers scheduled
	// messages at fire time. Stays within the workflow schedule stream.
	// A dedicated ScheduleFire consumer picks these up and re-publishes to
	// the resume stream.
	ScheduleTargetSubject = config.Subject("workflow", "schedule.fired")

	// ScheduleFiredSubject is the filter for the ScheduleFire consumer.
	ScheduleFiredSubject = config.Subject("workflow", "schedule.fired")

	// CodeStreamName is the JetStream stream for JS code execution requests.
	// Consumed by the js-workflow-executor service.
	// Cross-service contract — re-exported from packages/contracts/services/workflow/runtime.
	CodeStreamName = runtimeContract.StreamWorkflowJSCode

	// CodeSubject is the subject for code execution requests.
	// Cross-service contract — re-exported from packages/contracts/services/workflow/runtime.
	CodeSubject = runtimeContract.SubjectWorkflowJSCode

	// SignalStreamName is the JetStream stream for external signal delivery.
	// Consumed by the instances module to route signals to waiting instances.
	// Resolved at package init — e.g. "DEV-MAPEXOS-WORKFLOW-SIGNAL".
	SignalStreamName = config.StreamName("WORKFLOW", "SIGNAL")

	// SignalSubject is the subject pattern for signal messages.
	// Subject format: ${env}.mapexos.workflow.signal.{instanceId}
	SignalSubject = config.Subject("workflow", "signal") + ".>"

	// StatePatternSubject builds per-status state subjects.
	// Format: ${env}.mapexos.workflow.state.{status}
	StatePatternSubject = workflowSubjectBase + ".state.%s"

	// ResumeReenqueuePatternSubject builds per-instance reenqueue subjects.
	ResumeReenqueuePatternSubject = workflowSubjectBase + ".resume.reenqueue.%s"

	// ResumeCallbackPatternSubject builds per-instance callback subjects.
	ResumeCallbackPatternSubject = workflowSubjectBase + ".resume.callback.%s"

	// ResumeSignalPatternSubject builds per-instance signal resume subjects.
	ResumeSignalPatternSubject = workflowSubjectBase + ".resume.signal.%s"

	// ResumeTimerPatternSubject builds per-instance timer resume subjects.
	ResumeTimerPatternSubject = workflowSubjectBase + ".resume.timer.%s"

	// ExecutionSubworkflowPatternSubject builds per-instance subworkflow execution subjects.
	ExecutionSubworkflowPatternSubject = workflowSubjectBase + ".execution.subworkflow.%s"

	// TriggerWorkflowExecuteSubject is the fixed subject for workflow trigger dispatch.
	// Cross-service contract — re-exported from packages/contracts/services/workflow/runtime.
	TriggerWorkflowExecuteSubject = runtimeContract.SubjectTriggerWorkflowExecute
)

// CAS (Compare-And-Swap) configuration for concurrent callback handling.
const (
	// MaxCASRetries is the maximum number of CAS retry attempts for concurrent fanout callbacks.
	// After this many retries, the message is Nacked for NATS redelivery.
	MaxCASRetries = 5
)

/*
 * NodeState Map Keys
 * Application-layer keys used inside the per-node state map carried in
 * WorkflowExecution.NodeStates. Stable wire format consumed by the runtime
 * walker, the resume handler and external executors.
 */
const (
	NodeStateKeyExecutionToken = "executionToken"
	NodeStateKeyExpiresAt      = "expiresAt"
	NodeStateKeyWaitType       = "waitType"
	NodeStateKeyEnableOutput   = "enableOutput"
	NodeStateKeyNodeType       = "nodeType"
	NodeStateKeyRetryAttempt   = "retryAttempt"
	NodeStateKeyInternalRetry  = "__retryAttempt"
	NodeStateKeySignalName     = "signalName"
	NodeStateKeyLoopIndex      = "loop_index"
	NodeStateKeyTriggerID      = "triggerId"
	NodeStateKeyPayload        = "payload"
	NodeStateKeyPluginID       = "pluginId"
	NodeStateKeyOperation      = "operation"
	NodeStateKeyAction         = "action"
	NodeStateKeyHooks          = "hooks"
	NodeStateKeyMode           = "mode"
	NodeStateKeyStack          = "stack"
	NodeStateKeyFanoutMeta     = "__fanout_meta"
)

/*
 * Plugin Action Map Keys
 * Application-layer keys used inside the per-node action map carried in
 * NodeState[NodeStateKeyAction]. Stable contract negotiated with the
 * triggers / js-workflow-executor services.
 */
const (
	ActionKeyType = "type"
)

/*
 * Execution Message Data Keys
 * Application-layer keys used inside WorkflowExecutionMessage.Data and
 * inside fired schedule bodies. Stable wire format negotiated with publishers.
 */
const (
	ExecDataKeyInstanceID       = "instanceId"
	ExecDataKeyWorkflowUUID     = "workflowUUID"
	ExecDataKeyEventTrackerID   = "eventTrackerId"
	ExecDataKeyDefinitionID     = "definitionId"
	ExecDataKeyParentInstanceID = "parentInstanceId"
	ExecDataKeyParentNodeID     = "parentNodeId"
	ExecDataKeyDepth            = "depth"
	ExecDataKeyInputData        = "inputData"
	ExecDataKeyCallbackSubject  = "callbackSubject"
	ExecDataKeyExecutionToken   = "executionToken"
	ExecDataKeyRetryAttempt     = "retryAttempt"
	ExecDataKeySignalName       = "signalName"
	ExecDataKeySignalData       = "signalData"
)
