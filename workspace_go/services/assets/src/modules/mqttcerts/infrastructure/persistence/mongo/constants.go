package mongo

const (
	CollectionName    = "mqttRevokedCertificates"
	RevokedTTLSeconds = int32(30 * 24 * 60 * 60) // 30 days
)
