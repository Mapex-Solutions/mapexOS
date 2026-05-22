// Package asset_status holds the cross-service DTO contract for asset
// connectivity history (offline/online transitions).
//
// This package is the single source of truth for the wire shape of
// the asset_status_history feature. It is consumed by:
//   - events MS asset_status module (persistence consumer + query HTTP API)
//   - @mapexos/schemas (Zod schemas mirror these types)
//
// Contracts stay leaf-level — no imports from services/.
package asset_status

import (
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

// AssetConnectivityEvent is a single row in asset_status_history.
// It mirrors the ClickHouse columns produced by the healthmonitor
// alert publisher's persistence payload.
type AssetConnectivityEvent struct {
	Created          time.Time  `json:"created"`
	OrgId            string     `json:"orgId"`
	PathKey          string     `json:"pathKey"`
	AssetUUID        string     `json:"assetUUID"`
	AssetName        string     `json:"assetName,omitempty"`
	EventId          string     `json:"eventId"`
	EventType        string     `json:"eventType"`
	LastSeenAt       *time.Time `json:"lastSeenAt,omitempty"`
	ThresholdMinutes uint16     `json:"thresholdMinutes,omitempty"`
	MissCount        uint16     `json:"missCount,omitempty"`
}

// AssetConnectivityHistoryQuery are the query parameters for
// GET /api/v1/events/connectivity_history (listing across all assets for the
// org) AND GET /api/v1/events/assets/:assetUUID/connectivity_history (single
// asset timeline — assetUUID comes from the path and is ignored here).
//
// Embeds CursorQueryDTO for cursor-based pagination (Cursor, Direction,
// Limit, SortAsc, IncludeChildren).
type AssetConnectivityHistoryQuery struct {
	query.CursorQueryDTO

	// From/To bound the `created` column via BETWEEN.
	From *time.Time `query:"from" validate:"omitempty"`
	To   *time.Time `query:"to" validate:"omitempty"`

	// EventType filters offline vs online transitions.
	EventType *string `query:"eventType" validate:"omitempty,oneof=offline online"`

	// AssetUUID is an OPTIONAL filter used by the list endpoint (no assetUUID
	// in the path). When the path-scoped endpoint is hit the handler injects
	// the path param into this field before delegating to the service, so the
	// query layer is always identical.
	AssetUUID *string `query:"assetUUID" validate:"omitempty"`
}

// AssetConnectivityCursorResult is the cursor-paginated response for
// connectivity history. Shape mirrors EventsRawCursorResult so frontend
// pagination helpers work identically across events list endpoints.
type AssetConnectivityCursorResult struct {
	Items       []AssetConnectivityEvent `json:"items"`
	NextCursor  *time.Time               `json:"nextCursor,omitempty"`
	PrevCursor  *time.Time               `json:"prevCursor,omitempty"`
	HasNext     bool                     `json:"hasNext"`
	HasPrevious bool                     `json:"hasPrevious"`
}
