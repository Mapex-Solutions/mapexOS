// Package assets holds cross-service contract constants published by the
// assets service assets module and consumed by other services.
//
// These constants describe NATS subjects/streams that cross service
// boundaries and therefore cannot live inside a service-local
// application/constants file (see /go-arch §3 + §4 reciprocity).
//
// Ownership: assets service (publisher of asset cache invalidation events).
// Consumers (Go): router, js-executor, events.
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/assets/assets.
//
// Contracts stay leaf-level — no imports from services/.
package assets

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// FanoutStreamName is the platform-wide JetStream stream carrying FANOUT
// broadcast messages (cache invalidation). Shared across all publishers
// and consumers of any *.fanout.* subject (assets, workflow, future
// services). Resolved at package init from GO_ENV — e.g. "DEV-MAPEXOS-FANOUT".
var FanoutStreamName = config.StreamName("FANOUT", "")

// FanoutAssetSubject is the NATS subject published by the assets service
// whenever an asset is created, updated, or deleted, so that consuming
// services invalidate their TieredCache (L0/L1) for that asset. Resolved
// at package init from GO_ENV — e.g. "dev.mapexos.fanout.asset.invalidate".
//
// Published by: assets service (assets module).
// Consumed by: router, js-executor, events.
var FanoutAssetSubject = config.Subject("fanout", "asset.invalidate")

// FanoutAssetEventType is the DLQ event-type identifier for consumer
// failure routing of asset-invalidate FANOUT messages.
const FanoutAssetEventType = "fanout.asset.invalidate"

// MQTT auth-type enum values. The asset declares which credential the
// device is allowed to present at CONNECT — the broker plugin enforces
// the choice (cert-mode asset cannot bcrypt a password; password-mode
// asset cannot present a cert). Switching modes is an explicit
// operator action; the side-effect of switching wipes the unused
// credential (active cert revoked on password→cert, hash cleared on
// cert→password).
const (
	MqttAuthTypePassword = "password"
	MqttAuthTypeCert     = "cert"
)
