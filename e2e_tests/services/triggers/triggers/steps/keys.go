package steps

import (
	"net/http"
)

// Bag keys this package writes. Other packages reading these keys
// import the constants from here.
const (
	// BagKeyTriggerID is the Mongo ObjectID hex of the trigger created
	// by CreateXTrigger steps. Route groups of kind=trigger reference
	// it on Router.Trigger.TriggerId.
	BagKeyTriggerID = "triggers.triggerID"

	// HTTP sink — used by HTTP / Slack / Teams trigger smokes.

	// BagKeyTriggerSinkServer holds the *http.Server the HTTP sink
	// step started, so the Compensate path can shut it down. The
	// stored value type-asserts to *http.Server at the consumer site.
	BagKeyTriggerSinkServer = "triggers.sinkServer"

	// BagKeyTriggerSinkHits holds a *atomic.Int64 counter the HTTP
	// sink increments on each POST it receives. Asserts read this
	// directly to verify the trigger fired without round-tripping
	// through the events service.
	BagKeyTriggerSinkHits = "triggers.sinkHits"

	// SMTP sink — used by Email trigger smoke.

	// BagKeySmtpServer holds the *smtp.Server the SMTP sink step
	// started so Compensate can call Close() on it. Type asserts to
	// *smtp.Server at the consumer site.
	BagKeySmtpServer = "triggers.smtpServer"

	// BagKeySmtpHits is a *atomic.Int64 incremented on every
	// successful inbound message the SMTP sink finalizes.
	BagKeySmtpHits = "triggers.smtpHits"

	// BagKeySmtpLastMessage stores the captured fields of the most
	// recent SMTP message as a *SmtpCapturedMessage. Asserts read it
	// to validate subject/to/from/body without re-implementing
	// RFC 5322 parsing in every check.
	BagKeySmtpLastMessage = "triggers.smtpLastMessage"

	// WebSocket sink — used by WebSocket trigger smoke.

	// BagKeyWsServer holds the *http.Server backing the WS upgrade so
	// Compensate can stop it. Type asserts to *http.Server.
	BagKeyWsServer = "triggers.wsServer"

	// BagKeyWsHits is a *atomic.Int64 incremented on every successful
	// WebSocket handshake the sink accepts.
	BagKeyWsHits = "triggers.wsHits"

	// BagKeyWsLastMessage stores the first frame payload of the most
	// recent connection as a **string (pointer-to-pointer so the sink
	// can swap atomically without a mutex on the slot).
	BagKeyWsLastMessage = "triggers.wsLastMessage"

	// MQTT in-process broker (mochi-mqtt) — used by MQTT trigger smoke.

	// BagKeyMqttBroker holds the *mqtt.Server the saga started so
	// Compensate can Close() it.
	BagKeyMqttBroker = "triggers.mqttBroker"

	// BagKeyMqttBrokerHost is the bind host the trigger config's
	// broker field should target ("127.0.0.1").
	BagKeyMqttBrokerHost = "triggers.mqttBrokerHost"

	// BagKeyMqttBrokerPort is the OS-assigned port the broker is
	// listening on. The trigger config carries it under port.
	BagKeyMqttBrokerPort = "triggers.mqttBrokerPort"

	// NATS embedded server — used by NATS trigger smoke.

	// BagKeyNatsServer holds the *natsserver.Server the saga started
	// so Compensate can Shutdown() it.
	BagKeyNatsServer = "triggers.natsServer"

	// BagKeyNatsURL is the "nats://host:port" connect string the
	// trigger config's server field should target.
	BagKeyNatsURL = "triggers.natsURL"

	// RabbitMQ ephemeral container — used by RabbitMQ trigger smoke.

	// BagKeyRabbitmqContainer holds the *rabbitmq.RabbitMQContainer so
	// Compensate can Terminate() it.
	BagKeyRabbitmqContainer = "triggers.rabbitmqContainer"

	// BagKeyRabbitmqHost / Port / User / Pass are the connection
	// parameters the trigger config carries.
	BagKeyRabbitmqHost = "triggers.rabbitmqHost"
	BagKeyRabbitmqPort = "triggers.rabbitmqPort"
	BagKeyRabbitmqUser = "triggers.rabbitmqUser"
	BagKeyRabbitmqPass = "triggers.rabbitmqPass"
)

// Compile-time assertion: the sink server stored on the bag is the
// std-lib *http.Server type so consumers can stop it via its public
// API without depending on a custom struct.
var _ *http.Server = (*http.Server)(nil)
