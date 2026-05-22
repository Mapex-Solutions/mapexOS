package services

import (
	ctx "context"

	"assets/src/modules/assets/application/dtos"
	"assets/src/modules/healthmonitor/application/constants"
)

// enrichHealthStatus enriches a single asset response with real-time health data from Redis.
// Called by GetAssetById after building the response from MongoDB.
func (s *AssetService) enrichHealthStatus(c ctx.Context, response *dtos.AssetResponse, orgId string) {
	if response.HealthMonitor == nil || response.HealthMonitor.Enabled == nil || !*response.HealthMonitor.Enabled {
		return
	}

	assetUUID := ""
	if response.AssetUUID != nil {
		assetUUID = *response.AssetUUID
	}
	if assetUUID == "" {
		return
	}

	lastSeen, _ := s.deps.HealthRepo.GetLastSeen(c, orgId, assetUUID)
	isAlerted, _ := s.deps.HealthRepo.IsAlerted(c, orgId, assetUUID)

	response.LastSeenAt = lastSeen

	if isAlerted {
		status := constants.StatusOffline
		response.HealthStatus = &status
	} else if lastSeen != nil {
		status := constants.StatusOnline
		response.HealthStatus = &status
	} else {
		status := constants.StatusUnknown
		response.HealthStatus = &status
	}
}

// enrichHealthStatusBatch enriches a list of asset responses with real-time health data from Redis.
// Uses batch Redis commands (ZMSCORE + SMISMEMBER) for N items in 2 round-trips.
// Called by GetAssets after building the DTO list from MongoDB.
func (s *AssetService) enrichHealthStatusBatch(c ctx.Context, assets []dtos.AssetResponse, orgId string) {
	// Collect assetUUIDs where monitoring is enabled
	uuids := make([]string, 0)
	indexMap := make(map[string][]int) // assetUUID → indices in assets slice

	for i, a := range assets {
		if a.HealthMonitor == nil || a.HealthMonitor.Enabled == nil || !*a.HealthMonitor.Enabled {
			continue
		}
		if a.AssetUUID == nil || *a.AssetUUID == "" {
			continue
		}
		uuid := *a.AssetUUID
		if _, exists := indexMap[uuid]; !exists {
			uuids = append(uuids, uuid)
		}
		indexMap[uuid] = append(indexMap[uuid], i)
	}

	if len(uuids) == 0 {
		return
	}

	// 2 Redis commands for N items (batch)
	lastSeenMap, _ := s.deps.HealthRepo.GetLastSeenBatch(c, orgId, uuids)
	alertedMap, _ := s.deps.HealthRepo.IsAlertedBatch(c, orgId, uuids)

	for _, uuid := range uuids {
		indices := indexMap[uuid]
		for _, i := range indices {
			if ls, ok := lastSeenMap[uuid]; ok {
				assets[i].LastSeenAt = ls
			}

			if alertedMap[uuid] {
				status := constants.StatusOffline
				assets[i].HealthStatus = &status
			} else if lastSeenMap[uuid] != nil {
				status := constants.StatusOnline
				assets[i].HealthStatus = &status
			} else {
				status := constants.StatusUnknown
				assets[i].HealthStatus = &status
			}
		}
	}
}
