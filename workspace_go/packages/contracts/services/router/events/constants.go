// Package events holds the cross-service contract constants for the
// router service events module.
//
// These subject constants are the wire-level contract for messages published
// by the router service and consumed by other services:
//   - ${env}.mapexos.events.save              → events MS  (save_event router kind)
//   - ${env}.mapexos.events.lake_house        → events/lakehouse sink (lake_house router kind)
//   - ${env}.mapexos.events.notification      → events/notification consumer
//   - ${env}.mapexos.events.router            → events MS  (router execution history)
//   - ${env}.mapexos.workflow.execution.router → workflow MS (workflow router kind)
//   - ${env}.mapexos.trigger.router.execute   → triggers MS (trigger router kind)
//
// Ownership: router service (publisher).
// Consumers (Go): events, triggers, workflow.
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/router/events.
//
// Contracts stay leaf-level — no imports from services/.
package events

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// SubjectSaveEvent is the NATS subject for the save_event router kind.
// Consumed by the events service. Resolved at package init from GO_ENV —
// e.g. "dev.mapexos.events.save".
var SubjectSaveEvent = config.Subject("events", "save")

// StreamEvents is the NATS JetStream stream that carries save_event
// messages for the events service consumer (events_save). Resolved at
// package init — e.g. "DEV-MAPEXOS-EVENTS-SAVE".
var StreamEvents = config.StreamName("EVENTS", "SAVE")

// EventTypeSaveEvent tags DLQ messages produced by the events_save consumer.
const EventTypeSaveEvent = "events.save"

// SubjectLakeHouse is the NATS subject for the lake_house router kind.
var SubjectLakeHouse = config.Subject("events", "lake_house")

// SubjectNotification is the NATS subject for the notification router kind.
var SubjectNotification = config.Subject("events", "notification")

// SubjectRouterHistory is the NATS subject for router execution history
// events, consumed by the events service. Resolved at package init —
// e.g. "dev.mapexos.events.router".
var SubjectRouterHistory = config.Subject("events", "router")

// StreamRouterHistory is the NATS JetStream stream that carries router
// execution history events for the events service consumer (events_router).
// Resolved at package init — e.g. "DEV-MAPEXOS-EVENTS-ROUTER-LOGS".
var StreamRouterHistory = config.StreamName("EVENTS", "ROUTER-LOGS")

// EventTypeRouterHistory tags DLQ messages produced by the events_router
// consumer.
const EventTypeRouterHistory = "router"

// SubjectWorkflowExecution is the fixed subject for the workflow execution
// stream. Routing is done by mode in the payload, never by subject.
// Consumed by the workflow service.
var SubjectWorkflowExecution = config.Subject("workflow", "execution.router")

// SubjectTriggerRouterExecute is the NATS subject for trigger execution
// dispatched by the router service. Consumed by the triggers service.
// triggerId is sent in the payload.
var SubjectTriggerRouterExecute = config.Subject("trigger", "router.execute")
