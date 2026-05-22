package mqttcerts

// Reason enum for revoked rows. asset_deleted is NOT used here —
// asset deletion hard-deletes rows entirely (LGPD-friendly).
const (
	ReasonReplaced   = "replaced"
	ReasonUserAction = "user_action"
)

// EventTypeMqttCertIssued / Revoked are reserved DLQ labels in case the
// service ever publishes its own events; today the existing fanout suffices.
const (
	EventTypeMqttCertIssued  = "mqttcerts.issued"
	EventTypeMqttCertRevoked = "mqttcerts.revoked"
)
