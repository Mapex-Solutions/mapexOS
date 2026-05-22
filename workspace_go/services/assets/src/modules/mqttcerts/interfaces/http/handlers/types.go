package handlers

import mqttPorts "assets/src/modules/mqttcerts/application/ports"

// MqttCertsHandler bundles the service port for the external endpoints.
type MqttCertsHandler struct {
	service mqttPorts.MqttCertsServicePort
}

func NewMqttCertsHandler(s mqttPorts.MqttCertsServicePort) *MqttCertsHandler {
	return &MqttCertsHandler{service: s}
}
