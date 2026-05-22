// Package healthmonitor holds the cross-service contract constants emitted by
// the assets service healthmonitor module.
//
// These constants are the wire-level contract for messages published by
// healthmonitor and consumed by other services:
//   - mapexos.route.execute      → router MS (ROUTE-GROUPS stream)
//   - mapexos.events.asset_status_save → events MS (EVENTS-ASSET-STATUS stream)
//
// Ownership: assets service (publisher).
// Consumers (Go): router, events.
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/assets/healthmonitor.
//
// Contracts stay leaf-level — no imports from services/.
package healthmonitor

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Router stream and subject (for publishing route events). Stream and subject
// resolve at package init from GO_ENV — e.g. "DEV-MAPEXOS-ROUTER-EXECUTE" and
// "dev.mapexos.route.execute".
var RouterStream = config.StreamName("ROUTER", "EXECUTE")
var RouterSubject = config.Subject("route", "execute")

// AssetStatusSaveSubject is the subject for the dedicated asset connectivity
// persistence stream. The healthmonitor alert publisher ALWAYS publishes to
// this subject on a transition — regardless of whether RouteGroupIds are
// configured — so events MS can persist the history. Resolved at package
// init from GO_ENV — e.g. "dev.mapexos.events.asset_status_save".
var AssetStatusSaveSubject = config.Subject("events", "asset_status_save")

// AssetStatusSaveStream is the NATS JetStream stream that carries asset
// connectivity persistence events (offline/online transitions) for the
// events service consumer (asset_status_save). Resolved at package init —
// e.g. "DEV-MAPEXOS-EVENTS-ASSET-STATUS".
var AssetStatusSaveStream = config.StreamName("EVENTS", "ASSET-STATUS")

// AssetStatusSaveEventType tags DLQ messages produced by the
// asset_status_save consumer in the events service.
const AssetStatusSaveEventType = "asset-status"

// EventSourceHealthStatus is the top-level discriminator value placed on
// route-execute payloads produced by the healthmonitor publisher. Tells the
// router to read route groups from asset.HealthMonitor.{Offline,Online}RouteGroupIds
// instead of the default asset.RouteGroupIds (which is used for IoT events).
const EventSourceHealthStatus = "healthStatus"

// Event types published to Router.
const EventTypeOffline = "offline"
const EventTypeOnline = "online"

// EventSource identifies the emitter of the nested StandardizedPayload event.
// Placed in event.metadata.source so consumers can tell where the payload
// originated (healthmonitor scan vs js-executor script, etc.).
const EventSource = "healthmonitor"

// AssetHeartbeatStream is the JetStream stream that carries asset
// heartbeat events consumed by the assets/healthmonitor consumer.
// Subject pattern: ${env}.mapexos.asset.heartbeat.{orgId}. Resolved at
// package init from GO_ENV — e.g. "DEV-MAPEXOS-ASSETS-HEARTBEAT".
var AssetHeartbeatStream = config.StreamName("ASSETS", "HEARTBEAT")

// SubjectAssetHeartbeat is the base subject prefix; publishers append the
// orgId at publish time (${env}.mapexos.asset.heartbeat.{orgId}). Used by:
//   - js-executor implicit publish on each data event (when
//     heartbeatMode='implicit' and healthMonitor.enabled=true).
//   - http_gateway explicit publish on POST /api/v1/heartbeat with body
//     { assetUUID } (HTTP-protocol assets in heartbeatMode='explicit').
//
// MQTT-protocol assets in heartbeatMode='explicit' do NOT publish here —
// their liveness comes from NATS broker presence ($SYS.ACCOUNT.*.CONNECT +
// $SYS.ACCOUNT.*.DISCONNECT advisories) and is consumed directly by the
// healthmonitor presence path.
var SubjectAssetHeartbeat = config.Subject("asset", "heartbeat")
