package steps

// Bag keys this package writes. Other packages reading these keys import
// the constants from here.
const (
	// BagKeyAssetID is the asset id returned by CreateAsset. Asserts and
	// the Compensate path read it.
	BagKeyAssetID = "assets.assetID"

	// BagKeyAssetUUID is the public uuid the assets service generates
	// alongside the database id. The healthmonitor presence pipeline
	// looks the asset up by uuid, so saga journeys that exercise
	// telemetry read it from the bag.
	BagKeyAssetUUID = "assets.assetUUID"

	// BagKeyAssetMqttPassword is the plaintext MQTT password the saga
	// supplies on asset create — the operator copy of record. The
	// platform bcrypt-hashes it server-side and never returns it on
	// reads. The connect step presents it as the MQTT password against
	// the Mosquitto auth callout.
	BagKeyAssetMqttPassword = "assets.assetMqttPassword"

	// BagKeyMqttClient holds the live *mqttclient.Client established by
	// ConnectMqtt. Disconnect/Publish steps read it back; type-assert
	// to *mqttclient.Client at the consumer site.
	BagKeyMqttClient = "assets.mqttClient"

	// BagKeyMqttConnectedAt records the timestamp of the last
	// successful Connect, so polling asserts can scope event-search
	// windows to "after we connected" without depending on wall-clock
	// hacks.
	BagKeyMqttConnectedAt = "assets.mqttConnectedAt"

	// BagKeyTelemetrySentAt records the timestamp of the last
	// successful PublishTelemetry, used by the events assert to set
	// a "from" filter so the search ignores events older than the
	// publish.
	BagKeyTelemetrySentAt = "assets.telemetrySentAt"

	// BagKeyAssetDeleted is set true by DeleteAsset so the saga's
	// rollback chain knows the explicit-delete step already removed
	// the asset and CreateAsset.Compensate can no-op without
	// hitting Mongo with a 404 follow-up.
	BagKeyAssetDeleted = "assets.assetDeleted"

	// BagKeyAssetCertPEM holds the PEM-encoded device cert returned
	// by POST /api/v1/mqtt_certs. The ConnectMqttCert step reads it
	// to build the mTLS tls.Certificate before opening the broker
	// connection.
	BagKeyAssetCertPEM = "assets.assetCertPEM"

	// BagKeyAssetKeyPEM holds the PEM-encoded private key returned
	// by POST /api/v1/mqtt_certs. Paired with BagKeyAssetCertPEM.
	BagKeyAssetKeyPEM = "assets.assetKeyPEM"

	// BagKeyAssetCAChainPEM holds the PEM-encoded CA chain so the
	// mTLS client verifies the broker's server cert (issued by the
	// same platform PKI).
	BagKeyAssetCAChainPEM = "assets.assetCAChainPEM"

	// BagKeyAssetCertSerial holds the freshly-issued cert's serial
	// (uppercase hex). Used by asserts that need to correlate the
	// cert with mqttRevokedCertificates rows after a revoke flow.
	BagKeyAssetCertSerial = "assets.assetCertSerial"

	// BagKeyHeartbeatSentAt records the wall-clock timestamp of the
	// last successful POST /api/v1/heartbeat sent by the HTTP-protocol
	// connectivity journey. Asserts that observe online-side events
	// scope their search by this value.
	BagKeyHeartbeatSentAt = "assets.heartbeatSentAt"

	// BagKeyForceOfflineSentAt records the wall-clock timestamp of the
	// last POST /internal/health-monitor/:assetUUID/force-offline call.
	// Offline-action asserts scope their search by this value so they
	// only see executions triggered by the forced transition, never
	// older noise.
	BagKeyForceOfflineSentAt = "assets.forceOfflineSentAt"

	// BagKeyMqttDisconnectedAt records the wall-clock timestamp of the
	// last DisconnectMqtt call. Offline-action asserts use it to scope
	// the events search window to "after we disconnected".
	BagKeyMqttDisconnectedAt = "assets.mqttDisconnectedAt"
)
