package events_workflow

import (
	archiverContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/archiver"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for WorkflowEvent consumer.
 *
 * The wire-level contract (subject/stream) is owned by the workflow service
 * archiver module and declared in packages/contracts/services/workflow/archiver.
 * These locals are thin aliases kept so the consumer file stays stable when
 * subjects evolve. Durable is local to this consumer.
 */

// Stream name for workflow execution history events.
var Stream = archiverContract.StreamEventsWorkflow

// Subject pattern for workflow execution history events.
var Subject = archiverContract.SubjectEventsWorkflow

// Durable name for the events_workflow consumer.
var Durable = config.Durable("events", "workflow")

// EventType for DLQ metadata.
const EventType = archiverContract.EventTypeEventsWorkflow
