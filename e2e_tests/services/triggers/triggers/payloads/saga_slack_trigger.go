package payloads

import (
	"encoding/json"
	"fmt"
)

const sagaSlackTriggerJSON = `{
  "name": "Saga Slack Trigger",
  "description": "smoke test — webhook posts to in-process HTTP sink",
  "triggerType": "slack",
  "category": "communication",
  "enabled": true,
  "isSystem": false,
  "isTemplate": false,
  "config": {
    "slack": {
      "webhookUrl": "PLACEHOLDER",
      "message": "Saga slack smoke run=PLACEHOLDER"
    }
  }
}`

// SagaSlackTrigger returns the POST /api/v1/triggers body for the
// Slack smoke. webhookUrl is rewritten to sinkURL (the HTTP sink) and
// the message embeds the runID for content-key validation. Pure: the
// create step reads the sink's ephemeral address from the bag and
// passes it in.
func SagaSlackTrigger(runID, sinkURL string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaSlackTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaSlackTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-slack-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	slackCfg, _ := cfg["slack"].(map[string]any)
	slackCfg["webhookUrl"] = sinkURL
	slackCfg["message"] = fmt.Sprintf("Saga slack smoke run=%s", runID)
	return payload
}
