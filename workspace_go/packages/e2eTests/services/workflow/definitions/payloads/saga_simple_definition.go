// Package payloads holds DefinitionCreate fixtures for the workflow
// definitions module. The canonical SagaSimpleDefinition fixture is
// the smallest DAG the workflow runtime accepts that still exercises
// the start → state → code → end flow the saga's connectivity-action
// journey needs.
package payloads

import (
	"encoding/json"
	"fmt"
)

// sagaSimpleDefinitionJSON is the verbatim POST /api/v1/definitions
// body captured from the platform UI (DevTools → Network) for a
// "Device Status" definition. The saga keeps the payload literal so
// the shape stays in lockstep with what the operator creates; only
// the `name` field is parameterized on runID at call time so multiple
// concurrent saga runs do not collide on the name-unique index.
const sagaSimpleDefinitionJSON = `{"name":"Device Status","description":"Send the current status by Telegram","enabled":true,"isTemplate":false,"timezone":{"type":"literal","value":"UTC"},"retryPolicy":{"enabled":false,"maxAttempts":3,"initialInterval":"1s","backoffMultiplier":2,"maxInterval":"5m","nonRetryableErrors":[]},"states":[{"field":"counter","type":"number","defaultValue":0,"description":"counter","durable":false}],"captureFields":[],"externalInputs":[],"externalSignals":[{"name":"test","description":"sIMLE TEST"}],"nodes":[{"id":"__start__","type":"core/start","label":"Start","position":{"x":180,"y":-30},"config":{},"parentNodeId":""},{"id":"n_core_end_1777314105800","type":"core/end","label":"End","position":{"x":180,"y":300},"config":{"errorCode":"","errorMessage":{"type":"literal","value":""},"terminateWithError":false},"parentNodeId":""},{"id":"n_core_set_state_1778882565241","type":"core/set_state","label":"Set State","position":{"x":180,"y":90},"config":{"operation":"set","targetField":"counter","valueSource":{"type":"literal","value":"1"},"selectedTemplateIds":[]},"parentNodeId":""},{"id":"n_core_code_1778882580665","type":"core/code","label":"Code","position":{"x":150,"y":180},"config":{"script":"// Access: state, event, inputs, nodes\n\nreturn {\n    name: 'Thiago',\n    lastname: 'Anselmo'\n};","timeout":5000},"parentNodeId":""},{"id":"n_core_end_1778882628245_1","type":"core/end","label":"End","position":{"x":285,"y":255},"config":{"errorCode":"ERROR_TEST","errorMessage":{"type":"literal","value":"A simple error test"},"terminateWithError":true},"parentNodeId":""}],"edges":[{"id":"e_f31b93c1-dd51-4266-ad69-38b3c9f5f433","source":"__start__","sourceHandle":"out","target":"n_core_set_state_1778882565241","targetHandle":"in","label":"","pathOffsetX":0,"pathOffsetY":0},{"id":"e_5e89f582-6952-426a-9c11-10aeeba5f854","source":"n_core_set_state_1778882565241","sourceHandle":"out","target":"n_core_code_1778882580665","targetHandle":"in","label":"","pathOffsetX":0,"pathOffsetY":0},{"id":"e_c42a9077-624d-496d-ba03-596bfc8e0570","source":"n_core_code_1778882580665","sourceHandle":"success","target":"n_core_end_1777314105800","targetHandle":"in","label":"","pathOffsetX":0,"pathOffsetY":0},{"id":"e_70d55cf6-0747-4546-bb41-8fdb53d883f9","source":"n_core_code_1778882580665","sourceHandle":"error","target":"n_core_end_1778882628245_1","targetHandle":"in","label":"","pathOffsetX":0,"pathOffsetY":0}],"installedPlugins":["telegram"],"metadata":{"canvasViewport":{"x":-39.8385153215678,"y":41.6038371791167,"zoom":1.4999999999999987}}}`

// SagaSimpleDefinition returns the POST /api/v1/definitions body the
// connectivity-action journey uses. The body is the literal UI capture
// in sagaSimpleDefinitionJSON with `name` rewritten to embed runID
// (so concurrent runs do not collide on Mongo's name-unique index).
func SagaSimpleDefinition(runID string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaSimpleDefinitionJSON), &payload); err != nil {
		// The constant is hand-authored from a UI capture; a parse
		// error means a developer broke the literal during editing
		// and the saga can't recover at runtime.
		panic(fmt.Sprintf("SagaSimpleDefinition: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-workflow-def-%s", runID)
	return payload
}
