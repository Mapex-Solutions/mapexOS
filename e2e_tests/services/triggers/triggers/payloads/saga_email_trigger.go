// Package payloads holds CreateTrigger fixtures for the triggers
// module. SagaEmailTrigger wires the trigger config against the
// in-process SMTP sink the email smoke journey stands up.
package payloads

import (
	"encoding/json"
	"fmt"
)

const sagaEmailTriggerJSON = `{
  "name": "Saga Email Trigger",
  "description": "smoke test — delivers to in-process SMTP sink",
  "triggerType": "email",
  "category": "communication",
  "enabled": true,
  "isSystem": false,
  "isTemplate": false,
  "config": {
    "email": {
      "smtpHost": "PLACEHOLDER",
      "smtpPort": 0,
      "fromAddr": "saga-no-reply@mapex.test",
      "to": "saga-recipient@mapex.test",
      "subject": "Saga email smoke",
      "body": "placeholder body"
    }
  }
}`

// SagaEmailTrigger returns the POST /api/v1/triggers body the email
// smoke journey uses. smtpHost/smtpPort are rewritten to the saga
// sink's ephemeral address; subject and body interpolate the runID so
// the downstream assert can correlate the captured message with this
// run. Pure: the create step reads the sink's address from the bag and
// passes host + port in.
func SagaEmailTrigger(runID, smtpHost string, smtpPort int) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaEmailTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaEmailTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-email-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	emailCfg, _ := cfg["email"].(map[string]any)
	emailCfg["smtpHost"] = smtpHost
	emailCfg["smtpPort"] = smtpPort
	emailCfg["body"] = fmt.Sprintf("Saga email run=%s — delivered via in-process SMTP sink.", runID)
	emailCfg["subject"] = fmt.Sprintf("Saga email smoke run=%s", runID)
	return payload
}
