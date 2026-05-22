// Package assettemplates holds cross-service contract constants produced and
// consumed by the assets service assettemplates module across service
// boundaries.
//
// Two directions are covered:
//  1. Events produced OUTSIDE of assets (by the lists/mapexos service) and
//     consumed here to propagate classification name updates into
//     denormalized AssetTemplate documents.
//  2. Events produced by THIS module (template FANOUT cache invalidation)
//     and consumed by other services (router, js-executor, events).
//
// Ownership (per constant):
//   - ListNameUpdated*   : mapexos service (publisher of list name change events).
//   - Fanout*            : assets service (publisher of template invalidation).
//
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/assets/assettemplates.
//
// Contracts stay leaf-level — no imports from services/.
package assettemplates

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// ListNameUpdatedStream carries list name updated events. Resolved at
// package init from GO_ENV — e.g. "DEV-MAPEXOS-ASSETS-LISTS".
var ListNameUpdatedStream = config.StreamName("ASSETS", "LISTS")

// ListNameUpdatedSubject is the subject published by the mapexos service
// whenever a list's name changes (manufacturer, model, or category).
// Resolved at package init from GO_ENV — e.g. "dev.mapexos.lists.name_updated".
var ListNameUpdatedSubject = config.Subject("lists", "name_updated")

// ListNameUpdatedEventType is the DLQ event-type identifier for consumer
// failure routing.
const ListNameUpdatedEventType = "lists.name_updated"

// FanoutStreamName is the platform-wide JetStream stream carrying FANOUT
// broadcast messages (cache invalidation). Shared across all publishers
// and consumers of any *.fanout.* subject (assets, workflow, future
// services). Resolved at package init from GO_ENV — e.g. "DEV-MAPEXOS-FANOUT".
var FanoutStreamName = config.StreamName("FANOUT", "")

// FanoutTemplateSubject is the NATS subject published by the assets service
// whenever an asset template is created, updated, or deleted, so that
// consuming services invalidate their TieredCache (L0/L1) for that template.
// Resolved at package init — e.g. "dev.mapexos.fanout.template.invalidate".
//
// Published by: assets service (assettemplates module).
// Consumed by: router, js-executor, events.
var FanoutTemplateSubject = config.Subject("fanout", "template.invalidate")

// FanoutTemplateEventType is the DLQ event-type identifier for consumer
// failure routing of template-invalidate FANOUT messages.
const FanoutTemplateEventType = "fanout.template.invalidate"
