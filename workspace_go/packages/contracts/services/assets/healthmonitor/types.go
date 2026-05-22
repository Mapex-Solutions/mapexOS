// Package healthmonitor — cross-service payload contracts.
//
// This file holds the NATS payload structs that cross service boundaries
// for the assets/healthmonitor bounded context.
//
// Ownership: assets service (consumer).
// Publishers (polyglot): js-executor (TS), future IoT producers.
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/assets/healthmonitor.
//
// Contracts stay leaf-level — no imports from services/.
package healthmonitor

// HeartbeatEvent is the NATS payload received on the ASSET-HEARTBEAT stream.
// Published by JS Executor (and future producers) after successful event
// normalization, consumed by assets/healthmonitor to maintain the asset
// liveness state.
type HeartbeatEvent struct {
	OrgId     string `json:"orgId"`
	AssetUUID string `json:"assetUUID"`
	PathKey   string `json:"pathKey"`
	Timestamp int64  `json:"ts"`
}
