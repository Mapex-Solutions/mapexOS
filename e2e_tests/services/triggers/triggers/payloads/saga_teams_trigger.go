package payloads

import (
	"encoding/json"
	"fmt"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
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
// Teams smoke. webhookUrl is rewritten to the HTTP sink and the text
// embeds the runID.
func SagaTeamsTrigger(runID string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaTeamsTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaTeamsTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-teams-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	teamsCfg, _ := cfg["teams"].(map[string]any)
	teamsCfg["webhookUrl"] = constants.TriggerSinkURL
	teamsCfg["title"] = fmt.Sprintf("Saga teams smoke run=%s", runID)
	teamsCfg["text"] = fmt.Sprintf("Saga teams smoke text run=%s", runID)
	return payload
}
