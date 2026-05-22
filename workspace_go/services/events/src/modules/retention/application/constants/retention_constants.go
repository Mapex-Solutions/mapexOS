package constants

import "time"

// RetentionCacheTTL defines the TTL for retention policy cache entries.
// Set to 24 hours to balance freshness with performance.
const RetentionCacheTTL = 24 * time.Hour

// CacheKeyPrefix is the prefix used for retention policy cache keys.
// Format: RETENTION_POLICY:{orgId}:{type}
const CacheKeyPrefix = "RETENTION_POLICY"

// DefaultRetentionDays is the fallback retention when no policy is found.
const DefaultRetentionDays uint16 = 1

// asset_status_history bounds. These are platform-level (no orgId) — the
// UI exposes a 7–90 slider; the service enforces [min,max] and applies the
// TTL to ClickHouse on update.
const (
	TableAssetStatusHistory        = "asset_status_history"
	AssetStatusHistoryDefaultDays  = 7
	AssetStatusHistoryMinDays      = 1
	AssetStatusHistoryMaxDays      = 90
	AssetStatusHistoryPolicyName   = "Asset Connectivity History"
)
