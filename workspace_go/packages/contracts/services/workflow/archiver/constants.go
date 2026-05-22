// Package archiver holds the cross-service contract constants emitted by
// the workflow service archiver module.
//
// These subject constants are the wire-level contract for messages
// published by the workflow archiver and consumed by other services:
//   - mapexos.events.workflow → events MS (EVENTS-WORKFLOW stream) —
//     terminal workflow execution history sink (ClickHouse cold storage).
//
// Ownership: workflow service (publisher).
// Consumers (Go): events.
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/workflow/archiver.
//
// Contracts stay leaf-level — no imports from services/.
package archiver

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// SubjectEventsWorkflow is the NATS subject for publishing terminal
// workflow execution history events to the events service for ClickHouse
// cold storage. Resolved at package init — e.g. "dev.mapexos.events.workflow".
var SubjectEventsWorkflow = config.Subject("events", "workflow")

// StreamEventsWorkflow is the NATS JetStream stream that carries terminal
// workflow execution history events for the events service consumer
// (events_workflow). Resolved at package init —
// e.g. "DEV-MAPEXOS-EVENTS-WORKFLOW-LOGS".
var StreamEventsWorkflow = config.StreamName("EVENTS", "WORKFLOW-LOGS")

// EventTypeEventsWorkflow tags DLQ messages produced by the
// events_workflow consumer in the events service.
const EventTypeEventsWorkflow = "workflow"
