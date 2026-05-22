package heartbeat

import (
	hmMessage "assets/src/modules/healthmonitor/interfaces/message"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Consumer-local aliases for the heartbeat stream/subject/durable.
// Authoritative stream/subject values live in
// interfaces/message/constants.go so bootstrap and the service layer
// share the same source of truth. Durable is local to this consumer.
var (
	Stream  = hmMessage.HeartbeatStreamName
	Subject = hmMessage.HeartbeatSubject
	Durable = config.Durable("assets", "heartbeat")
)

const EventType = hmMessage.HeartbeatEventType
