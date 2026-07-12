// Package payloads holds CreateTrigger fixtures for the triggers
// module. The connectivity-action journey uses SagaSimpleTrigger to
// register an HTTP-kind trigger the route groups of kind=trigger
// reference and the saga later observes through /api/v1/events/trigger.
package payloads

import (
	"encoding/json"
	"fmt"
)

// sagaSimpleTriggerJSON is the verbatim POST /api/v1/triggers body
// captured from the platform UI for a "Generic HTTP" trigger. The
// saga keeps the payload literal so the shape stays in lockstep with
// what the operator creates; only `name` and the http endpoint URL
// are overridden at run time.
//
// The captured URL (https://localhost:1010) is overridden at run
// time with the sinkURL the create step passes in — the address of
// the local HTTP sink the journey stands up, on its ephemeral port.
const sagaSimpleTriggerJSON = `{"name":"Generic HTTP","description":"Test","triggerType":"http","category":"technical","enabled":true,"isSystem":false,"isTemplate":false,"config":{"http":{"endpoint":"https://localhost:1010","method":"POST","headers":{},"body":{"name":"test generic http"},"timeout":30000}}}`

// SagaSimpleTrigger returns the POST /api/v1/triggers body the
// connectivity-action journey uses. The endpoint URL is rewritten to
// sinkURL (the saga-managed test sink) so the trigger has a real HTTP
// listener to POST against (and the resulting events_trigger row
// carries success=true). Pure: the create step reads the sink's
// ephemeral address from the bag and passes it in.
func SagaSimpleTrigger(runID, sinkURL string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaSimpleTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaSimpleTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-trigger-%s", runID)

	// Override the captured HTTP endpoint with the saga-managed sink.
	cfg, _ := payload["config"].(map[string]any)
	if httpCfg, ok := cfg["http"].(map[string]any); ok {
		httpCfg["endpoint"] = sinkURL
	}
	return payload
}
