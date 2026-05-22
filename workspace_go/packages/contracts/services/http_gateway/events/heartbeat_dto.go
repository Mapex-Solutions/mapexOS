package events

// HeartbeatRequestDTO is the body shape of POST /api/v1/heartbeat?ds={dataSourceId}.
//
// Devices with HealthMonitorConfig.HeartbeatMode='explicit' on HTTP-protocol
// assets POST { "assetUUID": "<v>" } to keep their liveness fresh. orgId and
// pathKey are derived from the resolved DataSource via c.Locals server-side —
// never from the request body — so a compromised body cannot spoof a different
// tenant.
type HeartbeatRequestDTO struct {
	AssetUUID string `json:"assetUUID" validate:"required,min=1"`
}
