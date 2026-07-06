package permissions

// OTA Permissions
const (
	// OTAPlanList - Permission to list OTA plans
	OTAPlanList = "ota_plans.list"

	// OTAPlanCreate - Permission to create an OTA plan
	OTAPlanCreate = "ota_plans.create"

	// OTAPlanRead - Permission to read an OTA plan (and its executions)
	OTAPlanRead = "ota_plans.read"

	// OTAPlanUpdate - Permission to update an OTA plan
	OTAPlanUpdate = "ota_plans.update"

	// OTAPlanDelete - Permission to cancel/delete an OTA plan
	OTAPlanDelete = "ota_plans.delete"

	// OTAPlanAll - Wildcard permission for all OTA plan operations
	OTAPlanAll = "ota_plans.*"

	// OTAFirmwareUpload - Permission to upload firmware artifacts
	OTAFirmwareUpload = "ota_firmware.upload"
)
