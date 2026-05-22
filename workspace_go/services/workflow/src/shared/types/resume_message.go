package types

// ResumeMessage is the contract for resuming a waiting workflow instance.
// Published to WORKFLOW-RESUME stream.
// Producers: Instances (SendSignal), NATS Schedule (timer expired), Runtime (re-enqueue).
// Consumer: Runtime (HandleResume).
type ResumeMessage struct {
	InstanceID     string                 `json:"instanceId"`
	NodeID         string                 `json:"nodeId"`
	ExecutionToken string                 `json:"executionToken,omitempty"`
	Status         string                 `json:"status"`
	OutputHandle string                 `json:"outputHandle,omitempty"` // DAG handle to follow: "out", "success", "error", etc. Default: "out"
	Output       interface{}            `json:"output,omitempty"`
	StatePatch   map[string]interface{} `json:"statePatch,omitempty"`
	Error        *ExecutionError        `json:"error,omitempty"`
	SignalData   map[string]interface{} `json:"signalData,omitempty"`
	IsTimeout    bool                   `json:"isTimeout,omitempty"`    // True when resume is from NATS Schedule timer expiry
	EnableOutput bool                   `json:"enableOutput,omitempty"` // When true + isTimeout: route to "timeout" handle instead of failing
}
