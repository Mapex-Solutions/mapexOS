package scan

import (
	hmMessage "assets/src/modules/healthmonitor/interfaces/message"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Consumer-local aliases for the scan stream/subject/durable.
// Authoritative stream/subject values live in
// interfaces/message/constants.go so bootstrap and the service layer
// share the same source of truth. Durable is local to this consumer.
var (
	Stream  = hmMessage.ScanStreamName
	Subject = hmMessage.ScanSubject
	Durable = config.Durable("assets", "scan")
)

const EventType = "healthmonitor.scan"
