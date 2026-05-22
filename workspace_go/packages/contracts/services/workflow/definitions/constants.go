// Package definitions holds the cross-service contract constants emitted by
// the workflow service definitions module.
//
// These subject/stream constants are the wire-level contract for FANOUT
// cache-invalidation messages published by the workflow service when a
// definition's code nodes change:
//   - mapexos.fanout.workflow.definition.invalidate (FANOUT stream).
//
// Ownership: workflow service (publisher).
// Consumers (Go/TS): js-workflow-executor (TS, future) and any pod that
// caches workflow definition scripts.
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/workflow/definitions.
//
// Contracts stay leaf-level — no imports from services/.
package definitions

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// FanoutStreamName is the platform-wide JetStream stream carrying FANOUT
// broadcast messages (cache invalidation). Shared across all publishers
// and consumers of any *.fanout.* subject (assets, workflow, future
// services). Resolved at package init from GO_ENV — e.g.
// "DEV-MAPEXOS-FANOUT".
var FanoutStreamName = config.StreamName("FANOUT", "")

// FanoutDefinitionSubject is the NATS subject published by the workflow
// service whenever a workflow definition's code nodes are created,
// updated, or deleted, so consuming services invalidate their TieredCache
// (L0/L1) for that definition. Resolved at package init from GO_ENV —
// e.g. "dev.mapexos.fanout.workflow.definition.invalidate".
//
// Published by: workflow service (definitions module).
// Consumed by: js-workflow-executor (and future workflow pods).
var FanoutDefinitionSubject = config.Subject("fanout", "workflow.definition.invalidate")
