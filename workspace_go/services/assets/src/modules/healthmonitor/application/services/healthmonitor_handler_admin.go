package services

import (
	"context"
	"fmt"

	assetPorts "assets/src/modules/assets/application/ports"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// loadAssetForAdmin fetches the asset by UUID and gates on
// HealthMonitor.IsActive. Admin transitions are addressable by device
// id only — orgId is read from the loaded entity, so the wire never
// has to carry it. Returns (nil, false) on any drop and logs the
// reason at INFO level; admin transitions are explicit operator
// actions, so drops are visible.
func (s *HealthMonitorService) loadAssetForAdmin(ctx context.Context, assetUUID string) (*assetPorts.Asset, bool) {
	asset, err := s.deps.AssetRepo.FindByAssetUUID(ctx, &assetUUID)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:HealthMonitor] [ADMIN] asset lookup failed: uuid=%s err=%v",
			assetUUID, err))
		return nil, false
	}
	if asset == nil {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [ADMIN] drop=asset_not_found uuid=%s", assetUUID))
		return nil, false
	}
	if !asset.HealthMonitor.IsActive() {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [ADMIN] drop=asset_disabled uuid=%s orgId=%s",
			assetUUID, asset.OrgID.Hex()))
		return nil, false
	}
	return asset, true
}

// isAlreadyOffline gates the force-offline transition so repeated calls
// are no-ops once the asset is in the alerted (offline) state. The
// Redis read is best-effort — on error we treat the asset as online
// so the caller continues into the transition path and the publish
// is exercised end-to-end.
func (s *HealthMonitorService) isAlreadyOffline(ctx context.Context, orgId, assetUUID string) bool {
	alerted, _ := s.deps.HealthRepo.IsAlerted(ctx, orgId, assetUUID)
	return alerted
}
