package org_created

import (
	orgContract "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/organizations"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for RetentionOrgCreated consumer.
 *
 * The wire-level contract (subject/stream) is owned by the mapexIam service
 * organizations module and declared in
 * packages/contracts/services/mapexIam/organizations. These locals are
 * thin aliases kept so the consumer file stays stable when subjects evolve.
 * Durable is local to this consumer.
 */

// Stream name for organization events from mapexIam service.
var Stream = orgContract.StreamOrganizationEvents

// Subject for organization created events.
var Subject = orgContract.SubjectOrganizationCreated

// Durable name for the retention org_created consumer.
var Durable = config.Durable("events", "org-created")

// EventType for DLQ metadata.
const EventType = orgContract.EventTypeOrganizationCreated
