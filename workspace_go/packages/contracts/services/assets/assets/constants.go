// Package assets holds cross-service contract constants published by the
// assets service assets module and consumed by other services.
//
// These constants describe NATS subjects/streams that cross service
// boundaries and therefore cannot live inside a service-local
// application/constants file: a cross-service contract lives here and is
// mirrored on the TypeScript side.
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

// FanoutStreamName is the platform-wide JetStream stream carrying fanout
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
// failure routing of asset-invalidate fanout messages.
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

// LoRaWAN activation modes. OTAA derives the session at join from root keys;
// ABP carries a fixed DevAddr + session keys.
const (
	LorawanActivationOTAA = "otaa"
	LorawanActivationABP  = "abp"
)

// LoRaWAN asset kinds. A device is an end-device (identity + keys); a gateway is
// radio infrastructure (frequency plan + connection auth), never bound to a
// device.
const (
	LorawanKindDevice  = "device"
	LorawanKindGateway = "gateway"
)

// LoRaWAN gateway connection auth modes. "eui" gates a UDP gateway by its
// registered EUI (the Semtech UDP protocol allows nothing stronger); "cert" is
// mTLS for a Basics Station gateway via the platform PKI; "key" is a Basics
// Station bearer token validated by the LNS.
const (
	LorawanGatewayAuthModeEUI  = "eui"
	LorawanGatewayAuthModeCert = "cert"
	LorawanGatewayAuthModeKey  = "key"
)

// Asset custom-attribute kinds. Value is typed by kind and validated by the
// service (struct tags cannot type-switch on an any value).
const (
	AssetAttributeKindInteger = "integer"
	AssetAttributeKindString  = "string"
	AssetAttributeKindBoolean = "boolean"
	AssetAttributeKindDate    = "date"
	AssetAttributeKindGeo     = "geo"
)

// Asset custom-attribute limits and the reserved label prefix (platform-owned
// labels the operator may not set).
const (
	AssetAttributeMaxCount       = 20
	AssetAttributeLabelMaxLen    = 64
	AssetAttributeReservedPrefix = "mapex."
)
