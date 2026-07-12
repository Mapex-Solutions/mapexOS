package payloads

import (
	"encoding/json"
	"fmt"
)

const sagaTeamsTriggerJSON = `{
  "name": "Saga Teams Trigger",
  "description": "smoke test — webhook posts to in-process HTTP sink",
  "triggerType": "teams",
  "category": "communication",
  "enabled": true,
  "isSystem": false,
  "isTemplate": false,
  "config": {
    "teams": {
      "webhookUrl": "PLACEHOLDER",
      "title": "Saga teams smoke",
      "text": "Saga teams smoke run=PLACEHOLDER"
    }
  }
}`

// SagaTeamsTrigger returns the POST /api/v1/triggers body for the
// Teams smoke. webhookUrl is rewritten to sinkURL (the HTTP sink) and
// the text embeds the runID. Pure: the create step reads the sink's
// ephemeral address from the bag and passes it in.
func SagaTeamsTrigger(runID, sinkURL string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaTeamsTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaTeamsTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-teams-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	teamsCfg, _ := cfg["teams"].(map[string]any)
	teamsCfg["webhookUrl"] = sinkURL
	teamsCfg["title"] = fmt.Sprintf("Saga teams smoke run=%s", runID)
	teamsCfg["text"] = fmt.Sprintf("Saga teams smoke text run=%s", runID)
	return payload
}
