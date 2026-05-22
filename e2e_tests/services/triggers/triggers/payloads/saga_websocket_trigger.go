package payloads

import (
	"encoding/json"
	"fmt"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
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

// SagaWebsocketTrigger returns the POST /api/v1/triggers body for
// the WebSocket smoke. The URL points at the in-process HTTP sink's
// /ws path (a separate sink will be added when the assert needs to
// validate frames; until then the events-trigger oracle is enough
// since success=true requires the WS handshake to complete).
func SagaWebsocketTrigger(runID string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaWebsocketTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaWebsocketTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-websocket-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	wsCfg, _ := cfg["websocket"].(map[string]any)
	wsCfg["url"] = constants.WsSinkURL
	wsCfg["message"] = map[string]any{"saga": "websocket", "runID": runID}
	return payload
}
