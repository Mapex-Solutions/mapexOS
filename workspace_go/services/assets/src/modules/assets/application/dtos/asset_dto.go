package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
)

type (
	AssetCreateDTO                  = v1.AssetCreate
	AssetUpdateDTO                  = v1.AssetUpdate
	AssetQueryDTO                   = v1.AssetQuery
	AssetIdDto                      = v1.AssetId   // for params /:assetId (MongoDB _id)
	AssetUUIDDto                    = v1.AssetUUID // for field assetUUID in body (device identifier)
	AssetResponse                   = v1.AssetResponse
	GenerateMqttPasswordResponseDTO = v1.GenerateMqttPasswordResponse
	HealthMonitorConfig             = v1.HealthMonitorConfig
)
