package healthmonitor

import "time"

// PresenceAdvisory is the NATS payload every edge server publishes on each
// CONNECT and DISCONNECT, onto the shared subject mapexos.presence.advisory.
// This is the mapexOS-internal view of the cross-repo edge contract: the
// producers (the Mosquitto broker plugin for MQTT, the LNS Gateway Server for
// LoRaWAN gateways) marshal the same shape from mapexGoKit/contracts/presence.
// Field names + json tags MUST stay in sync with that canonical edge contract.
//
// Event is "connect" or "disconnect"; Protocol is "mqtt" or "lorawan";
// ReasonCode + ReasonText are populated only on an MQTT disconnect (clean=0,
// keepalive_timeout=4, session_taken_over=142, admin_action=152).
type PresenceAdvisory struct {
	Event      string    `json:"event"`
	Protocol   string    `json:"protocol,omitempty"`
	OrgID      string    `json:"orgId"`
	AssetUUID  string    `json:"assetUUID"`
	ClientID   string    `json:"clientId,omitempty"`
	SourceIP   string    `json:"sourceIp,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	ReasonCode int       `json:"reasonCode,omitempty"`
	ReasonText string    `json:"reasonText,omitempty"`
}
