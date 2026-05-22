package payloads

import (
	"encoding/json"
	"fmt"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
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
// Slack smoke. webhookUrl is rewritten to the HTTP sink and the
// message embeds the runID for content-key validation.
func SagaSlackTrigger(runID string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaSlackTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaSlackTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-slack-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	slackCfg, _ := cfg["slack"].(map[string]any)
	slackCfg["webhookUrl"] = constants.TriggerSinkURL
	slackCfg["message"] = fmt.Sprintf("Saga slack smoke run=%s", runID)
	return payload
}
