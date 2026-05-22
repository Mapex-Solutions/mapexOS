// Package payloads holds InstanceCreate fixtures for the workflow
// instances module. The connectivity-action journey uses
// SagaSimpleInstance to materialize a runnable workflow instance
// referenced by route groups of kind=workflow.
package payloads

import (
	"encoding/json"
	"fmt"
)

// sagaSimpleInstanceJSON is the verbatim POST /api/v1/instances body
// captured from the platform UI for a "Device Status" instance bound
// to the Device Status definition. The saga keeps the payload literal
// so the shape stays in lockstep with what the operator creates;
// definitionId, definitionVersion, definitionName and name are
// overridden by SagaSimpleInstance / CreateInstance at run time so
// the body lines up with the definition the saga just provisioned
// (id + version) and stays idempotent across runs (name + definitionName).
const sagaSimpleInstanceJSON = `{"definitionId":"69efa95b0bf1256ba56ff9de","definitionVersion":9,"definitionName":"Device Status","name":"Device Status","description":"A status test","pathKey":"","externalInputs":{},"isTemplate":false,"uniqueExecution":false,"workflowUUID":""}`

// SagaSimpleInstance returns the POST /api/v1/instances body the
// connectivity-action journey uses. Only `name` is overridden here
// (with the runID-derived value); `definitionId`, `definitionVersion`
// and `definitionName` are overridden by CreateInstance at run time
// from the bag values published by CreateDefinition.
func SagaSimpleInstance(runID string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaSimpleInstanceJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaSimpleInstance: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-workflow-inst-%s", runID)
	return payload
}
