package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/events"
	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"
)

type (
	EvenIdentificationDto = v1.EvenIdentification

	// HeartbeatRequestDTO is the body shape of POST /api/v1/heartbeat?ds={dataSourceId}.
	// Devices with HealthMonitorConfig.HeartbeatMode='explicit' on HTTP-protocol
	// assets POST { "assetUUID" } to keep their liveness fresh.
	HeartbeatRequestDTO = v1.HeartbeatRequestDTO

	// OTAStatusRequestDTO is the body shape of POST /api/v1/ota/status?ds={dataSourceId}
	// — a device's OTA progress report, normalized into the OTAStatusAdvisory.
	OTAStatusRequestDTO = otaDtos.OTAStatusRequestDTO
)
