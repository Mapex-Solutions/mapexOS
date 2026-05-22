package dagwalker

import (
	sharedTypes "workflow/src/shared/types"
)

// CallbackFunc generates a ResumeMessage when an async node suspends.
// Receives nodeID and the NodeState set by the executor.
type CallbackFunc func(nodeID string, nodeState map[string]interface{}) sharedTypes.ResumeMessage

// SuccessCallback returns a callback that resumes with success and the given output.
func SuccessCallback(output interface{}) CallbackFunc {
	return func(nodeID string, _ map[string]interface{}) sharedTypes.ResumeMessage {
		return sharedTypes.ResumeMessage{
			NodeID: nodeID,
			Status: "success",
			Output: output,
		}
	}
}

// CodeSuccessCallback returns a callback for code nodes with output and optional statePatch.
func CodeSuccessCallback(output map[string]interface{}, statePatch map[string]interface{}) CallbackFunc {
	return func(nodeID string, _ map[string]interface{}) sharedTypes.ResumeMessage {
		return sharedTypes.ResumeMessage{
			NodeID:     nodeID,
			Status:     "success",
			Output:     output,
			StatePatch: statePatch,
		}
	}
}

// CodeErrorCallback returns a callback for code nodes that fail.
func CodeErrorCallback(code, message string) CallbackFunc {
	return func(nodeID string, _ map[string]interface{}) sharedTypes.ResumeMessage {
		return sharedTypes.ResumeMessage{
			NodeID: nodeID,
			Status: "error",
			Error: &sharedTypes.ExecutionError{
				Code:    code,
				Message: message,
			},
		}
	}
}

// ErrorCallback returns a generic error callback.
func ErrorCallback(code, message string) CallbackFunc {
	return func(nodeID string, _ map[string]interface{}) sharedTypes.ResumeMessage {
		return sharedTypes.ResumeMessage{
			NodeID: nodeID,
			Status: "error",
			Error: &sharedTypes.ExecutionError{
				Code:    code,
				Message: message,
			},
		}
	}
}

// TimeoutCallback returns a callback that simulates timer expiry.
func TimeoutCallback(enableOutput bool) CallbackFunc {
	return func(nodeID string, _ map[string]interface{}) sharedTypes.ResumeMessage {
		return sharedTypes.ResumeMessage{
			NodeID:       nodeID,
			Status:       "timeout",
			IsTimeout:    true,
			EnableOutput: enableOutput,
		}
	}
}

// PluginSuccessCallback returns a callback for plugin nodes.
func PluginSuccessCallback(output interface{}) CallbackFunc {
	return SuccessCallback(output)
}

// SubworkflowSuccessCallback returns a callback for subworkflow nodes.
func SubworkflowSuccessCallback(output interface{}) CallbackFunc {
	return SuccessCallback(output)
}

// DefaultAsyncCallback returns a generic success callback with no output.
func DefaultAsyncCallback() CallbackFunc {
	return SuccessCallback(nil)
}

// SignalCallback returns a callback that simulates a signal resume with data.
func SignalCallback(signalData map[string]interface{}) CallbackFunc {
	return func(nodeID string, _ map[string]interface{}) sharedTypes.ResumeMessage {
		return sharedTypes.ResumeMessage{
			NodeID:     nodeID,
			Status:     "success",
			SignalData: signalData,
		}
	}
}
