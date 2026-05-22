// Package dlq holds the platform-level Dead Letter Queue NATS contract
// constants shared by every service that publishes to or consumes from
// the unified MAPEXOS DLQ.
//
// The DLQ is owned by no single service — ALL services that configure
// natsModel.DLQPolicy feed the same stream/subject, and the events
// service consumes it into ClickHouse (events_dlq table).
//
// Ownership: platform (infrastructure-wide), attributed to the gokit lib.
// Publishers: every service with a DLQ-enabled consumer.
// Consumers (Go): events (events_dlq consumer).
// Reciprocity: mirrored by workspace_js/packages/schemas/src/common/dlq.
//
// Contracts stay leaf-level — no imports from services/.
package dlq

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Subject is the unified NATS subject on which every service publishes
// DLQ messages via natsModel.DLQPolicy. Resolved at package init from
// GO_ENV — e.g. "dev.mapexos.mapexosgokit.dlq".
var Subject = config.Subject("mapexosgokit", "dlq")

// Stream is the unified NATS JetStream stream that retains DLQ messages
// for the events service consumer (events_dlq) to persist into
// ClickHouse. Resolved at package init from GO_ENV — e.g.
// "DEV-MAPEXOS-MAPEXOSGOKIT-DLQ".
var Stream = config.StreamName("MAPEXOSGOKIT", "DLQ")
