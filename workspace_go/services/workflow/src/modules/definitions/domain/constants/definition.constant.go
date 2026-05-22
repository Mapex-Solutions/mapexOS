// Package constants holds domain vocabulary for the definitions module.
//
// Node-type identifiers consumed by the validators (config validation,
// cycle detection, code-node extraction) live here as the single source of
// truth for the definitions module. The runtime module keeps its own copy
// for executor dispatch (see runtime/domain/constants/limits.go) — both
// MUST stay in sync because they describe the same workflow DSL.
package constants

// Node Type Identifiers — the workflow definition DSL vocabulary.

const (
	// NodeTypeStart is the entry point of every workflow definition.
	NodeTypeStart = "core/start"

	// NodeTypeEnd is a terminal node (succeeds or fails the run).
	NodeTypeEnd = "core/end"

	// NodeTypeLog emits a log entry from the running workflow.
	NodeTypeLog = "core/log"

	// NodeTypeCode runs user-supplied JavaScript via the JS executor.
	NodeTypeCode = "core/code"

	// NodeTypeCondition evaluates a condition group and routes to true/false outputs.
	NodeTypeCondition = "core/condition"

	// NodeTypeSetState mutates the instance state map.
	NodeTypeSetState = "core/set_state"

	// NodeTypeSwitch evaluates cases and routes to the matched output handle(s).
	NodeTypeSwitch = "core/switch"

	// NodeTypeSubworkflow invokes a child workflow and waits for it to complete.
	NodeTypeSubworkflow = "core/subworkflow"

	// NodeTypeDelay suspends execution for a configured duration.
	NodeTypeDelay = "core/delay"

	// NodeTypeWaitSignal suspends execution until an external signal arrives.
	NodeTypeWaitSignal = "core/wait_signal"

	// NodeTypeLoop iterates over an input source.
	NodeTypeLoop = "core/loop"

	// NodeTypeFanout spawns parallel branches.
	NodeTypeFanout = "core/fanout"

	// NodeTypeMerge joins parallel branches into a single path.
	NodeTypeMerge = "core/merge"

	// NodeTypeSequence walks a fixed number of steps in order.
	NodeTypeSequence = "core/sequence"

	// NodeTypeTriggerEvent fires a trigger event into the platform.
	NodeTypeTriggerEvent = "core/trigger_event"

	// NodeTypeWaitFor pauses the workflow until a field/operator predicate matches.
	NodeTypeWaitFor = "core/wait_for"

	// NodeTypeGoto jumps to the paired goto node by role.
	NodeTypeGoto = "core/goto"

	// NodeTypeTextNote is a visual-only annotation (no execution).
	NodeTypeTextNote = "core/text_note"

	// NodeTypeGroupFrame is a visual-only grouping container (no execution).
	NodeTypeGroupFrame = "core/group_frame"
)

// Node Config Keys — keys read out of WorkflowNode.Config maps.

const (
	// NodeConfigKeyScript is the script source field on a code node.
	NodeConfigKeyScript = "script"
)
