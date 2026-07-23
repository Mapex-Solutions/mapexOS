package mongo

const (
	CollectionName    = "mqtt_revoked_certificates"
	RevokedTTLSeconds = int32(30 * 24 * 60 * 60) // 30 days
)
