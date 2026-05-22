package presence

import (
	"encoding/json"
	"fmt"
	"time"

	"assets/src/modules/healthmonitor/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// presenceEnvelope is the slim shape used to discriminate the
// broker's presence advisory by Event before dispatching to the
// service handler. Keeping the parser inside the consumer (instead
// of duplicating the full PresenceAdvisory struct) avoids pulling
// the cross-service contract into the wiring layer.
type presenceEnvelope struct {
	Event string `json:"event"`
}

// NewConsumer starts ONE JetStream consumer that drives both
// connect and disconnect presence transitions for MQTT-protocol
// assets. The broker plugin publishes both shapes on the same
// subject, discriminated by the Event field; the stream uses
// WorkQueue retention, which only permits ONE consumer per subject
// filter — so we dispatch internally instead of registering two.
func NewConsumer(bus *natsModel.Bus, service ports.HealthMonitorServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")
	if natsFetchTimeout <= 0 {
		natsFetchTimeout = 1
	}

	logger.Info(fmt.Sprintf("[CONSUMER:MqttPresence] Starting %s", PresenceDurable))

	_, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:          Stream,
		Subject:         AdvisorySubject,
		Durable:         PresenceDurable,
		QueueGroup:      "assets-mqtt-presence",
		FetchTimeout:    time.Duration(natsFetchTimeout) * time.Second,
		DuplicateWindow: 1 * time.Minute,

		RetryPolicy: natsModel.DefaultRetryPolicy(),
		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "assets",
			EventType:   EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			dispatch(service, msg)
		},
	})
	if err != nil {
		logger.Error(err, "[CONSUMER:MqttPresence] Failed to start consumer")
		return
	}
	logger.Info("[CONSUMER:MqttPresence] Started successfully")
}

// dispatch peeks at the Event field and forwards the message to the
// matching service handler. Malformed payloads or unknown events are
// Ack-and-drop — the broker plugin only emits two event shapes and an
// outlier would mean a contract drift the consumer should not retry.
func dispatch(service ports.HealthMonitorServicePort, msg *natsModel.Message) {
	var env presenceEnvelope
	if err := json.Unmarshal(msg.Data, &env); err != nil {
		logger.Warn(fmt.Sprintf("[CONSUMER:MqttPresence] malformed payload: %v", err))
		msg.Ack()
		return
	}
	switch env.Event {
	case "connect":
		service.HandlePresenceConnect(msg)
	case "disconnect":
		service.HandlePresenceDisconnect(msg)
	default:
		msg.Ack()
	}
}
