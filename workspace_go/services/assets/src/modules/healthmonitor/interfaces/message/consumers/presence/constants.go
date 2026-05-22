package presence

import (
	hmMessage "assets/src/modules/healthmonitor/interfaces/message"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Consumer-local aliases for the MQTT presence stream. The broker
// plugin publishes both connect and disconnect advisories on the
// SAME subject, discriminated by the Event field. The stream uses
// WorkQueue retention, which only permits ONE consumer per subject
// filter — so we register ONE durable here and dispatch by Event
// inside the message handler.
var (
	Stream = hmMessage.MqttPresenceStreamName

	AdvisorySubject = hmMessage.MqttPresenceAdvisorySubject

	PresenceDurable = config.Durable("assets", "mqtt-presence")
)

const EventType = hmMessage.MqttPresenceEventType
