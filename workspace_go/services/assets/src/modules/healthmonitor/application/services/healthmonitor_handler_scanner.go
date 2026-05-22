package services

import (
	"context"
	"fmt"
	"time"

	assetPorts "assets/src/modules/assets/application/ports"
	"assets/src/modules/healthmonitor/application/constants"
	"assets/src/modules/healthmonitor/domain/entities"

	hmContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// fetchActiveOrgs returns the list of orgs that have active sensors.
// Returns (nil, false) after logging the error — callers should bail out.
func (s *HealthMonitorService) fetchActiveOrgs(ctx context.Context) ([]string, bool) {
	orgs, err := s.deps.HealthRepo.GetActiveOrgs(ctx)
	if err != nil {
		logger.Error(err, "[SERVICE:HealthMonitor] [SCAN] failed to get active orgs — loop continues")
		s.deps.Metrics.HealthRedisErrors.WithLabelValues("smembers").Inc()
		return nil, false
	}
	return orgs, true
}

// scanOrg scans a single org for stale sensors using paginated ZRANGEBYSCORE.
func (s *HealthMonitorService) scanOrg(ctx context.Context, orgId string) {
	// Pre-filter: skip assets that heartbeated very recently.
	// Kept minimal (1 min) so the real per-asset threshold in evaluateAsset drives offline detection.
	// Previously hardcoded 10 min — that was the bug causing assets with ThresholdMinutes<10 to never go offline.
	scanCutoffMinutes, _ := config.GetIntValue("health_monitor_scan_cutoff_minutes")
	if scanCutoffMinutes <= 0 {
		scanCutoffMinutes = 1
	}
	minThreshold := time.Duration(scanCutoffMinutes) * time.Minute
	cutoff := time.Now().Add(-minThreshold)
	offset := int64(0)
	totalStale := 0

	logger.Debug(fmt.Sprintf("[SERVICE:HealthMonitor] [SCAN] scanning org: orgId=%s cutoff=%s minThreshold=%s batchSize=%d",
		orgId, cutoff.Format(time.RFC3339), minThreshold, s.batchSize))

	for {
		logger.Debug(fmt.Sprintf("[SERVICE:HealthMonitor] [SCAN] fetching stale page: orgId=%s offset=%d limit=%d",
			orgId, offset, s.batchSize))

		staleUUIDs, err := s.deps.HealthRepo.FindStale(ctx, orgId, cutoff, offset, int64(s.batchSize))
		if err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [SCAN] find-stale failed: orgId=%s offset=%d", orgId, offset))
			s.deps.Metrics.HealthRedisErrors.WithLabelValues("zrangebyscore").Inc()
			break
		}

		if len(staleUUIDs) == 0 {
			logger.Debug(fmt.Sprintf("[SERVICE:HealthMonitor] [SCAN] page empty — done: orgId=%s offset=%d", orgId, offset))
			break
		}

		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [SCAN] stale assets found: orgId=%s offset=%d count=%d uuids=%v",
			orgId, offset, len(staleUUIDs), staleUUIDs))

		for _, assetUUID := range staleUUIDs {
			s.evaluateAsset(ctx, orgId, assetUUID)
		}

		totalStale += len(staleUUIDs)
		s.deps.Metrics.HealthSensorsScanned.Add(float64(len(staleUUIDs)))
		offset += int64(s.batchSize)
	}

	if totalStale == 0 {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [SCAN] orgId=%s stale=0 cutoff=%s → all healthy",
			orgId, cutoff.Format(time.RFC3339)))
	} else {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [SCAN] orgId=%s stale=%d cutoff=%s",
			orgId, totalStale, cutoff.Format(time.RFC3339)))
	}
}

// evaluateAsset checks if a single stale asset should be alerted.
func (s *HealthMonitorService) evaluateAsset(ctx context.Context, orgId string, assetUUID string) {
	asset, err := s.deps.AssetRepo.FindByAssetUUID(ctx, &assetUUID)
	if err != nil || asset == nil {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [EVAL] uuid=%s orgId=%s skip=asset_not_found", assetUUID, orgId))
		return
	}

	if asset.HealthMonitor == nil || !asset.HealthMonitor.Enabled {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [EVAL] uuid=%s name=%s skip=monitoring_disabled", assetUUID, asset.Name))
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [EVAL] evaluating: uuid=%s name=%s orgId=%s threshold=%dm requiredMisses=%d",
		assetUUID, asset.Name, orgId, asset.HealthMonitor.ThresholdMinutes, asset.HealthMonitor.RequiredMisses))

	alerted, _ := s.deps.HealthRepo.IsAlerted(ctx, orgId, assetUUID)
	if alerted {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [EVAL] uuid=%s name=%s skip=already_alerted(offline)", assetUUID, asset.Name))
		return
	}

	lastSeen, _ := s.deps.HealthRepo.GetLastSeen(ctx, orgId, assetUUID)
	threshold := time.Duration(asset.HealthMonitor.ThresholdMinutes) * time.Minute

	if lastSeen != nil && time.Since(*lastSeen) < threshold {
		elapsed := time.Since(*lastSeen).Round(time.Second)
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [EVAL] uuid=%s name=%s skip=within_threshold elapsed=%s threshold=%dm lastSeen=%s",
			assetUUID, asset.Name, elapsed, asset.HealthMonitor.ThresholdMinutes, lastSeen.Format(time.RFC3339)))
		return
	}

	missCount, err := s.deps.HealthRepo.IncrementMiss(ctx, orgId, assetUUID)
	if err != nil {
		s.deps.Metrics.HealthRedisErrors.WithLabelValues("hincrby").Inc()
		return
	}
	s.deps.Metrics.HealthMissIncrements.Inc()

	lastSeenStr := "never"
	if lastSeen != nil {
		lastSeenStr = lastSeen.Format(time.RFC3339)
	}

	if int(missCount) < asset.HealthMonitor.RequiredMisses {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [MISS] uuid=%s name=%s orgId=%s missCount=%d/%d lastSeen=%s",
			assetUUID, asset.Name, orgId, missCount, asset.HealthMonitor.RequiredMisses, lastSeenStr))
		return
	}

	s.deps.Metrics.HealthSensorsStale.Inc()

	if err := s.deps.HealthRepo.MarkAlerted(ctx, orgId, assetUUID); err != nil {
		s.deps.Metrics.HealthRedisErrors.WithLabelValues("sadd").Inc()
	}

	s.handleOfflineTransition(ctx, orgId, asset, lastSeen, int(missCount))
}

// handleOfflineTransition handles the online → offline state transition.
func (s *HealthMonitorService) handleOfflineTransition(ctx context.Context, orgId string, asset *assetPorts.Asset, lastSeen *time.Time, missCount int) {
	assetUUID := asset.AssetUUID

	lastSeenStr := "never"
	if lastSeen != nil {
		lastSeenStr = lastSeen.Format(time.RFC3339)
	}

	if err := s.deps.AssetRepo.UpdateHealthStatusWithChangedAt(ctx, &assetUUID, constants.StatusOffline, time.Now().UTC()); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [OFFLINE→] uuid=%s name=%s orgId=%s failed to update DB status",
			assetUUID, asset.Name, orgId))
		return
	}

	// Always call the publisher — persistence fires unconditionally inside the
	// publisher; the router publish is conditional on RouteGroupIds being set.
	if err := s.deps.AlertPublisher.PublishOffline(ctx, entities.AlertEvent{
		Type:             hmContract.EventTypeOffline,
		OrgId:            orgId,
		AssetUUID:        assetUUID,
		AssetName:        asset.Name,
		PathKey:          asset.PathKey,
		LastSeenAt:       lastSeen,
		ThresholdMinutes: asset.HealthMonitor.ThresholdMinutes,
		MissCount:        missCount,
		RouteGroupIds:    asset.HealthMonitor.OfflineRouteGroupIds,
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [OFFLINE→] uuid=%s name=%s orgId=%s miss=%d/%d lastSeen=%s nats=error routes=%d",
			assetUUID, asset.Name, orgId, missCount, asset.HealthMonitor.RequiredMisses, lastSeenStr, len(asset.HealthMonitor.OfflineRouteGroupIds)))
	} else {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [OFFLINE→] uuid=%s name=%s orgId=%s miss=%d/%d lastSeen=%s nats=sent-persistence routes=%d",
			assetUUID, asset.Name, orgId, missCount, asset.HealthMonitor.RequiredMisses, lastSeenStr, len(asset.HealthMonitor.OfflineRouteGroupIds)))
	}

	s.deps.Metrics.HealthAlertsPublished.WithLabelValues("offline").Inc()
}
