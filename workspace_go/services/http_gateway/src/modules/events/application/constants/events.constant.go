package constants

import (
	httpEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/events"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// EventsRawSubject is the subject for raw event publishing (ClickHouse events_raw table).
// Resolved at package init from GO_ENV — e.g. "dev.mapexos.events.raw".
var EventsRawSubject = config.Subject("events", "raw")

// HttpDataSubject is the STATIC subject http_gateway publishes HTTP device
// telemetry to for js-executor — mirrors mqtt.data / lorawan.data (identity in
// the payload, never in the subject). The authoritative declaration lives in
// packages/contracts/services/http_gateway/events.SubjectHTTPData; this constant
// is a local alias.
var HttpDataSubject = httpEvents.SubjectHTTPData
