package dtos

import (
	hm "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
)

// Module-local aliases over the canonical healthmonitor contracts.
// Handlers and services depend on these names so the call sites do not
// import the contracts package directly.
type (
	AdminAssetUUIDDto           = hm.AdminAssetUUID
	AdminForceOfflineRequestDto = hm.AdminForceOfflineRequest
)
