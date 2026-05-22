// Package plugins holds the cross-service contract constants emitted by
// the workflow service plugins module.
//
// These subject/stream constants are the wire-level contract for FANOUT
// cache-invalidation messages published by the workflow service when a
// plugin manifest changes:
//   - mapexos.fanout.workflow.plugin.invalidate (FANOUT stream).
//
// Ownership: workflow service (publisher).
// Consumers: all workflow pods (self-subscribe via NATS Fanout). The
// stream FANOUT is shared across the whole platform — declared once here
// for workflow consumers, and once per service that participates in the
// FANOUT pattern (assets, definitions, plugins).
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/workflow/plugins
// when the JS side gains a fanout consumer.
//
// Contracts stay leaf-level — no imports from services/.
package plugins

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// FanoutStreamName is the platform-wide JetStream stream carrying FANOUT
// broadcast messages (cache invalidation). Shared across all publishers
// and consumers of any *.fanout.* subject (assets, workflow, future
// services). Resolved at package init from GO_ENV — e.g.
// "DEV-MAPEXOS-FANOUT".
var FanoutStreamName = config.StreamName("FANOUT", "")

// FanoutPluginSubject is the NATS subject published by the workflow
// plugins module whenever a plugin manifest is created, updated, or
// deleted, so all workflow pods invalidate their TieredCache (L0/L1).
// Resolved at package init from GO_ENV — e.g.
// "dev.mapexos.fanout.workflow.plugin.invalidate".
//
// Published by: workflow service (plugins module).
// Consumed by: all workflow pods (self-subscribe via NATS Fanout).
var FanoutPluginSubject = config.Subject("fanout", "workflow.plugin.invalidate")
