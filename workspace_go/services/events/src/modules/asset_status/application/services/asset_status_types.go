package services

import (
	"time"

	"events/src/modules/asset_status/application/di"
)

// AssetStatusService orchestrates the asset connectivity history use cases:
// NATS batch persistence (ProcessAssetStatusBatch) and cursor-paginated HTTP
// queries (ListAssetConnectivityHistory).
//
// This file holds ONLY the struct + constructor + compile-time port check.
// Handler methods live in:
//   - asset_status_handler_save.go   (batch consumer handler)
//   - asset_status_handler_query.go  (query handler)
//
// Shared private helpers live in asset_status_helpers.go.
type AssetStatusService struct {
	deps di.AssetStatusServiceDependenciesInjection
}

// persistencePayload mirrors the FLAT shape produced by the healthmonitor
// alert publisher (subject mapexos.events.asset_status_save). Kept here
// (not in the public DTO alias) because it's a wire-only intermediate
// type — the public DTO exposed via HTTP uses the contract alias.
type persistencePayload struct {
	OrgId            string     `json:"orgId"`
	PathKey          string     `json:"pathKey"`
	AssetUUID        string     `json:"assetUUID"`
	AssetName        string     `json:"assetName"`
	EventId          string     `json:"eventId"`
	EventType        string     `json:"eventType"`
	Created          string     `json:"created"`
	LastSeenAt       *time.Time `json:"lastSeenAt,omitempty"`
	ThresholdMinutes uint16     `json:"thresholdMinutes,omitempty"`
	MissCount        uint16     `json:"missCount,omitempty"`
}
