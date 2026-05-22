package payloads

import (
	"encoding/json"
	"fmt"
)

const sagaMqttTriggerJSON = `{
  "name": "Saga MQTT Trigger",
  "description": "smoke test — publishes to the saga-managed broker",
  "triggerType": "mqtt",
  "category": "technical",
  "enabled": true,
  "isSystem": false,
  "isTemplate": false,
  "config": {
    "mqtt": {
      "broker": "PLACEHOLDER",
      "port": 0,
      "topic": "mapex-saga/mqtt/PLACEHOLDER",
      "qos": 1,
      "message": { "saga": "mqtt", "runID": "PLACEHOLDER" }
    }
  }
}`

// SagaMqttTrigger returns the POST /api/v1/triggers body for the MQTT
// smoke. broker / port are rewritten to the saga-managed in-process
// broker (mochi-mqtt) the journey started; topic embeds the runID so
// the events-trigger oracle can validate the resolved config.
func SagaMqttTrigger(runID, brokerHost string, brokerPort int) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(sagaMqttTriggerJSON), &payload); err != nil {
		panic(fmt.Sprintf("SagaMqttTrigger: literal payload is not valid JSON: %v", err))
	}
	payload["name"] = fmt.Sprintf("saga-mqtt-%s", runID)

	cfg, _ := payload["config"].(map[string]any)
	mqttCfg, _ := cfg["mqtt"].(map[string]any)
	mqttCfg["broker"] = brokerHost
	mqttCfg["port"] = brokerPort
	mqttCfg["topic"] = fmt.Sprintf("mapex-saga/mqtt/%s", runID)
	mqttCfg["message"] = map[string]any{"saga": "mqtt", "runID": runID}
	return payload
}
