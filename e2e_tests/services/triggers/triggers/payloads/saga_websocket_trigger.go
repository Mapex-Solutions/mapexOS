package payloads

import (
	"encoding/json"
	"fmt"
)

const sagaWebsocketTriggerJSON = `{
  "name": "Saga WebSocket Trigger",
  "description": "smoke test — opens a WS connection to the HTTP sink upgrade endpoint",
  "triggerType": "websocket",
  "category": "technical",
  "enabled": true,
  "isSystem": false,
  "isTemplate": false,
  "config": {
    "websocket": {
      "url": "PLACEHOLDER",
      "message": { "saga": "websocket", "runID": "PLACEHOLDER" }
    }
  }
}`

// SagaWebsocketTrigger returns the POST /api/v1/triggers body for the
// WebSocket smoke. The URL is rewritten to wsURL — the in-process WS
// sink's /ws upgrade endpoint on its ephemeral port. Pure: the create
// step reads the sink's ephemeral address from the bag and passes it in.
func SagaWebsocketTrigger(runID, wsURL string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaWebsocketTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaWebsocketTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-websocket-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	wsCfg, _ := cfg["websocket"].(map[string]any)
	wsCfg["url"] = wsURL
	wsCfg["message"] = map[string]any{"saga": "websocket", "runID": runID}
	return payload
}
