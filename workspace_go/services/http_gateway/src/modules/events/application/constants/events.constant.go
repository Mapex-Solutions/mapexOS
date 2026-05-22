package constants

import (
	processorContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/events"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// EventsRawSubject is the subject for raw event publishing (ClickHouse events_raw table).
// Resolved at package init from GO_ENV — e.g. "dev.mapexos.events.raw".
var EventsRawSubject = config.Subject("events", "raw")

// ProcessorJsExecuteSubject is the subject for dispatching events to js-executor.
//
// This is a cross-service NATS subject (http_gateway -> js-executor). The
// authoritative declaration lives in
// packages/contracts/services/http_gateway/events.SubjectProcessorJSExecute;
// this constant is a local alias kept to avoid churn at existing call sites.
var ProcessorJsExecuteSubject = processorContracts.SubjectProcessorJSExecute
