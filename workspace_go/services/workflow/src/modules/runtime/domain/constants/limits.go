package constants

/*
 * Domain Execution Limits
 * Used by executors to validate input and prevent resource exhaustion.
 */

const (
	// MaxInlineSteps is the maximum number of inline steps per execution cycle.
	// Prevents infinite loops from consuming unbounded resources.
	MaxInlineSteps = 300

	// MaxSubworkflowDepth is the maximum subworkflow nesting depth.
	MaxSubworkflowDepth = 10

	// MaxLoopIterations is the maximum number of iterations per loop node.
	MaxLoopIterations = 10000

	// MaxFanoutBranches is the maximum number of fanout branches.
	MaxFanoutBranches = 20

	// InlineTimeoutSeconds is the default timeout for inline execution cycles.
	InlineTimeoutSeconds = 30

	// LoopStackKey is the NodeStates key used to track the loop stack (LIFO).
	// When a loop emits "body", its nodeId is pushed. On terminal, popped to return.
	LoopStackKey = "__loop_stack"

	// MaxRetryAttempts is the hard cap on retry attempts regardless of node config.
	MaxRetryAttempts = 10

	// BranchCountKey is the NodeStates key used by the merge executor to count
	// completed parallel branches. Domain vocabulary — shared between the merge
	// executor (writer/reader) and the walker (writer of the initial value).
	BranchCountKey = "branchCount"
)

/*
 * Merge Strategies (domain vocabulary)
 * Determines when a merge executor proceeds.
 */

const (
	MergeStrategyAll   = "all"
	MergeStrategyAny   = "any"
	MergeStrategyFirst = "first"
)

/*
 * Default Output Handle (domain vocabulary)
 * Used when a node emits a single unnamed output edge.
 */

const (
	OutputHandleOut = "out"
)

/*
 * Retry Backoff Interval Units (domain vocabulary)
 * Accepted values for ErrorHandlerConfig.IntervalUnit. Default is seconds
 * when omitted.
 */

const (
	IntervalUnitMinutes = "minutes"
	IntervalUnitHours   = "hours"
)

/*
 * Retry Backoff Limits (domain limits)
 * Hard cap on the computed retry delay regardless of node config.
 */

const (
	// MaxRetryDelaySeconds caps any computed retry delay at one hour to
	// keep schedule timers bounded and avoid runaway exponential backoff.
	MaxRetryDelaySeconds = 3600

	// SecondsPerMinute is used to convert IntervalUnitMinutes to seconds.
	SecondsPerMinute = 60

	// SecondsPerHour is used to convert IntervalUnitHours to seconds.
	SecondsPerHour = 3600
)

/*
 * Node Type Identifiers (domain vocabulary)
 * Identifies the kind of workflow node and drives dispatch routing.
 */

const (
	NodeTypeStart        = "core/start"
	NodeTypeCode         = "core/code"
	NodeTypeSubworkflow  = "core/subworkflow"
	NodeTypeTriggerEvent = "core/trigger_event"
	NodeTypeWaitSignal   = "core/wait_signal"
	NodeTypeWaitFor      = "core/wait_for"
	NodeTypeDelay        = "core/delay"
	NodeTypeLoop         = "core/loop"
	NodeTypeSwitch       = "core/switch"
	NodeTypeMerge        = "core/merge"
	NodeTypeRetry        = "retry"
)

/*
 * Plugin Action Types (domain vocabulary)
 * Routes plugin nodes to the proper external executor.
 */

const (
	ActionTypeHTTP      = "http"
	ActionTypeMQTT      = "mqtt"
	ActionTypeNATS      = "nats"
	ActionTypeEmail     = "email"
	ActionTypeRabbitMQ  = "rabbitmq"
	ActionTypeWebsocket = "websocket"
	ActionTypeScript    = "script"
)

/*
 * Wait Types (domain vocabulary)
 * Identifies why a node is suspended.
 */

const (
	WaitTypeRetryTimer = "retryTimer"
	WaitTypeSignal     = "signal"
)

/*
 * Output Handles (domain vocabulary)
 * Named outputs of structural/control-flow nodes.
 */

const (
	OutputHandleBody    = "body"
	OutputHandleDone    = "done"
	OutputHandleError   = "error"
	OutputHandleTimeout = "timeout"
)

/*
 * Path Entry / State Event Status Values (domain vocabulary)
 * Used for execution path entries and state lifecycle events.
 */

const (
	StatusCreated   = "created"
	StatusRunning   = "running"
	StatusWaiting   = "waiting"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
	StatusError     = "error"
	StatusRetrying  = "retrying"
	StatusTimeout   = "timeout"
	StatusResumed   = "resumed"
)

/*
 * Execution Modes (domain vocabulary)
 * Routing key for incoming WORKFLOW-EXECUTION messages.
 */

const (
	ExecutionModeNewInstance   = "newInstance"
	ExecutionModeSignal        = "signal"
	ExecutionModeSignalOrStart = "signalOrStart"
	ExecutionModeSubworkflow   = "subworkflow"
)

/*
 * Fanout Modes (domain vocabulary)
 * Determines how a fanout merges its branch outcomes.
 */

const (
	FanoutModeWaitAll        = "waitAll"
	FanoutModeFirstCompleted = "firstCompleted"
)

/*
 * Trigger Sources (domain vocabulary)
 * Identifies the origin that started a workflow execution.
 */

const (
	TriggerSourceWorkflow    = "workflow"
	TriggerSourceSubworkflow = "subworkflow"
	TriggerSourceHTTP        = "http"
)

/*
 * Dispatch Types (domain vocabulary)
 * Metric label that classifies how a suspended node was dispatched.
 */

const (
	DispatchTypeCode        = "code"
	DispatchTypeSubworkflow = "subworkflow"
	DispatchTypeTrigger     = "trigger"
	DispatchTypePlugin      = "plugin"
)

/*
 * Dispatch Outcome Labels (domain vocabulary)
 * Metric label for dispatch success/error counters.
 */

const (
	DispatchOutcomeSuccess = "success"
	DispatchOutcomeError   = "error"
)

/*
 * Resume Types (domain vocabulary)
 * Identifies the kind of resume message published to WORKFLOW-RESUME.
 */

const (
	ResumeTypeReenqueue = "reenqueue"
)

/*
 * Execution Error Codes (domain vocabulary)
 * Stable error codes attached to ExecutionError records.
 */

const (
	ErrCodeNodeNotFound       = "NODE_NOT_FOUND"
	ErrCodeExecutorNotFound   = "EXECUTOR_NOT_FOUND"
	ErrCodeExecutionError     = "EXECUTION_ERROR"
	ErrCodeMaxFanoutExceeded  = "MAX_FANOUT_EXCEEDED"
	ErrCodeFanoutInconsistent = "FANOUT_INCONSISTENT"
	ErrCodeBranchMaxSteps     = "BRANCH_MAX_STEPS"
	ErrCodeTimeoutExceeded    = "TIMEOUT_EXCEEDED"
	ErrCodeExternalError      = "EXTERNAL_ERROR"
)
