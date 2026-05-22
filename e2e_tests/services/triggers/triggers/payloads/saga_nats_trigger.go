package payloads

import (
	"encoding/json"
	"fmt"
)

const sagaNatsTriggerJSON = `{
  "name": "Saga NATS Trigger",
  "description": "smoke test — publishes to the saga-managed NATS",
  "triggerType": "nats",
  "category": "technical",
  "enabled": true,
  "isSystem": false,
  "isTemplate": false,
  "config": {
    "nats": {
      "server": "PLACEHOLDER",
      "subject": "mapex.saga.nats.PLACEHOLDER",
      "message": { "saga": "nats", "runID": "PLACEHOLDER" }
    }
  }
}`

// SagaNatsTrigger returns the POST /api/v1/triggers body for the NATS
// smoke. server is rewritten to the saga-managed embedded server;
// subject embeds the runID.
func SagaNatsTrigger(runID, serverURL string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaNatsTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaNatsTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-nats-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	natsCfg, _ := cfg["nats"].(map[string]any)
	natsCfg["server"] = serverURL
	natsCfg["subject"] = fmt.Sprintf("mapex.saga.nats.%s", runID)
	natsCfg["message"] = map[string]any{"saga": "nats", "runID": runID}
	return payload
}
