// Package steps holds saga steps that exercise the assets OTA module: the
// operator firmware/plan lifecycle and the device sims (MQTT + HTTP). Both OTA
// journeys (ota_mqtt, ota_http) import these shared building blocks.
package steps

// Bag keys this package writes. Other packages reading these keys import the
// exported constants from here, never a string literal.
const (
	// BagKeyFirmwareID is the firmware artifact id returned by InitFirmware.
	// CompleteFirmware and CreatePlan read it.
	BagKeyFirmwareID = "assets/ota.firmwareId"

	// BagKeyFirmwareUploadURL is the short-TTL presigned PUT URL returned by
	// InitFirmware. UploadFirmwareBinary PUTs the bytes to it.
	BagKeyFirmwareUploadURL = "assets/ota.firmwareUploadUrl"

	// BagKeyFirmwareSHA256 is the sha256 InitFirmware declared for the artifact,
	// stashed so a later step/assert can cross-check without rebuilding the payload.
	BagKeyFirmwareSHA256 = "assets/ota.firmwareSha256"

	// BagKeyFirmwareSize is the byte size InitFirmware declared for the artifact.
	BagKeyFirmwareSize = "assets/ota.firmwareSize"

	// BagKeyFirmwareVersion is the version InitFirmware declared for the artifact.
	BagKeyFirmwareVersion = "assets/ota.firmwareVersion"

	// BagKeyPlanID is the OTA plan id returned by CreatePlan. The asserts poll
	// GET /ota/plans/{id} and /executions with it; Compensate cancels via it.
	BagKeyPlanID = "assets/ota.planId"

	// BagKeyExecutionID is the per-asset execution id the device sim reads from
	// the OTA command and reports its status against.
	BagKeyExecutionID = "assets/ota.executionId"
)

// bagKeyOtaCmdChan holds the buffered channel SubscribeOtaCommand writes the
// received OTA command to and RunMqttOtaDevice reads from. Unexported: it is an
// internal handoff within this package, never read across packages.
const bagKeyOtaCmdChan = "assets/ota.cmdChan"
