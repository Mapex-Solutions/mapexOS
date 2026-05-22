package types

// StateEvent is the contract for instance lifecycle events.
// Published to WORKFLOW-STATE stream. Consumed by the Archiver.
// Producers: Runtime (created, waiting, resumed, completed, failed), Instances (cancelled).
type StateEvent struct {
	InstanceID     string   `json:"instanceId"`
	ExecutionId    string   `json:"executionId,omitempty"`
	WorkflowID     string   `json:"workflowId,omitempty"`
	OrgID          string   `json:"orgId,omitempty"`
	WorkflowName   string   `json:"workflowName,omitempty"`
	InstanceName   string   `json:"instanceName,omitempty"`
	DefinitionName string   `json:"definitionName,omitempty"`
	InstanceObjID  string   `json:"instanceObjId,omitempty"`
	Status         string   `json:"status"`
	ActiveNodeIDs  []string `json:"activeNodeIds,omitempty"`
	Version        int      `json:"version,omitempty"`
	TriggerSource  string   `json:"triggerSource,omitempty"`
}
