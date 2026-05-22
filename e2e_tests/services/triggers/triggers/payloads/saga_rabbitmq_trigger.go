package payloads

import (
	"encoding/json"
	"fmt"
)

const sagaRabbitmqTriggerJSON = `{
  "name": "Saga RabbitMQ Trigger",
  "description": "smoke test — publishes to the saga-managed RabbitMQ",
  "triggerType": "rabbitmq",
  "category": "technical",
  "enabled": true,
  "isSystem": false,
  "isTemplate": false,
  "config": {
    "rabbitmq": {
      "host": "PLACEHOLDER",
      "port": 0,
      "vhost": "/",
      "username": "PLACEHOLDER",
      "password": "PLACEHOLDER",
      "publishMode": "queue",
      "queue": "saga-mq-PLACEHOLDER",
      "message": { "saga": "rabbitmq", "runID": "PLACEHOLDER" }
    }
  }
}`

// SagaRabbitmqTrigger returns the POST /api/v1/triggers body for the
// RabbitMQ smoke. host / port / username / password are rewritten to
// the saga-managed ephemeral container; queue embeds the runID.
func SagaRabbitmqTrigger(runID, host string, port int, user, pass string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaRabbitmqTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaRabbitmqTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-rabbitmq-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	rabbitCfg, _ := cfg["rabbitmq"].(map[string]any)
	rabbitCfg["host"] = host
	rabbitCfg["port"] = port
	rabbitCfg["username"] = user
	rabbitCfg["password"] = pass
	rabbitCfg["queue"] = fmt.Sprintf("saga-mq-%s", runID)
	rabbitCfg["message"] = map[string]any{"saga": "rabbitmq", "runID": runID}
	return payload
}
