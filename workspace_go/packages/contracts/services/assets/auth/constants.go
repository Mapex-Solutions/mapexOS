package auth

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// L2WritesStreamName carries durable retry messages for L2 (MinIO)
// writes that failed on the synchronous happy path. Stream is a NATS
// JetStream work queue with native Msg-Id dedup (5s window) — multiple
// rapid CRUDs to the same entity coalesce into a single retry message.
// Resolved at package init from GO_ENV — e.g. "DEV-MAPEXOS-L2-WRITES".
var L2WritesStreamName = config.StreamName("L2-WRITES", "")

// L2WritesAssetSubject is the NATS subject the assets module publishes
// to when the synchronous L2 write of an asset fails. The consumer
// inside the same module drains the subject and retries until success
// or DLQ exhaustion.
// Resolved at package init — e.g. "dev.mapexos.l2_writes.asset".
var L2WritesAssetSubject = config.Subject("l2_writes", "asset")

// L2WritesTemplateSubject mirrors L2WritesAssetSubject for the
// assettemplates module.
// Resolved at package init — e.g. "dev.mapexos.l2_writes.template".
var L2WritesTemplateSubject = config.Subject("l2_writes", "template")

// L2WritesAssetEventType is the DLQ event-type identifier for failure
// routing of asset L2 retry messages.
const L2WritesAssetEventType = "l2_writes.asset"

// L2WritesTemplateEventType mirrors L2WritesAssetEventType for templates.
const L2WritesTemplateEventType = "l2_writes.template"
