package healthmonitor

import "time"

// PresenceAdvisory is the NATS payload published by the Mosquitto broker
// plugin (mapex-broker-plugin) on every device CONNECT and DISCONNECT.
// One subject only:
//
//	mapexos.mqtt.presence.advisory
//
// The plugin parses the device's MQTT username — encoded as
// "{orgId}:{assetUUID}" by the platform on asset create — and emits the
// advisory pre-resolved, so consumers no longer need a username→asset
// lookup. Event is "connect" or "disconnect"; ReasonCode + ReasonText are
// populated only on disconnect (clean=0, keepalive_timeout=4,
// session_taken_over=142, admin_action=152).
type PresenceAdvisory struct {
	Event      string    `json:"event"`
	OrgID      string    `json:"orgId"`
	AssetUUID  string    `json:"assetUUID"`
	ClientID   string    `json:"clientId,omitempty"`
	SourceIP   string    `json:"sourceIp,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	ReasonCode int       `json:"reasonCode,omitempty"`
	ReasonText string    `json:"reasonText,omitempty"`
}
