package entities

import "time"

// AlertEvent represents a sensor state transition (offline/online) to be routed.
type AlertEvent struct {
	Type             string     `bson:"type"`
	OrgId            string     `bson:"orgId"`
	AssetUUID        string     `bson:"assetUUID"`
	AssetName        string     `bson:"assetName"`
	PathKey          string     `bson:"pathKey"`
	LastSeenAt       *time.Time `bson:"lastSeenAt,omitempty"`
	ThresholdMinutes int        `bson:"thresholdMinutes,omitempty"`
	MissCount        int        `bson:"missCount,omitempty"`
	RouteGroupIds    []string   `bson:"routeGroupIds"`
}
