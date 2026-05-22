package dtos

import (
	contract "github.com/Mapex-Solutions/MapexOS/contracts/services/events/asset_status"
)

// Application-layer DTO aliases. Definitions live in the contracts package —
// this file is a thin re-export so service code doesn't import contracts
// directly everywhere.
type (
	AssetConnectivityEvent         = contract.AssetConnectivityEvent
	AssetConnectivityHistoryQuery  = contract.AssetConnectivityHistoryQuery
	AssetConnectivityCursorResult  = contract.AssetConnectivityCursorResult
)
