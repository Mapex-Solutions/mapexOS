// Package message holds intra-service NATS constants for the events module.
//
// Subjects crossing service boundaries live in packages/contracts/services/...
// and are imported by the consumer files directly. Only values scoped to this
// service (e.g., the shared JetStream stream name consumed by both trigger and
// workflow execution consumers) belong here.
package message

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// StreamTriggers is the JetStream stream consumed by the triggers service for
// trigger and workflow execution dispatch. Shared by the trigger_execute and
// plugin_execute consumers. Resolved at package init —
// e.g. "DEV-MAPEXOS-TRIGGERS-EXECUTE".
var StreamTriggers = config.StreamName("TRIGGERS", "EXECUTE")
